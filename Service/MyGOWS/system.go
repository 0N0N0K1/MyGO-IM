package MyGOWS

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
)

type SystemMQChan struct {
	MQCh *amqp091.Channel
}

// SystemMQ 系统使用的 AMQP 通道实例
var SystemMQ SystemMQChan

// SystemMQChanInit 初始化系统使用的 AMQP 通道
func SystemMQChanInit() (err error) {
	SystemMQ.MQCh, _, err = NewChannel()
	if err != nil {
		return err
	}
	return nil
}

// NewServerACK 返回ACK消息实例
func NewServerACK(ack, desc string, replyID int64, resend bool) *ServerACK {
	return &ServerACK{
		Cmd:     ack,
		ReplyID: replyID,
		MsgID:   int64(SnowID.Generate()),
		Desc:    desc,
		Resend:  resend,
	}
}

// SendFrdStatus 发送好友申请相关系统通知
func SendFrdStatus(status string, fromID, toID uint) error {
	var msg *ServerACK
	switch status {
	case "pending":
		content := fmt.Sprintf("用户ID: %d 申请成为你的好友", fromID)
		msg = NewServerACK("notice", content, int64(toID), false)
	case "reject":
		content := fmt.Sprintf("用户ID: %d 拒绝了你的好友申请", fromID)
		msg = NewServerACK("notice", content, int64(toID), false)
	case "accept":
		content := fmt.Sprintf("用户ID: %d 接收了你的好友申请", fromID)
		msg = NewServerACK("notice", content, int64(toID), false)
	}
	err := SystemMQ.PublishServerACK(msg, toID)
	if err != nil {
		return err
	}
	return nil
}

// SendGrpStatus 发送进群相关系统通知
func SendGrpStatus(status string, fromID, toID, groupID uint) error {
	var msg *ServerACK
	switch status {
	case "pending":
		content := fmt.Sprintf("用户ID: %d 申请进入群聊", fromID)
		msg = NewServerACK("notice", content, int64(toID), false)
	case "reject":
		content := fmt.Sprintf("用户ID: %d  拒绝了你加入群聊ID: %d", fromID, groupID)
		msg = NewServerACK("notice", content, int64(toID), false)
	case "accept":
		content := fmt.Sprintf("用户ID: %d  同意了你加入群聊ID: %d", fromID, groupID)
		msg = NewServerACK("notice", content, int64(toID), false)
	}
	err := SystemMQ.PublishServerACK(msg, toID)
	if err != nil {
		return err
	}
	return nil
}

// UnbindExchanger 用于当退出群聊时解绑
func UnbindExchanger(userID, groupID uint) error {
	err := SystemMQ.MQCh.QueueUnbind(strconv.Itoa(int(userID)), strconv.Itoa(int(groupID)), "group", nil)
	if err != nil {
		return err
	}
	return nil
}
