package handler

import (
	"errors"
	"fmt"
	"github.com/TPM-Project-Larces/agent.git/schemas"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

// @Summary Login user
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body schemas.LoginRequest true "Request body"
// @Success 200 {object} schemas.AuthResponse
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_error"
// @Router /auth/login [post]
func Login(ctx *gin.Context) {

	request := schemas.LoginRequest{}
	ctx.BindJSON(&request)

	_ = bson.M{"email": request.Email, "password": request.Password}

	err := sendLoginRequest(request.Email, request.Password)
	if err != nil {
		response(ctx, 500, "login_not_successful", err)
	}

	file, err := os.Create("./config/config.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo config", err)
		return
	}
	defer file.Close()

	fmt.Fprintln(file, request.Email)
	fmt.Fprintln(file, request.Password)

	ctx.JSON(200, gin.H{"message": "user_logged"})
}

func sendLoginRequest(email string, password string) error {
	url := "http://localhost:5000/auth/login"

	token := authLogin(email, password, url)
	if token == "" {
		return errors.New("invalid token")
	}

	return nil
}
