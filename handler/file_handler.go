package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

// @BasePath /
// @Summary Delete file
// @Description deletes a file
// @Tags File
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_deleted"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /file/delete_file [post]
func DeleteFile(ctx *gin.Context) {

	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("file")
	if err != nil {
		response(ctx, 400, "bad_request", err)
		return
	}

	tempDir := "./files"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	token, err := Login()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file/" + file.Filename
	if err := sendFileDeleteRequest(tempFilePath, token, url); err != nil {
		response(ctx, 500, "file_not_deleted", nil)
		return
	}

	response(ctx, 200, "file_deleted", nil)
}
