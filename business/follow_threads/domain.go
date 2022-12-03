package follow_threads

import (
	"charum/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	UserID       primitive.ObjectID `json:"userID" bson:"userID"`
	ThreadID     primitive.ObjectID `json:"threadID" bson:"threadID"`
	Notification int                `json:"notification" bson:"notification"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetAllByUserID(userID primitive.ObjectID) ([]Domain, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByUserIDAndThreadID(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	// Update
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetAllByUserID(userID primitive.ObjectID) ([]Domain, error)
	DomainToResponse(domain Domain) (dto.ResponseFollowThread, error)
	DomainToResponseArray(domain []Domain) ([]dto.ResponseFollowThread, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	// Update
	// Delete
	Delete(domain *Domain) (Domain, error)
}
