package comments

import (
	"charum/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Id        primitive.ObjectID `json:"_id"`
	ThreadID  primitive.ObjectID `json:"threadID"`
	User      users.Domain       `json:"user"`
	Comment   string             `json:"comment"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt"`
}
