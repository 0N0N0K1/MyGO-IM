package DB

func GetPrivateSeq(formID, toID uint) (ReadSeq, WriteSeq uint64) {
	var user UserUser
	MySQL.Table("user_users").
		Select("active_read_seq,active_write_seq,passive_read_seq,passive_write_seq").
		Where("active_id=? and passive_id=? and status='accept'", formID, toID).First(&user)
	if user.ActiveWriteSeq != 0 {
		return user.ActiveReadSeq, user.ActiveWriteSeq
	}
	MySQL.Table("user_users").
		Select("active_read_seq,active_write_seq,passive_read_seq,passive_write_seq").
		Where("active_id=? and passive_id=? and status='accept'", toID, formID).First(&user)
	if user.PassiveReadSeq != 0 {
		return user.PassiveReadSeq, user.PassiveWriteSeq
	}
	return 0, 0
}
func GetGroupWriterSeq(gID uint) (writeSeq uint64) {
	var group Group
	MySQL.Table("groups").Select("write_seq").Where("id=?", gID).First(&group)
	return group.WriteSeq
}
func GetGroupReadSeq(gID uint, userID uint) (ReadSeq uint64) {
	var group GroupUser
	MySQL.Table("group_users").Select("read_seq").Where("group_id=? and user_id=? and status='accept'", gID, userID).First(&group)
	return group.ReadSeq
}
func IncrPrivateWriteSeq(formID, toID uint, seq uint64) {
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", formID, toID).Update("active_write_seq", seq)
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", toID, formID).Update("passive_write_seq", seq)

}

func IncrPrivateReadSeq(formID, toID uint, seq uint64) {
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", formID, toID).Update("active_read_seq", seq)
	MySQL.Table("user_users").Where("active_id=? and passive_id=? and status='accept'", toID, formID).Update("passive_read_seq", seq)

}
func IncrGroupWriteSeq(gID uint, seq uint64) {
	MySQL.Table("groups").Where("id=? ", gID).Update("write_seq", seq)

}

func IncrGroupReadSeq(gID, userID uint, seq uint64) {
	MySQL.Table("group_users").Where("group_id=? and user_id=? and status='accept'", gID, userID).Update("read_seq", seq)

}
