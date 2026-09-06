package Service

import (
	"github.com/gin-gonic/gin"
	"mygoim/DB"
	"mygoim/Utils"
	"net/http"
	"strconv"
)

func Login(c *gin.Context) {
	var postform Utils.LoginForm
	if err := c.BindJSON(&postform); err != nil || postform.Username == "" || postform.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again1!"})
		return
	}
	if !postform.VerifyPwd() {
		c.JSON(http.StatusOK, gin.H{"code": "100", "msg": "Try again2!"})
		return
	}
	token := postform.CreateJWT()
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"code": "100", "msg": "Try again3!"})
		return
	}
	user, err := DB.QueryUser(postform.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "100", "msg": "Try again3!"})
		return
	}
	if ban, reason := DB.QueryIfBan(strconv.Itoa(int(user.ID))); ban {
		c.JSON(http.StatusOK, gin.H{"code": "100", "msg": "Banned!!!", "reason": reason})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "001", "msg": "Login successfully!", "token": token})

}

func Register(c *gin.Context) {
	var postform Utils.LoginForm
	var repeat Utils.LoginForm
	if err := c.BindJSON(&postform); err != nil || postform.Username == "" || postform.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again!"})
		return
	}

	DB.MySQL.Table("users").Select("name").Where("name = ?", postform.Username).First(&repeat)
	if repeat.Username != "" {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again1!"})
		return
	}
	hashPwd, err := postform.CreateHashPwd()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again2!"})
		return
	}
	postform.Password = hashPwd
	err = DB.InsertUser(postform.Username, postform.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again3!"})
		return
	}
	user, err := DB.QueryUser(postform.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again4!"})
		return
	}
	err = DB.InsertUserInfo(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "100", "error": "Try again5!"})
		return
	}
	token := postform.CreateJWT()
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"code": "100", "msg": "Try again6!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "000", "msg": "Register successfully!", "token": token})

}
