package main

import (
	"github.com/TPM-Project-Larces/agent.git/router"
)

// @title Agent API
// @description Agent Operations
// @version 1.0.0
// @contact {
//   name: "Computer Networks and Security Laboratory (LARCES)",
//   url: "https://larces.uece.br/",
//   email: "larces@uece.br
// }
// @BasePath /
func main() {
	router.InitializeRoutes()
}
