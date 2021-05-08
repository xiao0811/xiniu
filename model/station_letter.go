package model

type StationLetter struct {
	ID          uint     `gorm:"primarykey" json:"id"`
	ContractID  uint     `json:"contract_id"`                    // 关联合约ID
	Contract    Contract `json:"contract"`                       // 合约信息
	SenderID    uint     `json:"sender_id"`                      // 发送者ID
	Sender      User     `json:"sender"`                         // 发送者
	RecipientID uint     `json:"recipient_id"`                   // 接收者ID
	Recipient   User     `json:"recipient"`                      // 接收者
	Title       string   `json:"title" gorm:"type:varchar(100)"` // 消息标题
	Content     string   `json:"content" gorm:"type:text"`       // 消息内容
	Status      uint8    `json:"status"`                         // 状态
	Reply       uint     `json:"reply"`                          // 回复哪条的
	CreatedAt   MyTime   `json:"created_at"`
	UpdatedAt   MyTime   `json:"updated_at"`
}
