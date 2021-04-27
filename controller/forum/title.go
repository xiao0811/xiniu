package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateTitle 发布新内容
func CreateTitle(c *gin.Context) {
	var r struct {
		Title            string `json:"title"`                    // 标题
		Content          string `json:"content" gorm:"type:text"` // 内容
		Type             uint8  `json:"type"`                     // 分类
		Recommended      bool   `json:"recommended"`              // 是否推荐
		TotalIntegral    uint8  `json:"total_integral"`           // 总积分
		OriginalPosterID uint   `json:"original_poster_id"`       // 楼主ID
		OriginalPoster   string `json:"original_poster"`          // 楼主
		Contract         uint   `json:"contract"`                 // 关联合约
		Images           string `json:"images"`                   // 图片
		Label            uint   `json:"label"`                    // 标签
		IsCarousel       bool   `json:"is_carousel"`              // 是否轮播
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()

	ft := model.ForumTitle{
		Title:          r.Title,
		Content:        r.Content,
		Type:           r.Type,
		Recommended:    r.Recommended,
		TotalIntegral:  r.TotalIntegral,
		OriginalPoster: r.OriginalPoster,
		Contract:       r.Contract,
		Images:         r.Images,
		Label:          r.Label,
		IsCarousel:     r.IsCarousel,
	}

	if err := db.Create(&ft).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "发布失败", c)
		return
	}
	handle.ReturnSuccess("ok", ft, c)
}

// UpdateTitle 更新主题
func UpdateTitle(c *gin.Context) {
	var r model.ForumTitle
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()

	var ft model.ForumTitle
	if err := db.Where("id = ?", r.ID).First(&ft).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "主题ID不存在", c)
		return
	}

	if err := db.Save(&r).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "更新失败", c)
		return
	}
	handle.ReturnSuccess("ok", r, c)
}

// DeleteTitle 删除主题
func DeleteTitle(c *gin.Context) {
	var r struct {
		ID uint `gorm:"primarykey" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var ft model.ForumTitle

	if err := db.Where("id = ?", r.ID).First(&ft).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "主题ID不存在", c)
		return
	}

	if err := db.Delete(&ft).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "删除失败", c)
		return
	}

	handle.ReturnSuccess("ok", ft, c)
}

// GetTitleList 获取主题列表
func GetTitleList(c *gin.Context) {
	var r struct {
		Type  uint8 `json:"type"`
		Label uint  `json:"label"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var fts []model.ForumTitle
	sql := db
	if r.Type != 0 {
		sql = sql.Where("type = ?", r.Type)
	}

	if r.Label != 0 {
		sql = sql.Where("label = ?", r.Label)
	}

	sql.Order("id").Find(&fts)

	handle.ReturnSuccess("ok", fts, c)
}

// TitleDetails 获取主题详情
func TitleDetails(c *gin.Context) {
	var r struct {
		ID uint `gorm:"primarykey" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.GetMysql()
	var ft model.ForumTitle

	if err := db.Where("id = ?", r.ID).Preload("Comment").First(&ft).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "主题ID不存在", c)
		return
	}

	ft.Pageviews++
	db.Save(&ft)

	handle.ReturnSuccess("ok", ft, c)
}
