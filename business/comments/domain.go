package comments

import (
	dtoComment "charum/dto/comments"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	ThreadID  primitive.ObjectID `json:"threadID" bson:"threadID"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	ParentID  primitive.ObjectID `json:"parentID,omitempty" bson:"parentID,omitempty"`
	ImageURL  string             `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
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
	GetByIDAndThreadID(id primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
	GetAllByUserID(userID primitive.ObjectID) ([]Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
	DeleteAllByUserID(userID primitive.ObjectID) error
	DeleteAllByThreadID(threadID primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain, image *multipart.FileHeader) (Domain, error)
	// Read
	GetByThreadID(threadID primitive.ObjectID) ([]Domain, error)
	DomainToResponse(comment Domain) (dtoComment.Response, error)
	DomainToResponseArray(comments []Domain) ([]dtoComment.Response, error)
	CountByThreadID(threadID primitive.ObjectID) (int, error)
	// Update
	Update(domain *Domain, image *multipart.FileHeader) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID, userID primitive.ObjectID) (Domain, error)
	DeleteAllByUserID(userID primitive.ObjectID) error
	DeleteAllByThreadID(threadID primitive.ObjectID) error
}
