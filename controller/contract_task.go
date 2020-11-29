package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateContractTask 添加合约任务记录
func CreateContractTask(c *gin.Context) {
	var r struct {
		Type            uint8    `json:"type"`
		ContractID      uint     `json:"contract_id"`
		OperationsStaff string   `json:"operations_staff" gorm:"type:varchar(20)"`
		TaskCount       uint8    `json:"task_count"`
		CompleteTime    string   `json:"complete_time"`
		Images          []string `json:"images"`
		Status          uint8    `json:"status"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	_CompleteTime, err := time.ParseInLocation("2006-01-02", r.CompleteTime, time.Local)
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "时间格式不正确", c)
		return
	}
	ij, _ := json.Marshal(r.Images)
	ct := model.ContractTask{
		Type:            r.Type,
		ContractID:      r.ContractID,
		OperationsStaff: r.OperationsStaff,
		TaskCount:       r.TaskCount,
		CompleteTime:    model.MyTime{Time: _CompleteTime},
		Images:          string(ij),
		Status:          1,
	}
	if err := db.Create(&ct).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "任务记录创建失败", c)
		return
	}
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Create Contract Task",
		Contract: r.ContractID,
		Remarks:  token.FullName + "创建合约任务: " + strconv.Itoa(int(r.ContractID)),
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", ct, c)
}

// DeleteContractTask 删除合约任务
func DeleteContractTask(c *gin.Context) {
	var r struct {
		ID uint `json:"id"  binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var ct model.ContractTask
	if err := db.Where("id = ?", r.ID).First(&ct).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "任务ID错误", c)
		return
	}
	ct.Status = 0
	if err := db.Save(&ct).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "删除失败", c)
		return
	}
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Create Contract Task",
		Contract: r.ID,
		Remarks:  token.FullName + "删除合约任务: " + strconv.Itoa(int(r.ID)),
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", ct, c)
}

// GetContractTaskList 获取任务列表
func GetContractTaskList(c *gin.Context) {
	var r struct {
		Type            uint8  `json:"type"`
		ContractID      uint   `json:"contract_id"`
		OperationsStaff string `json:"operations_staff"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	sql := db
	var cts []model.ContractTask
	if r.Type != 0 {
		sql = sql.Where("type = ?", r.Type)
	}
	if r.ContractID != 0 {
		sql = sql.Where("contract_id = ?", r.ContractID)
	}
	if r.OperationsStaff != "" {
		sql = sql.Where("operations_staff = ?", r.OperationsStaff)
	}
	sql.Order("created_at desc").Find(&cts)

	handle.ReturnSuccess("ok", cts, c)
}
