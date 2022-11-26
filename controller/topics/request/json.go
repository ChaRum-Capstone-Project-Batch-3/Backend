package request

import (
	"charum/business/topics"
	"charum/helper"
	"errors"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Create struct {
	Topic       string `json:"topic" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (req *Create) ToDomain() *topics.Domain {
	return &topics.Domain{
		Topic:       req.Topic,
		Description: req.Description,
	}
}

func (req *Create) Validate() []helper.ValidationError {
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
	Topic       string `json:"topic"`
	Description string `json:"description"`
}

func (req *Update) ToDomain() *topics.Domain {
	return &topics.Domain{
		Topic:       req.Topic,
		Description: req.Description,
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
