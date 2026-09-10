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
	if actorID.(uint) != ownerID.(uint) {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "You aren't the group owner"})
		return
	}
	err := DB.DropGroup(gID.(uint))
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
	groupName := c.Query("groupname")
	if groupName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "No group name input"})
		return
	}
	groupID, err := DB.QueryGroupID(actorID.(uint), groupName)
	if groupID != 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "Repeat create"})
		return
	}
	err = DB.InsertGroup(actorID.(uint), actorName.(string), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "InsertGroup: " + err.Error()})
		return
	}
	groupID, err = DB.QueryGroupID(actorID.(uint), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "QueryGroupID: " + err.Error()})
		return
	}
	err = DB.UpdateGrpStatus(actorID.(uint), groupID, "pending")
	err = DB.UpdateGrpStatus(actorID.(uint), groupID, "accept")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "InsertMember: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Create Group successfully"})
}

// EnterGroup 进入群的处理函数
func EnterGroup(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	actorName, _ := c.Get("actorName")
	gID, _ := c.Get("gID")
	gName, _ := c.Get("gName")
	ownerName, _ := c.Get("ownerName")
	ownerID, _ := c.Get("ownerID")
	applicantID := c.Query("applicant")
	status := c.Query("status")
	switch status {
	case "pending":
		if DB.QueryMember(actorID.(uint), gID.(uint)).ID != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "You are the member already"})
			return
		}
		err := DB.UpdateGrpStatus(actorID.(uint), gID.(uint), status)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "UpdateGrpStatus: " + err.Error()})
			return
		}
		MyGOWS.SendGrpStatus(status, actorName.(string), ownerName.(string), gName.(string), actorID.(uint), ownerID.(uint), gID.(uint))
	case "reject", "accept":
		if actorID.(uint) != ownerID.(uint) {
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
		err = DB.UpdateGrpStatus(uint(aplID), gID.(uint), status)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "UpdateGrpStatus: " + err.Error()})
			return
		}
		MyGOWS.SendGrpStatus(status, ownerName.(string), user[0].Name, gName.(string), ownerID.(uint), user[0].ID, gID.(uint))
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Enter Group successfully"})
}

// ExitGroup 退出群的处理函数
func ExitGroup(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	gID, _ := c.Get("gID")
	ownerID, _ := c.Get("ownerID")
	if actorID.(uint) == ownerID.(uint) {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "Group owner can't exit "})
		return
	}
	err := DB.NoMemberAnymore(actorID.(uint), gID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "NoMemberAnymore: " + err.Error()})
		return
	}
	err = MyGOWS.UnbindExchanger(actorID.(uint), gID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "UnbindExchanger: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Exit Group successfully"})
}

// GetMembers 查看所有群成员的处理函数
func GetMembers(c *gin.Context) {
	gID, _ := c.Get("gID")
	result := DB.QueryMemberAll(gID.(uint))
	c.JSON(http.StatusOK, gin.H{"code": "005", "msg": "Successful!", "members": result})
}

// KickOutMember 踢人出群的处理函数
func KickOutMember(c *gin.Context) {
	actorID, _ := c.Get("actorID")
	gID, _ := c.Get("gID")
	ownerID, _ := c.Get("ownerID")
	outerID := c.Query("ID")
	if actorID.(uint) != ownerID.(uint) {
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
	err = DB.NoMemberAnymore(uint(outID), gID.(uint))
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
	if actorID.(uint) != ownerID.(uint) {
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
	err = DB.OwnerBanUserHelper(uint(outID), gID.(uint), time.Duration(t)*time.Minute)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "OwnerBanUserHelper: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "Silence member successfully"})
}
