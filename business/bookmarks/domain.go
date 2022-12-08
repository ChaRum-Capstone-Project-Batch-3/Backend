package bookmarks

import (
	"charum/dto"
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
	// Update
	UpdateBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error)
	// Delete
	DeleteBookmark(userID primitive.ObjectID, threadID primitive.ObjectID) error
}

type UseCase interface {
	// Create
	AddBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error)
	// Read
	GetByID(UserID primitive.ObjectID, ThreadID primitive.ObjectID) (Domain, error)
	GetAllBookmark(userID primitive.ObjectID) ([]Domain, error)
	DomainToResponse(domain Domain) (dto.ResponseBookmark, error)
	DomainsToResponseArray(domains []Domain) ([]dto.ResponseBookmark, error)
	// Update
	UpdateBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error)
	// Delete
	DeleteBookmark(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
}
