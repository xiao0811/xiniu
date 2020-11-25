package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// GetUserLogs 用户记录
func GetUserLogs(c *gin.Context) {
	var r struct {
		UserID   int    `json:"user_id"`
		Action   string `json:"action"`
		Member   int    `json:"member"`
		Contract int    `json:"contract"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var logs model.UserLog
	sql := db.Where("1 = 1")
	if r.UserID != 0 {
		sql = sql.Where("operator = ?", r.UserID)
	}
	if r.Action != "" {
		sql = sql.Where("action = ?", r.Action)
	}
	if r.Member != 0 {
		sql = sql.Where("member = ?", r.Action)
	}
	if r.Contract != 0 {
		sql = sql.Where("contract = ?", r.Action)
	}

	if err := db.Order("id desc").Find(&logs).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "暂无记录", c)
		return
	}

	handle.ReturnSuccess("ok", logs, c)
}
