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

func QueryMemberAll(groupName string) []User {
	var result []User

	groupID := QueryGroupID(groupName).ID
	if groupID == 0 {
		return result
	}
	MySQL.
		Raw("select a.name, a.id from users a,`groups` b,group_users c where c.break=0 and b.id=? and b.id=c.group_id and a.id=c.user_id ",
			groupID).
		Find(&result)
	return result
}

func QueryMember(userName, groupName string) User {
	var result User

	groupID := QueryGroupID(groupName).ID
	if groupID == 0 {
		return result
	}
	MySQL.
		Raw("select a.name, a.id from users a,group_users c where c.break=0 and c.group_id=? and a.id=c.user_id and a.name=?",
			groupID, userName).
		Find(&result)
	return result
}

// NoMemberAnymore 删除成员表记录
func NoMemberAnymore(memberID, groupID uint) error {
	return MySQL.Table("group_users").Where("group_id=? and user_id=?", groupID, memberID).Delete(&GroupUser{}).Error
}
