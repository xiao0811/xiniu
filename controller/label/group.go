package label

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// GroupListRequest 标签分组列表
type GroupListRequest struct {
}

// GroupList 标签类型列表
func GroupList(c *gin.Context) {
	db := config.GetMysql()
	var groups []model.LabelGroup
	db.Where("status = 1").Find(&groups)
	handle.ReturnSuccess("ok", groups, c)
}

// CreateGroup 创建标签类型
func CreateGroup(c *gin.Context) {
	var r struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}

	db := config.GetMysql()
	group := model.LabelGroup{
		Name:   r.Name,
		Status: 1,
		Color:  r.Color,
	}
	if err := db.Create(&group).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "创建失败", c)
		return
	}
	handle.ReturnSuccess("ok", group, c)
}

// UpdateGroup 更新标签详情
func UpdateGroup(c *gin.Context) {
	var r model.LabelGroup

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var group model.LabelGroup

	if err := db.Where("id = ?", r.ID).First(&group).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该分组不存在", c)
		return
	}

	if err := db.Model(&group).Updates(r).Error; err != nil {
		handle.ReturnError(http.StatusServiceUnavailable, "标签信息更新失败", c)
		return
	}
	handle.ReturnSuccess("ok", group, c)
}

// DeleteGroup 删除标签类型
func DeleteGroup(c *gin.Context) {
	var r model.LabelGroup

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var group model.LabelGroup

	if err := db.Where("id = ?", r.ID).First(&group).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该标签不存在", c)
		return
	}
	group.Status = 0
	if err := db.Save(&group).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该标签不存在", c)
		return
	}
	handle.ReturnSuccess("ok", group, c)
}

// GroupDetails 标签详情
func GroupDetails(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var group model.LabelGroup
	db.Where("id = ?", r.ID).First(&group)
	handle.ReturnSuccess("ok", group, c)
}
