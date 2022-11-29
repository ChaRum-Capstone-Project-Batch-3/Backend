package bookmarks

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id        primitive.ObjectID
	UserID    primitive.ObjectID
	Threads   []Thread
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
}

type Thread struct {
	Id    primitive.ObjectID `json:"ThreadId" bson:"ThreadId"`
	Title string             `json:"title" bson:"title"`
}

type Repository interface {
	// Create
	AddBookmark(domain *Domain) (Domain, error)
	GetByID(id primitive.ObjectID) (Domain, error)
}

type UseCase interface {
	AddBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error)
	GetByID(userID primitive.ObjectID) (Domain, error)
}
