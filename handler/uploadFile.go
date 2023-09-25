package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// @Summary upload file

// @Description upload a file to encrypt
// @Tags User operations
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /upload_file [post]
func UploadFile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "file uploaded")
}
