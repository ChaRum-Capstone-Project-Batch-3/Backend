package comments

import (
	"charum/business/comments"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	ThreadID  primitive.ObjectID `json:"threadID" bson:"threadID"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	ParentID  primitive.ObjectID `json:"parentID,omitempty" bson:"parentID,omitempty"`
	Comment   string             `json:"comment" bson:"commment"`
	ImageURL  string             `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *comments.Domain) *Model {
	return &Model{
		Id:        domain.Id,
		ThreadID:  domain.ThreadID,
		UserID:    domain.UserID,
		ParentID:  domain.ParentID,
		Comment:   domain.Comment,
		ImageURL:  domain.ImageURL,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (comment *Model) ToDomain() comments.Domain {
	return comments.Domain{
		Id:        comment.Id,
		ThreadID:  comment.ThreadID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Comment:   comment.Comment,
		ImageURL:  comment.ImageURL,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func ToDomainArray(data []Model) []comments.Domain {
	var result []comments.Domain
	for _, v := range data {
		result = append(result, v.ToDomain())
	}
	return result
}
