package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Phone          string    `gorm:"unique;type:char(11)" json:"phone"`
	Password       string    `json:"-"`
	RealName       string    `gorm:"type:varchar(10)" json:"real_name"`
	Gender         uint8     `json:"gender"`
	Birthday       time.Time `gorm:"default:null" json:"birthday"`
	Identification string    `gorm:"type:char(18)" json:"identification"`
	Role           uint8     `json:"role"`
	Marshalling    uint8     `json:"marshalling"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
