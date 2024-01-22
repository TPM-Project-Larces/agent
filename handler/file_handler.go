package handler

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /
// @Summary Get encrypted files by username
// @Description Get a list of encrypted files by username
// @Tags File
// @Accept json
// @Produce json
// @Success 200 {object} schemas.ListFilesResponse
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_error"
// @Router /file/by_username [get]
func GetFilesByUsername(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file/by_username"
	files, err := sendGetFileRequestByUsername(token, url)
	if err != nil {
		response(ctx, 500, "user_not_found", nil)
		return
	}

	ctx.JSON(200, gin.H{"files": files})
}
