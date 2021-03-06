package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
	"gorm.io/gorm"
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
	db := config.GetMysql()
	db.Where("id = ?", token.UserID).First(&user)
	var users []model.User
	var names []string
	var userID []uint
	if user.Duty == 1 || user.Role == 1 {
		db.Where("status = 1").Find(&users)
	} else if user.Role == 2 {
		db.Where("duty = 2 AND status = 1 AND marshalling_id = ?", user.MarshallingID).Find(&users)
	} else {
		db.Where("id = ?", user.ID).Find(&users)
	}
	for _, u := range users {
		names = append(names, u.RealName)
		userID = append(userID, u.ID)
	}
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
	handle.ReturnSuccess("ok", gin.H{
		"newly": gin.H{
			"this_month": GetNewly(thisMonthStart, thisMonthEnd, names),
			"last_month": GetNewly(lastMonthStart, lastMonthEnd, names),
		},
		"renewal": gin.H{
			"this_month": GetRenewal(thisMonthStart, thisMonthEnd, names),
			"last_month": GetRenewal(lastMonthStart, lastMonthEnd, names),
		},
		"break": gin.H{
			"this_month": GetBreak(thisMonthStart, thisMonthEnd, names),
			"last_month": GetBreak(lastMonthStart, lastMonthEnd, names),
		},
		"refund": gin.H{
			"this_month": GetRefund(thisMonthStart, thisMonthEnd, names),
			"last_month": GetRefund(lastMonthStart, lastMonthEnd, names),
		},
		"client": gin.H{
			"this_month": GetClint(thisMonthEnd, userID),
			"last_month": GetClint(lastMonthEnd, userID),
		},
	}, c)
}

// MyContract 首页-我的合约
func MyContract(c *gin.Context) {
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db := config.GetMysql()
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
	db := config.GetMysql()
	db.Where("id = ?", token.UserID).First(&user)
	end, _ := time.ParseInLocation("2006-01-02", r.Date, time.Local)
	// start := end.AddDate(0, 0, -30)
	var count int64
	db.Model(&model.Contract{}).
		Where("status = 1 AND operations_staff = ?", token.FullName).
		Where("delay_time >= ?", end).Count(&count)
	handle.ReturnSuccess("ok", count, c)
}

// GetNewly 获取新签
func GetNewly(start, end time.Time, names []string) int {
	db := config.GetMysql()
	var contracts []model.Contract
	db.Where("status = 1 AND sort = 0").Where("operations_staff IN ?", names).
		Where("cooperation_time >= ? AND cooperation_time < ?", start, end).
		Find(&contracts)
	return len(contracts)
}

// GetRenewal 获取续约
func GetRenewal(start, end time.Time, names []string) int {
	db := config.GetMysql()
	var contracts []model.Contract
	db.Where("status = 1 AND sort > 0").Where("operations_staff IN ?", names).
		Where("cooperation_time >= ? AND cooperation_time < ?", start, end).
		Find(&contracts)
	return len(contracts)
}

// GetBreak 获取断约
func GetBreak(start, end time.Time, names []string) int {
	db := config.GetMysql()
	var contracts []model.Contract
	if start.Format("2006-01") == time.Now().Format("2006-01") {
		end = time.Now()
	}
	db.Where("status = 1").Where("operations_staff IN ?", names).
		Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Where("expire_time >= ? AND expire_time < ?", start, end)
		}).Where("delay_time >= ? AND delay_time < ?", start, end).
		Where("refund IS NULL").Find(&contracts)
	return len(contracts)
}

// GetRefund 获取退款
func GetRefund(start, end time.Time, names []string) int {
	db := config.GetMysql()
	var contracts []model.Contract
	db.Where("status = 1").Where("operations_staff IN ?", names).
		Where("refund >= ? AND refund < ?", start, end).Find(&contracts)
	return len(contracts)
}

// GetClint 获取客户总数
func GetClint(end time.Time, names []uint) int {
	db := config.GetMysql()
	var members []model.Member
	db.Where("operations_staff in ? AND created_at <= ?", names, end).Find(&members)
	return len(members)
}
