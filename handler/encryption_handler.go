package handler

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

	token, err := Auth()
	if err != nil || token == "" {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/encryption/upload_key"
	if err := sendFile(filePath, token, url); err != nil {
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
	fmt.Println(tempFile.Name())

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	fileChallenge, FileSign, FileAttestationKey := AttestationTPM()

	urlattestation := "http://localhost:5000/attestation/upload_challenge/"
	if err := sendFile(fileChallenge, token, urlattestation); err != nil {
		response(ctx, 500, "challenge_failed", err)
		return
	}

	urlattestation2 := "http://localhost:5000/attestation/upload_signature/"
	if err := sendFile(FileSign, token, urlattestation2); err != nil {
		response(ctx, 500, "signature_failed", err)
		return
	}

	urlattestation3 := "http://localhost:5000/attestation/upload_attestation_key/"
	if err := sendFile(FileAttestationKey, token, urlattestation3); err != nil {
		response(ctx, 500, "attestation_key_failed", err)
		return
	}

	url := "http://localhost:5000/file/upload_file"
	if err := sendFile(tempFile.Name(), token, url); err != nil {
		response(ctx, 500, "file_not_uploaded", err)
		return
	}

	response(ctx, 200, "file_uploaded", err)
}

// @BasePath /
// @Summary Upload encrypted file
// @Description Upload a file
// @Tags Encryption
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_uploaded"
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_rror"
// @Router /encryption/upload_encrypted_file [post]
func UploadEncryptedFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("file")
	if err != nil {
		response(ctx, 400, "bad_request", err)
		return
	}

	tempDir := "./files"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file/upload_encrypted_file/"
	if err := sendFile(tempFilePath, token, url); err != nil {
		response(ctx, 500, "file_not_uploaded", nil)
		return
	}

	response(ctx, 200, "file_uploaded", nil)
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
	fmt.Println(file.Filename)

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
	fmt.Println("tst3")
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

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/file/upload_encrypted_file"

	if err := sendFile(decryptedFile.Name(), token, url); err != nil {
		response(ctx, 500, "file_not_sent_to_server", nil)
		return
	}
	response(ctx, 200, "file_decrypted", nil)
}

// @BasePath /
// @Summary Decrypt a file
// @Description Decrypt a file stored in server
// @Tags File
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_decrypted"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /encryption/decrypt_saved_file [post]
func DecryptServerFile(ctx *gin.Context) {

	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("file")
	if err != nil {
		response(ctx, 400, "bad_request", err)
		return
	}

	tempDir := "./files"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	tempFilePath := filepath.Join(tempDir, file.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer tempFile.Close()

	uploadedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/encryption/decrypt_file"
	if err := sendFile(tempFilePath, token, url); err != nil {
		response(ctx, 500, "file_not_decrypted", nil)
		return
	}

	response(ctx, 200, "file_decrypted", nil)
}
