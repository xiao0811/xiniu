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
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	sl := model.StationLetter{
		SenderID:    r.SenderID,
		RecipientID: r.RecipientID,
		Title:       r.Title,
		Content:     r.Content,
		Status:      1,
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
