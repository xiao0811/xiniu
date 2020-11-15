package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/handle"
)

// VerifyToken 验证头部token
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			handle.ReturnError(http.StatusUnauthorized, "请求头中auth为空", c)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			handle.ReturnError(http.StatusUnauthorized, "请求头中auth格式有误", c)
			c.Abort()
			return
		}
		strToken := parts[1]
		claim, err := handle.VerifyAction(strToken)
		if err != nil {
			handle.ReturnError(http.StatusUnauthorized, err.Error(), c)
			c.Abort()
			return
		}
		c.Set("token", claim)
		c.Next()
	}
}
