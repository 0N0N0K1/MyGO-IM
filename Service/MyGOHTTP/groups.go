package MyGOHTTP

import (
	"MyGO-IM/DB"
	"MyGO-IM/Service/MyGOWS"

	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"time"
)

// DropMyGroup 销毁群的处理函数
func DropMyGroup(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	ownerID, _ := c.Get("ownerID")
	gID, _ := c.Get("gID")
	if actorID.(int64) != ownerID.(int64) {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "You aren't the group owner"})
		return
	}
	err := DB.DropGroup(gID.(int64))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "DropGroup: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Drop Group successfully"})

}

// CreateGroup 创建群的处理函数
func CreateGroup(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	actorName, _ := c.Get("actorName")
	groupName := c.Query("name")
	if groupName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "No group name input"})
		return
	}
	GID, err := DB.InsertGroup(actorID.(int64), actorName.(string), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "InsertGroup: " + err.Error()})
		return
	}
	group := DB.QueryGroup(GID)
	err = DB.UpdateGrpStatus(actorID.(int64), GID, "pending")
	err = DB.UpdateGrpStatus(actorID.(int64), GID, "accept")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "InsertMember: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "group": group})
}

// EnterGroup 进入群的处理函数
func EnterGroup(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	gID, _ := c.Get("gID")
	ownerID, _ := c.Get("ownerID")
	applicantID := c.Query("id")
	status := c.Query("status")
	switch status {
	case "pending":
		if DB.QueryMember(actorID.(int64), gID.(int64)).ID != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "You are the member already"})
			return
		}
		err := DB.UpdateGrpStatus(actorID.(int64), gID.(int64), status)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "UpdateGrpStatus: " + err.Error()})
			return
		}

		err = MyGOWS.SendGrpStatus(status, actorID.(int64), ownerID.(int64), gID.(int64))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "notice error"})
			return

		}
	case "reject", "accept":
		if actorID.(int64) != ownerID.(int64) {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "You aren't the group owner"})
			return
		}
		aplID, err := strconv.Atoi(applicantID)
		if err != nil || aplID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
			return
		}
		user, err := DB.QueryUser(aplID, DB.ByID)
		if err != nil || len(user) == 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "QueryUser error"})
			return
		}
		err = DB.UpdateGrpStatus(int64(aplID), gID.(int64), status)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "UpdateGrpStatus: " + err.Error()})
			return
		}

		err = MyGOWS.SendGrpStatus(status, ownerID.(int64), user[0].ID, gID.(int64))
		if err != nil {

			c.JSON(http.StatusOK, gin.H{"code": 102, "error": "notice error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "successfully " + status})
}

// ExitGroup 退出群的处理函数
func ExitGroup(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	gID, _ := c.Get("gID")
	ownerID, _ := c.Get("ownerID")
	if actorID.(int64) == ownerID.(int64) {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "Group owner can't exit "})
		return
	}
	err := DB.NoMemberAnymore(actorID.(int64), gID.(int64))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "NoMemberAnymore: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Exit Group successfully"})
}

// GetMembers 查看所有群成员的处理函数
func GetMembers(c *gin.Context) {
	gID, _ := c.Get("gID")
	limit := c.Query("limit")
	page := c.Query("page")
	l, _ := strconv.Atoi(limit)
	p, _ := strconv.Atoi(page)
	result, err := DB.QueryMemberAll(gID.(int64), p, l)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1005, "error": "QueryMember error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "result": result})
}

// KickOutMember 踢人出群的处理函数
func KickOutMember(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	gID, _ := c.Get("gID")
	ownerID, _ := c.Get("ownerID")
	outerID := c.Query("ID")
	if actorID.(int64) != ownerID.(int64) {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "You aren't the group owner"})
		return
	}
	if outerID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "No outerID input"})
		return
	}
	outID, err := strconv.Atoi(outerID)
	if err != nil || outID <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
		return
	}
	err = DB.NoMemberAnymore(int64(outID), gID.(int64))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "NoMemberAnymore: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Kick out member successfully"})

}

// DoNotSpeak 禁言成员的处理函数
func DoNotSpeak(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	gID, _ := c.Get("gID")
	ownerID, _ := c.Get("ownerID")
	outerID := c.Query("ID")
	timer := c.Query("time")
	if actorID.(int64) != ownerID.(int64) {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "You aren't the group owner"})
		return
	}
	if outerID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "No groupID input"})
		return
	}
	outID, err := strconv.Atoi(outerID)
	if err != nil || outID <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
		return
	}
	t, err := strconv.Atoi(timer)
	if err != nil || t <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
		return
	}
	err = DB.OwnerBanUserHelper(int64(outID), gID.(int64), time.Duration(t)*time.Minute)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "OwnerBanUserHelper: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "Silence member successfully"})
}

// GetGroups 查询已加入的群聊
func GetGroups(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	limit := c.Query("limit")
	page := c.Query("page")
	l, _ := strconv.Atoi(limit)
	p, _ := strconv.Atoi(page)
	result, err := DB.QueryJoinGroup(actorID.(int64), p, l)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "Query error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "result": result})
	return
}
