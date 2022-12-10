package threads

import (
	"charum/business/threads"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	TopicID       primitive.ObjectID `json:"topicId" bson:"topicId"`
	CreatorID     primitive.ObjectID `json:"creatorId" bson:"creatorId"`
	Title         string             `json:"title" bson:"title"`
	Description   string             `json:"description" bson:"description"`
	Likes         []threads.Like     `json:"likes" bson:"likes"`
	ImageURL      string             `json:"imageURL" bson:"imageURL"`
	SuspendStatus string             `json:"suspendStatus,omitempty" bson:"suspendStatus,omitempty"`
	SuspendDetail string             `json:"suspendDetail,omitempty" bson:"suspendDetail,omitempty"`
	CreatedAt     primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt     primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *threads.Domain) *Model {
	return &Model{
		Id:            domain.Id,
		TopicID:       domain.TopicID,
		CreatorID:     domain.CreatorID,
		Title:         domain.Title,
		Description:   domain.Description,
		Likes:         domain.Likes,
		ImageURL:      domain.ImageURL,
		SuspendStatus: domain.SuspendStatus,
		SuspendDetail: domain.SuspendDetail,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}
}

func (thread *Model) ToDomain() threads.Domain {
	return threads.Domain{
		Id:            thread.Id,
		TopicID:       thread.TopicID,
		CreatorID:     thread.CreatorID,
		Title:         thread.Title,
		Description:   thread.Description,
		Likes:         thread.Likes,
		ImageURL:      thread.ImageURL,
		SuspendStatus: thread.SuspendStatus,
		SuspendDetail: thread.SuspendDetail,
		CreatedAt:     thread.CreatedAt,
		UpdatedAt:     thread.UpdatedAt,
	}
}

func ToArrayDomain(data []Model) []threads.Domain {
	var result []threads.Domain
	for _, v := range data {
		result = append(result, v.ToDomain())
	}
	return result
}
