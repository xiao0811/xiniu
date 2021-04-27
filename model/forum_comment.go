package model

type ForumComment struct {
	ID         uint   `gorm:"primarykey" json:"id" binding:"required"`
	TitleID    uint   `json:"title_id"`    // 主题ID
	Content    string `json:"content"`     // 评论内容
	OperatorID string `json:"operator_id"` // 发表者ID
	Operator   string `json:"operator"`    // 发表者
	Integral   uint8  `json:"integral"`    // 获得积分
	Reply      uint   `json:"reply"`       // 回复楼层
	Adoption   bool   `json:"adoption"`    // 是否采纳
	CreatedAt  MyTime `json:"created_at"`
	UpdatedAt  MyTime `json:"updated_at"`
}
