package Service

import (
	"MyGO-IM/Chat"
	"MyGO-IM/Utils"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"MyGO-IM/DB"
	"net/http"
	"strconv"
	"time"
)

// BasicInfo 定义基本信息结构
type BasicInfo struct {
	UserNum      uint `json:"register_num"`
	OnlineNum    uint `json:"online_num"`
	MsgToday     uint `json:"msg_today"`
	FresherToday uint `json:"fresher_today"`
}

func AdminInit() {
	go func() {
		for {
			<-time.Tick(time.Hour * 24)
			DB.RDB.Set(context.TODO(), "messageNum", 0, 0)
		}
	}()
}

// GetBasicInfo 用于admin获取基本信息API处理函数
func GetBasicInfo(c *gin.Context) {
	var basicInfo BasicInfo
	basicInfo.UserNum = DB.UserNum()
	basicInfo.OnlineNum = DB.OnlineNum()
	basicInfo.MsgToday = DB.GetMsgNum()
	basicInfo.FresherToday = DB.GetFresherToday()
	c.JSON(http.StatusOK, gin.H{"code": "0", "basicinfo": basicInfo})
}

// QueryMessage 用于admin获取聊天记录API的处理函数
func QueryMessage(c *gin.Context) {
	userid := c.Query("id")
	msgType := c.Query("type")
	parseUint, err := strconv.ParseUint(userid, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1001", "error": "ID can't parse"})
		return
	}
	msgs, err := DB.QueryUserMsg(uint(parseUint), msgType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1001", "error": "Type can't search"})
		return
	}
	var messages = make([]Chat.Message, 1, 1)
	var message Chat.Message
	for _, msg := range msgs {
		_ = json.Unmarshal([]byte(msg.Content), &message)
		messages = append(messages, message)
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "messages": messages, "err": err})
}

// ChangePassword 用于管理员修改User密码API的处理函数
func ChangePassword(c *gin.Context) {
	userID := c.Query("id")
	userNewPwd := c.Query("password")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1002", "err": err.Error()})
		return
	}
	user, err := DB.QueryUser(uint(id), DB.ByID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1002", "err": "Query user by ID error" + err.Error()})
		return
	}

	pwd, err := Utils.CreateHashPwd(userNewPwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1002", "err": "Update password error"})
		return
	}
	user.Password = pwd
	err = DB.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1002", "err": "Update password error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "message": "Updated!", "user": user})
}

// BanUser 用于admin封禁/禁言用户
func BanUser(c *gin.Context) {
	userID := c.Query("id")
	banTime := c.Query("time")
	reason := c.Query("reason")

	t, err := strconv.Atoi(banTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1003", "err": err.Error()})
		return
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1003", "err": err.Error()})
		return
	}
	user, err := DB.QueryUser(uint(id), DB.ByID)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": "1002", "err": "Query user by ID error"})
		return
	}
	err = DB.BanUserHelper(uint(id), t, reason)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1003", "err": "Ban error" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "message": "Banned!", "time": banTime + " * Hour", "user": user})
}

func QueryUser(c *gin.Context) {
	userID := c.Query("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1003", "err": err.Error()})
		return
	}
	user, err := DB.QueryUserWithInfo(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "1002", "err": "Query user by ID error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "message": "Queried!", "user": user})
}
