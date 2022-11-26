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
	GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	// Update
	Update(id primitive.ObjectID, domain *Domain) (Domain, error)
	Suspend(id primitive.ObjectID) (Domain, error)
	Unsuspend(id primitive.ObjectID) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) (Domain, error)
}
