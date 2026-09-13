package DB

import (
	"context"
	"errors"
)

// InsertMsg 插入一条消息
func InsertMsg(from, to uint, seq uint64, fromname, toname, method string, Content string) (err error) {
	var msgs = PrivateMessage{
		ToID:    to,
		FromID:  from,
		Content: Content,
		Seq:     seq,
	}

	switch method {
	case "private":

		err = MySQL.Table("private_messages").Create(&msgs).Error

	case "group":

		err = MySQL.Table("group_messages").Create(&msgs).Error

	default:
		err = errors.New("message type error")
	}
	return err
}

//TODO 优化查找逻辑

// QueryUserMsg 查找聊天记录
func QueryUserMsg(ID uint, method string) (msgs []PrivateMessage, err error) {
	switch method {
	case "private":
		err = MySQL.Table("private_messages").Select("*").Where("from_id=? or to_id=?", ID, ID).Find(&msgs).Error

	case "group":
		err = MySQL.Table("group_messages").Select("*").Where("to_id=?", ID).Find(&msgs).Error
	default:
		err = errors.New("message type error")
	}
	return
}

// MsgNumPlus 今日聊天记录数+1
func MsgNumPlus() {
	RDB.Incr(context.TODO(), "messageNum")
}
