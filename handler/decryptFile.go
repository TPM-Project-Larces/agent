package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// @Summary decrypt file

// @Description decrypt an encrypted file
// @Tags User operations
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /decrypt_file [post]
func DecryptFile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "file decrypted")
}
