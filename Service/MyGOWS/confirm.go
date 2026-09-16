package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"context"
	"strconv"
)

// ConfirmHandler 异步处理发送消息后RabbitMQ的confirm
func (c *Client) ConfirmHandler() {
	for ack := range c.Confirm {
		cmd := DB.RDB.Get(context.TODO(), Utils.UsersIdSentTag(ack.DeliveryTag, c.ID))
		result, err := cmd.Result()
		if err != nil {
			return
		}
		replyID, _ := strconv.Atoi(result)
		if !ack.Ack {
			ACK := NewServerACK("nack", "broker broke", int64(replyID), false)
			SystemMQ.PublishServerACK(ACK, c.ID)
		}
		DB.RDB.Del(context.TODO(), Utils.UsersIdSentTag(ack.DeliveryTag, c.ID))

	}
}
