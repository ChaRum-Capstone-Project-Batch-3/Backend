package response

import (
	"charum/business/comments"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	ThreadID  primitive.ObjectID `json:"threadID" bson:"threadID"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	ParentID  primitive.ObjectID `json:"parentID,omitempty" bson:"parentID,omitempty"`
	Comment   string             `json:"comment" bson:"commment"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain comments.Domain) Comment {
	return Comment{
		Id:        domain.Id,
		ThreadID:  domain.ThreadID,
		UserID:    domain.UserID,
		ParentID:  domain.ParentID,
		Comment:   domain.Comment,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func FromDomainArray(domain []comments.Domain) []Comment {
	var result []Comment
	for _, v := range domain {
		result = append(result, FromDomain(v))
	}
	return result
}
