package model

// Label .
type Label struct {
	ID           uint   `gorm:"primarykey" json:"id" binding:"required"`
	Name         string `json:"name" gorm:"type:varchar(20)"`
	Status       int8   `json:"status"`
	LabelGroupID uint   `json:"label_group_id"`
	CreatedAt    MyTime `json:"created_at"`
	UpdatedAt    MyTime `json:"updated_at"`
}
