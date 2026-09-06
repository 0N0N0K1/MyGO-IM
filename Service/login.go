package Service

import (
	"MyGO-IM/DB"
	"MyGO-IM/Utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LoginForm struct {
	Password string `gorm:"column:password" json:"password"`
	Email    string `gorm:"column:email" json:"email"`
	Code     string `json:"code"`
	Mode     string `json:"mode"`
	Nickname string `json:"nickname"`
}

// Login 获取验证码后登录
func Login(c *gin.Context) {
	var postform LoginForm
	if err := c.BindJSON(&postform); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "BindJSON: " + err.Error()})
		return
	}
	if postform.Password == "" || postform.Email == "" || postform.Code == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Your input can't be NULL"})
		return
	}
	cmd := DB.RDB.Get(context.TODO(), "auth:code:"+postform.Code)
	result, err := cmd.Result()
	if err != nil || result != postform.Email {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Invalid auth-code"})
		return
	}
	user, err := DB.QueryUser(postform.Email, DB.ByEmail)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "QueryUser: " + err.Error()})
		return
	}
	if !Utils.VerifyPwd(user[0].Password, postform.Password) {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Incorrect password or email address"})
		return
	}
	if ban, reason, t := DB.QueryIfBan(strconv.Itoa(int(user[0].ID))); ban {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "YOU ARE BANNED!!!", "reason": reason, "remain-time": t})
		return
	}
	token, err := Utils.CreateJWT(user[0].ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "CreateJWT: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "Login successfully!", "token": token})

}

// Register 获取验证码后注册
func Register(c *gin.Context) {
	var postform LoginForm

	if err := c.BindJSON(&postform); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "BindJSON: " + err.Error()})
		return
	}
	if postform.Password == "" || postform.Email == "" || postform.Code == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Your input can't be NULL"})
		return
	}
	cmd := DB.RDB.Get(context.TODO(), "auth:code:"+postform.Code)
	result, err := cmd.Result()
	if err != nil || result != postform.Email {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Invalid auth-code"})
		return
	}
	user, err := DB.QueryUser(postform.Email, DB.ByEmail)
	if len(user) != 0 || err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "This email address is already registered"})
		return
	}
	hashPwd, err := Utils.CreateHashPwd(postform.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "CreateHashPwd: " + err.Error()})
		return
	}
	postform.Password = hashPwd
	err = DB.InsertUser(postform.Nickname, postform.Email, postform.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "InsertUser: " + err.Error()})
		return
	}
	user, err = DB.QueryUser(postform.Email, DB.ByEmail)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "QueryUser: " + err.Error()})
		return
	}
	err = DB.InsertUserInfo(user[0].ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "InsertUserInfo: " + err.Error()})
		return
	}
	token, err := Utils.CreateJWT(user[0].ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "CreateJWT: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Register successfully!", "token": token})

}

// AuthCode 获取验证码
func AuthCode(c *gin.Context) {
	var postform LoginForm

	if err := c.BindJSON(&postform); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "BindJSON: " + err.Error()})
		return
	}
	if postform.Password == "" || postform.Email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Your input can't be NULL"})
		return
	}
	switch postform.Mode {
	case "login":
		if !Utils.VerifyPwd(postform.Email, postform.Password) {
			c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Incorrect password or email address"})
			return
		}
		err := Utils.SendAuthMail(postform.Email)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 100, "error": "SendAuthMail: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Auth code is sent!"})

	case "register":

		user, err := DB.QueryUser(postform.Email, DB.ByEmail)
		if len(user) != 0 || err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 100, "error": "This email address is already registered"})
			return
		}
		err = Utils.SendAuthMail(postform.Email)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "SendAuthMail: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Auth code is sent!"})
	}
}
