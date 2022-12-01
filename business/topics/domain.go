package topics

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Topic       string             `json:"topic" bson:"topic"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Repository interface {
	// Create
	CreateTopic(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetAll() ([]Domain, error)
	GetByTopic(topic string) (Domain, error)
	// Update
	UpdateTopic(domain *Domain) (Domain, error)
	// Delete
	DeleteTopic(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	CreateTopic(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetAll() ([]Domain, error)
	GetByTopic(topic string) (Domain, error)
	// Update
	UpdateTopic(id primitive.ObjectID, domain *Domain) (Domain, error)
	// Delete
	DeleteTopic(id primitive.ObjectID) (Domain, error)
}
