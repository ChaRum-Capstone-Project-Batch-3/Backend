package reports

import (
	"charum/business/reports"
	"charum/business/threads"
	"charum/business/users"
	"charum/helper"
	"charum/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (ctrl *ReportController) ReportUser(c echo.Context) error {
	ReportedID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
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

func (ctrl *ReportController) ReportThread(c echo.Context) error {
	ReportedID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
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

func (ctrl *ReportController) GetUserReportedID(c echo.Context) error {
	ReportedID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    nil,
		})
	}

	reportData, err := ctrl.ReportUseCase.GetByReportedID(ReportedID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get report",
		Data: map[string]interface{}{
			"total_reports": reportData,
		},
	})
}
func (ctrl *ReportController) GetThreadReportedID(c echo.Context) error {
	ReportedID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    nil,
		})
	}

	reportData, err := ctrl.ReportUseCase.GetByReportedID(ReportedID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get report",
		Data: map[string]interface{}{
			"total_reports": reportData,
		},
	})
}
func (ctrl *ReportController) GetAll(c echo.Context) error {
	report, err := ctrl.ReportUseCase.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	reportedUsers, err := ctrl.ReportUseCase.GetAllReportedUsers()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	reportedThreads, err := ctrl.ReportUseCase.GetAllReportedThreads()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get report",
		Data: map[string]interface{}{
			"total_reports":          report,
			"total_reported_users":   reportedUsers,
			"total_reported_threads": reportedThreads,
		},
	})
}

func (ctrl *ReportController) GetAllReportedUsers(c echo.Context) error {
	reportData, err := ctrl.ReportUseCase.GetAllReportedUsers()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get report",
		Data: map[string]interface{}{
			"total_reported_users": reportData,
		},
	})
}

func (ctrl *ReportController) GetAllReportedThreads(c echo.Context) error {
	reportData, err := ctrl.ReportUseCase.GetAllReportedThreads()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get report",
		Data: map[string]interface{}{
			"total_reported_threads": reportData,
		},
	})
}

func (ctrl *ReportController) CountAll(c echo.Context) error {
	reportData, err := ctrl.ReportUseCase.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get all report",
		Data: map[string]interface{}{
			"total_reports": reportData,
		},
	})
}

func (ctrl *ReportController) CountAllData(c echo.Context) error {
	usersData, err := ctrl.userUseCase.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadsData, err := ctrl.threadUseCase.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	reportsData, err := ctrl.ReportUseCase.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success get all statistics",
		Data: map[string]interface{}{
			"total_users":   usersData,
			"total_threads": threadsData,
			"total_reports": reportsData,
		},
	})
}
