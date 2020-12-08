package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// SendMessage 发送信息
func SendMessage(c *gin.Context) {
	var r struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	var contract model.Contract
	var member model.Member
	db := config.GetMysql()
	if err := db.Where("id = ?", r.ID).First(&contract).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "合约ID不存在", c)
		return
	}
	db.Where("id = ?", contract.MemberID).First(&member)
	fmt.Println(member.Name, member.ExpireTime)
	m := handle.Info{
		MsgText:    member.Name + "即将到期:" + member.ExpireTime.String(),
		Destmobile: member.Phone,
	}
	id, err := m.Send()
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "短信发送失败", c)
		return
	}
	handle.ReturnSuccess("ok", id, c)
}
