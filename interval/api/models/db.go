package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

type Config struct {
	User   string
	Pass   string
	Adrr   string
	Port   string
	Dbname string
}

func init() {
	var err error
	// 配置可用
	// f := &config{
	// 	user:   "root",                  // 用户名
	// 	pass:   "",            			 // 密码
	// 	adrr:   localhost, 				 // 地址
	// 	port:   "3306",                  // 端口
	// 	dbname: "user",                  // 数据库名称
	// }

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", "root", "Aa123456", "118.89.199.105", "3366", "device")

	Db, err = gorm.Open("mysql", dsn)
	// Db.AutoMigrate(&User{}, &Tag{}, &Device{})
	Db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	if err != nil {

		zap.L().Error("打开数据库失败,error:" + err.Error())
	} else {

		zap.L().Info("连接数据库成功")
	}
}
