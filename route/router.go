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
	token.POST("/get_user_info", controller.GetUserInfo)
	return app
}
