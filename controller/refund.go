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

// CreateRefund 创建一个申请退款
func CreateRefund(c *gin.Context) {
	var r struct {
		ID     uint    `json:"id" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
		Remark string  `json:"remark"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var contract model.Contract
	if err := db.Where("id = ?", r.ID).First(&contract).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约ID不存在", c)
		return
	}
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	refund := model.Refund{
		ContractID:      r.ID,
		Amount:          r.Amount,
		Applicant:       token.UserID,
		Remark:          r.Remark,
		OperationsStaff: contract.OperationsStaff,
		BusinessPeople:  contract.BusinessPeople,
	}
	if err := db.Create(&refund).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "创建退款申请失败", c)
		return
	}
	handle.ReturnSuccess("ok", refund, c)
}

// ReviewRefund 审核退款
func ReviewRefund(c *gin.Context) {
	var r struct {
		ID     uint   `json:"id" binding:"required"`
		Status int8   `json:"status"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var refund model.Refund
	var contract model.Contract
	var member model.Member
	if err := db.Where("id = ?", r.ID).First(&refund).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "退款申请不存在", c)
		return
	}
	db.Where("id = ?", refund.ContractID).First(&contract)
	db.Where("id = ?", contract.MemberID).First(&member)
	refund.Status = r.Status
	refund.Reason = r.Reason
	contract.Refund = model.MyTime{Time: time.Now()}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Save(&refund).Error; err != nil {
			return err
		}

		if err := db.Save(&contract).Error; err != nil {
			return err
		}

		member.Refund = model.MyTime{Time: time.Now()}
		if err := db.Save(&member).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "审核失败", c)
		return
	}

	handle.ReturnSuccess("ok", refund, c)
}

// GetRefundList 获取审核列表
func GetRefundList(c *gin.Context) {
	var r struct {
		Contract uint   `json:"contract"`
		Key      string `json:"key"`
		Status   int8   `json:"status"`
		Page     int    `json:"page"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var refunds []model.Refund
	var count int64
	var pages int
	sql := db
	if r.Contract != 0 {
		sql.Where("contract_id = ?", r.Contract)
	}
	var page int
	_count := sql
	_count.Find(&refunds).Order("id desc").Count(&count)
	if r.Page == 0 {
		page = 1
	} else {
		page = r.Page
	}

	if r.Key != "" {
		_db := config.GetMysql()
		var members []model.Member
		var ids []uint
		var contracts []model.Contract
		var _ids []uint
		_db.Where("name LIKE ?", "%"+r.Key+"%").Find(&members)
		for _, member := range members {
			ids = append(ids, member.ID)
		}
		_db.Where("member_id IN ?", ids).Find(&contracts)
		for _, contract := range contracts {
			_ids = append(_ids, contract.ID)
		}

		sql = sql.Where("contract_id IN ?", _ids)
	}

	sql.Limit(10).Offset((page - 1) * 10).Find(&refunds)

	if int(count)%10 != 0 {
		pages = int(count)/10 + 1
	} else {
		pages = int(count) / 10
	}
	handle.ReturnSuccess("ok", gin.H{"refunds": refunds, "pages": pages, "currPage": page}, c)
}

// GetRefundDetails 获取退款详情
func GetRefundDetails(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	var refund model.Refund
	db := config.GetMysql()
	if err := db.Where("id = ?", r.ID).First(&refund).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "退款申请不存在", c)
		return
	}
	handle.ReturnSuccess("ok", refund, c)
}
