package threads

import (
	"charum/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	TopicID       primitive.ObjectID `json:"topicID" bson:"topicID"`
	CreatorID     primitive.ObjectID `json:"creatorID" bson:"creatorID"`
	Title         string             `json:"title" bson:"title"`
	Description   string             `json:"description" bson:"description"`
	Likes         []Like             `json:"likes" bson:"likes"`
	SuspendStatus string             `json:"suspendStatus,omitempty" bson:"suspendStatus"`
	SuspendDetail string             `json:"suspendDetail,omitempty" bson:"suspendDetail"`
	CreatedAt     primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt     primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
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
	SuspendByUserID(domain *Domain) error
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	DomainToResponse(domain Domain) (dto.ResponseThread, error)
	DomainsToResponseArray(domains []Domain) ([]dto.ResponseThread, error)
	// Update
	Update(domain *Domain) (Domain, error)
	SuspendByUserID(userID primitive.ObjectID) error
	// Delete
	Delete(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
}
