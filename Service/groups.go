package Service

import (
	"MyGO-IM/DB"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DropMyGroup 销毁群的处理函数
func DropMyGroup(c *gin.Context) {
	ownerID, _ := c.Get("ID")
	gID := c.Param("gID")
	if gID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "No groupID input"})
		return
	}
	groupID, err := strconv.Atoi(gID)
	if err != nil || groupID <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "String2int error"})
		return
	}
	if !DB.IsOwner(ownerID.(uint), uint(groupID)) {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "You aren't the group owner"})
		return
	}
	err = DB.DropGroup(uint(groupID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "DropGroup: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Drop Group successfully"})

}

// CreateGroup 创建群聊的处理函数
func CreateGroup(c *gin.Context) {
	ownerID, _ := c.Get("ID")
	ownerName, _ := c.Get("name")
	groupName := c.Query("groupname")
	if groupName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1002, "error": "No group name input"})
		return
	}
	err := DB.InsertGroup(ownerID.(uint), ownerName.(string), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "InsertGroup: " + err.Error()})
		return
	}
	groupID, err := DB.QueryGroupID(ownerID.(uint), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "QueryGroupID: " + err.Error()})
		return
	}
	err = DB.InsertMember(ownerID.(uint), groupID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1004, "error": "InsertMember: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Create Group successfully"})
}

// EnterGroup 进入群聊的处理函数
func EnterGroup(c *gin.Context) {
	memberID, _ := c.Get("ID")
	gID := c.Param("gID")
	if gID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "No groupID input"})
		return
	}
	groupID, err := strconv.Atoi(gID)
	if err != nil || groupID <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
		return
	}
	err = DB.InsertMember(memberID.(uint), uint(groupID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "InsertMember: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Enter Group successfully"})
}
func GetMembers(c *gin.Context) {
	groupName := c.Query("groupname")
	result := DB.QueryMemberAll(groupName)
	c.JSON(http.StatusOK, gin.H{"code": "005", "msg": "Successful!", "members": result})
}

// ExitGroup 退出群聊的处理函数
func ExitGroup(c *gin.Context) {
	memberID, _ := c.Get("ID")
	gID := c.Param("gID")
	if gID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "No groupID input"})
		return
	}
	groupID, err := strconv.Atoi(gID)
	if err != nil || groupID <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "String2int error"})
		return
	}
	err = DB.NoMemberAnymore(memberID.(uint), uint(groupID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1003, "error": "NoMemberAnymore: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Exit Group successfully"})
}
