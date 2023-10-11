package handler

import (
	"io"
	"os"

	"github.com/google/go-tpm/legacy/tpm2"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// @Summary upload file

// @Description upload a file to encrypt
// @Tags Server operations
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_decrypted"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /upload_file [post]
func DecryptFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("arquivo")
	if err != nil {
		response(ctx, 400, "bad_request", err)
	}

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
	keyHandle, _, err := tpm2.CreatePrimary(tpm, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", keyTemplate)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer tpm2.FlushContext(tpm, keyHandle)

	tempDir := "./decrypted_files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	// Abra o arquivo cifrado
	encryptedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer encryptedFile.Close()

	decryptedFile, err := os.Create(tempDir + "/" + file.Filename)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer decryptedFile.Close()
	buffer := make([]byte, 256)

	for {
		// Leia exatamente 256 bytes do arquivo cifrado
		n, err := io.ReadFull(encryptedFile, buffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// Se chegamos ao final do arquivo, saia do loop
			break
		}
		if err != nil {
			response(ctx, 500, "internal_server_error", err)
		}

		// Descriptografar o bloco
		decData, err := tpm2.RSADecrypt(tpm, keyHandle, "", buffer[:n], nil, "")
		if err != nil {
			response(ctx, 500, "internal_server_error", err)
		}

		_, err = decryptedFile.Write(decData[11:])
		if err != nil {
			response(ctx, 500, "internal_server_error", err)
		}
	}

	url := "http://localhost:3000/saved_file/"
	sendFile(decryptedFile.Name(), url)
	response(ctx, 200, "file_decrypted", err)
}
