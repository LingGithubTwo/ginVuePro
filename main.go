package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()

	//defer db.Close()//用完之后关闭数据库连接  可能版本问题，close没有用

	r := gin.Default()
	r.GET("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		//数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}

		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}

		//判断名称是否为空，如果为空给一个默认的10位随机数
		if len(name) == 0 {
			name = RandomString(10)
		}

		fmt.Printf("名称:%v--手机号:%v--密码:%v\n", name, telephone, password)
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
			return
		}

		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	panic(r.Run())
	//r.Run()
}

func RandomString(n int) string {
	var letters = []byte("123456789qwertyuiopjljljglfdsgdfsgfdsgdfg")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

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
	db.AutoMigrate(&User{})
	return db
}
