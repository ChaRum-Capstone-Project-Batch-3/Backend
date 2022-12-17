package response

import (
	"charum/business/reports"
	"charum/helper"
	"errors"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Create struct {
	ReportedID primitive.ObjectID `json:"reportedID"`
}

func (req *Create) ToDomain() *reports.Domain {
	return &reports.Domain{
		ReportedID: req.ReportedID,
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
