package DB

import (
	"context"
	"errors"
)

type MsgQueryMode string

const (
	Recv MsgQueryMode = "receive"
	Send MsgQueryMode = "send"
)

// InsertMsg 插入一条消息
func InsertMsg(from, to int64, writerseq uint64, cID, method string, Content string) (err error) {
	var msgs = PrivateMessage{
		ToID:           to,
		FromID:         from,
		Content:        Content,
		Seq:            writerseq,
		ConversationId: cID,
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

// QueryPrivateMsg 查找 ID 与 ID2 的聊天
func QueryPrivateMsg(ID1, ID2 int64, page, limit int, mode MsgQueryMode) (Pagination, error) {
	var msgs []PrivateMessage
	var err error
	switch mode {
	case Send:
		err = MySQL.Table("private_messages").Select("*").Where("from_id=? and to_id", ID1, ID2).Order("created_at DECR").Scopes(Paginate(page, limit)).Find(&msgs).Error
	case Recv:
		err = MySQL.Table("private_messages").Select("*").Where("(from_id=? and to_id=?) or (from_id=? or to_id=?) ", ID1, ID2, ID2, ID1).Order("created_at DECR").Scopes(Paginate(page, limit)).Find(&msgs).Error
	default:
		return Pagination{}, nil
	}
	return Pagination{
		Page:  page,
		Limit: limit,
		List:  msgs,
	}, err
}

// QueryGroupMsg 查找 ID 群的聊天记录
func QueryGroupMsg(gID, mID int64, page, limit int, mode MsgQueryMode) (Pagination, error) {
	var msgs []GroupMessage
	var err error
	switch mode {
	case Recv:
		err = MySQL.Table("group_messages").Select("*").Where("to_id=?", gID).Order("created_at DECR").Scopes(Paginate(page, limit)).Find(&msgs).Error
	case Send:
		err = MySQL.Table("group_messages").Select("*").Where("to_id=? and from_id", gID, mID).Order("created_at DECR").Scopes(Paginate(page, limit)).Find(&msgs).Error
	default:
		return Pagination{}, nil
	}
	return Pagination{
		Page:  page,
		Limit: limit,
		List:  msgs,
	}, err
}

// MsgNumPlus 今日聊天记录数+1
func MsgNumPlus() {
	RDB.Incr(context.TODO(), "messageNum")
}

// PullGroupMsg  拉取群的聊天记录
func PullGroupMsg(conversationID string, limit int, endseq uint64) ([]GroupMessage, error) {
	var msgs []GroupMessage
	var err error
	err = MySQL.Table("group_messages").
		Select("*").
		Where("conversation_id=? and seq<=?", conversationID, endseq).
		Order("seq desc").
		Limit(limit).Find(&msgs).Error

	return msgs, err
}

// PullPrivateMsg 拉取私聊的聊天记录
func PullPrivateMsg(conversationID string, limit int, endseq uint64) ([]GroupMessage, error) {
	var msgs []GroupMessage
	var err error
	err = MySQL.Table("private_messages").
		Select("*").
		Where("conversation_id=? and seq<=?", conversationID, endseq).
		Order("seq desc").
		Limit(limit).Find(&msgs).Error

	return msgs, err
}

// InsertNotice 插入系统通知防止申请/同意/拒绝……等消息丢失
func InsertNotice(ID, toid int64, content string) (NoticeMessage, error) {
	var notice = NoticeMessage{
		ID:      ID,
		Read:    false,
		ToID:    toid,
		Content: content,
	}
	err := MySQL.Table("notice_messages").Create(&notice).Error
	return notice, err
}

// QueryMyNewNoticeNum 获取新的系统通知数
func QueryMyNewNoticeNum(uid int64) ([]int64, int64, error) {
	var IDs []int64
	var notices int64
	err := MySQL.Table("notice_messages").
		Where("to_id=? and `read`=?", uid, false).
		Count(&notices).Error
	err = MySQL.Table("notice_messages").
		Select("id").
		Where("to_id=? and `read`=?", uid, false).
		Find(&IDs).
		Error
	return IDs, notices, err
}

// QueryMyNotice 获取系统通知
func QueryMyNotice(uid int64, lm int) ([]NoticeMessage, error) {
	var notices []NoticeMessage
	err := MySQL.Table("notice_messages").
		Select("*").
		Where("to_id=? ", uid).
		Order("id desc").Limit(lm).
		Find(&notices).Error
	MySQL.Table("notice_messages").
		Where("to_id=? and id<=?", uid, notices[0].ID).
		UpdateColumn("read", true)
	return notices, err
}
