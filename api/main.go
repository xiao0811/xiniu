package main

import (
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

func main() {
	password, _ := handle.HashPassword("xiaosha")
	var xiao = model.User{
		Phone:    "18949883585",
		Password: password,
	}
	db := config.GetMysql()
	db.Create(&xiao)
}
