package comments

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	// Update
	// Delete
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByThreadID(threadID primitive.ObjectID) ([]Domain, error)
	// Update
	// Delete
}
