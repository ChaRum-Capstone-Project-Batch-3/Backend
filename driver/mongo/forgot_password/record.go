package forgot_password

import (
	"charum/business/forgot_password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Email     string             `json:"email"`
	Token     string             `json:"token"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt"`
	ExpiredAt primitive.DateTime `json:"expiredAt"`
	IsUsed    bool               `json:"isUsed"`
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
