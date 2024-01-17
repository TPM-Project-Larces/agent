package handler

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
// @Description Search for a file to decrypt
// @Tags Encryption
// @Accept multipart/form-data
// @Produce json
// @Param filename query string true "Filename"
// @Success 200 {string} string "file_decrypted"
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_error"
// @Router /encryption/decrypt_file [post]
func DecryptFile(ctx *gin.Context) {
	filename := ctx.Query("filename")
	url := "http://localhost:5000/encryption/search_file"

	err := sendString(filename, url)
	if err != nil {
		response(ctx, http.StatusInternalServerError, "Error sending filename", err)
		return
	}

	// Exemplo de retorno de sucesso
	response(ctx, http.StatusOK, "String sent successfully", nil)
}

// @BasePath /
// @Summary Save file
// @Description Save file
// @Tags Encryption
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "file_saved"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /encryption/save_file [post]
func SaveFile(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	file, err := ctx.FormFile("arquivo")
	if err != nil {
		response(ctx, 400, "bad_request", err)
	}

	//Abra o arquivo diretamente sem salvá-lo no disco
	uploadedFile, err := file.Open()
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}
	defer uploadedFile.Close()

	data, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	// Split data into smaller blocks
	maxBlockSize := 245
	var encryptedBlocks []byte
	for len(data) > 0 {
		blockSize := len(data)
		if blockSize > maxBlockSize {
			blockSize = maxBlockSize
		}

		//Writes the file in blocks of bytes
		encryptedBlock := data[:blockSize]

		encryptedBlocks = append(encryptedBlocks, encryptedBlock...)
		data = data[blockSize:]
	}

	tempDir := "./server_files"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	tempfile, err := os.Create(tempDir + "/" + file.Filename)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	defer tempfile.Close()

	err = ioutil.WriteFile(tempfile.Name(), encryptedBlocks, 0644)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
	}

	response(ctx, 200, "file_saved", err)
}

// @BasePath /
// @Summary Get size and decrypt
// @Description Get size and decrypt
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body StringData true "Request body"
// @Success 200 {object} string "file_decrypted"
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /encryption/size_and_decrypt [post]
func SizeAndDecrypt(ctx *gin.Context) {
	request := StringData{}
	ctx.BindJSON(&request)

	size := StringData{
		Data: request.Data,
	}

	fmt.Println("aquiiii tamanho", size.Data)
	sizeInt, _ := strconv.Atoi(size.Data)
	rest := sizeInt % 258
	total := sizeInt / 258

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
	fmt.Println("tst2")
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
	serverFilesDir := "./server_files"

	err = filepath.Walk(serverFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		encryptedFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer encryptedFile.Close()

		decryptedFilePath := tempDir + "/" + info.Name()
		decryptedFile, err := os.Create(decryptedFilePath)
		if err != nil {
			return err
		}
		defer decryptedFile.Close()

		buffer := make([]byte, 256)

		if rest != 0 {
			for i := 0; i < total; i++ {
				n, err := io.ReadFull(encryptedFile, buffer)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				if err != nil {
					response(ctx, 500, "internal_server_error", nil)
					return nil
				}

				// Decrypt block
				decData, err := tpm2.RSADecrypt(tpm, keyHandle, "", buffer[:n], nil, "")
				if err != nil {
					response(ctx, 500, "internal_server_error", nil)
					return err
				}

				_, err = decryptedFile.Write(decData[11:])
				if err != nil {
					response(ctx, 500, "internal_server_error", nil)
					return err
				}
			}

			n, err := io.ReadFull(encryptedFile, buffer)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
			if err != nil {
				response(ctx, 500, "internal_server_error", nil)
				return err
			}

			// Decrypt block
			decData, err := tpm2.RSADecrypt(tpm, keyHandle, "", buffer[:n], nil, "")
			if err != nil {
				response(ctx, 500, "internal_server_error", nil)
				return err
			}

			_, err = decryptedFile.Write(decData[256-rest:])
			if err != nil {
				response(ctx, 500, "internal_server_error", nil)
				return err
			}
		} else {
			for {
				n, err := io.ReadFull(encryptedFile, buffer)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				if err != nil {
					response(ctx, 500, "internal_server_error", nil)
					return err
				}

				// Decrypt block
				decData, err := tpm2.RSADecrypt(tpm, keyHandle, "", buffer[:n], nil, "")
				if err != nil {
					response(ctx, 500, "internal_server_error", nil)
					return err
				}

				_, err = decryptedFile.Write(decData[11:])
				if err != nil {
					response(ctx, 500, "internal_server_error", nil)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}
	os.RemoveAll(serverFilesDir)
}

/*

// @BasePath /
// @Summary Decrypt a file
// @Description Decrypt a file stored in server
// @Tags Encryption
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
*/
