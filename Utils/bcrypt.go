package Utils

import (
	"MyGO-IM/DB"
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
func VerifyPwd(userEmail, password string) (OK bool) {

	user, _ := DB.QueryUser(userEmail, DB.ByEmail)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
