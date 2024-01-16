package schemas

import "github.com/TPM-Project-Larces/agent.git/model"

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LoginRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type AuthResponse struct {
	Message string `bson:"message"`
}

type AuthRequest struct {
	Token model.Token `bson:"token"`
}
