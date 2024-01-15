package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

// @BasePath /
// @Summary Get all encrypted files
// @Description Get a list of all encrypted files
// @Tags File
// @Accept json
// @Produce json
// @Success 200 {object} schemas.ListFilesResponse
// @Failure 500 {string} string "internal_server_error"
// @Router /file [get]
func GetFiles(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Login()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file"
	files, err := sendGetFileRequestByUsername(token, url)
	if err != nil {
		response(ctx, 500, "files_not_get", nil)
		return
	}

	ctx.JSON(200, gin.H{"message": "all_files", "files": files})
}

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

	token, err := Login()
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

	token, err := Login()
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
