package route

import (
	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/controller"
	"github.com/xiao0811/xiniu/controller/forum"
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
	app.Use(middleware.Cors())
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
		user.POST("get_member", controller.GetMember)
		user.POST("/user_list_all", controller.UserListAll)
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
		member.POST("/delete", controller.DeleteMember)
		member.POST("/review", controller.MemberReview)
		member.POST("/member_list", controller.MemberList)
		member.POST("/get_member_details", controller.GetMemberDetails)
		member.POST("/export", controller.ExportMembers)
	}
	contract := token.Group("/contract")
	{
		contract.POST("/create", controller.CreateContract)
		contract.POST("/update", controller.UpdateContract)
		contract.POST("/delete", controller.DeleteContract)
		contract.POST("/review", controller.ContractReview)
		contract.POST("/export", controller.ExportContract)
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
		contract.POST("/batch_change_management", controller.BatchChangeManagement)
	}

	cts := token.Group("/contract_task")
	{
		cts.POST("/create", controller.CreateContractTask)
		cts.POST("/delete", controller.DeleteContractTask)
		cts.POST("/get_list", controller.GetContractTaskList)
		cts.POST("/update", controller.UpdateContractTask2)

		cts.POST("/create_details", controller.CreateContractTaskDetails)
		cts.POST("/get_tesk_details", controller.GetContractTasKDetails)
		cts.POST("/export", controller.ExportContractTask)
		cts.POST("/delete_details", controller.DeleteContractTasKDetails)
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

	cls := token.Group("/contract_log")
	{
		// 创建合约操作记录
		cls.POST("/create", controller.CreateContractLog)
		// 删除合约操作记录
		cls.POST("/delate", controller.DeleteContractLog)
		cls.POST("/get_logs_by_contrat_id", controller.GetLogsByContratID)
	}

	// 论坛主题
	ft := token.Group("/forum_title")
	{
		// 新建主题
		ft.POST("/create", forum.CreateTitle)
		ft.POST("/update", forum.UpdateTitle)
		ft.POST("/delete", forum.DeleteTitle)
		ft.POST("/get_title_list", forum.GetTitleList)
		ft.POST("/details", forum.TitleDetails)
		ft.POST("/get_title_by_user", forum.GetForumTitleByUser)
		ft.POST("/carousel_or_recommended", forum.CarouselOrRecommended)
	}

	// 论坛评论
	fc := token.Group("/forum_comment")
	{
		fc.POST("/create", forum.CreateComment)
		fc.POST("/update", forum.UpdateComment)
		fc.POST("/delete", forum.DeleteComment)
		fc.POST("/adoption", controller.CommentAdoption)
	}

	// 点赞
	like := token.Group("/forum_like")
	{
		like.POST("/like", forum.Like)
		like.POST("/unlike", forum.Unlike)
	}

	// 积分
	integra := token.Group("/integra")
	{
		integra.POST("/change", controller.IntegralChange)
		integra.POST("/rank", controller.IntegralRank)
	}

	// 站内信
	stationLetter := token.Group("/station_letter")
	{
		stationLetter.POST("/create", controller.CreateStationLetter)
		stationLetter.POST("/update", controller.UpdateStationLetter)
		stationLetter.POST("/delete", controller.DeleteStationLetter)
		stationLetter.POST("/get", controller.GetStationLetter)
	}

	// 组队
	team := token.Group("/team")
	{
		// 创建队伍
		team.POST("/create", controller.CreateTeam)
		// 编辑队伍
		team.POST("/update", controller.UpdateTeam)
		//  删除分组
		team.POST("/delete", controller.DeleteTeam)
		// 获取组队积分
		team.POST("/team_rank", controller.GetTeamRank)
		// 获取队伍
		team.POST("/get_teams", controller.GetTeams)
	}
	return app
}
