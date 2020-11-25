package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateContract 创建新的合约
func CreateContract(c *gin.Context) {
	var r struct {
		MemberID                 uint   `json:"member_id"`                   // 门店名称
		CooperationTime          string `json:"cooperation_time"`            // 合作时间
		ExpireTime               string `json:"expire_time"`                 // 到期时间
		IsStartService           bool   `json:"is_start_service"`            // 是否开始服务
		DelayTime                string `json:"delay_time"`                  // 延期后到期时间
		ContractAmount           uint   `json:"contract_amount"`             // 签约金额
		Arrives                  bool   `json:"arrives"`                     // 是否到账
		Arrears                  uint   `json:"arrears"`                     // 有无欠款
		CurrentStoreCollections  uint   `json:"current_store_collections"`   // 目前门店收藏量
		CurrentNumber            uint   `json:"current_number"`              // 目前评价数
		CurrentStar              uint   `json:"current_star"`                // 目前星级
		CurrentLeaderboard       string `json:"current_leaderboard"`         // 目前排行榜
		StoreCollections         uint   `json:"store_collections"`           // 门店收藏量
		InformationFlow          uint   `json:"information_flow"`            // 信息流
		BigVReview               uint   `json:"big_v_review"`                // 大V评论
		GroupBuyingVolume        uint   `json:"group_buying_volume"`         // 团购卖量
		Like                     uint   `json:"like"`                        // 点赞
		FollowPeers              string `json:"follow_peers"`                // 关注同行
		CurrentStatusOfPromotion int8   `json:"current_status_of_promotion"` // 推广通现状
		Upgrade                  bool   `json:"upgrade"`                     // 是否提升金牌店铺
		IncludeDetailsPage       bool   `json:"include_details_page"`        // 是否包含详情页
		Remarks                  string `json:"remarks"`                     // 备注
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "门店创建失败", c)
		return
	}
	db := config.GetMysql()
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	_CooperationTime, _ := time.ParseInLocation(model.TimeFormat, r.CooperationTime, time.Local)
	_ExpireTime, _ := time.ParseInLocation(model.TimeFormat, r.ExpireTime, time.Local)
	_DelayTime, _ := time.ParseInLocation(model.TimeFormat, r.DelayTime, time.Local)
	con := model.Contract{
		UUID:                     "XINIU-ORD-" + time.Now().Format("200601021504") + strconv.Itoa(handle.RandInt(1000, 9999)),
		MemberID:                 r.MemberID,
		CooperationTime:          model.MyTime{Time: _CooperationTime},
		ExpireTime:               model.MyTime{Time: _ExpireTime},
		IsStartService:           r.IsStartService,
		DelayTime:                model.MyTime{Time: _DelayTime},
		ContractAmount:           r.ContractAmount,
		Arrives:                  r.Arrives,
		Arrears:                  r.Arrears,
		CurrentStoreCollections:  r.CurrentStoreCollections,
		CurrentNumber:            r.CurrentNumber,
		CurrentStar:              r.CurrentStar,
		CurrentLeaderboard:       r.CurrentLeaderboard,
		StoreCollections:         r.StoreCollections,
		InformationFlow:          r.InformationFlow,
		BigVReview:               r.BigVReview,
		GroupBuyingVolume:        r.GroupBuyingVolume,
		Like:                     r.Like,
		FollowPeers:              r.FollowPeers,
		CurrentStatusOfPromotion: r.CurrentStatusOfPromotion,
		Upgrade:                  r.Upgrade,
		IncludeDetailsPage:       r.IncludeDetailsPage,
		Remarks:                  r.Remarks,
		Status:                   2,
	}
	db.Create(&con)
	// 创建用户记录
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Create Contact",
		Contract: con.ID,
		Remarks:  token.FullName + "创建合约: " + con.UUID,
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", con, c)
}
