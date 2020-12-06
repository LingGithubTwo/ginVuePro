package main

import (
	"github.com/LingGithubTwo/ginVuePro/controller"
	"github.com/LingGithubTwo/ginVuePro/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/getInfo", middleware.AuthMiddleware(), controller.GetInfo)
	return r
}
