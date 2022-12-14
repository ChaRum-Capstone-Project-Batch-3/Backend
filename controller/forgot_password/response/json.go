package response


import (
	"charum/business/forgot_password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForgotPassword struct {
	Id 	  primitive.ObjectID `json:"_id" bson:"_id"`
	Email string `json:"email" bson:"email"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	ExpiredAt primitive.DateTime `json:"expiredAt" bson:"expiredAt"`
	IsUsed bool `json:"isUsed" bson:"isUsed"`
}

func FromDomain(domain forgot_password.Domain) ForgotPassword {
	return ForgotPassword{
		Id: domain.Id,
		Email: domain.Email,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		ExpiredAt: domain.ExpiredAt,
		IsUsed: domain.IsUsed,
	}
}