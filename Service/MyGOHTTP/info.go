package MyGOHTTP

import (
	"MyGO-IM/DB"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UpdateInfo 更新用户信息
func UpdateInfo(c *gin.Context) {
	userID, _ := c.Get("actorID")
	var postform DB.UserInfo
	err := c.BindJSON(postform)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "BindJSON: " + err.Error()})
		return
	}
	postform.Age = time.Now().Year() - postform.Birthday.Year()
	err = DB.UpdateUserInfo(userID.(uint), postform)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "UpdateUserInfo: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "success": "UpdateUserInfo Successfully", "new_info": postform})
}

// UpdateName 更新用户昵称
func UpdateName(c *gin.Context) {
	userID, _ := c.Get("actorID")
	newName := c.Query("new")
	if len(newName) >= 10 {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "nickname too long "})
		return
	}
	err := DB.UpdateUserName(userID.(uint), newName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100, "error": "UpdateUserName: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "success": "UpdateName Successfully", "new_nickname": newName})
}
