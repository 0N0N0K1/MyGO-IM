package MyGOWS

import (
	"MyGO-IM/Conf"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
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

// RegisterNewBind 用户注册后创建自己的队列与绑定交换机
func RegisterNewBind(userID uint) error {

	//声明收消息的队列
	q, err := SystemMQ.MQCh.QueueDeclare(strconv.Itoa(int(userID)), true, false, false, false, nil)
	if err != nil {
		return err
	}
	//将队列绑定到交换机上
	err = SystemMQ.MQCh.QueueBind(q.Name, strconv.Itoa(int(userID)), "private", false, nil)
	if err != nil {
		return err
	}

	err = SystemMQ.MQCh.QueueBind(q.Name, strconv.Itoa(int(userID)), "system", false, nil)
	if err != nil {
		return err
	}
	return nil
}

// GroupNewBind 用户成功加入群聊后绑定对应群聊
func GroupNewBind(userID, groupID uint) error {

	err := SystemMQ.MQCh.QueueBind(strconv.Itoa(int(userID)), strconv.Itoa(int(groupID)), "group", false, nil)
	if err != nil {
		return err
	}
	return nil
}
