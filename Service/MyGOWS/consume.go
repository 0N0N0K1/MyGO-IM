package MyGOWS

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func (c *Client) ConsumeMyQueue() error {
	//消费就绪
	msgs, err := c.MQCh.Consume(c.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println(" c.MQCh.Consume:", err)
		return err
	}
	var msg amqp091.Delivery
	for {
		select {
		//消费消息,转发到客户端
		case msg = <-msgs:
			c.Send <- msg.Body
			msg.Ack(false)
		//监听客户端连接是否关闭
		case <-c.Close:
			return nil
		//返回NACK 处理失败，重新入队
		default:
			if msg.DeliveryTag == 0 {
				msg.Nack(false, true)
			}
			continue
		}
	}
}
