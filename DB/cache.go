package DB

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// UsersIdSent 返回 users:<id>:sent 键名
func UsersIdSent(userID int64) string {
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
	_, err := RDB.ZScore(context.TODO(), CacheDedupMsgID(ID), strconv.FormatInt(msgID, 10)).Result()
	if err == nil {
		log.Println("err==nil")
		return false // 已存在
	}
	if !errors.Is(err, redis.Nil) {
		log.Println("not nil err:" + err.Error())
		return false
	}
	// 2. 插入
	if err = RDB.ZAdd(context.TODO(), CacheDedupMsgID(ID), redis.Z{Score: float64(time.Now().Unix()), Member: strconv.FormatInt(msgID, 10)}).Err(); err != nil {
		log.Println("ZAdd err:" + err.Error())
		return false
	}

	RDB.Expire(context.TODO(), CacheDedupMsgID(ID), 24*time.Hour)

	// 3. 裁剪到 20 条
	result, _ := RDB.ZCard(context.TODO(), CacheDedupMsgID(ID)).Result()
	if result < 100 {
		return true
	}
	if err = RDB.ZRemRangeByRank(context.TODO(), CacheDedupMsgID(ID), 0, -21).Err(); err != nil {
		log.Println("ZRemRangeByRank err:" + err.Error())
	}
	return true

}

// SendAuthMail 发送验证码到指定邮箱并缓存
func SendAuthMail(to string) error {
	qqEmail := "3397545837@qq.com"
	authCode := "xsxdeqfjsuoncigi"
	smtpHost := "smtp.qq.com"
	smtpPort := 587

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", rnd.Intn(1000000))

	m := gomail.NewMessage()
	m.SetHeader("From", qqEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "MyGO IM验证码")
	m.SetBody("text/plain", fmt.Sprintf("欢迎使用MyGo IM!\n您的本次操作的验证码是：\n\n%s\n\n请在2分钟内使用", code))

	d := gomail.NewDialer(smtpHost, smtpPort, qqEmail, authCode)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	RDB.Set(context.TODO(), "auth:code:"+code, to, time.Minute*2)
	return nil
}
