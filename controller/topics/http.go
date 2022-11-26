package topics

import (
	"charum/business/topics"
	"charum/controller/topics/request"
	"charum/controller/topics/response"
	"charum/helper"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type TopicController struct {
	TopicUseCase topics.UseCase
}

func NewTopicController(topicUC topics.UseCase) *TopicController {
	return &TopicController{
		TopicUseCase: topicUC,
	}
}

func (topicCtrl *TopicController) CreateTopic(c echo.Context) error {
	userInput := request.Create{}

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
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
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
	topicID, err := primitive.ObjectIDFromHex(c.Param("id"))
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

/*
Update
*/

func (topicCtrl *TopicController) UpdateTopic(c echo.Context) error {
	topicID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    nil,
		})
	}

	userInput := request.Create{}
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

	statusCode := http.StatusInternalServerError
	if err != nil {
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

func (topicCtrl *TopicController) DeleteTopic(c echo.Context) error {
	topicID, err := primitive.ObjectIDFromHex(c.Param("id"))
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

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success delete topic",
		Data: map[string]interface{}{
			"topic": response.FromDomain(deletedTopic),
		},
	})
}
