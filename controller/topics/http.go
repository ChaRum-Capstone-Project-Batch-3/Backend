package topics

import (
	"charum/business/bookmarks"
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/business/threads"
	"charum/business/topics"
	"charum/controller/topics/request"
	"charum/controller/topics/response"
	dtoPagination "charum/dto/pagination"
	"charum/helper"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicController struct {
	TopicUseCase        topics.UseCase
	ThreadUseCase       threads.UseCase
	CommentUseCase      comments.UseCase
	FollowThreadUseCase followThreads.UseCase
	bookmarkUseCase     bookmarks.UseCase
}

func NewTopicController(topicUC topics.UseCase, threadUC threads.UseCase, commentUC comments.UseCase, followThreadUC followThreads.UseCase, bookmarkUC bookmarks.UseCase) *TopicController {
	return &TopicController{
		TopicUseCase:        topicUC,
		ThreadUseCase:       threadUC,
		CommentUseCase:      commentUC,
		FollowThreadUseCase: followThreadUC,
		bookmarkUseCase:     bookmarkUC,
	}
}

func (topicCtrl *TopicController) Create(c echo.Context) error {
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

	userInput := request.Topic{}
	c.Bind(&userInput)

	inputErr := userInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	topic, err := topicCtrl.TopicUseCase.Create(userInput.ToDomain(), image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errors.New("topic already exist") {
			statusCode = http.StatusConflict
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "success create topic",
		Data: map[string]interface{}{
			"topic": response.FromDomain(topic),
		},
	})
}

func (topicCtrl *TopicController) GetByID(c echo.Context) error {
	topicID, err := primitive.ObjectIDFromHex(c.Param("topic-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	topic, err := topicCtrl.TopicUseCase.GetByID(topicID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get topic by id",
		Data: map[string]interface{}{
			"topic": response.FromDomain(topic),
		},
	})
}

func (topicCtrl *TopicController) GetManyWithPagination(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "page must be a number",
			Data:       nil,
			Pagination: helper.Page{},
		})
	} else if page < 1 {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "page must be greater than 0",
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "25"
	}
	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "limit must be a number and greater than 0",
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "createdAt"
	} else if !(sort == "_id" || sort == "topic" || sort == "createdAt" || sort == "updatedAt") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "sort must be _id, topic, createdAt, or updatedAt",
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	order := c.QueryParam("order")
	if order == "" {
		order = "desc"
	} else if !(order == "asc" || order == "desc") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "order must be asc or desc",
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	userInputDomain := topics.Domain{
		Topic: c.QueryParam("topic"),
	}

	pagination := dtoPagination.Request{
		Page:  page,
		Limit: limitNumber,
		Sort:  sort,
		Order: order,
	}

	users, totalPage, totalData, err := topicCtrl.TopicUseCase.GetManyWithPagination(pagination, &userInputDomain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get topics",
		Data: map[string]interface{}{
			"topics": response.FromDomainArray(users),
		},
		Pagination: helper.Page{
			Size:        limitNumber,
			TotalData:   totalData,
			TotalPage:   totalPage,
			CurrentPage: page,
		},
	})
}

// get by topic name
func (topicCtrl *TopicController) GetByTopic(c echo.Context) error {
	topicName := c.Param("topic")

	topic, err := topicCtrl.TopicUseCase.GetByTopic(topicName)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get topic by name",
		Data: map[string]interface{}{
			"topic": response.FromDomain(topic),
		},
	})
}

/*
Update
*/

func (topicCtrl *TopicController) Update(c echo.Context) error {
	topicID, err := primitive.ObjectIDFromHex(c.Param("topic-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
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

	userInput := request.Topic{}
	c.Bind(&userInput)

	inputErr := userInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	userInputDomain := userInput.ToDomain()
	userInputDomain.Id = topicID

	topic, err := topicCtrl.TopicUseCase.Update(userInputDomain, image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errors.New("topic already exist") {
			statusCode = http.StatusConflict
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success update topic",
		Data: map[string]interface{}{
			"topic": response.FromDomain(topic),
		},
	})
}

/*
Delete
*/

func (topicCtrl *TopicController) Delete(c echo.Context) error {
	topicID, err := primitive.ObjectIDFromHex(c.Param("topic-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	deletedTopic, err := topicCtrl.TopicUseCase.Delete(topicID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threads, err := topicCtrl.ThreadUseCase.GetAllByTopicID(topicID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	for _, thread := range threads {
		err := topicCtrl.ThreadUseCase.DeleteByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		err = topicCtrl.CommentUseCase.DeleteAllByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		err = topicCtrl.FollowThreadUseCase.DeleteAllByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		err = topicCtrl.bookmarkUseCase.DeleteAllByThreadID(thread.Id)
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
		Message: "success delete topic",
		Data: map[string]interface{}{
			"topic": response.FromDomain(deletedTopic),
		},
	})
}
