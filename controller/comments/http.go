package comments

import (
	"charum/business/comments"
	"charum/controller/comments/request"
	"charum/controller/comments/response"
	"charum/helper"
	"charum/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController struct {
	CommentUseCase comments.UseCase
}

func NewCommentController(commentUC comments.UseCase) *CommentController {
	return &CommentController{
		CommentUseCase: commentUC,
	}
}

/*
Create
*/

func (cc *CommentController) Create(c echo.Context) error {
	uid, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	commentInput := request.Comment{}
	c.Bind(&commentInput)
	commentInput.UserID = uid
	commentInput.ThreadID = threadID

	if err := commentInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	comment, err := cc.CommentUseCase.Create(commentInput.ToDomain())
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "failed to get thread" {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    err,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success create comment",
		Data: map[string]interface{}{
			"comment": response.FromDomain(comment),
		},
	})
}

/*
Read
*/

func (cc *CommentController) GetByThreadID(c echo.Context) error {
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	comments, err := cc.CommentUseCase.GetByThreadID(threadID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "failed to get thread" {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    err,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get comments",
		Data: map[string]interface{}{
			"comments": response.FromDomainArray(comments),
		},
	})
}

/*
Update
*/

/*
Delete
*/
