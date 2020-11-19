package model

// LabelGroup .
type LabelGroup struct {
	ID        uint    `gorm:"primarykey" json:"id" binding:"required"`
	Name      string  `json:"name" gorm:"type:varchar(20)"`
	Status    int8    `json:"status"`
	Color     string  `json:"color" gorm:"varchar(10)"`
	Labels    []Label `json:"labels"`
	CreatedAt MyTime  `json:"created_at"`
	UpdatedAt MyTime  `json:"updated_at"`
}
