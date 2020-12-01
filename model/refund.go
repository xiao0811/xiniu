package model

// Refund 退款结构
type Refund struct {
	ID              uint    `gorm:"primarykey" json:"id" binding:"required"`
	ContractID      uint    `json:"contract_id"`
	Amount          float64 `json:"amount" gorm:"type:DECIMAL(8,2)"`
	Status          int8    `json:"status"`           // 状态 0 待审核, 1 通过, 2 拒绝
	Applicant       uint    `json:"applicant"`        // 申请人
	Reviewer        uint    `json:"reviewer"`         // 审核人
	Remark          string  `json:"remark"`           // 备注
	Reason          string  `json:"reason"`           // 审核备注
	OperationsStaff string  `json:"operations_staff"` // 运营人员
	BusinessPeople  string  `json:"business_people"`  // 业务人员
	CreatedAt       MyTime  `json:"created_at"`
	UpdatedAt       MyTime  `json:"updated_at"`
}
