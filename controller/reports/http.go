package reports

import (
	"charum/business/reports"
	"charum/business/threads"
	"charum/business/users"
	"charum/helper"
	"charum/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ReportController struct {
	ReportUseCase reports.UseCase
	userUseCase   users.UseCase
	threadUseCase threads.UseCase
}

func NewReportController(rUC reports.UseCase, userUC users.UseCase, threadUC threads.UseCase) *ReportController {
	return &ReportController{
		ReportUseCase: rUC,
		userUseCase:   userUC,
		threadUseCase: threadUC,
	}
}

func (ctrl *ReportController) Create(c echo.Context) error {
	ReportedID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    nil,
		})
	}

	user, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Data:    nil,
		})
	}

	domain := reports.Domain{
		ReportedID: ReportedID,
		UserID:     user,
	}

	report, err := ctrl.ReportUseCase.Create(&domain)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "already reported" {
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
		Message: "success create report",
		Data:    report,
	})
}
