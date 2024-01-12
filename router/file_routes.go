package router

import (
	"github.com/TPM-Project-Larces/agent.git/handler"
	"github.com/gin-gonic/gin"
)

func fileRoutes(router *gin.Engine, basePath string, pathResource string) {

	encryption := router.Group(basePath + pathResource)
	{
		encryption.POST("/delete_file", handler.DeleteFile)
	}
}
