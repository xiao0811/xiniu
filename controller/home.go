package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CountData 首页统计数据
func CountData(c *gin.Context) {
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db := config.MysqlConn
	db.Where("id = ?", token.UserID).First(&user)
	handle.ReturnSuccess("ok", user, c)
}
