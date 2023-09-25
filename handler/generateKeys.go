package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// @Summary generate keys

// @Description generat a pair of keys
// @Tags User operations
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /generate_keys [post]
func GenerateKeys(ctx *gin.Context) {
	// Gerar um par de chaves RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Erro ao gerar a chave privada:", err)
		return
	}

	// Codificar a chave pública em formato PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Erro ao codificar a chave pública em formato PEM:", err)
		return
	}

	encodedPub := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	if err := ioutil.WriteFile("public_key.pem", encodedPub, 0600); err != nil {
		log.Fatalf("failed to write PEM data to file: %v", err)
	}

	// Codificar a chave privada em formato PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	encodedPriv := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	if err := ioutil.WriteFile("private_key.pem", encodedPriv, 0600); err != nil {
		log.Fatalf("failed to write PEM data to file: %v", err)
	}

	ctx.JSON(http.StatusOK, LerArquivo("public_key.pem"))

}
