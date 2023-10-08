package handler

import (
	"io"
	"os"
	_ "path/filepath"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// @Summary upload file

// @Description upload a file to encrypt
// @Tags User operations
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_uploaded"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /upload_file [post]
func UploadFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("file")
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	tempDir := "./files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	tempFile, err := os.Create(tempDir + "/" + file.Filename)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	// Get the path of the temporary file
	tempFilePath := tempFile.Name()

	url := "http://localhost:3000/upload_file/"

	sendFile(tempFilePath, url)

	response(ctx, 200, "file_uploaded", err)
}
