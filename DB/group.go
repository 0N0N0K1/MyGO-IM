package DB

// InsertGroup 创建群聊
func InsertGroup(ownerID uint, ownerName, groupName string) error {
	var group = Group{
		GroupName: groupName,
		OwnerID:   ownerID,
		OwnerName: ownerName,
	}

	err := MySQL.Table("groups").Create(&group).Error
	return err
}

// QueryGroup 查找群聊
func QueryGroup(groupID uint) Group {
	var result Group
	MySQL.Table("groups").Where("id=?", groupID).First(&result)
	return result
}

// QueryMyGroup 查找自己加入的群聊
func QueryMyGroup(userID uint) []Group {
	var result []Group
	MySQL.
		Raw("select b.group_name, b.id from users a,`groups` b,group_users c where a.id=? and b.id=c.group_id and a.id=c.user_id ",
			userID).
		Find(&result)
	return result
}

// DropGroup 销毁一个群
func DropGroup(groupID uint) error {
	return MySQL.Table("groups").Where("id=?", groupID).Delete(&Group{}).Error
}

// QueryGroupID 查询用户名下指定名字的群聊
func QueryGroupID(ownerID uint, groupName string) (uint, error) {
	var group Group
	err := MySQL.
		Raw("select id from `groups` where owner_id=? and group_name = ?", ownerID, groupName).
		First(&group).Error
	if err != nil || group.ID == 0 {
		return 0, err
	}
	return group.ID, nil
}
