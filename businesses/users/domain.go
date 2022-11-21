package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id          primitive.ObjectID
	Email       string
	UserName    string
	DisplayName string
	Password    string `json:"-"`
	IsActive    bool
	Role        string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetUserByEmail(email string) (Domain, error)
	GetUserByUsername(username string) (Domain, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	UserRegister(domain *Domain) (Domain, string, error)
	// Read
	Login(domain *Domain) (string, error)
	// Update
	// Delete
}
