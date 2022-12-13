package request

import (
	"charum/business/users"
	"charum/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Register struct {
	Email       string `json:"email" validate:"required,email" bson:"email" form:"email"`
	UserName    string `json:"userName" validate:"required" bson:"userName" form:"userName"`
	DisplayName string `json:"displayName" validate:"required" bson:"displayName" form:"displayName"`
	Biodata     string `json:"biodata" bson:"biodata" form:"biodata"`
	SocialMedia string `json:"socialMedia" bson:"socialMedia" form:"socialMedia"`
	Password    string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" bson:"password" form:"password"`
}

func (req *Register) ToDomain() *users.Domain {
	return &users.Domain{
		Email:       req.Email,
		UserName:    req.UserName,
		DisplayName: req.DisplayName,
		Biodata:     req.Biodata,
		SocialMedia: req.SocialMedia,
		Password:    req.Password,
	}
}

func (req *Register) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
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
	Key      string `json:"key" validate:"required" bson:"key" form:"key"`
	Password string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" bson:"password" form:"password"`
}

func (req *Login) ToDomain() *users.Domain {
	return &users.Domain{
		Email:    req.Key,
		Password: req.Password,
	}
}

func (req *Login) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
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

type Update struct {
	Email       string `json:"email" validate:"required,email" bson:"email" form:"email"`
	UserName    string `json:"userName" validate:"required" bson:"userName"  form:"userName"`
	DisplayName string `json:"displayName" validate:"required" bson:"displayName"  form:"displayName"`
	Biodata     string `json:"biodata" bson:"biodata"  form:"Biodata"`
	SocialMedia string `json:"socialMedia" bson:"socialMedia"  form:"socialMedia"`
}

func (req *Update) ToDomain() *users.Domain {
	return &users.Domain{
		Email:       req.Email,
		UserName:    req.UserName,
		DisplayName: req.DisplayName,
		Biodata:     req.Biodata,
		SocialMedia: req.SocialMedia,
	}
}

func (req *Update) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
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

type ChangePassword struct {
	OldPassword string `json:"oldPassword" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" bson:"password"`
	NewPassword string `json:"newPassword" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" bson:"password"`
}

func (req *ChangePassword) ToDomain() *users.Domain {
	return &users.Domain{
		NewPassword: req.NewPassword,
		OldPassword: req.OldPassword,
	}
}

func (req *ChangePassword) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
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
