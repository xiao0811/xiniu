package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// LoginRequest 请求登录结构体
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var (
	// ExpireTime token过期时间
	ExpireTime = 3600
)

// Login 用户登录
func Login(c *gin.Context) {
	db := config.GetMysql()
	var lq LoginRequest
	var user model.User
	if err := c.ShouldBind(&lq); err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户名密码输入不正确", c)
		return
	}

	if err := db.Where("phone = ?", lq.Phone).First(&user).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户名密码输入不正确", c)
		return
	}

	if !handle.CheckPasswordHash(lq.Password, user.Password) {
		handle.ReturnError(http.StatusBadRequest, "密码错误", c)
		return
	}
	claims := &handle.JWTClaims{
		UserID:      user.ID,
		Username:    lq.Phone,
		FullName:    user.RealName,
		Permissions: []string{},
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	signedToken, err := handle.GetToken(claims)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	// c.String(http.StatusOK, signedToken)
	handle.ReturnSuccess("ok", signedToken, c)
}


