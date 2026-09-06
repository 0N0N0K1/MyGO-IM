package Utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CreateHashPwd(password string) (hashPwd string, err error) {
	hashPwdTep, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword: ", err)
		return "", err
	}
	return string(hashPwdTep), nil
}
func VerifyPwd(hashPwd, password string) (OK bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(password))
	if err != nil {
		return false
	}
	return true
}
