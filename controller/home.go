package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CountData 首页统计数据
func CountData(c *gin.Context) {
	var r struct {
		Date string `json:"date"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db := config.MysqlConn
	db.Where("id = ?", token.UserID).First(&user)
	var users []model.User
	var names []string
	if user.Duty == 1 || user.Role == 1 {
		db.Where("duty = 2 AND status = 1").Find(&users)
	} else if user.Role == 2 {
		db.Where("duty = 2 AND status = 1 AND MarshallingID = ?", user.MarshallingID).Find(&users)
	} else {
		db.Where("id = ?", user.ID).Find(&users)
	}
	for _, u := range users {
		names = append(names, u.RealName)
	}
	var contracts []model.Contract
	db.Where("status = 1").Where("operations_staff in ?", names).Find(&contracts)
	handle.ReturnSuccess("ok", contracts, c)
}
