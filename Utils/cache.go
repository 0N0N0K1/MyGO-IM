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
func UsersIdSentTag(deliverTag uint64, userID uint) string {
	return fmt.Sprintf("user:%s:sent:%s", strconv.Itoa(int(userID)), strconv.Itoa(int(deliverTag)))
}

// CacheDedupMsgID 返回最近收到消息缓存msgID的键名
func CacheDedupMsgID(userID uint, msgID int64) string {
	return fmt.Sprintf("user:%s:dedup:%s", strconv.Itoa(int(userID)), strconv.Itoa(int(msgID)))
}

// Dedup 用于使用 msgID 判断是否已经处理过
func Dedup(userID uint, msgID int64, seq uint64) bool {

	// 1. 判重
	_, err := DB.RDB.ZScore(context.TODO(), CacheDedupMsgID(userID, msgID), strconv.FormatInt(msgID, 10)).Result()
	log.Println(CacheDedupMsgID(userID, msgID) + "    " + strconv.FormatInt(msgID, 10))
	if err == nil {
		log.Println("err==nil")
		return false // 已存在
	}
	if !errors.Is(err, redis.Nil) {
		log.Println("not nil err:" + err.Error())
		return false
	}

	// 2. 插入
	if err = DB.RDB.ZAdd(context.TODO(), CacheDedupMsgID(userID, msgID), redis.Z{Score: float64(seq), Member: strconv.FormatInt(msgID, 10)}).Err(); err != nil {
		log.Println("ZAdd err:" + err.Error())
		return false
	}

	// 3. 裁剪到 20 条
	if err = DB.RDB.ZRemRangeByRank(context.TODO(), CacheDedupMsgID(userID, msgID), 0, -21).Err(); err != nil {
		log.Println("ZRemRangeByRank err:" + err.Error())
		return false
	}

	// 4. 刷新 TTL
	DB.RDB.Expire(context.TODO(), CacheDedupMsgID(userID, msgID), 24*time.Hour)
	log.Println("return Ture")
	return true
}
