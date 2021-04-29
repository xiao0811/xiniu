package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateComment 发表评论
func CreateComment(c *gin.Context) {
	var r struct {
		TitleID    uint   `json:"title_id"`    // 主题ID
		Content    string `json:"content"`     // 评论内容
		OperatorID uint   `json:"operator_id"` // 发表者ID
		Operator   string `json:"operator"`    // 发表者
		Reply      uint   `json:"reply"`       // 回复楼层
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()

	fc := model.ForumComment{
		TitleID:    r.TitleID,
		Content:    r.Content,
		OperatorID: r.OperatorID,
		Operator:   r.Operator,
		Reply:      r.Reply,
	}

	if err := db.Create(&fc).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "发表失败", c)
		return
	}
	handle.ReturnSuccess("ok", fc, c)
}

// UpdateComment 编辑评论
func UpdateComment(c *gin.Context) {
	var r model.ForumComment

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var fc model.ForumComment
	if err := db.Where("id = ?", r.ID).First(&fc).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "该条品论不存在", c)
		return
	}

	if err := db.Save(&fc).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "编辑失败", c)
		return
	}
	handle.ReturnSuccess("ok", fc, c)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	var r struct {
		ID uint `gorm:"primarykey" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	var fc model.ForumComment
	if err := db.Where("id = ?", r.ID).First(&fc).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该条评论不存在", c)
		return
	}

	if err := db.Delete(&fc).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "删除失败", c)
		return
	}

	handle.ReturnSuccess("ok", fc, c)
}
