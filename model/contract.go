package model

// Contract 客户合约结构
type Contract struct {
	ID                       uint           `gorm:"primarykey" json:"id" binding:"required"`
	UUID                     string         `json:"uuid" gorm:"type:char(26);unique"` // 门店编号
	MemberID                 uint           `json:"member_id"`                        // 门店名称
	Member                   Member         `json:"member"`
	CooperationTime          MyTime         `json:"cooperation_time"`                      // 合作时间
	ExpireTime               MyTime         `json:"expire_time"`                           // 到期时间
	IsStartService           bool           `json:"is_start_service"`                      // 是否开始服务
	DelayTime                MyTime         `json:"delay_time"`                            // 延期后到期时间
	ContractAmount           uint           `json:"contract_amount"`                       // 签约金额
	Arrives                  bool           `json:"arrives"`                               // 是否到账
	Arrears                  uint           `json:"arrears"`                               // 有无欠款
	CurrentStoreCollections  uint           `json:"current_store_collections"`             // 目前门店收藏量
	CurrentNumber            uint           `json:"current_number"`                        // 目前评价数
	CurrentStar              float32        `json:"current_star" gorm:"type:DECIMAL(3,2)"` // 目前星级
	CurrentLeaderboard       string         `json:"current_leaderboard"`                   // 目前排行榜
	StoreCollections         uint           `json:"store_collections"`                     // 门店收藏量
	InformationFlow          uint           `json:"information_flow"`                      // 信息流
	BigVReview               uint           `json:"big_v_review"`                          // 大V评论
	GroupBuyingVolume        uint           `json:"group_buying_volume"`                   // 团购卖量
	Like                     uint           `json:"like"`                                  // 点赞
	FollowPeers              string         `json:"follow_peers"`                          // 关注同行
	CurrentStatusOfPromotion int8           `json:"current_status_of_promotion"`           // 推广通现状
	Upgrade                  bool           `json:"upgrade"`                               // 是否提升金牌店铺
	IncludeDetailsPage       bool           `json:"include_details_page"`                  // 是否包含详情页
	Status                   int8           `json:"status"`
	Type                     int16          `json:"type"`
	Task                     string         `json:"task"`
	ContractTask             []ContractTask `json:"contract_task"`
	Remarks                  string         `json:"remarks"`          // 备注
	OperationsStaff          string         `json:"operations_staff"` // 运营人员
	BusinessPeople           string         `json:"business_people"`  // 业务人员
	Refund                   MyTime         `json:"refund"`
	Reason                   string         `json:"reason"` // 审核备注
	Sort                     int64          `json:"sort"`
	ContractData             string         `json:"contract_data"` // 附加备注
	BuildPage                string         `json:"build_page"`
	IsBuild                  bool           `json:"is_build"`
	InitialLevel             uint8          `json:"initial_level"`
	ContractLogs             []ContractLog  `json:"contract_logs"`
	CreatedAt                MyTime         `json:"created_at"`
	UpdatedAt                MyTime         `json:"updated_at"`
}
