package MyGOWS

import (
	"MyGO-IM/Conf"
	"MyGO-IM/DB"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

// Upgrader HTTP -> WS 升级器
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  Conf.Conf.W.ReadBufferSize,
	WriteBufferSize: Conf.Conf.W.WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应校验域名
	},
}

// InitChat 升级为 WS 后负责注册等初始化操作
func InitChat(userName string, userID uint, conn *websocket.Conn) error {

	//创建AMQP信道
	ch, confirms, err := NewChannel()
	if err != nil {
		return err
	}
	//声明收消息的队列
	q, err := ch.QueueDeclare(strconv.Itoa(int(userID)), true, false, false, false, nil)
	if err != nil {
		return err
	}
	//将队列绑定到交换机上
	err = ch.QueueBind(q.Name, strconv.Itoa(int(userID)), "private", false, nil)
	if err != nil {
		return err
	}
	groups := DB.QueryMyGroup(userID)
	for _, group := range groups {
		if group.ID != 0 {
			err = ch.QueueBind(q.Name, strconv.Itoa(int(group.ID)), "group", false, nil)
			if err != nil {
				return err
			}
		}
	}
	err = ch.QueueBind(q.Name, strconv.Itoa(int(userID)), "system", false, nil)
	if err != nil {
		return err
	}

	// 从对象池中Get一个Client实例注册，加入Hub管理
	client := H.ClientPool.Get().(*Client)
	client.Close = make(chan struct{})
	client.ID = userID
	client.Name = userName
	client.MQCh = ch
	client.Confirm = confirms
	client.Queue = q
	client.Hub = H
	client.Conn = conn
	client.Send = make(chan []byte, 256)
	client.Hub.Register <- client
	client.SendFunc = client.MsgSender()
	return nil
}
