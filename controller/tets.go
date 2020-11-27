package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/model"
)

func Test(c *gin.Context) {
	db := config.GetMysql()
	var member model.Member
	db.First(&member)
	c.JSON(http.StatusOK, member.FirstCreate.String() == "0001-01-01 00:00:00 +0000 UTC")
}
