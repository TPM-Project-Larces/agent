package handler

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func response(ctx *gin.Context, code int, message string, err error) {
	response := gin.H{
		"code": code,
	}

	if message != "" {
		response["message"] = message
	}

	if err != nil {
		response["error"] = err.Error()
	}

	ctx.JSON(code, response)
}

func sendPutRequest(token string, url string) error {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func sendDeleteRequest(token string, url string) error {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func sendFile(fileName string, token string, url string) error {

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("arquivo", fileName)
	if err != nil {
		fmt.Println("Erro ao criar o campo de arquivo:", err)
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Erro ao copiar o conteúdo do arquivo para o formulário:", err)
		return err
	}

	err = writer.WriteField("token", token)
	if err != nil {
		fmt.Println("Erro ao adicionar o campo de token ao formulário:", err)
		return err
	}

	writer.Close()

	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Erro ao criar a requisição HTTP:", err)
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação HTTP para a outra API:", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("file not sent")
	}
}

func sendFileDeleteRequest(fileName string, token string, url string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("arquivo", filepath.Base(fileName))
	if err != nil {
		fmt.Println("Erro ao criar o campo de arquivo:", err)
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Erro ao copiar o conteúdo do arquivo para o formulário:", err)
		return err
	}

	err = writer.WriteField("token", token)
	if err != nil {
		fmt.Println("Erro ao adicionar o campo de token ao formulário:", err)
		return err
	}

	request, err := http.NewRequest("DELETE", url, &requestBody)
	if err != nil {
		fmt.Println("Erro ao criar a requisição HTTP:", err)
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação HTTP para a outra API:", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("request not successful")
	}
}

// Challenge represents the structure of a TPM challenge
type Challenge struct {
	Nonce []byte // Random value to prevent replay attacks
}

func AttestationTPM() (string, string, string) {
	// Open TPM
	rw, err := tpm2.OpenTPM("/dev/tpmrm0")
	if err != nil {
		log.Fatalf("Error opening TPM: %v", err)
	}
	defer rw.Close()

	HandleAttestation, outPublic, err := createAttestationKey(rw)
	if err != nil {
		log.Fatalf("Erro creating attestation key: %v", err)
	}

	PublicKey, err := converttoRSAKey(outPublic)
	if err != nil {
		log.Fatalf("Error convert to RSA Key: %v", err)
	}

	challenge, err := generateChallenge()
	if err != nil {
		log.Fatalf("Error generating challenge: %v", err)
	}

	// Serialize the challenge to a byte slice
	serializedChallenge, err := SerializeChallenge(challenge)
	if err != nil {
		log.Fatalf("Error serializing challenge: %v", err)
	}

	// TPM-specific operation using the serialized challenge
	sigScheme := tpm2.SigScheme{
		Alg:  tpm2.AlgRSASSA,
		Hash: tpm2.AlgSHA256,
	}

	// TPM-specific operation using the serialized challenge
	signature, err := tpm2.Sign(rw, HandleAttestation, "", serializedChallenge, nil, &sigScheme)
	if err != nil {
		log.Fatalf("Error signing challenge: %v", err)
	}

	rawSignature := signature.RSA.Signature

	fileChallenge, FileSign := uploadingfiles(rawSignature, serializedChallenge)

	return fileChallenge, FileSign, PublicKey
}

func uploadingfiles(rawSignature tpmutil.U16Bytes, challenge []byte) (string, string) {
	// Creates block challenge
	blockChallenge := &pem.Block{
		Type:  "CHALLENGE",
		Bytes: challenge,
	}

	filePath := "./challenge/message.pem"

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("error creating directory: %v", err)
	}

	fileChallenge, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("error creating challenge file: %v", err)
	}
	defer fileChallenge.Close()

	// Loads challenge in file
	err = pem.Encode(fileChallenge, blockChallenge)

	blockSignature := &pem.Block{
		Type:  "SIGNATURE",
		Bytes: rawSignature,
	}

	fileSign := "./Signature/Signature_key.pem"

	dir2 := filepath.Dir(fileSign)

	if err := os.MkdirAll(dir2, os.ModePerm); err != nil {
		log.Fatalf("error creating directory: %v", err)
	}

	fileSignature, err := os.Create(fileSign)
	if err != nil {
		log.Fatalf("error creating Signature file: %v", err)
	}
	defer fileSignature.Close()

	// Loads challenge in file
	err = pem.Encode(fileSignature, blockSignature)

	return filePath, fileSign

}

func converttoRSAKey(outPublic crypto.PublicKey) (string, error) {
	// Converts outPublic type to bytes
	publicKey, err := x509.MarshalPKIXPublicKey(outPublic)
	if err != nil {
		return "", err
	}
	// Creates block public key
	blockPublicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}

	fileKeySign := "./key/public_attestation_key.pem"

	dir := filepath.Dir(fileKeySign)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	filePublicKey, err := os.Create(fileKeySign)
	if err != nil {
		return "", err
	}
	defer filePublicKey.Close()

	// Loads public key in file
	err = pem.Encode(filePublicKey, blockPublicKey)

	return fileKeySign, nil
}

// serializeChallenge serializes the Challenge struct to a byte slice
func SerializeChallenge(challenge *Challenge) ([]byte, error) {
	return challenge.Nonce, nil // Adjust serialization based on your requirements
}

func createAttestationKey(rw io.ReadWriter) (tpmutil.Handle, crypto.PublicKey, error) {
	KeyParams := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedParent | tpm2.FlagFixedTPM | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth | tpm2.FlagSign,
		AuthPolicy: nil,
		RSAParameters: &tpm2.RSAParams{
			KeyBits:     2048,
			ExponentRaw: 65537,
			ModulusRaw:  make([]byte, 256),
		},
	}

	// Crie uma nova chave no TPM para atestação.
	handle, pub, err := tpm2.CreatePrimary(rw, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", KeyParams)
	if err != nil {
		log.Fatalf("Error creating attestation key: %v", err)
	}

	return handle, pub, nil
}

// createChallenge generates a random nonce for the challenge
func generateChallenge() (*Challenge, error) {
	nonce := make([]byte, 32) // Adjust the size based on your requirements
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return &Challenge{Nonce: nonce}, nil
}

func authLogin(email string, password string, url string) string {

	// Construir o payload do corpo da solicitação
	payload := map[string]string{"email": email, "password": password}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	// Criar uma solicitação HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return ""
	}

	// Configurar os cabeçalhos da solicitação
	req.Header.Set("Content-Type", "application/json")

	// Executar a solicitação
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Verificar o código de status da resposta
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	// Ler o corpo da resposta
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return ""
	}

	// Extrair o token da resposta
	token, ok := responseBody["token"].(string)
	if !ok {
		return ""
	}

	return token
}

func Login() (string, error) {

	f := "./config/config.txt"

	file, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err
	}

	// Extrair email e senha do conteúdo do arquivo
	lines := bytes.Split(file, []byte("\n"))
	if len(lines) < 2 {
		return "", err
	}

	email := string(lines[0])
	password := string(lines[1])

	url := "http://localhost:5000/auth/login"

	token := authLogin(email, password, url)
	if token == "" {
		return "", fmt.Errorf("invalid_token")
	}
	fmt.Println(token)
	return token, nil
}

/*func sendString(stringToSend string, token string, url string) error {

	requestBody := fmt.Sprintf(`{"string": "%s"}`, stringToSend)

	// Faça a solicitação POST para a segunda API com o token no cabeçalho
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	// Envie a solicitação
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Verifique o código de status da resposta
	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("request not successful")
	}
}*/
