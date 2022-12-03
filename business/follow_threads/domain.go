package follow_threads

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByUserIDAndThreadID(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
	// Update
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	// Update
	// Delete
	Delete(domain *Domain) (Domain, error)
}
