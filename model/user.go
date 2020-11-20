package model

// User 用户模型
type User struct {
	ID             uint        `gorm:"primarykey" json:"id" binding:"required"`
	Phone          string      `gorm:"unique;type:char(11)" json:"phone"`
	Password       string      `gorm:"type:varchar(100)" json:"-"`
	RealName       string      `gorm:"type:varchar(10)" json:"real_name"`
	Gender         uint8       `json:"gender"`
	Birthday       MyTime      `gorm:"default:null" json:"birthday"`
	Identification string      `gorm:"type:char(18)" json:"identification"`
	Role           uint8       `json:"role"`
	MarshallingID  uint        `json:"marshalling_id"`
	Marshalling    Marshalling `json:"marshalling,omitempty"`
	Status         int8        `json:"status"`
	Duty           int8        `json:"duty"`
	CreatedAt      MyTime      `json:"created_at"`
	UpdatedAt      MyTime      `json:"updated_at"`
}
