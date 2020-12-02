package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/handle"
)

func SendMessage(c *gin.Context) {
	var r handle.Info
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	if id, err := r.Send(); err != nil {
		handle.ReturnError(http.StatusBadRequest, "短信发送失败", c)
		return
	} else {
		handle.ReturnSuccess("ok", id, c)
	}
}
