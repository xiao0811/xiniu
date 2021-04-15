package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
	"gorm.io/gorm"
)

// CreateContractTaskDetails 创建任务详情
func CreateContractTaskDetails(c *gin.Context) {
	var r struct {
		TaskID    uint   `json:"task_id" binding:"required"` // 合约 ID
		Completed uint16 `json:"completed"`                  // 完成量
		Operator  string `json:"operator"`                   // 操作人员
		Image     string `json:"image"`                      // 完成图片
		Remark    string `json:"remark"`                     // 备注
		DoneTime  string `json:"done_time"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	var task model.ContractTask
	var operator string
	var contract model.Contract
	db.Where("id = ?", r.TaskID).First(&task)
	db.Where("id = ?", task.ContractID).First(&contract)
	if r.Operator == "" {
		_token, _ := c.Get("token")
		token, _ := _token.(*handle.JWTClaims)
		operator = token.FullName
	} else {
		operator = r.Operator
	}

	d := model.TaskDetails{
		TaskID:    r.TaskID,
		Completed: r.Completed,
		Operator:  operator,
		Image:     r.Image,
		Remark:    r.Remark,
		DoneTime:  r.DoneTime,
	}

	if contract.Type == 3 || contract.Type == 33 || contract.Type == 333 {
		task.Status = 10
	} else {
		task.CompletedCount += uint(r.Completed)

		if task.CompletedCount >= uint(task.TaskCount) {
			task.CompleteTime = model.MyTime{Time: time.Now()}
			task.Status = 10
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&d)

		tx.Save(&task)

		return nil
	})

	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约任务详情创建失败", c)
		return
	}

	handle.ReturnSuccess("ok", d, c)
}

// GetContractTasKDetails 获取合约任务详情
func GetContractTasKDetails(c *gin.Context) {
	var ct model.ContractTask
	var td []model.TaskDetails
	if err := c.ShouldBind(&ct); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	db.Where("contract_id = ?", ct.ContractID).Find(&td)
	handle.ReturnSuccess("ok", td, c)
}
