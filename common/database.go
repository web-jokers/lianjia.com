package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"os"
	"xietong.me/LianjiaSpider/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//fmt.Println(viper.GetString("datasource.driverName"))
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	if len(os.Args) < 2 {
		fmt.Println("请输入数据库密码")
		return nil
	}
	password := os.Args[1]
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("fail to connect databse,err:" + err.Error())
	}
	db.AutoMigrate(&model.Selling{}, &model.Sold{})
	DB = db
	return DB
}
func GetDB() *gorm.DB {
	return DB
}
