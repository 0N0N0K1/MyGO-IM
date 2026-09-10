package API

import (
	"MyGO-IM/Service/MyGOHTTP"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ServiceEngine *gin.Engine
var AdminEngine *gin.Engine

func InitRouters() {

	ServiceEngine = gin.Default()
	AdminEngine = gin.Default()

	ServiceEngine.POST("/login", MyGOHTTP.Login)
	ServiceEngine.POST("/register", MyGOHTTP.Register)
	ServiceEngine.POST("/send-code", MyGOHTTP.AuthCode)

	// 用户组 API
	Users := ServiceEngine.Group("/users/:ID")
	Users.Use(MyGOHTTP.VerifyJWT)
	{
		// 通过 Query 传ID增删查
		Users.GET("/friends", MyGOHTTP.GetFriends)
		Users.POST("/friends", MyGOHTTP.AddFriend)
		Users.DELETE("/friends", MyGOHTTP.DeleteFriend)

		Users.GET("/chat", MyGOHTTP.WSUpgrade)

		//  test API
		Users.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": "999", "msg": "Hello!"})
		})

		// 群组 API
		Groups := Users.Group("/groups")
		{
			// 创建
			Groups.POST("/", MyGOHTTP.CreateGroup)

			GroupsWithMW := Groups.Group("/:gID")
			GroupsWithMW.Use(MyGOHTTP.GroupMiddleware)
			{
				// 销毁
				GroupsWithMW.DELETE("/", MyGOHTTP.DropMyGroup)
				// 进入/退出
				GroupsWithMW.DELETE("/members", MyGOHTTP.ExitGroup)
				GroupsWithMW.POST("/members", MyGOHTTP.EnterGroup)

				// 踢出成员
				GroupsWithMW.DELETE("/owner", MyGOHTTP.KickOutMember)
				// 禁言成员
				GroupsWithMW.PATCH("/owner", MyGOHTTP.DoNotSpeak)

				// 查询
				GroupsWithMW.GET("/:gID/members", MyGOHTTP.GetMembers)
			}
		}
	}

	// 管理员 API
	Admin := AdminEngine.Group("/admin")

	{
		Admin.GET("/chat", MyGOHTTP.WSUpgrade)
		Admin.GET("/basicinfo", MyGOHTTP.GetBasicInfo)
		Admin.GET("/message", MyGOHTTP.QueryMessage)
		// 用户管理API
		UserManamge := Admin.Group("/users")
		{
			UserManamge.GET("/password", MyGOHTTP.ChangePassword)
			UserManamge.GET("/ban", MyGOHTTP.BanUser)
			UserManamge.GET("/info", MyGOHTTP.QueryUser)
		}
	}
}
