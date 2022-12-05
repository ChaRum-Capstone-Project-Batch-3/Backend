package threads

import (
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	dtoThread "charum/dto/threads"

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
	GetManyWithPagination(query dtoQuery.Request, domain *Domain) ([]Domain, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetAllByTopicID(topicID primitive.ObjectID) ([]Domain, error)
	GetAllByUserID(userID primitive.ObjectID) ([]Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	SuspendByUserID(domain *Domain) error
	// Delete
	Delete(id primitive.ObjectID) error
	DeleteAllByUserID(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetAllByTopicID(topicID primitive.ObjectID) ([]Domain, error)
	GetAllByUserID(userID primitive.ObjectID) ([]Domain, error)
	DomainToResponse(domain Domain) (dtoThread.Response, error)
	DomainsToResponseArray(domains []Domain) ([]dtoThread.Response, error)
	// Update
	UserUpdate(domain *Domain) (Domain, error)
	AdminUpdate(domain *Domain) (Domain, error)
	SuspendByUserID(userID primitive.ObjectID) error
	// Delete
	Delete(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error)
	DeleteAllByUserID(id primitive.ObjectID) error
	DeleteByThreadID(threadID primitive.ObjectID) error
	AdminDelete(threadID primitive.ObjectID) (Domain, error)
}
