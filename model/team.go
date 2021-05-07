package model

type Team struct {
	ID        uint   `gorm:"primarykey" json:"id" binding:"required"`
	Name      string `json:"name" gorm:"type:varchar(10)"`   // 组名
	Users     string `json:"users" gorm:"type:varchar(512)"` // 用户
	Month     uint8  `json:"month"`                          // 月份
	Key       string `json:"key" gorm:"varchar(100)"`        // 标识
	CreatedAt MyTime `json:"created_at"`
	UpdatedAt MyTime `json:"updated_at"`
}
