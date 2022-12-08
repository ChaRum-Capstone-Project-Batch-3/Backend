package bookmarks

import (
	"charum/dto/bookmarks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id        primitive.ObjectID
	UserID    primitive.ObjectID
	ThreadID  primitive.ObjectID
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
}

type Repository interface {
	// Create
	AddBookmark(domain *Domain) (Domain, error)
	// Read
	GetByID(UserID primitive.ObjectID, ThreadID primitive.ObjectID) (Domain, error)
	GetAllBookmark(UserID primitive.ObjectID) ([]Domain, error)
	// Delete
	DeleteBookmark(userID primitive.ObjectID, threadID primitive.ObjectID) error
}

type UseCase interface {
	// Create
	AddBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error)
	// Read
	GetByID(UserID primitive.ObjectID, ThreadID primitive.ObjectID) (Domain, error)
	GetAllBookmark(userID primitive.ObjectID) ([]Domain, error)
	DomainToResponse(domain Domain) (bookmarks.Response, error)
	DomainsToResponseArray(domains []Domain) ([]bookmarks.Response, error)
	// Delete
	DeleteBookmark(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
}
