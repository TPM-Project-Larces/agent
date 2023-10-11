package router

import (
	docs "github.com/TPM-Project-Larces/agent.git/docs"
	"github.com/TPM-Project-Larces/agent.git/handler"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine) {
	basePath := "/"
	docs.SwaggerInfo.BasePath = basePath
	v1 := router.Group(basePath)
	{
		v1.POST("generate_keys/", handler.GenerateKeys)
		v1.POST("upload_file/", handler.UploadFile)
		v1.POST("decrypt_file/", handler.DecryptFile)
	}

	// initialize swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
