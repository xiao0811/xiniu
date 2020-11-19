package label

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// Create 新增标签
func Create(c *gin.Context) {
	var r struct {
		Name         string `json:"name" binding:"required"`
		LabelGroupID uint   `json:"label_group_id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	label := model.Label{Name: r.Name, Status: 1, LabelGroupID: r.LabelGroupID}
	if err := db.Create(&label).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "标签创建失败", c)
		return
	}
	handle.ReturnSuccess("ok", label, c)
}

// Update 标签更新
func Update(c *gin.Context) {
	var r model.Label

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var label model.Label

	if err := db.Where("id = ?", r.ID).First(&label).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该标签不存在", c)
		return
	}

	if err := db.Model(&label).Updates(r).Error; err != nil {
		handle.ReturnError(http.StatusServiceUnavailable, "标签信息更新失败", c)
		return
	}
	handle.ReturnSuccess("ok", label, c)
}

// Delete 标签删除 status = 0
func Delete(c *gin.Context) {
	var r model.Label

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var label model.Label

	if err := db.Where("id = ?", r.ID).First(&label).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该标签不存在", c)
		return
	}
	label.Status = 0
	if err := db.Save(&label).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该标签不存在", c)
		return
	}
	handle.ReturnSuccess("ok", label, c)
}

// IndexRequest 标签查询结构
type IndexRequest struct {
	Group  uint `json:"group"`
	Status int8 `json:"status"`
}

// Index 标签列表
func Index(c *gin.Context) {
	var r IndexRequest
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var labels []model.Label
	sql := db.Where("status = ?", r.Status)
	if r.Group != 0 {
		sql = sql.Where("group = ?", r.Group)
	}
	sql.Find(&labels)
	handle.ReturnSuccess("ok", labels, c)
}

// Info 获取一个标签详情
func Info(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var label model.Label
	db.Where("id = ?", r.ID).First(&label)
	handle.ReturnSuccess("ok", label, c)
}
