package controller

import (
	"fmt"
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
	var thisMonthRenewal int64
	var lastMonthRenewal int64
	var thisMonthBreak int64
	var lastMonthBreak int64
	var thisMonthRefund int64
	var lastMonthRefund int64
	var thisMonthClient int64
	var lastMonthClient int64
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
	fmt.Println(names)
	sql := db.Model(&model.Contract{}).Where("status = 1").Where("operations_staff IN ?", names)

	var _thisMonthNewly = sql
	var _lastMonthNewly = sql
	var _thisMonthRenewal = sql
	var _lastMonthRenewal = sql
	var _thisMonthBreak = sql
	var _lastMonthBreak = sql
	var _thisMonthRefund = sql
	var _lastMonthRefund = sql
	// 新增
	_thisMonthNewly.Where("cooperation_time >= ? AND cooperation_time < ?", thisMonthStart, thisMonthEnd).Where("sort = 0").Count(&thisMonthNewly)
	_lastMonthNewly.Where("cooperation_time >= ? AND cooperation_time < ?", lastMonthStart, lastMonthEnd).Where("sort = 0").Count(&lastMonthNewly)

	// 续签
	_thisMonthRenewal.Where("cooperation_time >= ? AND cooperation_time < ?", thisMonthStart, thisMonthEnd).Where("sort > 0").Count(&thisMonthRenewal)
	_lastMonthRenewal.Where("cooperation_time >= ? AND cooperation_time < ?", lastMonthStart, lastMonthEnd).Where("sort > 0").Count(&lastMonthRenewal)

	// 断约
	_thisMonthBreak.Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Where("expire_time >= ? AND expire_time < ?", thisMonthStart, thisMonthEnd)
	}).Where("delay_time >= ? AND delay_time < ?", thisMonthStart, thisMonthEnd).Where("refund IS NULL").Count(&thisMonthBreak)
	_lastMonthBreak.Preload("Member", func(db *gorm.DB) *gorm.DB {
		return db.Where("expire_time >= ? AND expire_time < ?", lastMonthStart, lastMonthEnd)
	}).Where("delay_time >= ? AND delay_time < ?", lastMonthStart, lastMonthEnd).Where("refund IS NULL").Count(&lastMonthBreak)

	// 退款
	_thisMonthRefund.Where("refund >= ? AND refund < ?", thisMonthStart, thisMonthEnd).Count(&thisMonthRefund)
	_lastMonthRefund.Where("refund >= ? AND refund < ?", lastMonthStart, lastMonthEnd).Count(&lastMonthRefund)

	// 总服务数
	// _thisMonthClient.Where("delay_time < ?", thisMonthEnd).Count(&thisMonthClient)
	// _lastMonthClient.Where("delay_time < ?", lastMonthEnd).Count(&lastMonthClient)
	db.Model(model.Member{}).Where("operations_staff = ? AND created_at <= ?", user.ID, thisMonthEnd).Count(&thisMonthClient)
	db.Model(model.Member{}).Where("operations_staff = ? AND created_at <= ?", user.ID, lastMonthEnd).Count(&lastMonthClient)
	// handle.ReturnSuccess("ok", contracts, c)
	handle.ReturnSuccess("ok", gin.H{
		"newly":   gin.H{"this_month": thisMonthNewly, "last_month": lastMonthNewly},
		"renewal": gin.H{"this_month": thisMonthRenewal, "last_month": lastMonthRenewal},
		"break":   gin.H{"this_month": thisMonthBreak, "last_month": lastMonthBreak},
		"refund":  gin.H{"this_month": thisMonthRefund, "last_month": lastMonthRefund},
		"client":  gin.H{"this_month": thisMonthClient, "last_month": lastMonthClient},
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
	db.Preload("ContractTask").Where("operations_staff = ?", user.RealName).
		Where("status = 1 AND refund IS NULL").Find(&contracts)
	handle.ReturnSuccess("ok", contracts, c)
}

// ServiceDays30 30天服务数
func ServiceDays30(c *gin.Context) {
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
	end, _ := time.ParseInLocation("2006-01-02", r.Date, time.Local)
	// start := end.AddDate(0, 0, -30)
	var count int64
	db.Model(&model.Contract{}).
		Where("status = 1 AND operations_staff = ?", token.FullName).
		Where("delay_time >= ?", end).Count(&count)
	handle.ReturnSuccess("ok", count, c)
}
