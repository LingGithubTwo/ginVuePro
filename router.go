package main

import (
	"github.com/LingGithubTwo/ginVuePro/controller"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/api/auth/register", controller.Register)
	return r
}
