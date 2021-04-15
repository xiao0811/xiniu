package model

// ContractTask 合约任务
type ContractTask struct {
	ID               uint   `gorm:"primarykey" json:"id" binding:"required"`
	Type             uint8  `json:"type"`                                     // 任务类型
	ContractID       uint   `json:"contract_id"`                              // 合约ID
	Member           string `json:"member" gorm:"type:varchar(20)"`           // 门店名称
	OperationsStaff  string `json:"operations_staff" gorm:"type:varchar(20)"` // 运营人员
	TaskCount        uint8  `json:"task_count"`                               // 总任务量
	CompletedCount   uint   `json:"completed_count"`                          // 完成任务量
	Initial          uint   `json:"initial"`                                  // 初始值
	CompleteTime     MyTime `json:"complete_time"`                            // 完成时间
	ActualCompletion MyTime `json:"actual_completion"`                        // 实际完成时间
	StoreLink        string `json:"store_link"`                               // 门店链接
	Mediator         string `json:"mediator" gorm:"type:varchar(20)"`         // 媒介人员
	Requirements     string `json:"requirements"`                             // 任务要求
	Images           string `json:"images"`                                   // 图片
	Status           uint8  `json:"status"`                                   // 状态
	Remark           string `json:"remark"`                                   // 备注
	CreatedAT        MyTime `json:"created_at"`
	UpdatedAt        MyTime `json:"updated_at"`
}
