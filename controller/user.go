package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

type CreateUserRequest struct {
	Phone          string       `json:"phone" binding:"required"`
	Password       string       `json:"password"`
	RealName       string       `json:"real_name"`
	Gender         uint8        `json:"gender"`
	Birthday       model.MyTime `json:"birthday"`
	Identification string       `json:"identification"`
	Role           uint8        `json:"role"`
	Marshalling    uint8        `json:"marshalling"`
}

// CreateUser 创建新的管理员
func CreateUser(c *gin.Context) {
	var cr CreateUserRequest
	var user model.User
	db := config.GetMysql()
	if err := c.ShouldBind(&cr); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db.Where("phone = ?", cr.Phone).First(&user)
	if user.ID > 0 {
		handle.ReturnError(http.StatusBadRequest, "用户名重复", c)
		return
	}
	if cr.Password == "" {
		cr.Password, _ = handle.HashPassword("123456")
	}
	user = model.User{
		Phone:          cr.Phone,
		Password:       cr.Password,
		RealName:       cr.RealName,
		Gender:         cr.Gender,
		Birthday:       cr.Birthday,
		Identification: cr.Identification,
		Role:           cr.Role,
		Marshalling:    cr.Marshalling,
	}

	if err := db.Create(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "添加用户失败", c)
		return
	}
	handle.ReturnSuccess("ok", user, c)
}
