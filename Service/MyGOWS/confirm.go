package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"context"
	"encoding/json"
)

// ConfirmHandler 异步处理发送消息后RabbitMQ的confirm
func (c *Producer) ConfirmHandler() {
	for {
		select {
		case ack := <-c.Cfm:
			cmd := DB.RDB.Get(context.TODO(), Utils.UsersIdSentTag(ack.DeliveryTag, c.ID))
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
			DB.RDB.Del(context.TODO(), Utils.UsersIdSentTag(ack.DeliveryTag, c.ID))
		case <-c.Stop:
			return
		}

	}
}
