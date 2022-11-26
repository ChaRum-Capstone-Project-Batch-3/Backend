package topics

import (
	"charum/business/topics"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Topic       string             `json:"topic" bson:"topic"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *topics.Domain) *Model {
	return &Model{
		Id:          domain.Id,
		Topic:       domain.Topic,
		Description: domain.Description,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func (topic *Model) ToDomain() topics.Domain {
	return topics.Domain{
		Id:          topic.Id,
		Topic:       topic.Topic,
		Description: topic.Description,
		CreatedAt:   topic.CreatedAt,
		UpdatedAt:   topic.UpdatedAt,
	}
}
