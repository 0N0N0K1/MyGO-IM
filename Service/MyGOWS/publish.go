package MyGOWS

import (
	"MyGO-IM/Conf"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Ch   *amqp091.Channel
	Cfm  chan amqp091.Confirmation
	ID   int
	Stop chan struct{}
}

type ProducerChan struct {
	Chan chan *Producer
}

var ProducerPool ProducerChan

// ProducerPoolChanInit 初始化系统使用的 AMQP 通道
func ProducerPoolChanInit() {
	ProducerPool.Chan = make(chan *Producer, 32)
	var producer Producer
	var err error
	for i := 0; i < Conf.Conf.Worker.Producer; i++ {
		producer.Ch, producer.Cfm, err = NewChannel()
		producer.ID = i
		producer.Stop = make(chan struct{})
		if err != nil {
			i--
			continue
		}
		ProducerPool.Chan <- &producer
	}
	return
}
func (pool *ProducerChan) Put(producer *Producer) {
	if !producer.Ch.IsClosed() {
		producer.Ch, producer.Cfm, _ = NewChannel()
	}
	ProducerPool.Chan <- producer
}
func (pool *ProducerChan) Get() (producer *Producer) {
	producer = <-ProducerPool.Chan
	if !producer.Ch.IsClosed() {
		producer.Ch, producer.Cfm, _ = NewChannel()
	}
	go producer.ConfirmHandler()
	return producer
}

// PublishHandler 服务端 Publish 处理函数
func (c *Producer) PublishHandler(msg any, conversationID string) error {
	data, _ := json.Marshal(&msg)
	err := c.Ch.Publish("MyGO", conversationID, false, false,
		amqp091.Publishing{
			Body: data,
		})
	if err != nil {
		return err
	}
	return nil
}
