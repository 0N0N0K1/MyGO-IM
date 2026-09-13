package DB

// GetPrivateSeq 获得 ID1 的ReadSeq, WriteSeq
func GetPrivateSeq(ID1, ID2 uint) (ReadSeq, WriteSeq uint64) {
	var user UserUser
	MySQL.Table("user_users").
		Select("active_read_seq,active_write_seq,passive_read_seq,passive_write_seq").
		Where("active_id=? and passive_id=? and status='accept'", ID1, ID2).First(&user)
	if user.ActiveWriteSeq != 0 {
		return user.ActiveReadSeq, user.ActiveWriteSeq
	}
	MySQL.Table("user_users").
		Select("active_read_seq,active_write_seq,passive_read_seq,passive_write_seq").
		Where("active_id=? and passive_id=? and status='accept'", ID2, ID1).First(&user)
	if user.PassiveReadSeq != 0 {
		return user.PassiveReadSeq, user.PassiveWriteSeq
	}
	return 0, 0
}

// GetGroupWriterSeq 获取群聊的writeSeq
func GetGroupWriterSeq(gID uint) (writeSeq uint64) {
	var group Group
	MySQL.Table("groups").Select("write_seq").Where("id=?", gID).First(&group)
	return group.WriteSeq
}

// GetGroupReadSeq 获取群聊的ReadSeq
func GetGroupReadSeq(gID uint, userID uint) (ReadSeq uint64) {
	var group GroupUser
	MySQL.Table("group_users").Select("read_seq").Where("group_id=? and user_id=? and status='accept'", gID, userID).First(&group)
	return group.ReadSeq
}

// IncrPrivateWriteSeq 更新私聊中 ID1 的WriteSeq
func IncrPrivateWriteSeq(ID1, ID2 uint, seq uint64) {
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID1, ID2).Update("active_write_seq", seq)
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID2, ID1).Update("passive_write_seq", seq)

}

// IncrPrivateReadSeq 更新私聊中 ID1 的readSeq
func IncrPrivateReadSeq(ID1, ID2 uint, seq uint64) {
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID1, ID2).Update("active_read_seq", seq)
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", ID2, ID1).Update("passive_read_seq", seq)

}

// IncrGroupWriteSeq 更新群聊的 WriteSeq
func IncrGroupWriteSeq(gID uint, seq uint64) {
	MySQL.Table("groups").Where("id=? ", gID).Update("write_seq", seq)

}

// IncrGroupReadSeq 更新群聊成员的 ReadSeq
func IncrGroupReadSeq(gID, userID uint, seq uint64) {
	MySQL.Table("group_users").Where("group_id=? and user_id=? and status='accept'", gID, userID).Update("read_seq", seq)
}
