package controller

import (
	"net/http"
	"sort"

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

	handle.ReturnSuccess("ok", title, c)
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

type Rank struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Integral uint8  `json:"integral"`
}

type Ranks []Rank

// IntegralRank 获取积分排名
func IntegralRank(c *gin.Context) {
	ranks := GetIntegral()
	handle.ReturnSuccess("ok", ranks, c)
}

// getIntegral 积分排名
func GetIntegral() Ranks {
	var users []model.User
	var ranks Ranks
	db := config.GetMysql()
	db.Where("duty = 2").Find(&users)
	for _, user := range users {
		var integrals []model.IntegralLog
		db.Where("user_id = ?", user.ID).Find(&integrals)
		var total uint8
		for _, integral := range integrals {

			if integral.IsIncrease {
				total += integral.Quantity
			} else {
				total -= integral.Quantity
			}
		}

		ranks = append(ranks, Rank{
			ID:       user.ID,
			Name:     user.RealName,
			Integral: total,
		})
	}
	sort.Sort(ranks)
	return ranks
}

// Len()
func (s Ranks) Len() int {
	return len(s)
}

// Less():成绩将有低到高排序
func (s Ranks) Less(i, j int) bool {
	return s[i].Integral > s[j].Integral
}

// Swap()
func (s Ranks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
