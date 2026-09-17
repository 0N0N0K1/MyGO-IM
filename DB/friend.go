package DB

import (
	"errors"
	"fmt"
)

// UpdateFrdStatus  更新好友表中好友状态
func UpdateFrdStatus(activeID, passiveID int64, status string) error {
	var CID string
	if activeID < passiveID {
		CID = fmt.Sprintf("p:%d-%d", activeID, passiveID)
	} else {
		CID = fmt.Sprintf("p:%d-%d", passiveID, activeID)
	}
	var friendship = UserUser{
		ActiveID:       activeID,
		PassiveID:      passiveID,
		Status:         status,
		ConversationId: CID,
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

// QueryFrd  查找好友
func QueryFrd(ID int64, input any, option int) ([]User, error) {
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
	case ByCID:
		MySQL.
			Raw("select id from user_users  where status='accept'  and conversation_id=?", input).
			Find(&result)
		return result, nil
	case ALL:
		MySQL.
			Raw("select b.name, b.id, b.last_online from users a,users b,user_users c where c.status='accept' and a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id))", ID).
			Find(&result)
		return result, nil
	default:
		return []User{}, errors.New("option error")
	}
}

func QueryFrdConversationID(ID int64) ([]UserUser, error) {
	var err error
	var result = make([]UserUser, 0)
	err = MySQL.
		Raw("select b.id,c.conversation_id from users a,users b,user_users c where c.status='accept' and a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id)) ", ID).
		Find(&result).Error
	return result, err
}

// QueryFrdLimit  分页查找好友
func QueryFrdLimit(ID int64, page, limit int) (Pagination, error) {
	var result = make([]User, 0)
	var err error
	err = MySQL.
		Raw("select b.* from users a,users b,user_users c where c.status='accept' and a.id=? and ((a.id=c.active_id and b.id =c.passive_id) or  (b.id=c.active_id and a.id =c.passive_id)) limit ? offset ?", ID, limit, page*limit).
		Find(&result).Error
	return Pagination{
		Page:  page,
		Limit: limit,
		List:  result,
	}, err
}

// DeleteFrd  删除好友
func DeleteFrd(ID, frdID int64) error {
	err := MySQL.Table("user_users").Where("((active_id=? and passive_id=?) or (active_id=? and passive_id=?)) and status='accept' ", ID, frdID, frdID, ID).Delete(&UserUser{}).Error
	return err
}
