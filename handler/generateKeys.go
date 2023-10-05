package handler

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/go-tpm/legacy/tpm2"
)

// @BasePath /

// @Summary generate keys

// @Description generat a pair of keys
// @Tags User operations
// @Accept json
// @Produce json
// @Success 200 {string} Keys_generated
// @Router /generate_keys [post]
func GenerateKeys(ctx *gin.Context) {

	// Open the TPM device.
	tpmDevice := "/dev/tpmrm0"
	tpm, err := tpm2.OpenTPM(tpmDevice)
	handleError("Error opening TPM device", err)
	defer tpm.Close()

	// Creates primary key template
	keyTemplate := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedParent | tpm2.FlagFixedTPM | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth | tpm2.FlagDecrypt,
		AuthPolicy: nil,
		RSAParameters: &tpm2.RSAParams{
			KeyBits:    2048,
			ModulusRaw: make([]byte, 256),
		},
	}

	// Creates the primary key in the TPM.
	keyHandle, outPublic, err := tpm2.CreatePrimary(tpm, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", keyTemplate)
	handleError("Error creating primary key", err)
	defer tpm2.FlushContext(tpm, keyHandle)

	// Converts outPublic type to bytes
	publicKey, err := x509.MarshalPKIXPublicKey(outPublic)
	handleError("Error marshaling primary key", err)

	// Creates block public key
	blockPublicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}

	filePath := "./key/public_key.pem"

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		handleError("Error creating directory", err)
		ctx.JSON(http.StatusInternalServerError, "Error creating directory")
		return
	}

	filePublicKey, err := os.Create(filePath)
	if err != nil {
		handleError("Error creating file public key in PEM format", err)
		ctx.JSON(http.StatusInternalServerError, "Error creating public key file")
		return
	}
	defer filePublicKey.Close()

	url := "http://localhost:5000/upload_key/"

	// Loads public key in file
	err = pem.Encode(filePublicKey, blockPublicKey)
	handleError("Error enconding block public key in PEM file", err)

	sendFile(filePath, url)

	ctx.JSON(http.StatusOK, "keys_generated")
}
