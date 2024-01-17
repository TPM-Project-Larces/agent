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
	CreatedAt        time.Time          `bson:"created_at"`
}
