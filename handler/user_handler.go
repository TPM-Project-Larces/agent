package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/TPM-Project-Larces/agent.git/schemas"
	"github.com/gin-gonic/gin"
)

// @BasePath /
// @Summary Create user
// @Description Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body schemas.CreateUserRequest true "Request body"
// @Success 200 {object} schemas.CreateUserResponse
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_error"
// @Router /user [post]
func CreateUser(ctx *gin.Context) {
	request := schemas.CreateUserRequest{}
	ctx.BindJSON(&request)

	jsonData, err := json.Marshal(request)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	resp, err := http.Post("http://localhost:5000/user",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	response(ctx, 200, "user_created", err)
}

// @BasePath /
// @Summary Update user
// @Description Updates a user
// @Tags User
// @Produce json
// @Param user body schemas.UpdateUserRequest true "User data to Update"
// @Success 200 {object} schemas.UpdateUserResponse
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "user_not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /user [put]
func UpdateUser(ctx *gin.Context) {
	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
	}

	request := schemas.UpdateUserRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	req, err := http.NewRequest("PUT",
		"http://localhost:5000/user",
		bytes.NewBuffer(jsonData))
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response(ctx, 500, "internal_server_error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response(ctx, 500, "internal_server_error", err)
		return
	}

	ctx.JSON(200, gin.H{"message": "user_updated"})
}

// @BasePath /
// @Summary Find user by username
// @Description Provide the user data
// @Tags User
// @Produce json
// @Success 200 {object} schemas.ShowUserResponse
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /user/username [get]
func GetUserByUsername(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/user/username"
	user, err := sendGetRequestForUser(token, url)
	if err != nil {
		response(ctx, 500, "user_not_found", nil)
		return
	}

	ctx.JSON(200, gin.H{"message": "user", "user": user})
}

// @BasePath /
// @Summary Delete user
// @Description Delete a user
// @Tags User
// @Produce json
// @Success 200 {string} string "user_deleted"
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_rror"
// @Router /user [delete]
func DeleteUser(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Auth()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/user"
	if err := sendDeleteRequest(token, url); err != nil {
		response(ctx, 500, "file_not_uploaded", nil)
		return
	}

	response(ctx, 200, "user_deleted", nil)
}
