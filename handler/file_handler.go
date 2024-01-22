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
// @Param username query string true "Username"
// @Success 200 {object} schemas.ListFilesResponse
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_error"
// @Router /file/by_username [get]
func GetFilesByUsername(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	username := ctx.Query("username")

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file/by_username/" + username
	files, err := sendGetFileRequestByUsername(token, url)
	if err != nil {
		response(ctx, 500, "user_not_found", nil)
		return
	}

	ctx.JSON(200, gin.H{"message": "files_by_" + username, "files": files})
}

// @BasePath /
// @Summary Find file by name
// @Description Provide the file data
// @Tags File
// @Produce json
// @Param filename query string true "filename to find"
// @Success 200 {object} schemas.ShowFileResponse
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /file/by_name [get]
func GetFileByName(ctx *gin.Context) {

	filename := ctx.Query("filename")

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file/" + filename
	files, err := sendGetFileRequestByName(token, url)
	if err != nil {
		response(ctx, 500, "files_not_get", nil)
		return
	}

	ctx.JSON(200, gin.H{"message": "all_encrypted_files", "files": files})
}
