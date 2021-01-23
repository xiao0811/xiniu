package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateMember 创建一个新的客户
func CreateMember(c *gin.Context) {
	var r struct {
		Name              string `json:"name" binding:"required"` // 门店名称
		City              string `json:"city"`                    // 所在城市
		FirstCategory     string `json:"first_category"`          // 一级类目
		SecondaryCategory string `json:"secondary_category"`      // 二级类目
		BusinessScope     string `json:"business_scope"`          // 主营范围
		Stores            uint8  `json:"stores"`                  // 门店数量
		Accounts          uint8  `json:"accounts"`                // 账户数量
		Bosses            uint8  `json:"bosses"`                  // 老板人数
		Brands            uint8  `json:"brands"`                  // 品牌数量
		OperationsGroup   int    `json:"operations_group"`        // 运营组
		OperationsStaff   int    `json:"operations_staff"`        // 运营人员
		BusinessGroup     int    `json:"business_group"`          // 业务组
		BusinessPeople    int    `json:"business_people"`         // 业务人员
		ReviewAccount     string `json:"review_account"`          // 点评账号
		CommentPassword   string `json:"comment_password"`        // 点评密码
		Email             string `json:"email"`                   // 客户邮箱
		Phone             string `json:"phone"`                   // 客户手机号码
		OtherTags         string `json:"other_tags"`              // 其他标签
		Auditors          uint   `json:"auditors"`                // 审核人员
		Type              int8   `json:"type"`                    // 备注信息
		Status            int8   `json:"status"`
		Remarks           string `json:"remarks"`
	}
	var m model.Member
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	if err := db.Where("name = ?", r.Name).First(&m).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "门店已存在", c)
		return
	}
	m.UUID = "XINIU-CUS-" + time.Now().Format("200601021504") + strconv.Itoa(handle.RandInt(1000, 9999))
	m.Name = r.Name
	m.City = r.City
	m.FirstCategory = r.FirstCategory
	m.SecondaryCategory = r.SecondaryCategory
	m.BusinessScope = r.BusinessScope
	m.Stores = r.Stores
	m.Accounts = r.Accounts
	m.Bosses = r.Bosses
	m.Brands = r.Brands
	m.OperationsGroup = r.OperationsGroup
	m.OperationsStaff = r.OperationsStaff
	m.BusinessGroup = r.BusinessGroup
	m.BusinessPeople = r.BusinessPeople
	m.ReviewAccount = r.ReviewAccount
	m.CommentPassword = r.CommentPassword
	m.Email = r.Email
	m.Phone = r.Phone
	m.OtherTags = r.OtherTags
	m.Auditors = r.Auditors
	m.Type = r.Type
	m.Status = 0
	m.Remarks = r.Remarks
	if err := db.Create(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "门店创建失败", c)
		return
	}
	// 创建用户记录
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Create Member",
		Member:   m.ID,
		Contract: 0,
		Remarks:  token.FullName + "创建用户: " + m.Name,
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", m, c)
}

// UpdateMember 更新客户信息
func UpdateMember(c *gin.Context) {
	var r model.Member
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确"+err.Error(), c)
		return
	}
	db := config.GetMysql()
	var m model.Member
	if err := db.Where("id = ?", r.ID).First(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "客户ID不存在", c)
		return
	}
	if err := db.Updates(&r).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "门店更新失败", c)
		return
	}
	// 创建用户记录
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Review Member",
		Member:   m.ID,
		Remarks:  token.FullName + "更新用户: " + r.Name,
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", r, c)
}

// MemberList 客户列表
func MemberList(c *gin.Context) {
	var r struct {
		Name            string `json:"name"`
		Page            int    `json:"page"`
		Limit           int    `json:"limit"`
		Status          int    `json:"status"`
		OperationsStaff int    `json:"operations_staff"`
		BusinessPeople  int    `json:"business_people"`
	}
	var members []model.Member
	var count int64
	var pages int
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	sql := db

	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db.Where("id = ?", token.UserID).First(&user)
	var users []model.User
	var names []string
	var userID []uint
	_sql := db.Where("status = 1")
	if user.Duty > 1 {
		_sql.Where("duty = ?", user.Duty)
	}

	if user.Role > 1 {
		_sql.Where("marshalling_id = ?", user.MarshallingID)
	}

	if user.Role > 2 {
		_sql.Where("id = ?", user.ID)
	}
	_sql.Find(&users)

	for _, u := range users {
		names = append(names, u.RealName)
		userID = append(userID, u.ID)
	}

	if user.Duty == 2 { // 运营
		sql = sql.Where("operations_staff IN ?", userID)
	} else if user.Duty == 3 { // 业务
		sql = sql.Where("business_people IN ?", userID)
	}
	if r.Status != -1 {
		sql = sql.Where("status = ?", r.Status)
	}
	if r.Name != "" {
		sql = sql.Where("name like '%" + r.Name + "%'")
	}
	if r.BusinessPeople != 0 {
		sql = sql.Where("business_people = ?", r.BusinessPeople)
	}
	if r.OperationsStaff != 0 {
		sql = sql.Where("operations_staff = ?", r.OperationsStaff)
	}
	if r.Limit != 0 {
		sql = sql.Limit(r.Limit)
	} else {
		sql = sql.Limit(10)
	}
	var page int
	_count := sql
	_count.Find(&members).Order("id desc").Count(&count)
	if r.Page == 0 {
		page = 1
	} else {
		page = r.Page
	}
	sql.Limit(10).Offset((page - 1) * 10).Find(&members)

	if int(count)%10 != 0 {
		pages = int(count)/10 + 1
	} else {
		pages = int(count) / 10
	}
	handle.ReturnSuccess("ok", gin.H{"members": members, "pages": pages, "currPage": page}, c)
}

// MemberReview 客户审核
func MemberReview(c *gin.Context) {
	var r struct {
		ID     int    `json:"id" binding:"required"`
		Status int8   `json:"status" binding:"required"`
		Remark string `json:"remark"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var m model.Member
	if err := db.Where("id = ?", r.ID).First(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "客户ID不存在", c)
		return
	}
	m.Status = r.Status
	m.Reason = r.Reason
	if err := db.Save(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "审核失败", c)
		return
	}
	// 创建用户记录
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var msg string
	if r.Status == 1 {
		msg = token.FullName + "审核用户通过: " + m.Name
	} else if r.Status == 2 {
		msg = token.FullName + "审核用户拒绝: " + r.Remark
	}
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Review Member",
		Member:   m.ID,
		Remarks:  msg,
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", m, c)
}

// GetMemberDetails 获取客户详情
func GetMemberDetails(c *gin.Context) {
	var r struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	var member model.Member
	db := config.GetMysql()
	sql := db.Where("id = ?", r.ID)
	if err := sql.First(&member).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户不存在", c)
		return
	}
	var operation model.User
	var business model.User
	db.Where("id = ?", member.OperationsStaff).First(&operation)
	db.Where("id = ?", member.BusinessPeople).First(&business)
	handle.ReturnSuccess("ok", gin.H{
		"member":    member,
		"business":  business,
		"operation": operation,
	}, c)
}

func ExportMembers(c *gin.Context) {
	var r struct {
		Name            string `json:"name"`
		Status          int    `json:"status"`
		OperationsStaff int    `json:"operations_staff"`
		BusinessPeople  int    `json:"business_people"`
	}
	var members []model.Member
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	sql := db

	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db.Where("id = ?", token.UserID).First(&user)
	var users []model.User
	var names []string
	var userID []uint
	_sql := db.Where("status = 1")
	if user.Duty > 1 {
		_sql.Where("duty = ?", user.Duty)
	}

	if user.Role > 1 {
		_sql.Where("marshalling_id = ?", user.MarshallingID)
	}

	if user.Role > 2 {
		_sql.Where("id = ?", user.ID)
	}
	_sql.Find(&users)

	for _, u := range users {
		names = append(names, u.RealName)
		userID = append(userID, u.ID)
	}

	if user.Duty == 2 { // 运营
		sql = sql.Where("operations_staff IN ?", userID)
	} else if user.Duty == 3 { // 业务
		sql = sql.Where("business_people IN ?", userID)
	}
	if r.Status != -1 {
		sql = sql.Where("status = ?", r.Status)
	}
	if r.Name != "" {
		sql = sql.Where("name like '%" + r.Name + "%'")
	}
	if r.BusinessPeople != 0 {
		sql = sql.Where("business_people = ?", r.BusinessPeople)
	}
	if r.OperationsStaff != 0 {
		sql = sql.Where("operations_staff = ?", r.OperationsStaff)
	}

	sql.Find(&members)

	head := []string{"表头一", "表头二", "表头三"}
	var body [][]interface{}
	for _, member := range members {
		memberInfo := []interface{}{
			member.UUID,
			member.Accounts,
		}
		body = append(body, memberInfo)
	}

	// body := [][]interface{}{{1, "2020", ""}, {2, "2019", ""}, {3, "2018", ""}}
	filename := "test.xlsx"
	handle.ExcelExport(c, head, body, filename)
}
