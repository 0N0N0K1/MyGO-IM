package MyGOWS

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
)

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
func (c *Client) PublishGroup(msg *Message) error {
	data, _ := json.Marshal(msg)
	err := c.MQCh.Publish("group", strconv.Itoa(int(msg.ToID)), false, false,
		amqp091.Publishing{
			Body: data,
		})
	if err != nil {
		return err
	}
	fmt.Println("publish成功")
	return nil
}

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
