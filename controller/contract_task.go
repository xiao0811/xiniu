package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

var TaskType = map[uint16]string{
	1:   "收藏",
	2:   "浏览量",
	4:   "团购卖量",
	5:   "点赞",
	7:   "预约",
	8:   "访客",
	3:   "小红书",
	33:  "评价",
	333: "笔记",
}

// CreateContractTask 添加合约任务记录
func CreateContractTask(c *gin.Context) {
	var r struct {
		Type            uint8    `json:"type"`                                     // 任务类型
		ContractID      uint     `json:"contract_id"`                              // 合约ID
		OperationsStaff string   `json:"operations_staff" gorm:"type:varchar(20)"` // 运营人员
		TaskCount       uint8    `json:"task_count"`                               // 总任务量
		Initial         uint     `json:"initial"`                                  // 初始值
		CompleteTime    string   `json:"complete_time"`                            // 完成时间
		StoreLink       string   `json:"store_link"`                               // 门店链接
		Requirements    string   `json:"requirements"`                             // 任务要求
		Mediator        string   `json:"mediator"`                                 // 媒介人员
		Images          []string `json:"images"`                                   // 图片
		Status          uint8    `json:"status"`                                   // 状态 1: 新建/未完成 10: 完成 20: 取消
		Remark          string   `json:"remark"`                                   // 备注
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	_CompleteTime, err := time.ParseInLocation("2006-01-02", r.CompleteTime, time.Local)
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "时间格式不正确", c)
		return
	}

	var member model.Member
	var contreact model.Contract
	db.Where("id = ?", r.ContractID).First(&contreact)
	db.Where("id = ?", contreact.MemberID).First(&member)

	ij, _ := json.Marshal(r.Images)
	ct := model.ContractTask{
		Type:            r.Type,
		ContractID:      r.ContractID,
		Member:          member.Name,
		OperationsStaff: r.OperationsStaff,
		TaskCount:       r.TaskCount,
		CompleteTime:    model.MyTime{Time: _CompleteTime},
		Images:          string(ij),
		Status:          1,
		Remark:          r.Remark,
		Mediator:        r.Mediator,
		StoreLink:       r.StoreLink,
		Requirements:    r.Requirements,
		Initial:         r.Initial,
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
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
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
		Type            uint8        `json:"type"`
		ContractID      uint         `json:"contract_id"`
		OperationsStaff string       `json:"operations_staff"`
		StartTime       model.MyTime `json:"start_time"`
		EndTime         model.MyTime `json:"end_time"`
		Pagination      bool         `json:"pagination"`
		Page            int          `json:"page"`
		Limit           int          `json:"limit"`
		Status          uint8        `json:"status"`
		Member          string       `json:"member"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
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

	if r.Member != "" {
		sql = sql.Where("member LIKE ?", "?"+r.Member+"")
	}

	if r.Status != 0 {
		sql = sql.Where("status = ?", r.Status)
	}

	if r.StartTime.Format("2006-01-02 15:04:05") != "0001-01-01 00:00:00" {
		sql = sql.Where("complete_time > ?", r.StartTime)
	}

	if r.EndTime.Format("2006-01-02 15:04:05") != "0001-01-01 00:00:00" {
		sql = sql.Where("complete_time < ?", r.EndTime)
	}

	var data gin.H
	if r.Pagination {
		var count int64
		var pages int
		var page int
		_count := sql
		_count.Find(&cts).Order("id desc").Count(&count)
		if r.Page == 0 {
			page = 1
		} else {
			page = r.Page
		}
		sql.Limit(10).Offset((page - 1) * 10).Find(&cts)

		if int(count)%10 != 0 {
			pages = int(count)/10 + 1
		} else {
			pages = int(count) / 10
		}
		data = gin.H{"tasks": cts, "pages": pages, "currPage": page}
	} else {
		sql.Order("created_at desc").Find(&cts)
		data = gin.H{"tasks": cts}
	}

	handle.ReturnSuccess("ok", data, c)
}

// UpdateContractTask2 更新合约任务
func UpdateContractTask2(c *gin.Context) {
	var r model.ContractTask
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	var td model.ContractTask
	if err := db.Where("id = ?", r.ID).First(&td).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "任务ID不存在", c)
		return
	}
	if err := db.Updates(&r).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "任务更新失败", c)
		return
	}

	handle.ReturnSuccess("ok", r, c)
}

// ExportContractTask 导出合约列表
func ExportContractTask(c *gin.Context) {
	var r struct {
		Type            []uint8      `json:"type"`
		ContractID      uint         `json:"contract_id"`
		OperationsStaff string       `json:"operations_staff"`
		StartTime       model.MyTime `json:"start_time"`
		EndTime         model.MyTime `json:"end_time"`
		Pagination      bool         `json:"pagination"`
		Status          uint8        `json:"status"`
		Member          string       `json:"member"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	sql := db
	var cts []model.ContractTask
	if r.Type != nil {
		sql = sql.Where("type IN ?", r.Type)
	}
	if r.ContractID != 0 {
		sql = sql.Where("contract_id = ?", r.ContractID)
	}
	if r.OperationsStaff != "" {
		sql = sql.Where("operations_staff = ?", r.OperationsStaff)
	}
	if r.Member != "" {
		sql = sql.Where("member LIKE ?", "?"+r.Member+"")
	}
	if r.Status != 0 {
		sql = sql.Where("status = ?", r.Status)
	}

	if r.StartTime.Format("2006-01-02 15:04:05") != "0001-01-01 00:00:00" {
		sql = sql.Where("complete_time > ?", r.StartTime)
	}

	if r.EndTime.Format("2006-01-02 15:04:05") != "0001-01-01 00:00:00" {
		sql = sql.Where("complete_time < ?", r.EndTime)
	}

	sql.Find(&cts)

	head := []string{"下单日期", "运营名字", "门店", "下单项目", "门店链接", "安排数量",
		"初始值", "完成值", "特殊要求", "反馈"}

	var body [][]interface{}
	for _, ct := range cts {
		var pics []string
		var images string
		json.Unmarshal([]byte(ct.Images), &pics)

		for _, pic := range pics {
			if pic != "" && !strings.Contains(pic, "http") {
				pic = "http://8.136.135.212:8080" + pic
			}
			images = images + "   " + pic
		}
		ctInfo := []interface{}{
			ct.CreatedAt.Format("2006-01-02 15:04:05"), // 下单日期
			ct.OperationsStaff,                         // 运营名字
			ct.Member,                                  // 门店
			TaskType[uint16(ct.Type)],                  // 下单项目
			ct.StoreLink,                               // 门店链接
			ct.CompletedCount,                          // 安排数量
			ct.Initial,                                 // 初始值
			ct.CompletedCount,                          // 完成值
			ct.Requirements,                            // 要求
			images,                                     // 反馈
		}
		body = append(body, ctInfo)
	}

	filename := "任务列表" + time.Now().Format("20060102150405") + ".xlsx"
	handle.ExcelExport(c, head, body, filename)
}
