package threads

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id            primitive.ObjectID
	TopicID       primitive.ObjectID
	CreatorID     primitive.ObjectID
	Title         string
	Description   string
	Likes         []Like
	SuspendStatus string
	SuspendDetail string
	CreatedAt     primitive.DateTime
	UpdatedAt     primitive.DateTime
}

type Like struct {
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetWithSortAndOrder(skip int, limit int, sort string, order int) ([]Domain, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(creatorID primitive.ObjectID, topicName string, domain *Domain) (Domain, error)
	// Read
	GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	// Update
	Update(creatorID primitive.ObjectID, threadID primitive.ObjectID, topicName string, domain *Domain) (Domain, error)
	// Delete
	Delete(creatorID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
}
