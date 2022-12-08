package bookmarks

import (
	"charum/business/bookmarks"
	"charum/controller/bookmarks/request"
	"charum/controller/bookmarks/response"
	"charum/helper"
	"charum/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type BookmarkController struct {
	bookmarkUseCase bookmarks.UseCase
}

func NewBookmarkController(bu bookmarks.UseCase) *BookmarkController {
	return &BookmarkController{
		bookmarkUseCase: bu,
	}
}

func (bc *BookmarkController) AddBookmark(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("thread_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	bookmarkInput := request.Bookmark{}
	c.Bind(&bookmarkInput)

	if err := bookmarkInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	result, err := bc.bookmarkUseCase.AddBookmark(userID, threadID, bookmarkInput.ToDomain())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "success add to bookmark",
		Data: map[string]interface{}{
			"bookmark": response.FromDomain(result),
		},
	})
}

// get bookmark by id
func (bc *BookmarkController) GetByID(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}
	//
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread_id"))
	result, err := bc.bookmarkUseCase.GetByID(userID, threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get bookmark by id",
		Data: map[string]interface{}{
			"bookmark": response.FromDomain(result),
		},
	})
}

func (bc *BookmarkController) GetAllBookmark(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	result, err := bc.bookmarkUseCase.GetAllBookmark(userID)
	if err != nil {
		return c.JSON(http.StatusCreated, helper.BaseResponse{
			Status:  http.StatusCreated,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseBookmark, err := bc.bookmarkUseCase.DomainsToResponseArray(result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get all bookmark",
		Data: map[string]interface{}{
			"bookmarks": responseBookmark,
		},
	})
}

/*
Delete
*/
func (bc *BookmarkController) DeleteBookmark(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	_, err = bc.bookmarkUseCase.DeleteBookmark(userID, threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success delete bookmark",
		Data:    nil,
	})
}