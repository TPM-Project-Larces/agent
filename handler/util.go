package handler

import (
	"io/ioutil"
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
