package Service

import (
	"github.com/gin-gonic/gin"
	"mygoim/DB"
	"net/http"
	"time"
)

func GetFriends(c *gin.Context) {
	userID, _ := c.Get("ID")
	results := DB.QueryFrdAll(userID.(uint))
	for index := range results {
		if DB.CheckOnline(results[index].Name) {
			results[index].Online = 1
			results[index].LastOnline = new(time.Time)
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": "003", "msg": "successfully!", "friends": results})
}

func AddFriend(c *gin.Context) {
	userID, _ := c.Get("ID")
	frdName := c.Param("friendname")
	frd, err := DB.QueryUser(frdName, DB.ByName)
	//对方是否存在，是否为自己
	if frd.ID == 0 || userID == frd.ID || err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "Illegal Addition!"})
		return
	}
	//是否已添加
	result := DB.QueryFrd(userID.(uint), frd.ID)
	if result.Name != "" {
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "Repeated Addition!"})
		return
	}
	//插入到中间表
	err = DB.InsertFrd(userID.(uint), frd.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "102", "msg": "Try Again!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "002", "msg": "Successful Addition!"})
}
