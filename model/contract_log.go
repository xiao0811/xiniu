package model

type ContratLog struct {
	ID            uint   `gorm:"primarykey" json:"id" binding:"required"`
	ContratID     uint   `json:"contrat_id"`     // 合约ID
	OperatorID    uint   `json:"operator_id"`    // 操作人员ID
	Operator      string `json:"operator"`       // 操作人员名字
	Type          uint8  `json:"type"`           // 类型: 1 牌级 2 推广通
	OperatingTime MyTime `json:"operating_time"` // 操作时间
	GradeScore    uint8  `json:"grade_score"`    // 等级分数 - 牌级
	Spend         int    `json:"spend"`          // 花费 - 推广通
	GuestCapital  int    `json:"guest_capital"`  // 客资 - 推广通
	Week          uint8  `json:"week"`           // 第几周
	CreatedAt     MyTime `json:"created_at"`
	UpdatedAt     MyTime `json:"updated_at"`
}
