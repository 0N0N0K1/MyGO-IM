package DB

// GetPrivateSeq 获得 ID1、ID2 的ReadSeq, WriteSeq
func GetPrivateSeq(ID1, ID2 int64) (ID1ReadSeq, ID2ReadSeq, globalWriteSeq uint64) {
	var user UserUser
	MySQL.Table("user_users").
		Select("active_read_seq,passive_read_seq,global_write_seq").
		Where("active_id=? and passive_id=? and status='accept'", ID1, ID2).First(&user)
	if user.GlobalWriteSeq != 0 {
		return user.ActiveReadSeq, user.PassiveReadSeq, user.GlobalWriteSeq
	}
	MySQL.Table("user_users").
		Select("active_read_seq,passive_read_seq,global_write_seq").
		Where("active_id=? and passive_id=? and status='accept'", ID2, ID1).First(&user)
	if user.PassiveReadSeq != 0 {
		return user.PassiveReadSeq, user.ActiveReadSeq, user.GlobalWriteSeq
	}
	return 0, 0, 0
}
func GetPrivateSeqByCID(uID int64, conversationID string) (ReadSeq, globalWriteSeq uint64) {
	var user UserUser
	MySQL.Table("user_users").
		Select("active_read_seq,passive_read_seq,global_write_seq").
		Where("active_id=? and conversation_id=? and status='accept'", uID, conversationID).First(&user)
	if user.GlobalWriteSeq != 0 {
		return user.ActiveReadSeq, user.GlobalWriteSeq
	}
	MySQL.Table("user_users").
		Select("active_read_seq,passive_read_seq,global_write_seq").
		Where("conversation_id=? and passive_id=? and status='accept'", conversationID, uID).First(&user)
	if user.PassiveReadSeq != 0 {
		return user.PassiveReadSeq, user.GlobalWriteSeq
	}
	return 0, 0
}

// GetGroupWriterSeq 获取群聊的writeSeq
func GetGroupWriterSeq(gID int64) (writeSeq uint64) {
	var group Group
	MySQL.Table("groups").Select("write_seq").Where("id=?", gID).First(&group)
	return group.WriteSeq
}

// GetGroupWriterSeqByCID 获取群聊的writeSeq
func GetGroupWriterSeqByCID(conversationID string) (writeSeq uint64) {
	var group Group
	MySQL.Table("groups").Select("write_seq").Where("conversation_id=?", conversationID).First(&group)
	return group.WriteSeq
}

// GetGroupReadSeq 获取群聊的ReadSeq
func GetGroupReadSeq(gID int64, userID int64) (ReadSeq uint64) {
	var group GroupUser
	MySQL.Table("group_users").Select("read_seq").Where("group_id=? and user_id=? and status='accept'", gID, userID).First(&group)
	return group.ReadSeq
}

// GetGroupReadSeqByCID 获取群聊的ReadSeq
func GetGroupReadSeqByCID(conversationID string, userID int64) (ReadSeq uint64) {
	var group GroupUser
	MySQL.Table("group_users").Select("read_seq").Where("conversation_id=? and user_id=? and status='accept'", conversationID, userID).First(&group)
	return group.ReadSeq
}

//// IncrPrivateWriteSeq 更新私聊中 ID1 的WriteSeq
//func IncrPrivateWriteSeq(ID1, ID2 int64, seq uint64) {
//	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID1, ID2).Update("active_write_seq", seq)
//	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID2, ID1).Update("passive_write_seq", seq)
//
//}

// IncrPrivateGlobalWriteSeq 更新私聊中 globalWriteSeq
func IncrPrivateGlobalWriteSeq(ID1, ID2 int64, seq uint64) {
	MySQL.Table("user_users").Where("(active_id=? and passive_id=? and status='accept') or (active_id=? and passive_id=? and status='accept')", ID1, ID2, ID2, ID1).Update("global_write_seq", seq)

}

// IncrPrivateReadSeqByCID 更新私聊中readSeq
func IncrPrivateReadSeqByCID(uid int64, conversationID string, seq uint64) {
	MySQL.Table("user_users").Where("active_id=? and conversation_id=? and status='accept'", uid, conversationID).Update("active_read_seq", seq)
	MySQL.Table("user_users").Where("conversation_id=? and passive_id=? and status='accept'", conversationID, uid).Update("passive_read_seq", seq)

}

// IncrPrivateReadSeq 更新私聊中 ID1 的readSeq
func IncrPrivateReadSeq(ID1, ID2 int64, seq uint64) {
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID1, ID2).Update("active_read_seq", seq)
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID2, ID1).Update("passive_read_seq", seq)
}

// IncrGroupWriteSeq 更新群聊的 WriteSeq
func IncrGroupWriteSeq(gID int64, seq uint64) {
	MySQL.Table("groups").Where("id=? ", gID).Update("write_seq", seq)

}

// IncrGroupReadSeq 更新群聊成员的 ReadSeq
func IncrGroupReadSeq(gID, userID int64, seq uint64) {
	MySQL.Table("group_users").Where("group_id=? and user_id=? and status='accept'", gID, userID).Update("read_seq", seq)
}

// IncrGroupReadSeqByCID 更新群聊成员的 ReadSeq
func IncrGroupReadSeqByCID(conversationID string, userID int64, seq uint64) {
	MySQL.Table("group_users").Where("conversation_id=? and user_id=? and status='accept'", conversationID, userID).Update("read_seq", seq)
}

//func NewGroupMsgNum(gid, uid int64) (writeSeq, readSeq uint64) {
//	var groupSeq Group
//	var userSeq GroupUser
//	MySQL.Table("groups").Select("write_seq").Where("id=?", gid).First(&groupSeq)
//	MySQL.Table("group_users").Select("read_seq").Where("group_id=? and user_id=?", gid, uid).First(&userSeq)
//	return groupSeq.WriteSeq, userSeq.ReadSeq
//}
