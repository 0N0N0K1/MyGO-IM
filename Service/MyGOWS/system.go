package MyGOWS

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type SystemMQChan struct {
	MQCh *amqp091.Channel
}

var SystemMQ SystemMQChan

func SystemMQChanInit() (err error) {
	SystemMQ.MQCh, err = NewChannel()
	if err != nil {
		return err
	}
	return nil
}
func NewSystemMsg(content, fromName, toName string, fromID, toID uint) *Message {
	return &Message{
		Type:      "system",
		FromName:  fromName,
		FromID:    fromID,
		ToID:      toID,
		ToName:    toName,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
}

func SendFrdStatus(status, fromName, toName string, fromID, toID uint) {
	var msg *Message
	switch status {
	case "pending":
		content := fmt.Sprintf("用户: %s  ID: %d 申请成为你的好友", fromName, fromID)
		msg = NewSystemMsg(content, fromName, toName, fromID, toID)
	case "reject":
		content := fmt.Sprintf("用户: %s  ID: %d 拒绝了你的好友申请", fromName, fromID)
		msg = NewSystemMsg(content, fromName, toName, fromID, toID)
	case "accept":
		content := fmt.Sprintf("用户: %s  ID: %d 接收了你的好友申请", fromName, fromID)
		msg = NewSystemMsg(content, fromName, toName, fromID, toID)
	}
	err := SystemMQ.PublishSystem(msg)
	if err != nil {
		log.Println(err)
		return
	}

}

func SendGrpStatus(status, fromName, toName, groupName string, fromID, toID, groupID uint) {
	var msg *Message
	switch status {
	case "pending":
		content := fmt.Sprintf("用户: %s  ID: %d 申请进入群聊", fromName, fromID)
		msg = NewSystemMsg(content, fromName, toName, fromID, toID)
	case "reject":
		content := fmt.Sprintf("用户: %s  (%d) 拒绝了你加入群聊 %s  (%d)", fromName, fromID, groupName, groupID)
		msg = NewSystemMsg(content, fromName, toName, fromID, toID)
	case "accept":
		content := fmt.Sprintf("用户: %s  (%d) 同意了你加入群聊 %s  (%d)", fromName, fromID, groupName, groupID)
		msg = NewSystemMsg(content, fromName, toName, fromID, toID)
	}
	err := SystemMQ.PublishSystem(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
