package MyGOWS

import (
	"MyGO-IM/Conf"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

var Conn *amqp091.Connection

// InitMQ RabbitMQ连接与交换机/队列初始化
func InitMQ() {
	var err error
	Conn, err = amqp091.Dial(Conf.MQURL)
	if err != nil {
		log.Fatal("rabbitMQ连接失败")
	}
	log.Println("rabbitMQ连接成功")
	ProducerPoolChanInit()
	producer := ProducerPool.Get()
	defer ProducerPool.Put(producer)
	//声明不同种类消息对应的交换机
	err = producer.Ch.ExchangeDeclare("MyGO", "x-modulus-hash", true, false, false, false, nil)
	if err != nil {
		log.Fatal("x-modulus-hash交换机创建失败")
	}
	err = producer.Ch.ExchangeDeclare("broadcast", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal("x-modulus-hash交换机创建失败")
	}
}

// NewChannel 返回带confirm机制的 AMQP 通道
func NewChannel() (*amqp091.Channel, chan amqp091.Confirmation, error) {
	channel, err := Conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	err = channel.Confirm(false)
	confirms := channel.NotifyPublish(make(chan amqp091.Confirmation, 100))
	return channel, confirms, nil
}
