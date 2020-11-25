package config

import (
	"log"

	"github.com/xiao0811/xiniu/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	db := GetMysql()
	// 自动迁移
	db.AutoMigrate(
		&model.User{},
		&model.LabelGroup{},
		&model.Label{},
		&model.Member{},
		&model.UserLog{},
	)
}

// GetMysql 获取MySQL链接
func GetMysql() *gorm.DB {
	conf := Conf.MysqlConfig
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
