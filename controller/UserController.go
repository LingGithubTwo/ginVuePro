package controller

import (
	"fmt"
	"github.com/LingGithubTwo/ginVuePro/common"
	"github.com/LingGithubTwo/ginVuePro/model"
	"github.com/LingGithubTwo/ginVuePro/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
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

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "加密出错"})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	ctx.JSON(200, gin.H{
		"code":    http.StatusOK,
		"message": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	//获取数据库实例
	DB := common.GetDB()

	//获取参数
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

	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户密码错误"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统异常",
		})
		log.Printf("token generate error : %v", err)
		return
	}

	//返回结果
	ctx.JSON(200, gin.H{
		"code":    http.StatusOK,
		"token":   token,
		"message": "登录成功",
	})
}

func GetInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"user": user,
		},
		"message": "获取成功",
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
