package DB

// InsertMember 插入成员表记录
func InsertMember(memberID, groupID uint) (err error) {
	var member = GroupUser{
		GroupID: groupID,
		UserID:  memberID,
	}
	err = MySQL.Table("group_users").Create(&member).Error
	return err
}

func QueryMemberAll(groupID uint) []User {
	var result []User

	if groupID == 0 {
		return result
	}
	MySQL.
		Raw("select a.name, a.id from users a,`groups` b,group_users c where  b.id=? and b.id=c.group_id and a.id=c.user_id ",
			groupID).
		Find(&result)
	return result
}

func QueryMember(userID, groupID uint) User {
	var result User
	MySQL.
		Raw("select a.name, a.id from users a,group_users c where c.group_id=? and a.id=c.user_id and a.id=?",
			groupID, userID).
		First(&result)
	return result
}

// NoMemberAnymore 删除成员表记录
func NoMemberAnymore(memberID, groupID uint) error {
	return MySQL.Table("group_users").Where("group_id=? and user_id=?", groupID, memberID).Delete(&GroupUser{}).Error
}
