package dto

import (
	"charum/business/users"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseFollowThread struct {
	Id        primitive.ObjectID `json:"_id"`
	User      users.Domain       `json:"user"`
	Thread    ResponseThread     `json:"thread"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
