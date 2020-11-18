package route

import (
	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/controller"
	"github.com/xiao0811/xiniu/middleware"
)

// GetRouter 获取路由
func GetRouter() *gin.Engine {
	app := gin.Default()

	app.POST("/login", controller.Login)

	token := app.Group("/v1/")
	token.Use(middleware.VerifyToken())

	user := token.Group("user")
	{
		user.POST("/get_info", controller.GetUserInfo)
		user.POST("/create", controller.CreateUser)
		user.POST("/get_details", controller.GetUserDetails)
		user.POST("/update", controller.UpdateUser)
	}

	marshalling := token.Group("marshalling")
	{
		marshalling.POST("/create", controller.CreateMarshalling)
		marshalling.POST("/update", controller.UpdateMarshalling)
		marshalling.POST("/delete", controller.DeleteMarshalling)
	}

	return app
}
