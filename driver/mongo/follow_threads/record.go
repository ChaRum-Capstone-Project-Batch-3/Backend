package follow_threads

import (
	followThreads "charum/business/follow_threads"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	UserID       primitive.ObjectID `json:"userID" bson:"userID"`
	ThreadID     primitive.ObjectID `json:"threadID" bson:"threadID"`
	Notification int                `json:"notification" bson:"notification"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *followThreads.Domain) *Model {
	return &Model{
		Id:           domain.Id,
		UserID:       domain.UserID,
		ThreadID:     domain.ThreadID,
		Notification: domain.Notification,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}
}

func (m *Model) ToDomain() followThreads.Domain {
	return followThreads.Domain{
		Id:           m.Id,
		UserID:       m.UserID,
		ThreadID:     m.ThreadID,
		Notification: m.Notification,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToDomainArray(model []Model) []followThreads.Domain {
	var domain []followThreads.Domain
	for _, v := range model {
		domain = append(domain, v.ToDomain())
	}
	return domain
}
