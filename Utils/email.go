package Utils

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"mygoim/DB"
	"time"
)

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
	DB.RDB.Set(context.TODO(), "auth:code:"+code, to, time.Minute*2)
	return nil
}
