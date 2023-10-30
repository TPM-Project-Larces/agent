package handler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func response(ctx *gin.Context, code int, message string, err error) {
	response := gin.H{
		"code": code,
	}

	if message != "" {
		response["message"] = message
	}

	if err != nil {
		response["error"] = err.Error()
	}

	ctx.JSON(code, response)
}

func sendFile(fileName string, url string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("arquivo", fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if fileName != "./key/public_key.pem" {
		os.Remove(fileName)
	}

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("file not sent")
	}
}
