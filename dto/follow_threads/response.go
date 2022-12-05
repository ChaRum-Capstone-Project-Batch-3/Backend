package follow_threads

import (
	"charum/business/users"
	"charum/dto/threads"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Id           primitive.ObjectID `json:"_id"`
	User         users.Domain       `json:"user"`
	Thread       threads.Response   `json:"thread"`
	Notification int                `json:"notification"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
