package controller

import (
	"fmt"
	"github.com/LingGithubTwo/ginVuePro/common"
	"github.com/LingGithubTwo/ginVuePro/model"
	"github.com/LingGithubTwo/ginVuePro/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Register(ctx *gin.Context) {

	DB := common.GetDB()

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
		name = util.RandomString(10)
	}

	fmt.Printf("名称:%v--手机号:%v--密码:%v\n", name, telephone, password)
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)

	ctx.JSON(200, gin.H{
		"message": "注册成功",
	})
}

func isTelephoneExist(DB *gorm.DB, telephone string) bool {
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
