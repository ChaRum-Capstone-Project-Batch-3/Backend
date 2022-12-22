package forgot_password

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Email     string             `json:"email"`
	Token     string             `json:"token"`
	Password  string             `json:"-"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt"`
	ExpiredAt primitive.DateTime `json:"expiredAt"`
	IsUsed    bool               `json:"isUsed"`
}

type Repository interface {
	// Create
	Generate(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByToken(token string) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Generate(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByToken(token string) (Domain, error)
	ValidateToken(token string) (Domain, error)
	// Update
	UpdatePassword(domain *Domain) (Domain, error)
}
