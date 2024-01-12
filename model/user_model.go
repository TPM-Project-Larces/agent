package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Street  string
	City    string
	State   string
	ZIPCode string
}

type Contact struct {
	Celphone string
	Phone    string
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CPF         string             `bson:"cpf"`
	Name        string             `bson:"name"`
	Username    string             `bson:"username"`
	DateOfBirth string             `bson:"date_of_birth"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	Contact     Contact            `bson:"contact"`
	Address     Address            `bson:"address"`
}
