package MyGOWS

import (
	"MyGO-IM/DB"
	"encoding/json"
	"fmt"
	"log"
)

func (c *Client) SendPrivate(msg Message) {
	var frd []DB.User
	// 判断是否为好友
	userTo, _ := DB.QueryUser(msg.ToID, DB.ByID)
	if len(userTo) != 0 {
		frd, _ = DB.QueryFrd(msg.FromID, userTo[0].ID, DB.ByID)
	}
	// 对方不是好友返回系统消息
	if len(frd) == 0 {
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
		return
	}
	// Publish到私聊交换机
	err := c.PublishPrivate(msg)
	log.Println("published")
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
		return
	}
	//持久化到MySQL
	data, _ := json.Marshal(msg)
	err = DB.InsertMsg(msg.FromID, msg.ToID, msg.FromName, msg.ToName, msg.Type, string(data))
	if err != nil {
		//todo 持久化失败原因及处理
		return
	}
}

func (c *Client) SendGroup(msg Message) {
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
		return
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
		return
	}
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
		return
	}

	data, _ := json.Marshal(msg)
	err = DB.InsertMsg(msg.FromID, msg.ToID, msg.FromName, msg.ToName, msg.Type, string(data))
	if err != nil {
		//todo 持久化失败原因及处理
		return
	}
}
