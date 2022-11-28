package threads

import (
	"charum/business/threads"
	"charum/controller/threads/request"
	"charum/controller/threads/response"
	"charum/helper"
	"charum/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThreadController struct {
	threadUseCase threads.UseCase
}

func NewThreadController(threadUC threads.UseCase) *ThreadController {
	return &ThreadController{
		threadUseCase: threadUC,
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

	result, err := tc.threadUseCase.Create(userID, threadInput.Topic, threadInput.ToDomain())
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
			"thread": response.FromDomain(result),
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
	} else if !(sort == "id" || sort == "title" || sort == "createdAt" || sort == "updatedAt") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "sort must be id, title, createdAt, or updatedAt",
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

	threads, totalPage, err := tc.threadUseCase.GetWithSortAndOrder(page, limitNumber, sort, order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get threads",
		Data: map[string]interface{}{
			"totalPage": totalPage,
			"threads":   response.FromDomainArray(threads),
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

	result, err := tc.threadUseCase.GetByID(threadID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get thread",
		Data: map[string]interface{}{
			"thread": response.FromDomain(result),
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

	result, err := tc.threadUseCase.Update(userID, threadID, threadInput.Topic, threadInput.ToDomain())

	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "you are not the thread creator") {
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
			"thread": response.FromDomain(result),
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
		} else if strings.Contains(err.Error(), "you are not the thread creator") {
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
			"thread": response.FromDomain(deletedThread),
		},
	})
}
