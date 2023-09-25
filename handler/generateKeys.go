package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /

// @Summary generate keys

// @Description generat a pair of keys
// @Tags User operations
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /generate_keys [post]
func GenerateKeys(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "chaves_criadas")
}
