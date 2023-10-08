package handler

import (
	"crypto/x509"
	"encoding/pem"
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
// @Success 200 {string} string "keys_generated"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /generate_keys [post]
func GenerateKeys(ctx *gin.Context) {

	// Open the TPM device.
	tpmDevice := "/dev/tpmrm0"
	tpm, err := tpm2.OpenTPM(tpmDevice)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
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
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer tpm2.FlushContext(tpm, keyHandle)

	// Read key public part
	//fmt.Println(tpm2.ReadPublic(tpm, keyHandle))
	//fmt.Println("\nPublic part: \n", outPublic)

	// Converts outPublic type to bytes
	publicKey, err := x509.MarshalPKIXPublicKey(outPublic)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	// Creates block public key
	blockPublicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}

	filePath := "./key/public_key.pem"

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		if err != nil {
			response(ctx, 500, "internal_server_error", err)
		}

	}

	filePublicKey, err := os.Create(filePath)
	if err != nil {
		if err != nil {
			response(ctx, 500, "internal_server_error", err)
		}

	}
	defer filePublicKey.Close()

	url := "http://localhost:3000/upload_key/"

	// Loads public key in file
	err = pem.Encode(filePublicKey, blockPublicKey)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	sendFile(filePath, url)

	response(ctx, 200, "keys_generated", err)
}
