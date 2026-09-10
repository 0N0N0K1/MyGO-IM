package DB

func InsertGroup(ownerID uint, ownerName, groupName string) error {
	var group = Group{
		GroupName: groupName,
		OwnerID:   ownerID,
		OwnerName: ownerName,
	}

	err := MySQL.Table("groups").Create(&group).Error
	return err
}

func QueryGroup(groupID uint) Group {
	var result Group
	MySQL.Table("groups").Where("id=?", groupID).First(&result)
	return result
}
func QueryMyGroup(input uint) []Group {
	var result []Group
	MySQL.
		Raw("select b.group_name, b.id from users a,`groups` b,group_users c where a.id=? and b.id=c.group_id and a.id=c.user_id ",
			input).
		Find(&result)
	return result
}

// IsOwner 用于判断是否为群主
func IsOwner(ownerID, groupID uint) bool {
	var out []GroupUser
	err := MySQL.
		Raw("select * from `groups` a where a.id=? and a.owner_id=?", groupID, ownerID).
		First(&out).Error
	if err != nil || len(out) == 0 {
		return false
	}
	return true
}

// DropGroup 用于销毁一个群
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
