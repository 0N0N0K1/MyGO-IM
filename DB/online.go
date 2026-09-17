package DB

import (
	"context"
	"strconv"
	"time"
)

// SetOnline 缓存在线状态
func SetOnline(userID int64) {
	var online = 1

	RDB.HSetNX(context.TODO(), "users:online", strconv.Itoa(int(userID)), online)
}

// SetOffline 删除缓存的在线状态
func SetOffline(userID int64) {
	RDB.HDel(context.TODO(), "users:online", strconv.Itoa(int(userID)))
}

// CheckOnline 检查是否在线
func CheckOnline(userID int64) (online bool) {
	res := RDB.HGet(context.TODO(), "users:online", strconv.Itoa(int(userID)))
	_, err := res.Result()
	if err != nil {
		return false
	}
	return true
}

// SetLastOnline 更新上次在线时间
func SetLastOnline(userID int64) {
	MySQL.Table("users").Where("id=?", userID).Update("last_online", time.Now())
}
