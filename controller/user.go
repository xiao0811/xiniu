package controller

import (
	"net/http"

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
	Marshalling    uint         `json:"marshalling"`
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
	user = model.User{
		Phone:          r.Phone,
		Password:       r.Password,
		RealName:       r.RealName,
		Gender:         r.Gender,
		Birthday:       r.Birthday,
		Identification: r.Identification,
		Role:           r.Role,
		MarshallingID:  r.Marshalling,
	}

	if err := db.Create(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "添加用户失败", c)
		return
	}
	handle.ReturnSuccess("ok", user, c)
}

// GetUserDetailsRequest 获取用户详情结构体
type GetUserDetailsRequest struct {
	ID uint `json:"id" binding:"required"`
}

// GetUserDetails 获取用户详情
func GetUserDetails(c *gin.Context) {
	var r GetUserDetailsRequest
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
