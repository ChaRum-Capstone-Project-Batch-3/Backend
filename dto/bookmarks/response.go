package bookmarks

import (
	dtoThread "charum/dto/threads"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Id     primitive.ObjectID `json:"_id"`
	UserID primitive.ObjectID `json:"userID"`
	Thread dtoThread.Response `json:"thread"`
}
