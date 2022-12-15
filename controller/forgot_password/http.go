package forgot_password

import (
	"charum/business/forgot_password"
	"charum/business/users"
	"charum/controller/forgot_password/request"
	"charum/helper"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ForgotPasswordController struct {
	forgotPasswordUseCase forgot_password.UseCase
	userUseCase           users.UseCase
}

func NewForgotPasswordController(fpUC forgot_password.UseCase, userUC users.UseCase) *ForgotPasswordController {
	return &ForgotPasswordController{
		forgotPasswordUseCase: fpUC,
		userUseCase:           userUC,
	}
}

func (ctrl *ForgotPasswordController) Generate(c echo.Context) error {
	userInput := request.Generate{}
	c.Bind(&userInput)

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	_, err := ctrl.forgotPasswordUseCase.Generate(userInput.ToDomain())
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not registered") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "success to generate token",
		Data:    nil,
	})
}

func (ctrl *ForgotPasswordController) Update(c echo.Context) error {
	token := c.Param("token")
	userInput := request.Update{}
	userInput.Token = token
	c.Bind(&userInput)

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	_, err := ctrl.forgotPasswordUseCase.UpdatePassword(userInput.ToDomain())
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not registered") {
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
		Message: "success to update password",
		Data:    nil,
	})
}

func (ctrl *ForgotPasswordController) ValidateToken(c echo.Context) error {
	token := c.Param("token")
	_, err := ctrl.forgotPasswordUseCase.ValidateToken(token)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "token has") {
			statusCode = http.StatusUnauthorized
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "token valid",
		Data:    nil,
	})
}
