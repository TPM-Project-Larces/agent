package router

import (
	"github.com/TPM-Project-Larces/agent.git/handler"
	"github.com/gin-gonic/gin"
)

func encryptionRoutes(router *gin.Engine, basePath string, pathResource string) {

	encryption := router.Group(basePath + pathResource)
	{
		encryption.POST("/generate_keys", handler.GenerateKeys)
		encryption.POST("/upload_file", handler.UploadFile)
		encryption.POST("/upload_encrypted_file", handler.UploadEncryptedFile)
		encryption.POST("/decrypt_file", handler.DecryptFile)
	}
}
