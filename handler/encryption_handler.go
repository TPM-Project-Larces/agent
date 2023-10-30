package handler

import (
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/go-tpm/legacy/tpm2"
)

// @BasePath /
// @Summary Generate keys
// @Description Generate a pair of keys
// @Tags Encryption
// @Accept json
// @Produce json
// @Success 200 {string} string "keys_generated"
// @Failure 500 {string} string "internal_server_rror"
// @Router /encryption/generate_keys [post]
func GenerateKeys(ctx *gin.Context) {

	// Open the TPM device.
	tpmDevice := "/dev/tpmrm0"
	tpm, err := tpm2.OpenTPM(tpmDevice)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
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
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer tpm2.FlushContext(tpm, keyHandle)

	// Converts outPublic type to bytes
	publicKey, err := x509.MarshalPKIXPublicKey(outPublic)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	// Creates block public key
	blockPublicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}

	filePath := "./key/public_key.pem"
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	filePublicKey, err := os.Create(filePath)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer filePublicKey.Close()

	// Loads public key in file
	err = pem.Encode(filePublicKey, blockPublicKey)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/encryption/upload_key"
	if err := sendFile(filePath, url); err  != nil {
		response(ctx, 500, "keys_not_sent_to_server", nil)
		return
	}

	response(ctx, 200, "keys_generated", nil)
}

// @BasePath /
// @Summary Upload file
// @Description Upload a file to encrypt
// @Tags Encryption
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_uploaded"
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_rror"
// @Router /encryption/upload_file [post]
func UploadFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("file")
	if err != nil {
		response(ctx, 400, "bad_request", err)
	}

	tempDir := "./files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	tempFile, err := os.Create(tempDir + "/" + file.Filename)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/encryption/upload_file/"
	if err := sendFile(tempFile.Name(), url); err  != nil {
		response(ctx, 500, "file_not_uploaded", nil)
		return
	}

	response(ctx, 200, "file_uploaded", err)
}


// @BasePath /
// @Summary Decrypt file
// @Description Decrypt a file
// @Tags Encryption
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_decrypted"
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_error"
// @Router /encryption/decrypt_file [post]
func DecryptFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("arquivo")
	if err != nil {
		response(ctx, 400, "bad_request", err)
		return
	}

	// Open the TPM device.
	tpmDevice := "/dev/tpmrm0"
	tpm, err := tpm2.OpenTPM(tpmDevice)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
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
	keyHandle, _, err := tpm2.CreatePrimary(tpm, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", keyTemplate)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer tpm2.FlushContext(tpm, keyHandle)

	tempDir := "./decrypted_files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	encryptedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer encryptedFile.Close()

	decryptedFile, err := os.Create(tempDir + "/" + file.Filename)
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	defer decryptedFile.Close()
	buffer := make([]byte, 256)

	for {
		n, err := io.ReadFull(encryptedFile, buffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			response(ctx, 500, "internal_server_error", nil)
			return
		}

		// Decrypt block
		decData, err := tpm2.RSADecrypt(tpm, keyHandle, "", buffer[:n], nil, "")
		if err != nil {
			response(ctx, 500, "internal_server_error", nil)
			return
		}

		_, err = decryptedFile.Write(decData[11:])
		if err != nil {
			response(ctx, 500, "internal_server_error", nil)
			return
		}
	}

	url := "http://localhost:5000/encryption/saved_file/"
	if err := sendFile(decryptedFile.Name(), url); err  != nil {
		response(ctx, 500, "file_not_sent_to_server", nil)
		return
	}
	response(ctx, 200, "file_decrypted", nil)
}

