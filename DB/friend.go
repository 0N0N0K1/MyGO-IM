package DB

import (
	"errors"
)

// InsertFrd SQL 用于插入一条好友记录
func InsertFrd(ID, frdID uint) error {
	//删除标记恢复机制
	err := MySQL.Table("user_users").Create(&UserUser{
		ActiveID:  ID,
		PassiveID: frdID,
	}).Error
	return err
}

// QueryFrd SQL 查找好友（by id/nickname/all）
func QueryFrd(ID uint, input any, option int) ([]User, error) {
	var result = make([]User, 0)
	var tmp User
	switch option {
	case ByID:
		MySQL.
			Raw("select a.name, a.id from users a,user_users c where  a.id=? and ((c.active_id=? and c.passive_id=?) or  (c.active_id=? and c.passive_id=?))", input, ID, input, input, ID).
			Find(&result)
		return result, nil
	case ByName:
		frds, err := QueryUser(input, ByName)
		if err != nil {
			for _, frd := range frds {
				MySQL.
					Raw("select a.name, a.id from users a,user_users c where   a.id=? and ((c.active_id=? and c.passive_id=?) or  (c.active_id=? and c.passive_id=?))", frd.ID, ID, frd.ID, frd.ID, ID).
					Find(&tmp)
				if tmp.ID != 0 {
					result = append(result, tmp)
				}
			}
			return result, nil
		}
		return []User{}, err
	case ALL:
		MySQL.
			Raw("select b.name, b.id, b.last_online from users a,users b,user_users c where  a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id))", ID).
			Find(&result)
		return result, nil
	default:
		return []User{}, errors.New("option error")
	}

}

// QueryFrdLimit SQL 查找好友 分页
func QueryFrdLimit(ID uint, page, limit int) []User {
	var result = make([]User, 0)
	MySQL.
		Raw("select b.name, b.id, b.last_online from users a,users b,user_users c where  a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id)) limit ? offset ?", ID, limit, page*limit).
		Find(&result)
	return result
}

// DeleteFrd SQL 删除好友
func DeleteFrd(ID, frdID uint) error {
	err := MySQL.Table("user_users").Where("(active_id=? and passive_id=?) or (active_id=? and passive_id=?) ", ID, frdID, frdID, ID).Delete(&UserUser{}).Error
	return err
}
