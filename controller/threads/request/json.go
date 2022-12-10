package request

import (
	"charum/business/threads"
	"charum/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Thread struct {
	TopicID     string `json:"topicID" validate:"required" form:"topicID"`
	Title       string `json:"title" validate:"required" form:"title"`
	Description string `json:"description" validate:"required" form:"description"`
}

func (req *Thread) ToDomain() *threads.Domain {
	return &threads.Domain{
		Title:       req.Title,
		Description: req.Description,
	}
}

func (req *Thread) Validate() []helper.ValidationError {
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
