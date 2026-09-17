package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

// ReadPump 持续从客户端读消息
func (c *Client) ReadPump() {
	// defer 下线注销通知
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()
	c.Conn.SetReadLimit(512 * 1024) // 最大消息 512KB
	// 设置心跳检测
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		return nil
	})
	// 循环读取客户端消息
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
		}
		c.ReadHandler(message)
		select {
		case <-c.Close:
			return
		default:
		}
	}
}

// ReadHandler 从客户端收到消息的处理函数
func (c *Client) ReadHandler(message []byte) {
	clientMsg := c.Hub.ClientMsgPool.Get().(*ClientMessage)
	serverMsg := c.Hub.ServerMsgPool.Get().(*ServerMessage)
	defer c.Hub.ClientMsgPool.Put(clientMsg)
	defer c.Hub.ServerMsgPool.Put(serverMsg)
	if err := json.Unmarshal(message, &clientMsg); err != nil {
		log.Println("Unmarshal error " + err.Error())
		return
	}

	log.Println("ReadHandler 开始")
	switch clientMsg.Cmd {
	case "send":
		serverMsg.MsgID = int64(Utils.SnowID.Generate())
		serverMsg.FromID = c.ID
		serverMsg.ToID = clientMsg.ToId
		serverMsg.Type = clientMsg.Type
		serverMsg.Method = clientMsg.Method
		serverMsg.Payload = clientMsg.Payload
		serverMsg.Cmd = "send"
		serverMsg.ReplyID = clientMsg.MsgId
		serverMsg.ConversationId = clientMsg.ConversationId
		serverMsg.Timestamp = time.Now().Unix()
		DB.MsgNumPlus()
		c.SendFunc(serverMsg)
	case "ack":
	case "nack":
	}

	log.Println("ReadHandler 完成")
	return
}

// MsgSender 闭包返回 SendFunc ,用于发送消息，并保持 DeliveryTag 与 ServerMessage 的对应关系
func (c *Client) MsgSender() func(msg *ServerMessage) {
	var deliveryTag uint64 = 1
	return func(msg *ServerMessage) {
		cmd := DB.RDB.HGet(context.TODO(), DB.UsersIdSent(c.ID), strconv.FormatInt(msg.ReplyID, 10))
		if cmd.Err() == nil {
			return
		}
		// 根据消息类型路由
		switch msg.Method {
		case private:
			log.Println("SendFunc 开始 private")
			if c.SendPrivate(msg, deliveryTag) {
				deliveryTag++
			}
			log.Println("SendFunc 成功")
			return

		case group:
			if c.SendGroup(msg, deliveryTag) {
				deliveryTag++
				return
			}
		case system:
		}
	}
}

func (c *Client) SendPrivate(msg *ServerMessage, deliveryTag uint64) bool {
	producer := ProducerPool.Get()
	defer ProducerPool.Put(producer)
	// 判断是否为好友,拿到发送的Seq
	_, _, globalSeq := DB.GetPrivateSeq(msg.FromID, msg.ToID)
	if globalSeq != 0 {
		msg.Seq = globalSeq
	} else {
		ACK := NewSystemMsg("nack", "对方不是你的好友", msg.ReplyID, msg.FromID, false, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	marshal, _ := json.Marshal(&msg)

	//todo 事务保障两步原子性
	// 消息入库	// writeSeq+1
	err := DB.InsertMsg(msg.FromID, msg.ToID, globalSeq, msg.ConversationId, msg.Method, string(marshal))
	if err != nil {
		ACK := NewSystemMsg("nack", "持久化失败", msg.ReplyID, msg.FromID, true, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	DB.IncrPrivateGlobalWriteSeq(msg.FromID, msg.ToID, globalSeq+1)
	DB.IncrPrivateReadSeq(msg.FromID, msg.ToID, msg.Seq)
	// 在线 Publish 到私聊交换机
	err = producer.PublishHandler(msg, msg.ConversationId)
	// Publish失败通知
	if err != nil {
		ACK := NewSystemMsg("nack", "publish error", msg.ReplyID, msg.FromID, false, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	data, _ := json.Marshal(msg)
	DB.RDB.HSet(context.TODO(), DB.UsersIdSent(c.ID), msg.ReplyID, msg.ReplyID)
	DB.RDB.HSet(context.TODO(), DB.UsersIdSentTag(deliveryTag, producer.ID), msg.ReplyID, data)

	// ack
	ACK := NewSystemMsg("ack", "ok", msg.ReplyID, msg.FromID, false, nil)
	producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
	return true

}

// SendGroup publish到群聊交换机
func (c *Client) SendGroup(msg *ServerMessage, deliveryTag uint64) bool {
	producer := ProducerPool.Get()
	defer ProducerPool.Put(producer)
	// 判断发送者是否为群成员
	if DB.QueryMember(msg.FromID, msg.ToID).ID == 0 {
		ACK := NewSystemMsg("nack", "你不是群成员", msg.ReplyID, msg.FromID, false, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	// 判断发送者是否被禁言
	ban, t := DB.QueryIfSilence(msg.FromID, msg.ToID)
	if ban {
		content := fmt.Sprintf("你已被禁言，剩余解禁时间%v", t)
		ACK := NewSystemMsg("nack", content, msg.ReplyID, msg.FromID, false, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	marshal, _ := json.Marshal(&msg)

	//todo 事务保障三步原子性
	// 拿到群的WriteSeq  // 消息入库	// writeSeq+1
	msg.Seq = DB.GetGroupWriterSeq(msg.ToID)
	err := DB.InsertMsg(msg.FromID, msg.ToID, msg.Seq, msg.ConversationId, msg.Method, string(marshal))
	if err != nil {
		ACK := NewSystemMsg("nack", "持久化失败", msg.ReplyID, msg.FromID, true, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	DB.IncrGroupWriteSeq(msg.ToID, msg.Seq+1)
	DB.IncrGroupReadSeq(msg.ToID, msg.FromID, msg.Seq)
	// Publish到群聊交换机
	err = producer.PublishHandler(msg, msg.ConversationId)
	// Publish失败通知
	if err != nil {
		ACK := NewSystemMsg("nack", "publish error", msg.ReplyID, msg.FromID, false, nil)
		producer.PublishHandler(ACK, strconv.Itoa(int(msg.FromID)))
		return false
	}
	// 缓存到Redis
	data, _ := json.Marshal(msg)
	DB.RDB.HSet(context.TODO(), DB.UsersIdSent(c.ID), msg.ReplyID, msg.MsgID)
	DB.RDB.HSet(context.TODO(), DB.UsersIdSentTag(deliveryTag, producer.ID), msg.ReplyID, data)
	return true
}
