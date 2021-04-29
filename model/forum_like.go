package model

type ForumLike struct {
	ID         uint   `gorm:"primarykey" json:"id" binding:"required"`
	TitleID    uint   `json:"title_id"`    // 主题ID
	OperatorID uint   `json:"operator_id"` // 发表者ID
	Operator   string `json:"operator"`    // 发表者
	Status     bool   `json:"status"`      // 状态
	CreatedAt  MyTime `json:"created_at"`
	UpdatedAt  MyTime `json:"updated_at"`
}
