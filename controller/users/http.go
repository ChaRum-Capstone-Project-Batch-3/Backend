package users

import (
	"charum/business/users"
	"charum/controller/users/request"
	"charum/controller/users/response"
	"net/http"
	"strconv"

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

func (userCtrl *UserController) GetAllUser(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "page must be a number",
		})
	} else if page < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "page must be greater than 0",
		})
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "25"
	}
	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "limit must be more than 0",
		})
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "createdAt"
	} else if !(sort == "id" || sort == "email" || sort == "userName" || sort == "displayName" || sort == "createdAt" || sort == "updatedAt") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "sort must be id, email, userName, displayName, createdAt, or updatedAt",
		})
	}

	order := c.QueryParam("order")
	if order == "" {
		order = "desc"
	} else if !(order == "asc" || order == "desc") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "order must be asc or desc",
		})
	}

	users, totalPage, err := userCtrl.userUseCase.GetUsersWithSortAndOrder(page, limitNumber, sort, order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "success to get all users",
		"totalPage": totalPage,
		"users":     response.FromArrayDomain(users),
	})
}

/*
Update
*/

/*
Delete
*/
