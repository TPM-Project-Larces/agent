package handler

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /
// @Summary Find user by username
// @Description Provide the user data
// @Tags User
// @Produce json
// @Param username query string true "User`s username to find"
// @Success 200 {object} schemas.ShowUserResponse
// @Failure 400 {string} string "bad_request"
// @Failure 404 {string} string "not_found"
// @Failure 500 {string} string "internal_server_error"
// @Router /user/username [get]
func GetUserByUsername(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	username := ctx.Query("username")

	token, err := Login()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/user/username?username=" + username
	user, err := sendGetRequestForUser(token, url)
	if err != nil {
		response(ctx, 500, "user_not_found", nil)
		return
	}

	ctx.JSON(200, gin.H{"message": "user", "user": user})
}

// @BasePath /
// @Summary Get all users
// @Description Get all users
// @Tags User
// @Produce json
// @Success 200 {object} schemas.ListUsersResponse
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_rror"
// @Router /user [get]
func GetUsers(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Login()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/user"
	users, err := sendGetRequest(token, url)
	if err != nil {
		response(ctx, 500, "users_not_get", nil)
		return
	}

	ctx.JSON(200, gin.H{"message": "all_users", "users": users})
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

	token, err := Login()
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
