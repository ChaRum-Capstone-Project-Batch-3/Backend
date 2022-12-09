package bookmarks

import (
	"charum/business/bookmarks"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	ThreadID  primitive.ObjectID `json:"threadID" bson:"threadID"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *bookmarks.Domain) *Model {
	return &Model{
		Id:        domain.Id,
		UserID:    domain.UserID,
		ThreadID:  domain.ThreadID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (bookmark *Model) ToDomain() bookmarks.Domain {
	return bookmarks.Domain{
		Id:        bookmark.Id,
		UserID:    bookmark.UserID,
		ThreadID:  bookmark.ThreadID,
		CreatedAt: bookmark.CreatedAt,
		UpdatedAt: bookmark.UpdatedAt,
	}
}

func ToDomainArray(data []Model) []bookmarks.Domain {
	var result []bookmarks.Domain
	for _, v := range data {
		result = append(result, v.ToDomain())
	}
	return result
}
