package comments

import (
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/controller/comments/request"
	"charum/helper"
	"charum/util"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController struct {
	CommentUseCase      comments.UseCase
	FollowThreadUseCase followThreads.UseCase
}

func NewCommentController(commentUC comments.UseCase, followThreadUC followThreads.UseCase) *CommentController {
	return &CommentController{
		CommentUseCase:      commentUC,
		FollowThreadUseCase: followThreadUC,
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

	var validationErr []helper.ValidationError
	image, _ := c.FormFile("image")
	if image != nil {
		imageExt := filepath.Ext(image.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if imageExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if image.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	commentInput := request.Comment{}
	c.Bind(&commentInput)
	commentInput.UserID = uid
	commentInput.ThreadID = threadID

	inputErr := commentInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	comment, err := cc.CommentUseCase.Create(commentInput.ToDomain(), image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    err,
		})
	}

	err = cc.FollowThreadUseCase.UpdateNotification(comment.ThreadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    err,
		})
	}

	responseComment, err := cc.CommentUseCase.DomainToResponse(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    err,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success create comment",
		Data: map[string]interface{}{
			"comment": responseComment,
		},
	})
}

/*
Read
*/

/*
Update
*/

func (cc *CommentController) Update(c echo.Context) error {
	uid, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	commentID, err := primitive.ObjectIDFromHex(c.Param("comment-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid comment id",
			Data:    nil,
		})
	}

	var validationErr []helper.ValidationError
	image, _ := c.FormFile("image")
	if image != nil {
		imageExt := filepath.Ext(image.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if imageExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if image.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	commentInput := request.Comment{}
	c.Bind(&commentInput)

	inputErr := commentInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	commentDomain := commentInput.ToDomain()
	commentDomain.Id = commentID
	commentDomain.UserID = uid

	comment, err := cc.CommentUseCase.Update(commentDomain, image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "failed to get comment" {
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
		Message: "success update comment",
		Data: map[string]interface{}{
			"comment": comment,
		},
	})
}

/*
Delete
*/

func (cc *CommentController) Delete(c echo.Context) error {
	uid, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	commentID, err := primitive.ObjectIDFromHex(c.Param("comment-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid comment id",
			Data:    nil,
		})
	}

	comment, err := cc.CommentUseCase.Delete(commentID, uid)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "failed to get comment" {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    err,
		})
	}

	responseComment, err := cc.CommentUseCase.DomainToResponse(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    err,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success delete comment",
		Data: map[string]interface{}{
			"comment": responseComment,
		},
	})
}
