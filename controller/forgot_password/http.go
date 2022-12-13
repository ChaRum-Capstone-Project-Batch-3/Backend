package forgot_password

import (
	"charum/business/forgot_password"
	"charum/business/users"
	"charum/controller/forgot_password/request"
	"charum/helper"
	"github.com/labstack/echo/v4"
	"net/http"
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

	forgotPassword, err := ctrl.forgotPasswordUseCase.Generate(userInput.ToDomain())
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    forgotPassword,
	})
}

func (ctrl *ForgotPasswordController) Update(c echo.Context) error {
	// get token from params and validate
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

	user, err := ctrl.forgotPasswordUseCase.UpdatePassword(userInput.ToDomain())
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    user,
	})
}

func (ctrl *ForgotPasswordController) ValidateToken(c echo.Context) error {
	token := c.Param("token")
	forgotPassword, err := ctrl.forgotPasswordUseCase.ValidateToken(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "token invalid/not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "token valid",
		Data:    forgotPassword,
	})
}
