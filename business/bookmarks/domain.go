package bookmarks

import (
	"charum/dto/bookmarks"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	ThreadID  primitive.ObjectID `json:"threadID" bson:"threadID"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByUserIDAndThreadID(UserID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
	GetAllByUserID(UserID primitive.ObjectID) ([]Domain, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	// Delete
	Delete(domain *Domain) error
	DeleteAllByUserID(userID primitive.ObjectID) error
	DeleteAllByThreadID(threadID primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetAllByUserID(userID primitive.ObjectID) ([]Domain, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	CheckBookmarkedThread(userID primitive.ObjectID, threadID primitive.ObjectID) (bool, error)
	DomainToResponse(domain Domain, userID primitive.ObjectID) (bookmarks.Response, error)
	DomainsToResponseArray(domains []Domain, userID primitive.ObjectID) ([]bookmarks.Response, error)
	// Delete
	Delete(domain *Domain) (Domain, error)
	DeleteAllByUserID(userID primitive.ObjectID) error
	DeleteAllByThreadID(threadID primitive.ObjectID) error
}
