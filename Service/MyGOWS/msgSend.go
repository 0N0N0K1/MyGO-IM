package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (c *Client) MsgSender() func(msg *Message) bool {
	var deliveryTag uint64 = 1
	return func(msg *Message) bool {
		log.Printf("收到一条消息: %v", msg)
		// 根据消息类型路由
		switch msg.Method {
		case private:
			if c.SendPrivate(msg, deliveryTag) {
				deliveryTag++
				return true
			}
		case group:
			if c.SendGroup(msg, deliveryTag) {
				deliveryTag++
				return true
			}
		case system:
		}
		return false
	}
}
func (c *Client) SendPrivate(msg *Message, deliveryTag uint64) bool {

	// 判断是否为好友,拿到发送的Seq
	_, writeSeq := DB.GetPrivateSeq(msg.FromID, msg.ToID)
	if writeSeq != 0 {
		msg.Seq = writeSeq
	} else {
		c.Hub.mu.RLock()
		target, ok := c.Hub.Clients[msg.FromID]
		c.Hub.mu.RUnlock()
		data, _ := json.Marshal(NewSystemMsg("对方不是你的好友", "system", msg.FromName, 0, msg.FromID))
		if ok {
			select {
			case target.Send <- data:
			default:
				close(target.Send)
			}
		}
		return false
	}

	// Publish到私聊交换机
	err := c.PublishPrivate(msg)

	log.Printf("deliveryTag:%d is published", deliveryTag)

	// Publish失败通知
	if err != nil {
		c.Hub.mu.RLock()
		target, ok := c.Hub.Clients[msg.FromID]
		c.Hub.mu.RUnlock()
		data, _ := json.Marshal(NewSystemMsg("发送失败", "system", msg.FromName, 0, msg.FromID))
		if ok {
			select {
			case target.Send <- data:
			default:
				close(target.Send)
			}
		}
		return false
	}
	// 缓存到Redis
	data, _ := json.Marshal(msg)
	DB.RDB.Set(context.TODO(), Utils.CachePublishMsgName(deliveryTag, c.ID), data, time.Minute)

	log.Printf("deliveryTag:%d is cachaed", deliveryTag)

	return true

}

func (c *Client) SendGroup(msg *Message, deliveryTag uint64) bool {

	// 判断发送者是否为群成员
	if DB.QueryMember(msg.FromID, msg.ToID).ID == 0 {
		c.Hub.mu.RLock()
		target, ok := c.Hub.Clients[msg.FromID]
		c.Hub.mu.RUnlock()
		data, _ := json.Marshal(NewSystemMsg("你不是该群成员", "system", msg.FromName, 0, msg.FromID))
		if ok {
			select {
			case target.Send <- data:
			default:
				close(target.Send)
			}
		}
		return false
	}
	// 判断发送者是否被禁言
	ban, t := DB.QueryIfSilence(msg.FromID, msg.ToID)
	if ban {
		c.Hub.mu.RLock()
		target, ok := c.Hub.Clients[msg.FromID]
		c.Hub.mu.RUnlock()
		content := fmt.Sprintf("你已被禁言，剩余解禁时间%v", t)
		data, _ := json.Marshal(NewSystemMsg(content, "system", msg.FromName, 0, msg.FromID))
		if ok {
			select {
			case target.Send <- data:
			default:
				close(target.Send)
			}
		}
		return false
	}
	// 拿到群的WriteSeq
	msg.Seq = DB.GetGroupWriterSeq(msg.ToID)
	// Publish到群聊交换机
	err := c.PublishGroup(msg)
	// Publish失败通知
	if err != nil {
		c.Hub.mu.RLock()
		target, ok := c.Hub.Clients[msg.FromID]
		c.Hub.mu.RUnlock()
		data, _ := json.Marshal(NewSystemMsg("发送失败", "system", msg.FromName, 0, msg.FromID))
		if ok {
			select {
			case target.Send <- data:
			default:
				close(target.Send)
			}
		}
		return false
	}
	// 缓存到Redis
	data, _ := json.Marshal(msg)
	DB.RDB.Set(context.TODO(), Utils.CachePublishMsgName(deliveryTag, c.ID), data, time.Minute)

	return true

}
