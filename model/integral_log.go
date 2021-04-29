package model

type IntegralLog struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	UserID     uint   `json:"user_id"`                         // 用户ID
	IsIncrease bool   `json:"is_increase"`                     // 积分增减
	Quantity   uint8  `json:"quantity"`                        // 增减个数
	Remark     string `json:"remark" gorm:"type:varchar(100)"` // 备注
	CreatedAt  MyTime `json:"created_at"`
	UpdatedAt  MyTime `json:"updated_at"`
}
