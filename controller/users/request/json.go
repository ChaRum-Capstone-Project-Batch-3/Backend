package request

import (
	"charum/business/users"
	"charum/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type UserRegister struct {
	Email       string `json:"email" validate:"required,email" bson:"email"`
	UserName    string `json:"userName" validate:"required" bson:"userName"`
	DisplayName string `json:"displayName" validate:"required" bson:"displayName"`
	Password    string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" bson:"password"`
}

func (req *UserRegister) ToDomain() *users.Domain {
	return &users.Domain{
		Email:       req.Email,
		UserName:    req.UserName,
		DisplayName: req.DisplayName,
		Password:    req.Password,
	}
}

func (u *UserRegister) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(u); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(u)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				// Get field tag
				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

type Login struct {
	Email    string `json:"email" validate:"required,email" bson:"email"`
	Password string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" bson:"password"`
}

func (req *Login) ToDomain() *users.Domain {
	return &users.Domain{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (u *Login) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(u); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(u)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				// Get field tag
				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}
