package handler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func LerArquivo(nomeArquivo string) string {
	// Lê o conteúdo do arquivo
	conteudo, erro := ioutil.ReadFile(nomeArquivo)
	if erro != nil {
		return "erro na leitura"
	}

	// Converte os bytes lidos em uma string
	conteudoString := string(conteudo)

	return conteudoString
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func handleError(message string, err error) {
	if err != nil {
		fmt.Println(message+":", err)
		panic((err))
	}
}
