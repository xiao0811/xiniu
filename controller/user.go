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

// CreateUserRequest .
type CreateUserRequest struct {
	Phone          string `json:"phone" binding:"required"`
	Password       string `json:"password"`
	RealName       string `json:"real_name"`
	Gender         uint8  `json:"gender"`
	Birthday       string `json:"birthday"`
	Identification string `json:"identification"`
	Role           uint8  `json:"role"`
	Duty           int8   `json:"duty"`
	MarshallingID  uint   `json:"marshalling_id"`
}

// CreateUser 创建新的管理员
func CreateUser(c *gin.Context) {
	var r CreateUserRequest
	var user model.User
	db := config.GetMysql()
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db.Where("phone = ?", r.Phone).First(&user)
	if user.ID > 0 {
		handle.ReturnError(http.StatusBadRequest, "用户名重复", c)
		return
	}
	if r.Password == "" {
		r.Password, _ = handle.HashPassword("123456")
	}
	ps, _ := handle.HashPassword(r.Password)
	var m uint
	if r.MarshallingID == 0 {
		m = 1
	} else {
		m = r.MarshallingID
	}
	dt, _ := time.ParseInLocation(model.TimeFormat, r.Birthday, time.Local)
	user = model.User{
		Phone:          r.Phone,
		Password:       ps,
		RealName:       r.RealName,
		Gender:         r.Gender,
		Birthday:       model.MyTime{Time: dt},
		Identification: r.Identification,
		Role:           r.Role,
		MarshallingID:  m,
		Status:         1,
		Duty:           r.Duty,
	}

	if err := db.Create(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "添加用户失败", c)
		return
	}
	handle.ReturnSuccess("ok", user, c)
}

// GetUserDetails 获取用户详情
func GetUserDetails(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	var user model.User
	db := config.GetMysql()
	if err := db.Where("id = ?", r.ID).Preload("Marshalling").First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户未找到", c)
		return
	}
	handle.ReturnSuccess("ok", user, c)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	var r model.User
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}

	var user model.User
	db := config.GetMysql()
	if err := db.Where("id = ?", r.ID).Preload("Marshalling").First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户未找到", c)
		return
	}
	if err := db.Model(&user).Updates(r).Error; err != nil {
		handle.ReturnError(http.StatusServiceUnavailable, "用户信息更新失败", c)
		return
	}
	handle.ReturnSuccess("ok", user, c)
}

// SendChangePasswordMessage 发送修改密码短信
func SendChangePasswordMessage(c *gin.Context) {
	var r struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	rc := config.GetRedis()
	defer rc.Close()
	var user model.User
	if err := db.Where("phone = ?", r.Phone).First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "手机号码不正确", c)
		return
	}
	_, err := rc.Get("change_password_" + r.Phone).Result()
	if err == nil {
		handle.ReturnError(http.StatusBadRequest, "验证码已发送,请勿重复操作", c)
		return
	}
	randNum := handle.RandInt(100000, 999999)
	rc.Set("change_password_"+r.Phone, randNum, 10*time.Minute)
	sms := handle.Info{
		MsgText:    "修改密码验证码: " + strconv.Itoa(randNum),
		Destmobile: r.Phone,
	}
	task, err := sms.Send()
	if err != nil {
		handle.ReturnError(http.StatusInternalServerError, "短信发送失败", c)
		return
	}
	handle.ReturnSuccess("ok", "发送成功:"+task, c)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var r struct {
		Phone       string `json:"phone" binding:"required"`
		OldPassword string `json:"old_password"`
		Password    string `json:"password" binding:"required"`
		Compulsory  bool   `json:"compulsory"` // 忘记旧密码强制更新
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}

	var user model.User
	db := config.GetMysql()
	if err := db.Where("phone = ?", r.Phone).First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "手机号码不正确", c)
		return
	}

	if !r.Compulsory {
		if !handle.CheckPasswordHash(r.OldPassword, user.Password) {
			handle.ReturnError(http.StatusInternalServerError, "旧密码错误", c)
			return
		}
	}

	ps, _ := handle.HashPassword(r.Password)
	user.Password = ps
	if err := db.Save(&user).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "密码修改失败", c)
		return
	}
	handle.ReturnSuccess("ok", "密码修改成功", c)
}

// UserList 用户列表
func UserList(c *gin.Context) {
	var r struct {
		RealName      string `json:"real_name"`
		Status        int8   `json:"status"`
		MarshallingID uint   `json:"marshalling_id"`
		Role          uint8  `json:"role"`
		Limit         int    `json:"limit"`
		Offset        int    `json:"offset"`
		Duty          int8   `json:"duty"`
		Page          int    `json:"page"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()

	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db.Where("id = ?", token.UserID).First(&user)
	var _users []model.User
	// var names []string
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
	_sql.Find(&_users)

	for _, u := range _users {
		// names = append(names, u.RealName)
		userID = append(userID, u.ID)
	}

	var users []model.User
	var count int64
	var page int
	var pages int
	sql := db.Preload("Marshalling").Where("status = 1").Where("id IN ?", userID)
	if r.RealName != "" {
		sql = sql.Where("real_name like '%"+r.RealName+"%'").Or("phone", r.RealName)
	}
	if r.Status != 0 {
		sql = sql.Where("status = ?", r.Status)
	}
	if r.MarshallingID != 0 {
		sql = sql.Where("marshalling_id = ?", r.MarshallingID)
	}
	if r.Role != 0 {
		sql = sql.Where("role = ?", r.Role)
	}
	if r.Duty != 0 {
		sql = sql.Where("duty = ?", r.Duty)
	}
	if r.Limit != 0 {
		sql = sql.Limit(r.Limit)
	} else {
		sql = sql.Limit(10)
	}
	_count := sql
	_count.Find(&users).Order("id desc").Count(&count)
	if r.Page == 0 {
		page = 1
	} else {
		page = r.Page
	}
	sql.Limit(10).Offset((page - 1) * 10).Find(&users)

	if int(count)%10 != 0 {
		pages = int(count)/10 + 1
	} else {
		pages = int(count) / 10
	}
	handle.ReturnSuccess("ok", gin.H{"user": users, "pages": pages, "currPage": page}, c)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var user model.User
	if err := db.Where("id = ?", r.ID).First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户不存在", c)
		return
	}
	// user.Status = 0
	if err := db.Delete(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户删除失败", c)
		return
	}

	handle.ReturnSuccess("ok", user, c)
}

// UserBatchGroup 用户批量分组
func UserBatchGroup(c *gin.Context) {
	var r struct {
		GroupID uint  `json:"groupId" binding:"required"`
		UserIds []int `json:"userids"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var group model.Marshalling
	if err := db.Where("id = ?", r.GroupID).First(&group).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "组别不存在", c)
		return
	}
	db.Model(model.User{}).Where("id IN ?", r.UserIds).Updates(model.User{
		MarshallingID: r.GroupID,
	})
	handle.ReturnSuccess("ok", "批量修改成功", c)
}

func GetMember(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	var members []model.Member
	db := config.GetMysql()
	if err := db.Where("operations_staff = ? OR business_people = ?", r.ID, r.ID).Find(&members).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "", c)
	}
	handle.ReturnSuccess("ok", members, c)
}

// UserListAll 获取用户信息
func UserListAll(c *gin.Context) {
	var r struct {
		MarshallingID uint8 `json:"marshalling_id"`
		DutyID        uint8 `json:"duty_id"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var users []model.User
	if r.MarshallingID != 0 {
		db = db.Where("marshalling_id = ?", r.MarshallingID)
	}

	if r.DutyID != 0 {
		db = db.Where("duty = ?", r.DutyID)
	}

	db.Find(&users)
	handle.ReturnSuccess("ok", users, c)
}
