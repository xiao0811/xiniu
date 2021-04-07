package model

type TaskDetails struct {
	ID        uint   `gorm:"primarykey" json:"id" binding:"required"`
	CreatedAT MyTime `json:"created_at"`
	UpdatedAt MyTime `json:"updated_at"`
}
