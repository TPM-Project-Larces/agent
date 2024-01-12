package router

import (
	docs "github.com/TPM-Project-Larces/agent.git/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeRoutes() {
	basePath := "/"
	docs.SwaggerInfo.BasePath = basePath

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	encryptionRoutes(router, basePath, "encryption/")
	userRoutes(router, basePath, "user/")
	fileRoutes(router, basePath, "file/")

	router.Run(":3000")
}
