package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EncryptedFile struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Username         string             `bson:"username"`
	Name             string             `bson:"name"`
	Data             []byte             `bson:"data"`
	LocallyEncrypted bool               `bson:"locally_encrypted"`
	//AnonymizedFile AnonymizedFile     `bson:"anonymized_file"`
	CreatedAt time.Time `bson:"created_at"`
}
