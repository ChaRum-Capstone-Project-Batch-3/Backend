package topics

import (
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/business/threads"
	"charum/business/topics"
	"charum/controller/topics/request"
	"charum/controller/topics/response"
	"charum/helper"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicController struct {
	TopicUseCase        topics.UseCase
	ThreadUseCase       threads.UseCase
	CommentUseCase      comments.UseCase
	FollowThreadUseCase followThreads.UseCase
}

func NewTopicController(topicUC topics.UseCase, threadUC threads.UseCase, commentUC comments.UseCase, followThreadUC followThreads.UseCase) *TopicController {
	return &TopicController{
		TopicUseCase:        topicUC,
		ThreadUseCase:       threadUC,
		CommentUseCase:      commentUC,
		FollowThreadUseCase: followThreadUC,
	}
}

func (topicCtrl *TopicController) Create(c echo.Context) error {
	userInput := request.Topic{}

	if c.Bind(&userInput) != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "fill all the required fields and make sure data type is correct",
			Data:    nil,
		})
	}

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	topic, err := topicCtrl.TopicUseCase.CreateTopic(userInput.ToDomain())
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
			Message: "invalid id",
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

func (topicCtrl *TopicController) GetAll(c echo.Context) error {
	topic, err := topicCtrl.TopicUseCase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get all topics",
		Data: map[string]interface{}{
			"topics": response.FromDomainArray(topic),
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
			Message: "invalid id",
			Data:    nil,
		})
	}

	userInput := request.Topic{}
	if c.Bind(&userInput) != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "fill all the required fields and make sure data type is correct",
			Data:    nil,
		})
	}

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	topic, err := topicCtrl.TopicUseCase.UpdateTopic(topicID, userInput.ToDomain())
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
			Message: "invalid id",
			Data:    nil,
		})
	}

	deletedTopic, err := topicCtrl.TopicUseCase.DeleteTopic(topicID)
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
			return c.JSON(http.StatusNotFound, helper.BaseResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
				Data:    nil,
			})
		}

		err = topicCtrl.CommentUseCase.DeleteAllByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusNotFound, helper.BaseResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
				Data:    nil,
			})
		}

		err = topicCtrl.FollowThreadUseCase.DeleteAllByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusNotFound, helper.BaseResponse{
				Status:  http.StatusNotFound,
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
