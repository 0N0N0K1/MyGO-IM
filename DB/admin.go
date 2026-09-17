package DB

import (
	"MyGO-IM/Utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// InsertUser 插入用户/注册
func InsertUser(name, email string, pwd string) error {
	var user = User{
		ID:       int64(Utils.SnowID.Generate()),
		Name:     name,
		Email:    email,
		Password: pwd,
	}
	err := MySQL.Table("users").Create(&user).Error
	return err
}

// UpdateUser 用于更新用户信息
func UpdateUser(user *User) error {
	err := MySQL.Table("users").Select("*").Where("id=?", user.ID).Updates(user).Error
	return err
}

// QueryUser 查询用户
func QueryUser(input any, option int) ([]User, error) {
	var result = make([]User, 0)
	switch option {
	case ByID:
		MySQL.Table("users").Select("*").Where("id=?", input).First(&result)
	case ByEmail:
		MySQL.Table("users").Select("*").Where("email=?", input).First(&result)
	case ByName:
		MySQL.Table("users").Select("*").Where("name=?", input).Find(&result)
		log.Println(result)
	default:
		return []User{}, errors.New("option error")
	}
	return result, nil
}

// QueryUserWithInfo 查询用户关键与详细信息
func QueryUserWithInfo(input any, option int) (DetailInfo, error) {
	var result DetailInfo
	switch option {
	case ByID:
		MySQL.Raw("select * from users a,user_infos b where a.id=b.user_id and a.id =?", input).First(&result)
	case ByEmail:
		MySQL.Raw("select * from users a,user_infos b where a.id=b.user_id and a.email =?", input).First(&result)
	default:
		return DetailInfo{}, errors.New("input type error")
	}
	return result, nil
}

// UserNum 返回已注册用户总数
func UserNum() uint {
	var num uint
	MySQL.Table("users").Select("count(*)").Find(&num)
	return num
}

// OnlineNum 返回在线用户总数
func OnlineNum() uint {
	cmd := RDB.HLen(context.TODO(), "online")
	u, _ := cmd.Uint64()
	return uint(u)
}

// GetMsgNum 返回今日发帖量
func GetMsgNum() uint {
	cmd := RDB.Get(context.TODO(), "messageNum")
	num, _ := cmd.Int()
	return uint(num)
}

// GetFresherToday 返回今天新注册的用户数
func GetFresherToday() uint {
	var num uint
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	MySQL.Table("user_infos").Select("count(*)").Where("created_at between ? and ?", start, end).Find(&num)
	return num
}

// AdminBanUserHelper admin永久或者暂时全局禁言用户
func AdminBanUserHelper(userID int64, t int, reason string) error {
	if t <= 0 {
		err := MySQL.Table("users").Where("id=?", userID).Delete(&User{}).Error
		return err
	}
	user := fmt.Sprintf("admin:slience:%d", userID)
	err := RDB.SetNX(context.TODO(), user, reason, time.Hour*time.Duration(t)).Err()
	if err != nil {
		return err
	}
	return err
}

// QueryIfBan 查询是否被admin全局禁言
func QueryIfBan(userID int64) (ban bool, reason string, t time.Duration) {
	user := fmt.Sprintf("admin:slience:%d", userID)
	cmd := RDB.Get(context.TODO(), user)
	if cmd.Err() != nil {
		return false, "", 0
	}
	cmd2 := RDB.TTL(context.TODO(), user)
	t, _ = cmd2.Result()
	reason = cmd.String()
	return true, reason, t
}

// OwnerBanUserHelper 群主永久或者暂时全局禁言用户
func OwnerBanUserHelper(userID, groupID int64, t time.Duration) error {
	if t <= 0 {
		err := MySQL.Table("users").Where("id=?", userID).Delete(&User{}).Error
		return err
	}
	group := fmt.Sprintf("groups:%d:silence:%d", groupID, userID)
	err := RDB.SetNX(context.TODO(), group, userID, t).Err()
	if err != nil {
		return err
	}
	return err
}

// QueryIfSilence 查询是否被群主禁言
func QueryIfSilence(userID, groupID int64) (ban bool, t time.Duration) {
	user := fmt.Sprintf("groups:%d:silence:%d", groupID, userID)
	cmd := RDB.Get(context.TODO(), user)
	if cmd.Err() != nil {
		return false, 0
	}
	cmd2 := RDB.TTL(context.TODO(), user)
	t, _ = cmd2.Result()
	return true, t
}

// DeleteUser 删除用户
func DeleteUser(userID int64) error {
	err := MySQL.Table("user_infos").Where("user_id=?", userID).Delete(&UserInfo{}).Error
	if err != nil {
		return err
	}
	err = MySQL.Table("users").Where("id=?", userID).Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}
