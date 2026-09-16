package MyGOWS

import (
	"fmt"
	"strconv"
)

// NewServerACK 返回ACK消息实例
func NewServerACK(ack, desc string, replyID int64, ToID uint, resend bool) *ServerACK {
	return &ServerACK{
		Cmd:     ack,
		ReplyID: replyID,
		MsgID:   int64(SnowID.Generate()),
		Desc:    desc,
		Resend:  resend,
		Method:  "system",
		ToId:    ToID,
	}
}

// SendFrdStatus 发送好友申请相关系统通知
func SendFrdStatus(status string, fromID, toID uint) error {
	var msg *ServerACK
	switch status {
	case "pending":
		content := fmt.Sprintf("用户ID: %d 申请成为你的好友", fromID)
		msg = NewServerACK("notice", content, int64(toID), toID, false)
	case "reject":
		content := fmt.Sprintf("用户ID: %d 拒绝了你的好友申请", fromID)
		msg = NewServerACK("notice", content, int64(toID), toID, false)
	case "accept":
		content := fmt.Sprintf("用户ID: %d 接收了你的好友申请", fromID)
		msg = NewServerACK("notice", content, int64(toID), toID, false)
	}
	producer := ProducerPool.Get()
	defer ProducerPool.Put(producer)
	err := producer.PublishHandler(msg, strconv.Itoa(int(toID)))
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
		msg = NewServerACK("notice", content, int64(toID), toID, false)
	case "reject":
		content := fmt.Sprintf("用户ID: %d  拒绝了你加入群聊ID: %d", fromID, groupID)
		msg = NewServerACK("notice", content, int64(toID), toID, false)
	case "accept":
		content := fmt.Sprintf("用户ID: %d  同意了你加入群聊ID: %d", fromID, groupID)
		msg = NewServerACK("notice", content, int64(toID), toID, false)
	}
	producer := ProducerPool.Get()
	defer ProducerPool.Put(producer)
	err := producer.PublishHandler(msg, strconv.Itoa(int(toID)))
	if err != nil {
		return err
	}
	return nil
}
