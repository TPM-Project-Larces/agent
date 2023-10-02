package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
// @Success 200 {string} file_uploaded
// @Router /upload_file [post]
func UploadFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("file")
	handleError("Error retrieving file from the form", err)

	fmt.Print("File uploaded: " + file.Filename)

	tempDir := "./files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	handleError("Error creating 'files' directory", err)

	ext := filepath.Ext(file.Filename)

	tempFile, err := ioutil.TempFile(tempDir, "file_to_encrypt"+ext)
	handleError("Error creating temporary file", err)
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	handleError("Error opening the file", err)
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	handleError("Error copying the file to the destination", err)

	// Get the path of the temporary file
	tempFilePath := tempFile.Name()

	url := "http://localhost:3000/upload_file/"

	sendFile(tempFilePath, url)

	ctx.JSON(http.StatusOK, gin.H{"message": "File successfully sent to the server"})
}
