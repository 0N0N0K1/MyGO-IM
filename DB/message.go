package DB

import (
	"context"
	"errors"
)

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
func MsgNumPlus() {
	RDB.Incr(context.TODO(), "messageNum")
}
