package threads

import (
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/business/threads"
	"charum/controller/threads/request"
	"charum/helper"
	"charum/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThreadController struct {
	threadUseCase       threads.UseCase
	commentUseCase      comments.UseCase
	followThreadUseCase followThreads.UseCase
}

func NewThreadController(threadUC threads.UseCase, commentUC comments.UseCase, followThreadUC followThreads.UseCase) *ThreadController {
	return &ThreadController{
		threadUseCase:       threadUC,
		commentUseCase:      commentUC,
		followThreadUseCase: followThreadUC,
	}
}

/*
Create
*/

func (tc *ThreadController) Create(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadInput := request.Thread{}
	c.Bind(&threadInput)

	if err := threadInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	threadDomain := threadInput.ToDomain()
	threadDomain.CreatorID = userID
	threadDomain.TopicID, err = primitive.ObjectIDFromHex(threadInput.TopicID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	result, err := tc.threadUseCase.Create(threadDomain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "success to create thread",
		Data: map[string]interface{}{
			"thread": result,
		},
	})
}

/*
Read
*/

func (tc *ThreadController) GetManyWithPagination(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "page must be a number",
			Data:    nil,
		})
	} else if page < 1 {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "page must be greater than 0",
			Data:    nil,
		})
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "25"
	}
	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "limit must be a number and greater than 0",
			Data:    nil,
		})
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "createdAt"
	} else if !(sort == "id" || sort == "title" || sort == "createdAt" || sort == "updatedAt" || sort == "likes") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "sort must be id, title, createdAt, updatedAt, or likes",
			Data:    nil,
		})
	}

	order := c.QueryParam("order")
	if order == "" {
		order = "desc"
	} else if !(order == "asc" || order == "desc") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "order must be asc or desc",
			Data:    nil,
		})
	}

	threads, totalPage, totalData, err := tc.threadUseCase.GetWithSortAndOrder(page, limitNumber, sort, order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:     http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	responseThreads, err := tc.threadUseCase.DomainsToResponseArray(threads)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:     http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	for i, thread := range responseThreads {
		responseThreads[i].TotalFollow, err = tc.followThreadUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		responseThreads[i].TotalComment, err = tc.commentUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		responseThreads[i].TotalLike = len(thread.Likes)
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get threads",
		Data: map[string]interface{}{
			"threads": responseThreads,
		},
		Pagination: helper.Page{
			Size:        limitNumber,
			TotalData:   totalData,
			TotalPage:   totalPage,
			CurrentPage: page,
		},
	})
}

func (tc *ThreadController) GetByID(c echo.Context) error {
	threadID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	thread, err := tc.threadUseCase.GetByID(threadID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	comment, err := tc.commentUseCase.GetByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	totalFollow, err := tc.followThreadUseCase.CountByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseComment, err := tc.commentUseCase.DomainToResponseArray(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThread, err := tc.threadUseCase.DomainToResponse(thread)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	uid, err := util.GetUIDFromToken(c)
	if err == nil {
		err = tc.followThreadUseCase.ResetNotification(threadID, uid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}
	}

	responseThread.TotalComment = len(responseComment)
	responseThread.TotalFollow = totalFollow

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get thread",
		Data: map[string]interface{}{
			"thread":   responseThread,
			"comments": responseComment,
		},
	})
}

/*
Update
*/

func (tc *ThreadController) Update(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	threadInput := request.Thread{}
	c.Bind(&threadInput)

	if err := threadInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	topicID, err := primitive.ObjectIDFromHex(threadInput.TopicID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	threadDomain := threadInput.ToDomain()
	threadDomain.Id = threadID
	threadDomain.TopicID = topicID
	threadDomain.CreatorID = userID

	result, err := tc.threadUseCase.Update(threadDomain)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user are not the thread creator") {
			statusCode = http.StatusForbidden
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to update thread",
		Data: map[string]interface{}{
			"thread": result,
		},
	})
}

/*
Delete
*/

func (tc *ThreadController) Delete(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, helper.BaseResponse{
			Status:  http.StatusForbidden,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	deletedThread, err := tc.threadUseCase.Delete(userID, threadID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user are not the thread creator") {
			statusCode = http.StatusForbidden
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to delete thread",
		Data: map[string]interface{}{
			"thread": deletedThread,
		},
	})
}
