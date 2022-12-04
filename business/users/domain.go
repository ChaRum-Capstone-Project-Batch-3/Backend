package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Email       string             `json:"email" bson:"email"`
	UserName    string             `json:"userName" bson:"userName"`
	DisplayName string             `json:"displayName" bson:"displayName"`
	Biodata     string             `json:"biodata" bson:"biodata"`
	SocialMedia string             `json:"socialMedia" bson:"socialMedia"`
	Password    string             `json:"-"`
	IsActive    bool               `json:"isActive" bson:"isActive"`
	Role        string             `json:"role" bson:"role"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByEmail(email string) (Domain, error)
	GetByUsername(username string) (Domain, error)
	GetWithSortAndOrder(skip int, limit int, sort string, order int) ([]Domain, int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Register(domain *Domain) (Domain, string, error)
	// Read
	Login(domain *Domain) (Domain, string, error)
	GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	Suspend(id primitive.ObjectID) (Domain, error)
	Unsuspend(id primitive.ObjectID) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) (Domain, error)
}
