package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateContractLog 新建合约记录
func CreateContractLog(c *gin.Context) {
	var r struct {
		ContractID   uint   `json:"contrat_id"`    // 合约ID
		OperatorID   uint   `json:"operator_id"`   // 操作人员ID
		Operator     string `json:"operator"`      // 操作人员名字
		Type         uint8  `json:"type"`          // 类型: 1 牌级 2 推广通
		GradeScore   uint8  `json:"grade_score"`   // 等级分数 - 牌级
		Spend        int    `json:"spend"`         // 花费 - 推广通
		GuestCapital int    `json:"guest_capital"` // 客资 - 推广通
		Week         uint8  `json:"week"`          // 第几周
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()

	cl := model.ContratLog{
		ContractID:    r.ContractID,
		OperatorID:    r.OperatorID,
		Operator:      r.Operator,
		Type:          r.Type,
		GradeScore:    r.GradeScore,
		Spend:         r.Spend,
		GuestCapital:  r.GuestCapital,
		Week:          r.Week,
		OperatingTime: model.MyTime{Time: time.Now()},
	}

	if err := db.Create(&cl).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "合约记录创建失败", c)
		return
	}

	handle.ReturnSuccess("ok", cl, c)
}

// DeleteContractLog 删除合约操作记录
func DeleteContractLog(c *gin.Context) {
	var r struct {
		ID uint `gorm:"primarykey" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	var cl model.ContratLog
	if err := db.Where("id = ?", r.ID).First(&cl).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "合约操作记录不存在", c)
		return
	}

	if err := db.Delete(&cl).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约操作记录删除成功", c)
		return
	}

	handle.ReturnSuccess("ok", cl, c)
}

// GetLogsByContratID 根据合约ID获取操作记录
func GetLogsByContratID(c *gin.Context) {
	var r struct {
		ID uint `gorm:"primarykey" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	var cls []model.ContratLog

	db.Where("contract_id = ?", r.ID).Find(&cls)
	handle.ReturnSuccess("ok", cls, c)
}
