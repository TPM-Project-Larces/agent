package schemas

import (
	"github.com/TPM-Project-Larces/agent.git/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateUserRequest struct {
	CPF         string        `bson:"cpf"`
	Name        string        `bson:"name"`
	Username    string        `bson:"username"`
	DateOfBirth string        `bson:"date_of_birth"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
	Contact     model.Contact `bson:"contact"`
	Address     model.Address `bson:"address"`
}

type DeleteUserRequest struct {
	Username string `bson:"username"`
}

type UpdateUserRequest struct {
	CPF         string        `bson:"cpf"`
	Name        string        `bson:"name"`
	Username    string        `bson:"username"`
	DateOfBirth string        `bson:"date_of_birth"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
	Contact     model.Contact `bson:"contact"`
	Address     model.Address `bson:"address"`
}

type UserResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty"`
	CPF         string             `bson:"cpf"`
	Name        string             `bson:"name"`
	Username    string             `bson:"username"`
	DateOfBirth string             `bson:"date_of_birth"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	Contact     model.Contact      `bson:"contact"`
	Address     model.Address      `bson:"address"`
}

type CreateUserResponse struct {
	Message string       `bson:"message"`
	Data    UserResponse `bson:"data"`
}

type DeleteUserResponse struct {
	Message string       `bson:"message"`
	Data    UserResponse `bson:"data"`
}

type ShowUserResponse struct {
	Data UserResponse `bson:"data"`
}

type ListUsersResponse struct {
	Data []UserResponse `bson:"data"`
}

type UpdateUserResponse struct {
	Message string       `bson:"message"`
	Data    UserResponse `bson:"data"`
}
