package model

// UserLog 用户操作日志
type UserLog struct {
	ID        uint   `gorm:"primarykey" json:"id" binding:"required"`
	Operator  uint   `json:"operator"`
	Action    string `json:"action" gorm:"type:varchar(30)"`
	Member    uint   `json:"member"`
	Contract  uint   `json:"contract"`
	Remarks   string `json:"remarks" gorm:"type:varchar(200)"`
	CreatedAt MyTime `json:"created_at"`
	UpdatedAt MyTime `json:"updated_at"`
}
