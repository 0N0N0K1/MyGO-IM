package DB

import (
	"errors"
)

// UpdateFrdStatus  用于更新好友表中好友状态
func UpdateFrdStatus(activeID, passiveID uint, status string) error {
	var friendship = UserUser{
		ActiveID:  activeID,
		PassiveID: passiveID,
		Status:    status,
	}
	switch status {
	case "pending": //待处理
		err := MySQL.Table("user_users").Create(&friendship).Error
		return err
	case "reject": //拒绝
		err := MySQL.Table("user_users").Where("(active_id=? and passive_id=? and status='pending') ", activeID, passiveID).
			Delete(&UserUser{}).Error
		return err
	case "accept": //同意
		err := MySQL.Table("user_users").Updates(&friendship).Error
		return err
	}
	return errors.New("no this status")
}

// QueryFrd SQL 查找好友（by id/nickname/all）
func QueryFrd(ID uint, input any, option int) ([]User, error) {
	var result = make([]User, 0)
	var tmp User
	switch option {
	case ByID:
		MySQL.
			Raw("select a.name, a.id from users a,user_users c where c.status='accept'  and a.id=? and ((c.active_id=? and c.passive_id=?) or  (c.active_id=? and c.passive_id=?))", input, ID, input, input, ID).
			Find(&result)
		return result, nil
	case ByName:
		frds, err := QueryUser(input, ByName)
		if err != nil {
			for _, frd := range frds {
				MySQL.
					Raw("select a.name, a.id from users a,user_users c where  c.status='accept' and a.id=? and ((c.active_id=? and c.passive_id=?) or  (c.active_id=? and c.passive_id=?))", frd.ID, ID, frd.ID, frd.ID, ID).
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
			Raw("select b.name, b.id, b.last_online from users a,users b,user_users c where c.status='accept' and a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id))", ID).
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
		Raw("select b.name, b.id, b.last_online from users a,users b,user_users c where c.status='accept' and a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id)) limit ? offset ?", ID, limit, page*limit).
		Find(&result)
	return result
}

// DeleteFrd SQL 删除好友
func DeleteFrd(ID, frdID uint) error {
	err := MySQL.Table("user_users").Where("((active_id=? and passive_id=?) or (active_id=? and passive_id=?)) and status='accept' ", ID, frdID, frdID, ID).Delete(&UserUser{}).Error
	return err
}
