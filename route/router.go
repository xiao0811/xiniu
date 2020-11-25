package route

import (
	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/controller"
	"github.com/xiao0811/xiniu/controller/label"
	"github.com/xiao0811/xiniu/middleware"
)

// GetRouter 获取路由
func GetRouter() *gin.Engine {
	// Debug
	if !config.Conf.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	// 用户登录
	app.POST("/login", controller.Login)
	// 发送修改密码短信
	app.POST("/send_change_password_message", controller.SendChangePasswordMessage)
	// 修改密码
	app.POST("/change_password", controller.ChangePassword)

	token := app.Group("/v1/")
	token.Use(middleware.VerifyToken())

	user := token.Group("user")
	{
		user.POST("/get_info", controller.GetUserInfo)
		user.POST("/create", controller.CreateUser)
		user.POST("/get_details", controller.GetUserDetails)
		user.POST("/update", controller.UpdateUser)
		user.POST("/user_list", controller.UserList)
		user.POST("/delete", controller.DeleteUser)
		user.POST("/batch_group", controller.UserBatchGroup)
	}

	marshalling := token.Group("marshalling")
	{
		marshalling.POST("/create", controller.CreateMarshalling)
		marshalling.POST("/update", controller.UpdateMarshalling)
		marshalling.POST("/delete", controller.DeleteMarshalling)
		marshalling.POST("/marshalling_list", controller.MarshallingList)
	}

	l := token.Group("/label")
	{
		l.POST("/create", label.Create)
		l.POST("/update", label.Update)
		l.POST("/delete", label.Delete)
		l.POST("/index", label.Index)
		l.POST("/info", label.Info)
	}
	lg := token.Group("/label_group")
	{
		lg.POST("/create", label.CreateGroup)
		lg.POST("/update", label.UpdateGroup)
		lg.POST("/delete", label.DeleteGroup)
		lg.POST("/index", label.GroupList)
		lg.POST("/info", label.GroupDetails)
	}
	member := token.Group("/member")
	{
		member.POST("/create", controller.CreateMember)
		member.POST("/update", controller.UpdateMember)
		member.POST("/review", controller.MemberReview)
		member.POST("/member_list", controller.MemberList)
		member.POST("/get_member_details", controller.GetMemberDetails)
	}
	return app
}
