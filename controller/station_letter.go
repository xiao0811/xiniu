package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateStationLetter 创建一条信息
func CreateStationLetter(c *gin.Context) {
	var r struct {
		ContractID  uint   `json:"contract_id"`                     // 关联合约ID
		SenderID    uint   `json:"sender_id" binding:"required"`    // 发送者ID
		RecipientID uint   `json:"recipient_id" binding:"required"` // 接收者ID
		Title       string `json:"title"`                           // 消息标题
		Content     string `json:"content"`                         // 消息内容
		Reply       uint   `json:"reply"`                           // 回复哪条的
		Type        string `json:"type" binding:"required"`
		MemberName  string `json:"member_name"` // 店铺名称
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	sl := model.StationLetter{
		ContractID:  r.ContractID,
		SenderID:    r.SenderID,
		RecipientID: r.RecipientID,
		Title:       r.Title,
		Content:     r.Content,
		MemberName:  r.MemberName,
		Status:      1,
		Reply:       r.Reply,
	}
	if r.Type == "create" {
		sl.ContractID = r.ContractID
	} else if r.Type == "reply" {
		sl.Reply = r.Reply
	}
	db := config.GetMysql()
	if err := db.Create(&sl).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "消息发送失败", c)
		return
	}

	handle.ReturnSuccess("ok", sl, c)
}

// UpdateStationLetter 更新
func UpdateStationLetter(c *gin.Context) {
	var letter model.StationLetter
	if err := c.ShouldBind(&letter); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	if err := db.Where("id = ?", letter.ID).First(&letter).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "该条站内信不存在", c)
		return
	}

	if err := db.Save(&letter).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "更新失败", c)
		return
	}
	handle.ReturnSuccess("ok", letter, c)
}

// DeleteStationLetter 删除
func DeleteStationLetter(c *gin.Context) {
	var letter model.StationLetter
	if err := c.ShouldBind(&letter); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	if err := db.Where("id = ?", letter.ID).First(&letter).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "该条站内信不存在", c)
		return
	}

	if err := db.Delete(&letter).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "删除失败失败", c)
		return
	}
	handle.ReturnSuccess("ok", letter, c)
}

// GetStationLetter 获取站内信
func GetStationLetter(c *gin.Context) {
	var r struct {
		ContractID  uint   `json:"contract_id"`  // 关联合约ID
		SenderID    uint   `json:"sender_id"`    // 发送者ID
		RecipientID uint   `json:"recipient_id"` // 接收者ID
		MemberName  string `json:"member_name"`  // 店铺名称
		Title       string `json:"title"`        // 消息标题
		StartTime   string `json:"start_time"`   // 开始时间
		EndTime     string `json:"end_time"`     // 结束时间
		Reply       uint   `json:"reply"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	var letters []model.StationLetter
	db := config.GetMysql().Preload("Sender").Preload("Recipient")

	if r.ContractID != 0 {
		db = db.Where("contract_id = ?", r.ContractID)
	}

	if r.SenderID != 0 {
		db = db.Where("sender_id = ?", r.SenderID)
	}

	if r.RecipientID != 0 {
		db = db.Where("recipient_id = ?", r.RecipientID)
	}

	if r.Title != "" {
		db = db.Where("title LIKE ?", "%"+r.Title+"%")
	}

	if r.StartTime != "" {
		db = db.Where("created_at > ?", r.StartTime)
	}

	if r.EndTime != "" {
		db = db.Where("created_at < ?", r.EndTime)
	}

	if r.MemberName != "" {
		db = db.Where("member_name LIKE ?", "%"+r.MemberName+"%")
	}

	if r.Reply != 0 {
		db = db.Where("reply = ?", r.Reply)
	}

	db.Find(&letters)

	handle.ReturnSuccess("ok", letters, c)
}
