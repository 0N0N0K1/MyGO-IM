package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"

	"log"
	"time"
)

var H *Hub
var SnowID *snowflake.Node

// InitHub 初始化并启动 Hub
func InitHub() {
	H = NewHub()
	H.ClientPool.New = func() any {
		H.mu.Lock()
		defer H.mu.Unlock()
		c := new(Client)
		return c
	}
	H.MsgPool.New = func() any {
		H.mu.Lock()
		defer H.mu.Unlock()
		c := new(Message)
		return c
	}
	SnowID, _ = snowflake.NewNode(1)
	go H.Run()
	log.Println("WS集中管理器Hub创建成功")
}

// NewHub 返回一个全局 Hub 实例
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message, 256),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		//上线
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			OnlineHandler(client)
			log.Printf("用户%s上线", client.Name)

		//下线
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				OfflineHandler(client)
			}
			h.mu.Unlock()
			log.Printf("用户%s下线", client.Name)
		//广播
		case msg := <-h.Broadcast:
			data, _ := json.Marshal(msg)
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- data:
				default:
					// 发送缓冲区满，关闭连接
					close(client.Send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ReadPump 持续从客户端读取消息
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
		if c.ReadHandler(message) {
			select {
			case <-c.Close:
				return
			case <-c.SendReady:
			}
		} else {
			select {
			case <-c.Close:
				return
			default:

			}
		}
	}
}

// WritePump 持续向客户端写消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(27 * time.Second) // 心跳间隔
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case <-c.Close:
			return
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.WriteHandler(message)
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// OfflineHandler 下线后的处理
func OfflineHandler(client *Client) {
	close(client.Close)                   //关闭通道，通知子协程死亡
	client.MQCh.Close()                   //关闭AMQP信道，防止出现僵尸消费者
	DB.SetOffline(client.Name)            //删除Redis缓存中在线记录
	DB.SetLastOnline(client.ID)           //设置下线时间
	delete(client.Hub.Clients, client.ID) //从Hub中删除对应Client连接
	close(client.Send)                    //关闭发送消息发送通道
	H.ClientPool.Put(client)              //将client实例放回池
}

// OnlineHandler 注册上线后的处理
func OnlineHandler(client *Client) {
	//Redis进行缓存上线记录
	DB.SetOnline(client.Name)
	//开辟 读+写goroutine、消费消息的goroutine
	go client.WritePump()
	go client.ReadPump()
	go client.ConsumeMyQueue()
	go client.ConfirmHandler()
}

// ReadHandler 从客户端收到消息后的处理
func (c *Client) ReadHandler(message []byte) (sent bool) {
	msg := c.Hub.MsgPool.Get().(*Message)
	defer c.Hub.MsgPool.Put(msg)
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Println(err)
		return
	}
	msg.MsgID = int64(SnowID.Generate())
	msg.FromName = c.Name
	msg.FromID = c.ID
	msg.Timestamp = time.Now().Unix()
	DB.MsgNumPlus()
	return c.SendFunc(msg)
}

// ConfirmHandler 处理RabbitMQ的confirm
func (c *Client) ConfirmHandler() {
	for ack := range c.Confirm {

		log.Printf("收到confirm tag: %d ack: %v", ack.DeliveryTag, ack.Ack)

		if ack.Ack {
			cmd := DB.RDB.Get(context.TODO(), Utils.CachePublishMsgName(ack.DeliveryTag, c.ID))
			result, err := cmd.Result()
			if err != nil {
				return
			}
			DB.RDB.Del(context.TODO(), Utils.CachePublishMsgName(ack.DeliveryTag, c.ID))
			var msg Message
			err = json.Unmarshal([]byte(result), &msg)
			//持久化到MySQL
			err = DB.InsertMsg(msg.FromID, msg.ToID, msg.Seq, msg.FromName, msg.ToName, msg.Method, result)
			//发送者writeSeq+1
			switch msg.Method {
			case "group":
				DB.IncrGroupWriteSeq(msg.ToID, msg.Seq+1)
			case "private":
				DB.IncrPrivateWriteSeq(msg.FromID, msg.ToID, msg.Seq+1)
			}

			c.SendReady <- struct{}{}
			log.Printf("tag: %d  %v 已存入mysql", ack.DeliveryTag, msg)

		} else {
			cmd := DB.RDB.Get(context.TODO(), Utils.CachePublishMsgName(ack.DeliveryTag, c.ID))
			result, err := cmd.Result()
			if err != nil {
				return
			}
			DB.RDB.Del(context.TODO(), Utils.CachePublishMsgName(ack.DeliveryTag, c.ID))
			var msg Message
			_ = json.Unmarshal([]byte(result), &msg)
			c.Hub.mu.RLock()
			target, ok := c.Hub.Clients[msg.FromID]
			c.Hub.mu.RUnlock()

			data, _ := json.Marshal(NewSystemMsg(msg.Payload+"发送失败", "system", msg.FromName, 0, msg.FromID))
			if ok {
				select {
				case target.Send <- data:
				//监听客户端连接是否关闭
				case <-c.Close:
					return
				default:
					close(target.Send)
				}
			}
		}
	}
}

func (c *Client) WriteHandler(message []byte) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		c.AckReady <- false
		return
	}
	log.Printf("write begin %v", msg)
	switch msg.Method {
	case "system":
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
	default:
		if !Utils.Dedup(msg.ToID, msg.MsgID, msg.Seq) {
			log.Printf("repeat")
			c.AckReady <- true
			return
		}
		switch msg.Method {
		case "group":
			readSeq := DB.GetGroupReadSeq(msg.ToID, msg.FromID)
			if readSeq != msg.Seq {
				c.AckReady <- false
				return
			}
			err = c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				c.AckReady <- false
				return
			}
			DB.IncrGroupReadSeq(msg.ToID, c.ID, DB.GetGroupWriterSeq(msg.ToID)+1)
		case "private":
			_, readSeq := DB.GetPrivateSeq(msg.ToID, msg.FromID)
			if readSeq != msg.Seq {
				c.AckReady <- false
				return
			}
			err = c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				c.AckReady <- false
				return
			}
			DB.IncrPrivateReadSeq(msg.ToID, msg.FromID, msg.Seq+1)
		}
		c.AckReady <- true

	}

}
