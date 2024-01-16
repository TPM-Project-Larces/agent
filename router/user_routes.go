package router

import (
	"github.com/TPM-Project-Larces/agent.git/handler"
	"github.com/gin-gonic/gin"
)

func userRoutes(router *gin.Engine, basePath string, pathResource string) {

	user := router.Group(basePath + pathResource)
	{
		user.PUT("", handler.UpdateUser)
		user.POST("", handler.CreateUser)
		user.GET("/username", handler.GetUserByUsername)
		user.GET("", handler.GetUsers)
		user.DELETE("", handler.DeleteUser)
	}
}
