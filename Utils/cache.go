package Utils

import (
	"MyGO-IM/DB"

	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

// UsersIdSent 返回 users:<id>:sent 键名
func UsersIdSent(userID uint) string {
	return fmt.Sprintf("user:%s:sent", strconv.Itoa(int(userID)))
}

// UsersIdSentTag 返回 users:<id>:sent:<tag> 键名
func UsersIdSentTag(deliverTag uint64, workerID int) string {
	return fmt.Sprintf("tag:%s:worker:%s", strconv.Itoa(int(deliverTag)), strconv.Itoa(workerID))
}

// CacheDedupMsgID 返回最近收到消息缓存msgID的键名
func CacheDedupMsgID(userID uint) string {
	return fmt.Sprintf("consumed:%s", strconv.Itoa(int(userID)))
}

// Dedup 用于使用 msgID 判断是否已经处理过
func Dedup(ID uint, msgID int64) bool {

	// 1. 判重
	_, err := DB.RDB.ZScore(context.TODO(), CacheDedupMsgID(ID), strconv.FormatInt(msgID, 10)).Result()
	if err == nil {
		log.Println("err==nil")
		return false // 已存在
	}
	if !errors.Is(err, redis.Nil) {
		log.Println("not nil err:" + err.Error())
		return false
	}
	// 2. 插入
	if err = DB.RDB.ZAdd(context.TODO(), CacheDedupMsgID(ID), redis.Z{Score: float64(time.Now().Unix()), Member: strconv.FormatInt(msgID, 10)}).Err(); err != nil {
		log.Println("ZAdd err:" + err.Error())
		return false
	}

	DB.RDB.Expire(context.TODO(), CacheDedupMsgID(ID), 24*time.Hour)

	// 3. 裁剪到 20 条
	result, _ := DB.RDB.ZCard(context.TODO(), CacheDedupMsgID(ID)).Result()
	if result < 100 {
		return true
	}
	if err = DB.RDB.ZRemRangeByRank(context.TODO(), CacheDedupMsgID(ID), 0, -21).Err(); err != nil {
		log.Println("ZRemRangeByRank err:" + err.Error())
	}
	return true

}
