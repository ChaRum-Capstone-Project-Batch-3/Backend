package forgot_password

import (
	"charum/business/forgot_password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	ExpiredAt primitive.DateTime `json:"expiredAt" bson:"expiredAt"`
	IsUsed    bool               `json:"isUsed" bson:"isUsed"`
}

func FromDomain(domain *forgot_password.Domain) *Model {
	return &Model{
		Id:        domain.Id,
		Email:     domain.Email,
		Token:     domain.Token,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		ExpiredAt: domain.ExpiredAt,
		IsUsed:    domain.IsUsed,
	}
}

func (user *Model) ToDomain() forgot_password.Domain {
	return forgot_password.Domain{
		Id:        user.Id,
		Email:     user.Email,
		Token:     user.Token,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ExpiredAt: user.ExpiredAt,
		IsUsed:    user.IsUsed,
	}
}
