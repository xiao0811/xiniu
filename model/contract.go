package model

// Contract 客户合约结构
type Contract struct {
	ID                       uint   `gorm:"primarykey" json:"id" binding:"required"`
	UUID                     string `json:"uuid" gorm:"type:char(26);unique"` // 门店编号
	MemberID                 uint   `json:"member_id"`                        // 门店名称
	CooperationTime          MyTime `json:"cooperation_time"`                 // 合作时间
	ExpireTime               MyTime `json:"expire_time"`                      // 到期时间
	IsStartService           bool   `json:"is_start_service"`                 // 是否开始服务
	DelayTime                MyTime `json:"delay_time"`                       // 延期后到期时间
	ContractAmount           uint   `json:"contract_amount"`                  // 签约金额
	Arrives                  bool   `json:"arrives"`                          // 是否到账
	Arrears                  uint   `json:"arrears"`                          // 有无欠款
	CurrentStoreCollections  uint   `json:"current_store_collections"`        // 目前门店收藏量
	CurrentNumber            uint   `json:"current_number"`                   // 目前评价数
	CurrentStar              uint   `json:"current_star"`                     // 目前星级
	CurrentLeaderboard       string `json:"current_leaderboard"`              // 目前排行榜
	StoreCollections         uint   `json:"store_collections"`                // 门店收藏量
	InformationFlow          uint   `json:"information_flow"`                 // 信息流
	BigVReview               uint   `json:"big_v_review"`                     // 大V评论
	GroupBuyingVolume        uint   `json:"group_buying_volume"`              // 团购卖量
	Like                     uint   `json:"like"`                             // 点赞
	FollowPeers              string `json:"follow_peers"`                     // 关注同行
	CurrentStatusOfPromotion int8   `json:"current_status_of_promotion"`      // 推广通现状
	Upgrade                  bool   `json:"upgrade"`                          // 是否提升金牌店铺
	IncludeDetailsPage       bool   `json:"include_details_page"`             // 是否包含详情页
	Status                   int8   `json:"status"`
	Type                     int8   `json:"type"`
	Remarks                  string `json:"remarks"` // 备注
	CreatedAt                MyTime `json:"created_at"`
	UpdatedAt                MyTime `json:"updated_at"`
}