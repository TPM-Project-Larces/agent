package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func handleError(message string, err error) {
	if err != nil {
		fmt.Println(message+":", err)
		panic((err))
	}
}

func sendFile(fileName string, url string) {
	// Abra o arquivo que você deseja enviar
	arquivo, err := os.Open(fileName)
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

	// Excluir arquivo se nao for a chave
	if fileName != "./key/public_key.pem" {
		os.Remove(fileName)
	}

	// Verifique a resposta do servidor
	if response.StatusCode == http.StatusOK {
		return
	} else {
		return
	}
}
