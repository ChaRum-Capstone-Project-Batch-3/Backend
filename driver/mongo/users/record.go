package users

import (
	"charum/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Email       string             `json:"email" bson:"email"`
	UserName    string             `json:"userName" bson:"userName"`
	DisplayName string             `json:"displayName" bson:"displayName"`
	Password    string             `json:"password" bson:"password"`
	BookmarkId  primitive.ObjectID `json:"bookmarkId" bson:"bookmarkId"`
	IsActive    bool               `json:"isActive" bson:"isActive"`
	Role        string             `json:"role" bson:"role"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *users.Domain) *Model {
	return &Model{
		Id:          domain.Id,
		Email:       domain.Email,
		UserName:    domain.UserName,
		DisplayName: domain.DisplayName,
		Password:    domain.Password,
		BookmarkId:  domain.BookmarkID,
		IsActive:    domain.IsActive,
		Role:        domain.Role,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func (user *Model) ToDomain() users.Domain {
	return users.Domain{
		Id:          user.Id,
		Email:       user.Email,
		UserName:    user.UserName,
		DisplayName: user.DisplayName,
		Password:    user.Password,
		BookmarkID:  user.BookmarkId,
		IsActive:    user.IsActive,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func ToArrayDomain(data []Model) []users.Domain {
	var array []users.Domain
	for _, v := range data {
		array = append(array, v.ToDomain())
	}
	return array
}
