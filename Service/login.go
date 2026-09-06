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
	if !Utils.VerifyPwd(postform.Email, postform.Password) {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "Incorrect password or email address"})
		return
	}
	token, err := Utils.CreateJWT(postform.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "CreateJWT: " + err.Error()})
		return
	}
	user, err := DB.QueryUser(postform.Email, DB.ByEmail)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "QueryUser: " + err.Error()})
		return
	}
	if ban, reason, t := DB.QueryIfBan(strconv.Itoa(int(user.ID))); ban {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "YOU ARE BANNED!!!", "reason": reason, "remain-time": t})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "Login successfully!", "token": token})

}

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
	if user.Name != "" {
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
	err = DB.InsertUserInfo(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "InsertUserInfo: " + err.Error()})
		return
	}
	token, err := Utils.CreateJWT(postform.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "CreateJWT: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Register successfully!", "token": token})

}

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
		if user.Name != "" || err != nil {
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
