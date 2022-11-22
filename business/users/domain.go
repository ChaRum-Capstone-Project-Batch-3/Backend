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
	GetUserByID(id primitive.ObjectID) (Domain, error)
	GetUserByEmail(email string) (Domain, error)
	GetUserByUsername(username string) (Domain, error)
	GetUsersWithSortAndOrder(skip int, limit int, sort string, order int) ([]Domain, int, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	UserRegister(domain *Domain) (Domain, string, error)
	// Read
	Login(domain *Domain) (Domain, string, error)
	GetUsersWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, error)
	GetUserByID(id primitive.ObjectID) (Domain, error)
	// Update
	// Delete
}
