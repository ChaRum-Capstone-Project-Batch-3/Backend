package follow_threads

import (
	followthreads "charum/business/follow_threads"
	"charum/helper"
	"charum/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowThreadController struct {
	followThreadUseCase followthreads.UseCase
}

func NewFollowThreadController(followThreadUC followthreads.UseCase) *FollowThreadController {
	return &FollowThreadController{
		followThreadUseCase: followThreadUC,
	}
}

/*
Create
*/

func (ftc *FollowThreadController) Create(c echo.Context) error {
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	domain := followthreads.Domain{
		UserID:   userID,
		ThreadID: threadID,
	}

	result, err := ftc.followThreadUseCase.Create(&domain)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "Success to follow thread",
		Data:    result,
	})
}

/*
Read
*/

/*
Update
*/

/*
Delete
*/
