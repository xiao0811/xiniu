package model

type TaskDetails struct {
	ID        uint   `gorm:"primarykey" json:"id" binding:"required"`
	TaskID    uint   `json:"task_id"`                        // 合约 ID
	Completed uint16 `json:"completed"`                      // 完成量
	Operator  string `json:"operator"`                       // 操作人员
	Image     string `json:"image"`                          // 完成图片
	Remark    string `json:"remark"`                         // 备注
	DoneTime  string `json:"done_time" gorm:"type:char(10)"` // 完成时间 不做处理
	CreatedAT MyTime `json:"created_at"`
	UpdatedAt MyTime `json:"updated_at"`
}
