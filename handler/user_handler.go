package handler

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /
// @Summary Upload user
// @Description Upload a user
// @Tags User
// @Produce json
// @Success 200 {string} string "user_uploaded"
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_rror"
// @Router /user [put]
func UploadUser(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Login()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/user"
	if err := sendPutRequest(token, url); err != nil {
		response(ctx, 500, "file_not_uploaded", nil)
		return
	}

	response(ctx, 200, "file_uploaded", nil)
}

// @BasePath /
// @Summary Get all users
// @Description Get all users
// @Tags User
// @Produce json
// @Success 200 {object} schemas.ListUsersResponse
// @Failure 400 {string} string "bad_request"
// @Failure 500 {string} string "internal_server_rror"
// @Router /user [put]
func GetUsers(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	token, err := Login()
	if err != nil {
		response(ctx, 500, "internal_server_error", nil)
		return
	}

	url := "http://localhost:5000/user"
	if err := sendPutRequest(token, url); err != nil {
		response(ctx, 500, "file_not_uploaded", nil)
		return
	}

	response(ctx, 200, "file_uploaded", nil)
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
