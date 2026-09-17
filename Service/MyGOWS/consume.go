package MyGOWS

import (
	"MyGO-IM/DB"
	"encoding/json"

	"log"
	"strconv"
)

// ConsumerWorker 闭包返回worker持续从RabbitMQ拉取消息写入Send通道并等待ACK
func ConsumerWorker() func() error {
	var ID uint
	return func() error {
		workerChan, _, _ := NewChannel()
		//声明收消息的队列
		q, err := workerChan.QueueDeclare(strconv.Itoa(int(ID)), true, false, false, false, nil)
		if err != nil {
			return err
		}
		//将队列绑定到交换机上
		err = workerChan.QueueBind(q.Name, q.Name, "MyGO", false, nil)
		if err != nil {
			return err
		}
		err = workerChan.QueueBind(q.Name, q.Name, "broadcast", false, nil)
		if err != nil {
			return err
		}
		//消费就绪
		msgs, err := workerChan.Consume(q.Name, "", false, false, false, false, nil)
		if err != nil {
			log.Println(" c.MQCh.Consume:", err)
		}
		var data ServerMessage
		for {
			msg := <-msgs
			json.Unmarshal(msg.Body, &data)
			if !DB.Dedup(ID, data.MsgID) {
				log.Printf("repeat")
				goto loop
			}
			switch data.Method {
			case "private", "system":
				H.mu.Lock()
				cli, ok := H.Clients[data.ToID]
				H.mu.Unlock()
				if ok {
					cli.Send <- msg.Body
				}
				H.mu.Lock()
				cli, ok = H.Clients[data.FromID]
				H.mu.Unlock()
				if ok {
					cli.Send <- msg.Body
				}
			case "group":
				var result []DB.User
				H.mu.Lock()
				cli, ok := H.Clients[data.FromID]
				H.mu.Unlock()
				if ok {
					cli.Send <- msg.Body
				}
				err = DB.MySQL.
					Raw("select a.id from users a,`groups` b,group_users c where  b.id=? and b.id=c.group_id and a.id=c.user_id ",
						data.ToID).
					Find(&result).Error
				for _, user := range result {
					H.mu.Lock()
					cli, ok = H.Clients[user.ID]
					H.mu.Unlock()
					if ok {
						cli.Send <- msg.Body
					}
				}
			}
		loop:
			msg.Ack(false)
		}
	}
}
