package model

// Member .
type Member struct {
	ID                uint   `gorm:"primarykey" json:"id" binding:"required"`
	Name              string `gorm:"type:varchar(20)" json:"name"`
	City              string `json:"city" gorm:"type:varchar(10)"`
	FirstCategory     string `json:"first_category" gorm:"type:varchar(20)"`
	SecondaryCategory string `json:"secondary_category" gorm:"type:varchar(20)"`
	BusinessScope     string `json:"business_scope" gorm:"type:varchar(20)"`
	Stores            uint8  `json:"stores"`
	Accounts          uint8  `json:"accounts"`
	Bosses            uint8  `json:"bosses"`
	Brands            uint8  `json:"brands"`
	CreatedAt         MyTime `json:"created_at"`
	UpdatedAt         MyTime `json:"updated_at"`
}
