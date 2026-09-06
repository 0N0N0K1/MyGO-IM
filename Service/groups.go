package Service

import (
	"MyGO-IM/DB"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateGroup(c *gin.Context) {
	ownerID, _ := c.Get("ID")
	ownerName, _ := c.Get("name")
	groupName := c.Param("groupname")
	if DB.QueryGroupID(groupName).ID != 0 {
		c.JSON(http.StatusOK, gin.H{"code": "104", "msg": "Illegal Create!"})
		return
	}
	err := DB.InsertGroup(ownerID.(uint), ownerName.(string), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "104", "msg": "Try Again!"})
		return
	}
	err = DB.InsertMember(ownerID.(uint), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "104", "msg": "Try Again!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "004", "msg": "Successful Creation!"})
}
func EnterGroup(c *gin.Context) {
	memberID, _ := c.Get("ID")
	groupName := c.Param("groupname")
	err := DB.InsertMember(memberID.(uint), groupName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "105", "msg": "Try Again!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "004", "msg": "Successful Enter!"})
}
func GetMembers(c *gin.Context) {
	groupName := c.Param("groupname")
	result := DB.QueryMemberAll(groupName)
	c.JSON(http.StatusOK, gin.H{"code": "005", "msg": "Successful!", "members": result})
}
