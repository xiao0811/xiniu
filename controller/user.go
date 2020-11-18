package controller

import (
	"fmt"
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
	Phone          string       `json:"phone" binding:"required"`
	Password       string       `json:"password"`
	RealName       string       `json:"real_name"`
	Gender         uint8        `json:"gender"`
	Birthday       model.MyTime `json:"birthday"`
	Identification string       `json:"identification"`
	Role           uint8        `json:"role"`
	MarshallingID  uint         `json:"marshalling_id"`
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
	user = model.User{
		Phone:          r.Phone,
		Password:       ps,
		RealName:       r.RealName,
		Gender:         r.Gender,
		Birthday:       r.Birthday,
		Identification: r.Identification,
		Role:           r.Role,
		MarshallingID:  m,
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
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	rc := config.GetRedis()
	defer rc.Close()
	code, _ := rc.Get("change_password_" + r.Phone).Result()

	fmt.Println(code, r.Code)
	if code != r.Code {
		handle.ReturnError(http.StatusBadRequest, "验证码错误或过期", c)
		return
	}

	var user model.User
	db := config.GetMysql()
	if err := db.Where("phone = ?", r.Phone).First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "手机号码不正确", c)
		return
	}
	ps, _ := handle.HashPassword(r.Password)
	user.Password = ps
	if err := db.Save(&user).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "密码修改失败", c)
		return
	}
	handle.ReturnSuccess("ok", "密码修改成功", c)
}
