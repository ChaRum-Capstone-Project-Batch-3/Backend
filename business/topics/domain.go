package topics

import (
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Topic       string             `json:"topic" bson:"topic"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"imageURL" bson:"imageURL"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetManyWithPagination(query dtoQuery.Request, domain *Domain) ([]Domain, int, error)
	GetByTopic(topic string) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain, image *multipart.FileHeader) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error)
	GetByTopic(topic string) (Domain, error)
	// Update
	Update(domain *Domain, image *multipart.FileHeader) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) (Domain, error)
}
