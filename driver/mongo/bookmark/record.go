package bookmark

import (
	"charum/business/bookmarks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bookmark struct {
	Id        primitive.ObjectID   `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID   `json:"userId" bson:"userId"`
	Threads   []primitive.ObjectID `json:"threads" bson:"threads"`
	CreatedAt primitive.DateTime   `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime   `json:"updatedAt" bson:"updatedAt"`
}

type ThreadReturn struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Title string             `json:"title" bson:"title"`
}

func FromDomain(domain *bookmarks.Domain) *Bookmark {
	return &Bookmark{
		Id:        domain.Id,
		UserID:    domain.UserID,
		Threads:   domain.Threads,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (bookmark *Bookmark) ToDomain() bookmarks.Domain {
	return bookmarks.Domain{
		Id:        bookmark.Id,
		UserID:    bookmark.UserID,
		Threads:   bookmark.Threads,
		CreatedAt: bookmark.CreatedAt,
		UpdatedAt: bookmark.UpdatedAt,
	}
}