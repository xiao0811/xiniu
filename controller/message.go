package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/handle"
)

// SendMessage 发送信息
func SendMessage(c *gin.Context) {
	var r handle.Info
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	id, err := r.Send()
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "短信发送失败", c)
		return
	}
	handle.ReturnSuccess("ok", id, c)
}
