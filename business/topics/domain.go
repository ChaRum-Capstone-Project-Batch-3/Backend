package topics

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id          primitive.ObjectID
	Topic       string
	Description string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
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
