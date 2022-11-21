package request

import (
	"charum/business/users"

	"github.com/go-playground/validator/v10"
)

type UserRegister struct {
	Email       string `json:"email" validate:"required,email" bson:"email"`
	UserName    string `json:"userName" validate:"required" bson:"userName"`
	DisplayName string `json:"displayName" validate:"required" bson:"displayName"`
	Password    string `json:"password" validate:"required" bson:"password"`
}

func (req *UserRegister) ToDomain() *users.Domain {
	return &users.Domain{
		Email:       req.Email,
		UserName:    req.UserName,
		DisplayName: req.DisplayName,
		Password:    req.Password,
	}
}

func (req *UserRegister) Validate() error {
	validate := validator.New()
	err := validate.Struct(req)
	return err
}

type Login struct {
	Email    string `json:"email" validate:"required,email" bson:"email"`
	Password string `json:"password" validate:"required" bson:"password"`
}

func (req *Login) ToDomain() *users.Domain {
	return &users.Domain{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *Login) Validate() error {
	validate := validator.New()
	err := validate.Struct(req)
	return err
}
