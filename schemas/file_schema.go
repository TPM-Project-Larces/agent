package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EncryptedFileResponse struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Name     string             `bson:"name"`
	Data     []byte             `bson:"data"`
	//AnonymizedFile AnonymizedFile     `bson:"anonymized_file"`
	CreatedAt time.Time  `bson:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
}

type DeleteFileRequest struct {
	Filename string `bson:"filename"`
}

type ListFilesResponse struct {
	Message string                  `bson:"message"`
	Data    []EncryptedFileResponse `bson:"data"`
}
type ShowFileResponse struct {
	Message string                `bson:"message"`
	Data    EncryptedFileResponse `bson:"data"`
}
type DeleteFileResponse struct {
	Message string                `bson:"message"`
	Data    EncryptedFileResponse `bson:"data"`
}
