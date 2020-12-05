package controller

import (
	"net/http"
	"time"

	"gorm.io/gorm"

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
	// var contracts []model.Contract
	// db.Where("status = 1").Where("operations_staff in ?", names).Find(&contracts)
	var thisMonthNewly int64
	var lastMonthNewly int64
	var thisMonthBeexpire int64
	var lastMonthBeexpire int64
	var thisMonthBreak int64
	var lastMonthBreak int64
	var thisMonthStart time.Time
	// 查询月开始时间
	if r.Date != "" {
		thisMonthStart, _ = time.ParseInLocation("2006-01", r.Date, time.Local)
	} else {
		thisMonthStart, _ = time.ParseInLocation("2006-01", time.Now().Format("2006-01"), time.Local)
	}

	thisMonthEnd := thisMonthStart.AddDate(0, 1, 0)    // 查询月结束时间
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0) // 对比月开始时间
	lastMonthEnd := lastMonthStart.AddDate(0, 1, 0)    // 对比月结束时间

	sql := db.Model(&model.Contract{}).Where("sort = 0 AND status = 1").Where("operations_staff IN ?", names)

	// 新增
	sql.Where("cooperation_time >= ? AND cooperation_time < ?", thisMonthStart, thisMonthEnd).Count(&thisMonthNewly)
	sql.Where("cooperation_time >= ? AND cooperation_time < ?", lastMonthStart, lastMonthEnd).Count(&lastMonthNewly)

	// 续签
	sql.Where("cooperation_time >= ? AND cooperation_time < ?", thisMonthStart, thisMonthEnd).Where("sort > 0").Count(&thisMonthBeexpire)
	sql.Where("cooperation_time >= ? AND cooperation_time < ?", lastMonthStart, lastMonthEnd).Where("sort > 0").Count(&lastMonthBeexpire)

	// 断约
	sql.Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Where("expire_time >= ? AND expire_time < ?", thisMonthStart, thisMonthEnd)
	}).Where("delay_time >= ? AND delay_time < ?", thisMonthStart, thisMonthEnd).Where("refund IS NULL").Count(&thisMonthBreak)
	sql.Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Where("expire_time >= ? AND expire_time < ?", lastMonthStart, lastMonthEnd)
	}).Where("delay_time >= ? AND delay_time < ?", lastMonthStart, lastMonthEnd).Where("refund IS NULL").Count(&lastMonthBreak)

	// handle.ReturnSuccess("ok", contracts, c)
	handle.ReturnSuccess("ok", gin.H{
		"newly":    gin.H{"this_month": thisMonthNewly, "last_month": lastMonthNewly},
		"beexpire": gin.H{"this_month": thisMonthBeexpire, "last_month": lastMonthBeexpire},
		"break":    gin.H{"this_month": thisMonthBreak, "last_month": lastMonthBreak},
	}, c)
}

// MyContract 首页-我的合约
func MyContract(c *gin.Context) {
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db := config.MysqlConn
	db.Where("id = ?", token.UserID).First(&user)
	var contracts []model.Contract
	db.Where("operations_staff = ?", user.RealName).Where("status = 1").Where("refund IS NULL").Find(&contracts)
	handle.ReturnSuccess("ok", contracts, c)
}
