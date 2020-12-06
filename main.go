package main

import (
	"github.com/LingGithubTwo/ginVuePro/common"
	"github.com/gin-gonic/gin"
)

func main() {
	//db := common.InitDB()
	common.InitDB()
	//defer db.Close()//用完之后关闭数据库连接  可能版本问题，close没有用

	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run())
	//r.Run()
}
