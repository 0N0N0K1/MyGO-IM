package MyGOWS

import (
	"MyGO-IM/Conf"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

var Conn *amqp091.Connection

func InitMQ() {
	var err error
	Conn, err = amqp091.Dial(Conf.MQURL)
	if err != nil {
		log.Fatal("rabbitMQ连接失败")
	}
	log.Println("rabbitMQ连接成功")
	err = SystemMQChanInit()
	if err != nil {
		log.Fatal("systemMQCh创建失败")
	}
	//声明不同种类消息对应的交换机
	err = SystemMQ.MQCh.ExchangeDeclare("private", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("private交换机创建失败")
	}
	err = SystemMQ.MQCh.ExchangeDeclare("group", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("group交换机创建失败")
	}
	err = SystemMQ.MQCh.ExchangeDeclare("system", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("system交换机创建失败")
	}
}
func NewChannel() (*amqp091.Channel, error) {
	channel, err := Conn.Channel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}
