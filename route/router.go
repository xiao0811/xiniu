package route

import (
	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/controller"
)

// GetRouter 获取路由
func GetRouter() *gin.Engine {
	app := gin.Default()

	app.POST("/login", controller.Login)
	return app
}
