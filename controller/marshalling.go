package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

const (
	// YUNYING type:运营
	YUNYING = 1
	// YEWU type:业务
	YEWU = 2
)

// CreateMarshallingRequest .
type CreateMarshallingRequest struct {
	Name string `json:"name" gorm:"type:varchar(10)" binding:"required"`
	Type int8   `json:"type"`
}

// CreateMarshalling 创建一个新的组别
func CreateMarshalling(c *gin.Context) {
	var r CreateMarshallingRequest
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据不正确", c)
		return
	}

	m := model.Marshalling{
		Name:   r.Name,
		Type:   r.Type,
		Status: 1,
	}
	db := config.MysqlConn
	if err := db.Create(&m).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "输入数据格式不正确", c)
		return
	}

	handle.ReturnSuccess("ok", m, c)
}

// UpdateMarshalling 更新部门信息
func UpdateMarshalling(c *gin.Context) {
	var r model.Marshalling
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据不正确", c)
		return
	}
	var m model.Marshalling
	db := config.MysqlConn
	if err := db.Where("id = ?", r.ID).First(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据不正确", c)
		return
	}

	if err := db.Model(&m).Updates(r).Error; err != nil {
		handle.ReturnError(http.StatusServiceUnavailable, "用户信息更新失败", c)
		return
	}

	handle.ReturnSuccess("ok", m, c)
}

// DeleteMarshalling 删除部门
func DeleteMarshalling(c *gin.Context) {
	var r model.Marshalling
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据不正确", c)
		return
	}
	var m model.Marshalling
	db := config.MysqlConn
	if err := db.Where("id = ?", r.ID).First(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据不正确", c)
		return
	}

	m.Status = 0
	if err := db.Save(&m).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "输入数据不正确", c)
		return
	}

	handle.ReturnSuccess("ok", m, c)
}

// MarshallingList 分组列表
func MarshallingList(c *gin.Context) {
	var r struct {
		Type uint8 `json:"type"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据不正确", c)
		return
	}
	var marshallings []model.Marshalling
	db := config.MysqlConn

	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var user model.User
	db.Where("id = ?", token.UserID).First(&user)

	sql := db.Where("status = 1")
	if user.Role == 1 {
		if user.Duty == 1 {

		} else {
			sql = sql.Where("type = ?", r.Type)
		}
	} else if user.Role == 2 {
		sql = sql.Where("id = ?", user.MarshallingID)
	}

	sql.Find(&marshallings)
	handle.ReturnSuccess("ok", marshallings, c)
}
