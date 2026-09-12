package MyGOWS

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"time"
)

type SystemMQChan struct {
	MQCh *amqp091.Channel
}

var SystemMQ SystemMQChan

func SystemMQChanInit() (err error) {
	SystemMQ.MQCh, _, err = NewChannel()
	if err != nil {
		return err
	}
	return nil
}
func NewSystemMsg(payload, fromName, toName string, fromID, toID uint) *Message {
	return &Message{
		Type:      "system",
		FromName:  fromName,
		FromID:    fromID,
		ToID:      toID,
		ToName:    toName,
		Payload:   payload,
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

// UnbindExchanger 用于当
func UnbindExchanger(userID, groupID uint) error {
	err := SystemMQ.MQCh.QueueUnbind(strconv.Itoa(int(userID)), strconv.Itoa(int(groupID)), "group", nil)
	if err != nil {
		return err
	}
	return nil
}
