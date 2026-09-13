package DB

import (
	"context"
	"strconv"
	"time"
)

// SetOnline 缓存在线状态
func SetOnline(userID uint) {
	var online = 1

	RDB.HSetNX(context.TODO(), "online", strconv.Itoa(int(userID)), online)
}

// SetOffline 删除缓存的在线状态
func SetOffline(userID uint) {
	RDB.HDel(context.TODO(), "online", strconv.Itoa(int(userID)))
}

// CheckOnline 检查是否在线
func CheckOnline(userID uint) (online bool) {
	res := RDB.HGet(context.TODO(), "online", strconv.Itoa(int(userID)))
	val, err := res.Int()

	if err != nil || val == 0 {
		return false
	}
	return true
}

// SetLastOnline 更新上次在线时间
func SetLastOnline(userID uint) {
	MySQL.Table("users").Where("id=?", userID).Update("last_online", time.Now())
}
