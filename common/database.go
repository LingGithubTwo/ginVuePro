package common

import (
	"fmt"
	"github.com/LingGithubTwo/ginVuePro/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginvuepro"
	username := "ginVuePro"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, host, port, database, charset)

	//db, err := gorm.Open(mysql.Open(args),&gorm.Config{})

	//自定义驱动
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: driverName,
		DSN:        args, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database,err: " + err.Error())
	}
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
