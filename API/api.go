package API

import (
	"MyGO-IM/Service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var UserEngine *gin.Engine
var AdminEngine *gin.Engine

func InitRouters() {

	UserEngine = gin.Default()
	AdminEngine = gin.Default()

	UserEngine.POST("/login", Service.Login)
	UserEngine.POST("/register", Service.Register)
	UserEngine.POST("/send-code", Service.AuthCode)

	// 用户组 API
	User := UserEngine.Group("/users/:ID")
	User.Use(Service.VerifyJWT)
	{
		// 通过 Query 传ID增删查
		User.GET("/friends", Service.GetFriends)
		User.POST("/friends", Service.AddFriend)
		User.DELETE("/friends", Service.DeleteFriend)

		//todo 加好友申请与同意拒绝

		User.GET("/chat", Service.WSUpgrade)

		//  test API
		User.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": "999", "msg": "Hello!"})
		})

		// 群组 API
		Group := User.Group("/groups")
		{

			//todo 群主不能退群，群主转让功能，进/退群申请与同意拒绝，群主禁言与踢人功能

			// 创建/销毁
			Group.DELETE("/:gID", Service.DropMyGroup)
			Group.POST("/", Service.CreateGroup)

			// 进入/退出
			Group.DELETE("/:gID/members", Service.ExitGroup)
			Group.POST("/:gID/members", Service.EnterGroup)

			// 成员查看
			Group.GET("/:gID/members", Service.GetMembers)
		}
	}

	// 管理员 API
	Admin := AdminEngine.Group("/admin")

	{
		Admin.GET("/chat", Service.WSUpgrade)
		Admin.GET("/basicinfo", Service.GetBasicInfo)
		Admin.GET("/message", Service.QueryMessage)
		// 用户管理API
		UserManamge := Admin.Group("/users")
		{
			UserManamge.GET("/password", Service.ChangePassword)
			UserManamge.GET("/ban", Service.BanUser)
			UserManamge.GET("/info", Service.QueryUser)
		}
	}
}
