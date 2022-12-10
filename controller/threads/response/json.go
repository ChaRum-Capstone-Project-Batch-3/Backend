package response

import (
	"charum/business/threads"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thread struct {
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

func FromDomain(domain threads.Domain) Thread {
	return Thread{
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

func FromDomainArray(data []threads.Domain) []Thread {
	var array []Thread
	for _, v := range data {
		array = append(array, FromDomain(v))
	}
	return array
}
