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
	// 获取图片
	app.GET("/upload/images/:images_name", controller.ShowImage)
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
	contract := token.Group("/contract")
	{
		contract.POST("/create", controller.CreateContract)
		contract.POST("/update", controller.UpdateContract)
		contract.POST("/review", controller.ContractReview)
		contract.POST("/contract_list", controller.ContractList)
		contract.POST("/get_member_details", controller.GetContractDetails)
		// 合约延期
		contract.POST("/extension", controller.ContractExtension)
		// 获取合约七日任务
		contract.POST("/get_task", controller.GetContractTask)
		// 修改七日任务
		contract.POST("/update_task", controller.UpdateContractTask)
		contract.POST("/get_contract_by_status", controller.GetContractByStatus)
		// 退款
		contract.POST("/refund", controller.ContractRefund)
		contract.POST("/change_management", controller.ChangeManagement)
	}

	cts := token.Group("/contract_task")
	{
		cts.POST("/create", controller.CreateContractTask)
		cts.POST("/delete", controller.DeleteContractTask)
		cts.POST("/get_list", controller.GetContractTaskList)
	}
	upload := token.Group("/upload")
	{
		upload.POST("/images", controller.UploadImages)
		upload.POST("/image", controller.UploadImage)
	}

	refund := token.Group("/refund")
	{
		refund.POST("/create", controller.CreateRefund)
		refund.POST("/review", controller.ReviewRefund)
		refund.POST("/get_details", controller.GetRefundDetails)
		refund.POST("/list", controller.GetRefundList)
	}
	sms := token.Group("sms")
	{
		sms.POST("send_message", controller.SendMessage)
	}
	home := token.Group("/home")
	{
		home.POST("/count_data", controller.CountData)
		home.POST("/my_contract", controller.MyContract)
		home.POST("/30_days_service", controller.ServiceDays30)
	}
	return app
}
