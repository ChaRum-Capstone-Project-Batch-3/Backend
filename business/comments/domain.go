package comments

import (
	"charum/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	ThreadID  primitive.ObjectID `json:"threadID" bson:"threadID"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	Comment   string             `json:"comment" bson:"commment"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByThreadID(threadID primitive.ObjectID) ([]Domain, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByThreadID(threadID primitive.ObjectID) ([]Domain, error)
	DomainToResponse(comment Domain) (dto.ResponseComment, error)
	DomainToResponseArray(comments []Domain) ([]dto.ResponseComment, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID, userID primitive.ObjectID) (Domain, error)
}
