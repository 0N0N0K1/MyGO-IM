package DB

import (
	"errors"
)

// UpdateGrpStatus 插入成员表记录
func UpdateGrpStatus(memberID, groupID uint, status string) (err error) {
	var member = GroupUser{
		GroupID: groupID,
		UserID:  memberID,
		Status:  status,
		ReadSeq: GetGroupWriterSeq(groupID),
	}
	switch status {
	case "pending":
		err = MySQL.Table("group_users").Create(&member).Error
		return err
	case "reject":
		err = MySQL.Table("group_users").Where("group_id=? and user_id=?", groupID, memberID).Delete(&GroupUser{}).Error
		return err
	case "accept":
		err = MySQL.Table("group_users").Updates(&member).Error
		return err

	}
	return errors.New("no this status")
}

// QueryMemberAll 查找返回群成员列表
func QueryMemberAll(groupID uint, page, limit int) (Pagination, error) {
	var result []User

	err := MySQL.
		Raw("select a.*from users a,`groups` b,group_users c where  b.id=? and b.id=c.group_id and a.id=c.user_id limit ? offset ?",
			groupID, limit, page*limit).
		Find(&result).Error
	return Pagination{
		Page:  page,
		Limit: limit,
		List:  result,
	}, err
}

// QueryMember 查找返回一个群成员
func QueryMember(userID, groupID uint) User {
	var result User
	MySQL.
		Raw("select a.name, a.id from users a,group_users c where c.status='accept' and c.group_id=? and a.id=c.user_id and a.id=?",
			groupID, userID).
		First(&result)
	return result
}

// NoMemberAnymore 删除成员表记录
func NoMemberAnymore(memberID, groupID uint) error {
	return MySQL.Table("group_users").Where("group_id=? and user_id=?", groupID, memberID).Delete(&GroupUser{}).Error
}
