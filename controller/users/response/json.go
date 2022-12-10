package response

import (
	"charum/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	Email             string             `json:"email" bson:"email"`
	UserName          string             `json:"userName" bson:"userName"`
	DisplayName       string             `json:"displayName" bson:"displayName"`
	Biodata           string             `json:"biodata" bson:"biodata"`
	SocialMedia       string             `json:"socialMedia" bson:"socialMedia"`
	ProfilePictureURL string             `json:"profilePictureURL" bson:"profilePictureURL"`
	IsActive          bool               `json:"isActive" bson:"isActive"`
	Role              string             `json:"role" bson:"role"`
	CreatedAt         primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt         primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain users.Domain) User {
	return User{
		Id:                domain.Id,
		Email:             domain.Email,
		UserName:          domain.UserName,
		DisplayName:       domain.DisplayName,
		Biodata:           domain.Biodata,
		SocialMedia:       domain.SocialMedia,
		ProfilePictureURL: domain.ProfilePictureURL,
		IsActive:          domain.IsActive,
		Role:              domain.Role,
		CreatedAt:         domain.CreatedAt,
		UpdatedAt:         domain.UpdatedAt,
	}
}

func FromDomainArray(data []users.Domain) []User {
	var array []User
	for _, v := range data {
		array = append(array, FromDomain(v))
	}
	return array
}
