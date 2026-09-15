package MyGOHTTP

import (
	"MyGO-IM/DB"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetMessage 查询聊天记录
// 通过 query 指定
// 1.查询聊天对象 群/好友 ID
// 2.查询方式，双方消息/发送消息
func GetMessage(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	groupID := c.Query("group")
	friendID := c.Query("friend")
	page := c.Query("page")
	limit := c.Query("limit")
	mode := c.Query("mode")
	lm, _ := strconv.Atoi(limit)
	pg, _ := strconv.Atoi(page)
	if groupID != "" {
		gID, err := strconv.Atoi(groupID)
		if err != nil || gID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "ID invalid"})
			return
		}
		result, err := DB.QueryGroupMsg(uint(gID), actorID.(uint), pg, lm, DB.MsgQueryMode(mode))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "QueryGroupMemberMsg error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "messages": result})
		return
	}
	if friendID != "" {
		fID, err := strconv.Atoi(friendID)
		if err != nil || fID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "ID invalid"})
			return
		}
		result, err := DB.QueryPrivateMsg(uint(fID), actorID.(uint), pg, lm, DB.MsgQueryMode(mode))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1111, "error": "QueryGroupMemberMsg error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "messages": result})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1112, "error": "no appointed ID"})
	return
}
