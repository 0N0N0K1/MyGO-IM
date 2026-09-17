package DB

import (
	"MyGO-IM/Utils"
	"fmt"
)

// InsertGroup 创建群聊
func InsertGroup(ownerID int64, ownerName, groupName string) (int64, error) {
	ID := Utils.SnowID.Generate()
	var group = Group{
		MemberNum:      1,
		ID:             int64(ID),
		GroupName:      groupName,
		OwnerID:        ownerID,
		OwnerName:      ownerName,
		ConversationId: fmt.Sprintf("p:%d", ID),
	}
	err := MySQL.Table("groups").Create(&group).Error
	return int64(ID), err
}

// QueryGroup 查找一个群聊
func QueryGroup(groupID int64) Group {
	var result Group
	MySQL.Table("groups").Where("id=?", groupID).First(&result)
	return result
}

// QueryMyGroup 查找加入的群聊
func QueryMyGroup(userID int64) []Group {
	var result []Group
	MySQL.
		Raw("select b.* from users a,`groups` b,group_users c where a.id=? and b.id=c.group_id and a.id=c.user_id ",
			userID).
		Find(&result)
	return result
}

// DropGroup 销毁一个群
func DropGroup(groupID int64) error {
	return MySQL.Table("groups").Where("id=?", groupID).Delete(&Group{}).Error
}

// QueryGroupID 查询用户名下指定名字的群聊
func QueryGroupID(ownerID int64, groupName string) (int64, error) {
	var group Group
	err := MySQL.
		Raw("select id from `groups` where owner_id=? and group_name = ?", ownerID, groupName).
		First(&group).Error
	if err != nil || group.ID == 0 {
		return 0, err
	}
	return group.ID, nil
}

// QueryJoinGroup 查询加入的群聊
func QueryJoinGroup(uID int64, page, limit int) (Pagination, error) {
	var result []Group
	err := MySQL.
		Raw("select b.* from users a,`groups` b,group_users c where a.id=? and b.id=c.group_id and a.id=c.user_id ",
			uID).
		Find(&result).Error
	return Pagination{
		Page:  page,
		Limit: limit,
		List:  result,
	}, err
}
