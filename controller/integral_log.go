package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// IntegralChange 用户积分变动
func IntegralChange(c *gin.Context) {
	var r struct {
		UserID     uint   `json:"user_id" binding:"required"`  // 用户ID
		IsIncrease bool   `json:"is_increase"`                 // 积分增减
		Quantity   uint8  `json:"quantity" binding:"required"` // 增减个数
		Remark     string `json:"remark"`                      // 备注
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	var is_increase bool
	if r.IsIncrease {
		is_increase = true
	}

	if err := integralChange(r.UserID, r.Quantity, r.Remark, is_increase); err != nil {
		handle.ReturnError(http.StatusInternalServerError, "积分变更失败", c)
		return
	}

	handle.ReturnSuccess("ok", r, c)
}

// CommentAdoption 评论被采纳
func CommentAdoption(c *gin.Context) {
	var r struct {
		CommentID uint  `json:"comment_id" binding:"required"`
		Quantity  uint8 `json:"quantity"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var comment model.ForumComment
	var title model.ForumTitle
	if err := db.Where("id = ?", r.CommentID).First(&comment).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "该条评论不存在", c)
		return
	}

	var quantity uint8 = 2
	if r.Quantity != 0 {
		quantity = r.Quantity
	}
	db.Where("id = ?", comment.TitleID).First(&title)
	title.GivenIntegral += quantity
	comment.Integral = quantity
	if title.GivenIntegral > title.TotalIntegral {
		handle.ReturnError(http.StatusBadRequest, "所给积分超出上限", c)
		return
	}
	comment.Adoption = true
	db.Save(&comment)
	db.Save(&title)
	integralChange(comment.OperatorID, quantity, "评论被采纳", true)
}

func integralChange(user_id uint, quantity uint8, remark string, is_increase bool) error {
	ic := model.IntegralLog{
		UserID:     user_id,
		IsIncrease: is_increase,
		Quantity:   quantity,
		Remark:     remark,
	}

	db := config.GetMysql()
	return db.Create(&ic).Error
}
