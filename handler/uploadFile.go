package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

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
	handleError("Erro ao recuperar o arquivo do formulário", err)

	fmt.Print("Arquivo enviado: " + file.Filename)

	tempDir := "./files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	handleError("Erro ao criar diretório 'files'", err)

	ext := filepath.Ext(file.Filename)

	tempFile, err := ioutil.TempFile(tempDir, "file_to_encrypt_*"+ext)
	handleError("Erro ao criar arquivo temporário", err)
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	handleError("Erro ao abrir o arquivo", err)
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	handleError("Erro ao copiar arquivo para o destino", err)

	//respomse = sendFile()

	ctx.JSON(http.StatusOK, "file_uploaded_and_sent")
}

/*func sendFile(fileName)  {
	// Abra o arquivo que você deseja enviar
    arquivo, err := os.Open("file" + fileName)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer arquivo.Close()

    // Crie um buffer para a solicitação multipart/form-data
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Adicione o arquivo ao formulário
    part, err := writer.CreateFormFile("arquivo", fileName)
    if err != nil {
        fmt.Println(err)
        return
    }

    _, err = io.Copy(part, arquivo)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Finalize o formulário
    writer.Close()

    // Faça uma solicitação HTTP POST para o servidor
    url := "https://exemplo.com/upload" // Substitua pela URL real do servidor
    request, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Defina o cabeçalho Content-Type para multipart/form-data
    request.Header.Set("Content-Type", writer.FormDataContentType())

    // Faça a solicitação
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer response.Body.Close()

	// Excluir arquivo

    // Verifique a resposta do servidor
    if response.StatusCode == http.StatusOK {
        return 200
    } else {
        return 500
}*/
