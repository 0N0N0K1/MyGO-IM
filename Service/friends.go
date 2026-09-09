package Service

import (
	"MyGO-IM/DB"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetFriends 好友查询的处理函数，提供通过 nickname, id, limit+page 三种query方式
func GetFriends(c *gin.Context) {
	userID, _ := c.Get("ID")
	frdName := c.Query("nickname")
	frdID := c.Query("id")
	limit := c.Query("limit")
	page := c.Query("page")

	if frdID != "" {
		fID, err := strconv.Atoi(frdID)
		if err != nil || fID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
			return
		}
		frds, err := DB.QueryFrd(userID.(uint), fID, DB.ByID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "type": "ByID", "friends": frds})
		return
	}
	if frdName != "" {
		frds, err := DB.QueryFrd(userID.(uint), frdName, DB.ByName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ByName!", "friends": frds})
		return
	}
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil || l < 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
			return
		}
		if page == "" {
			frds := DB.QueryFrdLimit(userID.(uint), 0, l)
			c.JSON(http.StatusOK, gin.H{"code": 0, "type": "limit", "friends": frds})
			return
		}
		p, err := strconv.Atoi(page)
		if err != nil || p < 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "string2int error"})
			return
		}
		frds := DB.QueryFrdLimit(userID.(uint), p, l)
		c.JSON(http.StatusOK, gin.H{"code": 0, "type": "limit!", "friends": frds})
		return
	}
	frds, err := DB.QueryFrd(userID.(uint), nil, DB.ALL)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "type": "ALL", "friends": frds})
}

// AddFriend 添加好友的处理函数，提供通过 query id 添加好友
func AddFriend(c *gin.Context) {
	userID, _ := c.Get("ID")
	frdID := c.Query("id")
	var frd []DB.User
	var err error
	fID, err := strconv.Atoi(frdID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
		return
	}
	frd, err = DB.QueryUser(fID, DB.ByID)
	//对方是否存在，是否为自己
	if frd[0].ID == 0 || userID == frd[0].ID || err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "Illegal Addition!"})
		return
	}
	result, _ := DB.QueryFrd(userID.(uint), frd[0].ID, DB.ByID)
	//是否已添加
	if len(result) != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "Repeated addition"})
		return
	}
	//插入到中间表
	err = DB.InsertFrd(userID.(uint), frd[0].ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "InsertFrd: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Successful Addition!"})
}

// DeleteFriend 删除好友的处理函数，通过query id 删除
func DeleteFriend(c *gin.Context) {
	userID, _ := c.Get("ID")
	frdID := c.Query("id")
	var err error
	fID, err := strconv.Atoi(frdID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
		return
	}
	frd, err := DB.QueryFrd(userID.(uint), fID, DB.ByID)
	//是否已添加
	if frd[0].Name == "" || userID == fID || err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "Illegal delete"})
		return
	}
	//插入到中间表
	err = DB.DeleteFrd(userID.(uint), uint(fID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "DeleteFrd: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Successful delete!"})
}
