package router

import (
	"github.com/TPM-Project-Larces/agent.git/handler"
	"github.com/gin-gonic/gin"
)

func fileRoutes(router *gin.Engine, basePath string, pathResource string) {

	file := router.Group(basePath + pathResource)
	{
		file.GET("/by_username", handler.GetFilesByUsername)
		file.GET("/by_name", handler.GetFileByName)
	}
}
