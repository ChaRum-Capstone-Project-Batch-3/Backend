package users

import (
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	Email             string             `json:"email" bson:"email"`
	UserName          string             `json:"userName" bson:"userName"`
	DisplayName       string             `json:"displayName" bson:"displayName"`
	Biodata           string             `json:"biodata" bson:"biodata"`
	SocialMedia       string             `json:"socialMedia" bson:"socialMedia"`
	Password          string             `json:"-"`
	OldPassword       string             `json:"-"`
	NewPassword       string             `json:"-"`
	IsActive          bool               `json:"isActive" bson:"isActive"`
	Role              string             `json:"role" bson:"role"`
	CreatedAt         primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt         primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	ProfilePictureURL string             `json:"profilePictureURL" bson:"profilePictureURL"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByEmail(email string) (Domain, error)
	GetByUsername(username string) (Domain, error)
	GetManyWithPagination(query dtoQuery.Request, domain *Domain) ([]Domain, int, error)
	GetAll() ([]Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	UpdatePassword(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Register(domain *Domain, profilePicture *multipart.FileHeader) (Domain, string, error)
	// Read
	Login(key string, password string) (Domain, string, error)
	GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetAll() (int, error)
	// Update
	UpdatePassword(domain *Domain) (Domain, error)
	Update(domain *Domain, profilePicture *multipart.FileHeader) (Domain, error)
	Suspend(id primitive.ObjectID) (Domain, error)
	Unsuspend(id primitive.ObjectID) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) (Domain, error)
}
