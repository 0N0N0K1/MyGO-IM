package MyGOWS

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
)

func (c *Client) ConsumeMyQueue() {
	//消费就绪
	msgs, err := c.MQCh.Consume(strconv.Itoa(int(c.ID)), "", false, false, false, false, nil)
	if err != nil {
		log.Println(" c.MQCh.Consume:", err)
	}
	var msg amqp091.Delivery
	for {
		select {
		//监听客户端连接是否关闭
		case <-c.Close:
			return
		//消费消息,转发到客户端
		case msg = <-msgs:
			log.Printf("写入send消息%v", string(msg.Body))
			c.Send <- msg.Body
			ok := <-c.AckReady
			log.Printf("ackready %v", ok)
			if ok {
				msg.Ack(false)
			} else {
				msg.Nack(false, true)
			}
		//返回NACK 处理失败，重新入队
		default:

		}
	}
}
