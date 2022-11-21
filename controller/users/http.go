package users

import (
	"charum/businesses/users"
	"charum/controller/users/request"
	"charum/controller/users/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUseCase users.UseCase
}

func NewUserController(userUC users.UseCase) *UserController {
	return &UserController{
		userUseCase: userUC,
	}
}

/*
Create
*/

func (userCtrl *UserController) UserRegister(c echo.Context) error {
	userInput := request.UserRegister{}

	if c.Bind(&userInput) != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	if userInput.Validate() != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "validation failed",
		})
	}

	user, token, err := userCtrl.userUseCase.UserRegister(userInput.ToDomain())
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "success to register",
		"token":   token,
		"user":    response.FromDomain(user),
	})
}

/*
Read
*/

func (userCtrl *UserController) Login(c echo.Context) error {
	userInput := request.Login{}

	if c.Bind(&userInput) != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid request",
		})
	}

	if userInput.Validate() != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "validation failed",
		})
	}

	token, err := userCtrl.userUseCase.Login(userInput.ToDomain())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success to login",
		"token":   token,
	})
}

/*
Update
*/

/*
Delete
*/
