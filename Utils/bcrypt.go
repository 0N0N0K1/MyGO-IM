package Utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"mygoim/DB"
)

func (f *LoginForm) CreateHashPwd() (hashPwd string, err error) {
	hashPwdTep, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword: ", err)
		return "", err
	}
	return string(hashPwdTep), nil
}
func (f *LoginForm) VerifyPwd() (OK bool) {

	user, _ := DB.QueryUser(f.Username)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(f.Password))
	if err != nil {
		return false
	}
	return true
}
