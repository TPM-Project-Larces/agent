package router

import (
	"github.com/gin-gonic/gin"
)

func Initialize() {
	// initialize router
	router := gin.Default()

	//initialize routes
	initializeRoutes(router)

	// run the server
	router.Run(":3000")
}
