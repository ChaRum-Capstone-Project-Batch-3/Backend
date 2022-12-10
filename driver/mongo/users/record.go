package users

import (
	"charum/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	Email             string             `json:"email" bson:"email"`
	UserName          string             `json:"userName" bson:"userName"`
	DisplayName       string             `json:"displayName" bson:"displayName"`
	Biodata           string             `json:"biodata" bson:"biodata,omitempty"`
	SocialMedia       string             `json:"socialMedia" bson:"socialMedia,omitempty"`
	ProfilePictureURL string             `json:"profilePictureURL" bson:"profilePictureURL,omitempty"`
	Password          string             `json:"password" bson:"password"`
	IsActive          bool               `json:"isActive" bson:"isActive"`
	Role              string             `json:"role" bson:"role"`
	CreatedAt         primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt         primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *users.Domain) *Model {
	return &Model{
		Id:                domain.Id,
		Email:             domain.Email,
		UserName:          domain.UserName,
		DisplayName:       domain.DisplayName,
		Biodata:           domain.Biodata,
		SocialMedia:       domain.SocialMedia,
		ProfilePictureURL: domain.ProfilePictureURL,
		Password:          domain.Password,
		IsActive:          domain.IsActive,
		Role:              domain.Role,
		CreatedAt:         domain.CreatedAt,
		UpdatedAt:         domain.UpdatedAt,
	}
}

func (user *Model) ToDomain() users.Domain {
	return users.Domain{
		Id:                user.Id,
		Email:             user.Email,
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		Biodata:           user.Biodata,
		SocialMedia:       user.SocialMedia,
		ProfilePictureURL: user.ProfilePictureURL,
		Password:          user.Password,
		IsActive:          user.IsActive,
		Role:              user.Role,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

func ToArrayDomain(data []Model) []users.Domain {
	var array []users.Domain
	for _, v := range data {
		array = append(array, v.ToDomain())
	}
	return array
}
