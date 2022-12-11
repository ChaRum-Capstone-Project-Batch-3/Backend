package request

import (
	"charum/business/comments"
	"charum/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ThreadID primitive.ObjectID `json:"threadID" bson:"threadID" form:"threadID"`
	UserID   primitive.ObjectID `json:"userID" bson:"userID" form:"userID"`
	ParentID primitive.ObjectID `json:"parentID" bson:"parentID" form:"parentID"`
	Comment  string             `json:"comment" validate:"required" bson:"comment" form:"comment"`
}

func (req *Comment) ToDomain() *comments.Domain {
	return &comments.Domain{
		ThreadID: req.ThreadID,
		UserID:   req.UserID,
		ParentID: req.ParentID,
		Comment:  req.Comment,
	}
}

func (req *Comment) Validate() []helper.ValidationError {
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
