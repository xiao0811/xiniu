package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseData api接口返回数据
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ReturnJSON 返回json
func ReturnJSON(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(code, ResponseData{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ReturnSuccess 返回正确处理
func ReturnSuccess(message string, data interface{}, c *gin.Context) {
	ReturnJSON(http.StatusOK, message, data, c)
}

// ReturnError 返回错误处理
func ReturnError(code int, message string, c *gin.Context) {
	ReturnJSON(code, message, nil, c)
}
