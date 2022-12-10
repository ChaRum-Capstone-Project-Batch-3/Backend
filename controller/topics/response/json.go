package response

import (
	"charum/business/topics"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Topic       string             `json:"topic" bson:"topic"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"imageURL" bson:"imageURL"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain topics.Domain) Topic {
	return Topic{
		Id:          domain.Id,
		Topic:       domain.Topic,
		Description: domain.Description,
		ImageURL:    domain.ImageURL,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func FromDomainArray(data []topics.Domain) []Topic {
	var array []Topic
	for _, v := range data {
		array = append(array, FromDomain(v))
	}
	return array
}
