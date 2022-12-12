package bookmarks

import (
	"charum/business/bookmarks"
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/helper"
	"charum/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarkController struct {
	bookmarkUseCase     bookmarks.UseCase
	followThreadUseCase followThreads.UseCase
	commentUseCase      comments.UseCase
}

func NewBookmarkController(bUC bookmarks.UseCase, ftUC followThreads.UseCase, cUC comments.UseCase) *BookmarkController {
	return &BookmarkController{
		bookmarkUseCase:     bUC,
		followThreadUseCase: ftUC,
		commentUseCase:      cUC,
	}
}

func (bc *BookmarkController) Create(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
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

	userInputDomain := bookmarks.Domain{
		UserID:   userID,
		ThreadID: threadID,
	}

	result, err := bc.bookmarkUseCase.Create(&userInputDomain)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already") {
			statusCode = http.StatusConflict
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	response, err := bc.bookmarkUseCase.DomainToResponse(result, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "success to craete bookmark",
		Data: map[string]interface{}{
			"bookmark": response,
		},
	})
}

func (bc *BookmarkController) GetAllByToken(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	result, err := bc.bookmarkUseCase.GetAllByUserID(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseBookmark, err := bc.bookmarkUseCase.DomainsToResponseArray(result, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	for i, boomark := range responseBookmark {
		responseBookmark[i].Thread.TotalFollow, err = bc.followThreadUseCase.CountByThreadID(boomark.Thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseBookmark[i].Thread.TotalComment, err = bc.commentUseCase.CountByThreadID(boomark.Thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseBookmark[i].Thread.TotalBookmark, err = bc.bookmarkUseCase.CountByThreadID(boomark.Thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseBookmark[i].Thread.IsBookmarked = true
		responseBookmark[i].Thread.IsFollowed, err = bc.followThreadUseCase.CheckFollowedThread(userID, boomark.Thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get bookmarks",
		Data: map[string]interface{}{
			"bookmarks": responseBookmark,
		},
	})
}

/*
Delete
*/

func (bc *BookmarkController) Delete(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
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

	userInputDomain := bookmarks.Domain{
		UserID:   userID,
		ThreadID: threadID,
	}

	_, err = bc.bookmarkUseCase.Delete(&userInputDomain)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to delete bookmark",
		Data:    nil,
	})
}
