package Utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// CreateHashPwd 生成加密后的密码以保存
func CreateHashPwd(password string) (hashPwd string, err error) {
	hashPwdTep, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword: ", err)
		return "", err
	}
	return string(hashPwdTep), nil
}

// VerifyPwd 用原始密码与加密密码比较
func VerifyPwd(hashPwd, password string) (OK bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(password))
	if err != nil {
		return false
	}
	return true
}
