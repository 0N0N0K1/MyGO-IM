package MyGOHTTP

import (
	"MyGO-IM/DB"
	"MyGO-IM/Service/MyGOWS"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetFriends 好友查询的处理函数，提供通过 nickname, id, limit+page 三种query方式
func GetFriends(c *gin.Context) {
	userID, _ := c.Get("actorID")
	frdName := c.Query("nickname")
	frdID := c.Query("id")
	limit := c.Query("limit")
	page := c.Query("page")

	if frdID != "" {
		fID, err := strconv.Atoi(frdID)
		if err != nil || fID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "ID invalid"})
			return
		}
		frds, err := DB.QueryFrd(userID.(int64), fID, DB.ByID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "type": "ByID", "friends": frds})
		return
	}
	if frdName != "" {
		frds, err := DB.QueryFrd(userID.(int64), frdName, DB.ByName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "type": "ByName", "friends": frds})
		return
	}
	l, _ := strconv.Atoi(limit)
	p, _ := strconv.Atoi(page)
	result, err := DB.QueryFrdLimit(userID.(int64), p, l)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "type": "limit", "result": result})
	return
}

// AddFriend 添加好友的处理函数，提供通过 query id 添加好友
func AddFriend(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	frdID := c.Query("id")
	status := c.Query("status")
	var frd []DB.User
	var err error
	fID, err := strconv.Atoi(frdID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
		return
	}
	frd, err = DB.QueryUser(fID, DB.ByID)
	//对方是否存在，是否为自己
	if len(frd) == 0 || actorID == frd[0].ID || err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "Illegal Addition!"})
		return
	}
	result, _ := DB.QueryFrd(actorID.(int64), frd[0].ID, DB.ByID)
	//是否已添加
	if len(result) != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "Repeated addition"})
		return
	}
	switch status {
	case "pending": //待处理

		err = DB.UpdateFrdStatus(actorID.(int64), frd[0].ID, "pending")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "InsertFrd: " + err.Error()})
			return
		}
		err = MyGOWS.SendFrdStatus(status, actorID.(int64), int64(fID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "notice error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Successful Apply!"})
	case "reject": //拒绝
		err = DB.UpdateFrdStatus(frd[0].ID, actorID.(int64), "reject")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "InsertFrd: " + err.Error()})
			return
		}
		err = MyGOWS.SendFrdStatus(status, actorID.(int64), int64(fID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "notice error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Successful Reject!"})
	case "accept": //接受

		err = DB.UpdateFrdStatus(frd[0].ID, actorID.(int64), "accept")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "InsertFrd: " + err.Error()})
			return
		}
		err := MyGOWS.SendFrdStatus(status, actorID.(int64), int64(fID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "notice error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Successful accept!"})
	default:
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Input status error"})
		return
	}
}

// DeleteFriend 删除好友的处理函数，通过query id 删除
func DeleteFriend(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	frdID := c.Query("id")
	var err error
	fID, err := strconv.Atoi(frdID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
		return
	}
	frd, err := DB.QueryFrd(actorID.(int64), fID, DB.ByID)
	//是否已添加
	if len(frd) == 0 || actorID.(int64) == int64(fID) || err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "Illegal delete"})
		return
	}
	//删除
	err = DB.DeleteFrd(actorID.(int64), int64(fID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 102, "error": "DeleteFrd: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Successful delete!"})
}
