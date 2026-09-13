package MyGOWS

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
)

// PublishPrivate 私聊 Publish 处理函数
func (c *Client) PublishPrivate(msg *Message) error {
	data, _ := json.Marshal(msg)
	err := c.MQCh.Publish("private", strconv.Itoa(int(msg.ToID)), false, false,
		amqp091.Publishing{
			Body: []byte(data),
		})
	if err != nil {
		return err
	}
	return nil
}

// PublishGroup 群聊 Publish 处理函数
func (c *Client) PublishGroup(msg *Message) error {
	data, _ := json.Marshal(msg)
	err := c.MQCh.Publish("group", strconv.Itoa(int(msg.ToID)), false, false,
		amqp091.Publishing{
			Body: data,
		})
	if err != nil {
		return err
	}
	return nil
}

// PublishSystem 系统通知 Publish 处理函数
func (s *SystemMQChan) PublishSystem(msg *Message) error {
	data, _ := json.Marshal(&msg)
	err := s.MQCh.Publish("system", strconv.Itoa(int(msg.ToID)), false, false,
		amqp091.Publishing{
			Body: data,
		})
	if err != nil {
		return err
	}
	return nil
}
