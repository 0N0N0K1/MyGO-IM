package MyGOWS

import (
	"MyGO-IM/DB"
	"context"
	"encoding/json"
)

// ConfirmHandler 异步处理发送消息后RabbitMQ的confirm
func (c *Producer) ConfirmHandler() {
	for {
		select {
		case ack := <-c.Cfm:
			if ack.DeliveryTag == 0 {
				return
			}
			cmd := DB.RDB.Get(context.TODO(), DB.UsersIdSentTag(ack.DeliveryTag, c.ID))
			result, err := cmd.Result()
			if err != nil {
				return
			}
			var data ServerMessage
			json.Unmarshal([]byte(result), &data)
			if !ack.Ack {
				H.mu.Lock()
				cli, ok := H.Clients[data.FromID]
				H.mu.Unlock()
				if ok {
					cli.Send <- []byte(result)
				}

			}
			DB.RDB.Del(context.TODO(), DB.UsersIdSentTag(ack.DeliveryTag, c.ID))
		}

	}
}
