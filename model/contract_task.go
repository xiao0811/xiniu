package model

// ContractTask 合约任务
type ContractTask struct {
	ID              uint   `gorm:"primarykey" json:"id" binding:"required"`
	Type            uint8  `json:"type"`
	ContractID      uint   `json:"contract_id"`
	OperationsStaff string `json:"operations_staff" gorm:"type:varchar(20)"`
	TaskCount       uint8  `json:"task_count"`
	CompleteTime    MyTime `json:"complete_time"`
	Images          string `json:"images"`
	Status          uint8  `json:"status"`
	CreatedAt       MyTime `json:"created_at"`
	UpdatedAt       MyTime `json:"updated_at"`
}
