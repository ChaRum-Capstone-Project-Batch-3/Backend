package topics

import (
	"charum/business/topics"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Topic       string             `json:"topic" bson:"topic"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"imageURL" bson:"imageURL"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *topics.Domain) *Model {
	return &Model{
		Id:          domain.Id,
		Topic:       domain.Topic,
		Description: domain.Description,
		ImageURL:    domain.ImageURL,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func (topic *Model) ToDomain() topics.Domain {
	return topics.Domain{
		Id:          topic.Id,
		Topic:       topic.Topic,
		Description: topic.Description,
		ImageURL:    topic.ImageURL,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,
	}
}

func ToArrayDomain(data []Model) []topics.Domain {
	var result []topics.Domain
	for _, v := range data {
		result = append(result, v.ToDomain())
	}
	return result
}
