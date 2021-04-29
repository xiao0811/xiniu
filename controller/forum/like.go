package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// Like 点赞
func Like(c *gin.Context) {
	var r struct {
		TitleID    uint   `json:"title_id" binding:"required"`
		OperatorID uint   `json:"operator_id" binding:"required"`
		Operator   string `json:"operator"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var l model.ForumLike
	var f model.ForumTitle
	if err := db.Where("title_id = ? AND operator_id = ? AND status = 1").First(&l).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "已赞改主题", c)
		return
	}

	ll := model.ForumLike{
		TitleID:    r.TitleID,
		OperatorID: r.OperatorID,
		Operator:   r.Operator,
		Status:     true,
	}

	if err := db.Where("id = ?", r.TitleID).First(&f).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "点赞主题不存在", c)
		return
	}

	if err := db.Create(&ll).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "点赞失败", c)
		return
	}
	f.Like++
	db.Save(&f)
	handle.ReturnSuccess("ok", ll, c)
}

// Unlike 取消点赞
func Unlike(c *gin.Context) {
	var r struct {
		TitleID    uint   `json:"title_id" binding:"required"`
		OperatorID string `json:"operator_id" binding:"required"`
		Operator   string `json:"operator"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var l model.ForumLike
	var f model.ForumTitle
	if err := db.Where("title_id = ? AND operator_id = ? AND status = 1", r.TitleID, r.OperatorID).First(&l).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "未点赞改主题", c)
		return
	}
	if err := db.Where("id = ?", r.TitleID).First(&f).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "点赞主题不存在", c)
		return
	}

	l.Status = false
	f.Like--
	db.Save(&l)
	db.Save(&f)
	handle.ReturnSuccess("ok", l, c)
}
