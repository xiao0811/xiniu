package model

// Marshalling .
type Marshalling struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `json:"name" gorm:"type:varchar(10)"`
	Status    int8   `json:"status"`
	CreatedAt MyTime `json:"created_at"`
	UpdatedAt MyTime `json:"updated_at"`
}
