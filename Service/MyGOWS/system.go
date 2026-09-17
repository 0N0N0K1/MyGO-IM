package MyGOWS

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"fmt"
	"strconv"
)

// NewSystemMsg 返回ACK消息实例
func NewSystemMsg(ack, desc string, replyID int64, ToID int64, resend bool, extra any) *SystemMsg {
	return &SystemMsg{
		Cmd:     ack,
		ReplyID: replyID,
		MsgID:   int64(Utils.SnowID.Generate()),
		Desc:    desc,
		Resend:  resend,
		Method:  "system",
		ToId:    ToID,
		Extra:   extra,
	}
}

// SendFrdStatus 发送好友申请相关系统通知
func SendFrdStatus(status string, fromID, toID int64) error {
	var msg *SystemMsg
	var result DB.UserUser
	DB.MySQL.
		Raw("select * from user_users where  (active_id=? and passive_id=?) or  (active_id=? and passive_id=?)", fromID, toID, toID, fromID).
		Find(&result)

	switch status {
	case "pending":
		content := fmt.Sprintf("用户ID: %d 申请成为你的好友", fromID)
		msg = NewSystemMsg("notice", content, int64(toID), toID, false, result)
	case "reject":
		content := fmt.Sprintf("用户ID: %d 拒绝了你的好友申请", fromID)
		msg = NewSystemMsg("notice", content, int64(toID), toID, false, result)
	case "accept":
		content := fmt.Sprintf("用户ID: %d 接收了你的好友申请", fromID)
		msg = NewSystemMsg("notice", content, int64(toID), toID, false, result)
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
func SendGrpStatus(status string, fromID, toID, groupID int64) error {
	var msg *SystemMsg
	result := DB.QueryGroup(groupID)
	switch status {
	case "pending":
		content := fmt.Sprintf("用户ID: %d 申请进入群聊", fromID)
		msg = NewSystemMsg("notice", content, int64(toID), toID, false, result)
	case "reject":
		content := fmt.Sprintf("用户ID: %d  拒绝了你加入群聊ID: %d", fromID, groupID)
		msg = NewSystemMsg("notice", content, int64(toID), toID, false, result)
	case "accept":
		content := fmt.Sprintf("用户ID: %d  同意了你加入群聊ID: %d", fromID, groupID)
		msg = NewSystemMsg("notice", content, int64(toID), toID, false, result)
	}
	producer := ProducerPool.Get()
	defer ProducerPool.Put(producer)
	err := producer.PublishHandler(msg, strconv.Itoa(int(toID)))
	if err != nil {
		return err
	}
	return nil
}
