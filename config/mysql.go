package config

import (
	"log"

	"github.com/xiao0811/xiniu/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlConn mysql链接
var mysqlConn *gorm.DB

func init() {
	conf := Conf.MysqlConfig
	var err error
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db := mysqlConn
	// 自动迁移
	_ = db.AutoMigrate(
		&model.User{},
		&model.LabelGroup{},
		&model.Label{},
		&model.Member{},
		&model.UserLog{},
		&model.Contract{},
		&model.ContractTask{},
		&model.Refund{},
	)
}

// GetMysql 获取mysql链接
func GetMysql() *gorm.DB {
	return mysqlConn
}

// alter table contracts modify column current_star DECIMAL(3,2)
