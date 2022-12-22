package users

import (
	"charum/business/bookmarks"
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/business/threads"
	"charum/business/users"
	"charum/controller/users/request"
	"charum/controller/users/response"
	dtoPagination "charum/dto/pagination"
	"charum/helper"
	"charum/util"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userUseCase         users.UseCase
	threadUseCase       threads.UseCase
	commentUseCase      comments.UseCase
	followThreadUseCase followThreads.UseCase
	bookmarksUseCase    bookmarks.UseCase
}

func NewUserController(userUC users.UseCase, threadUC threads.UseCase, commentUC comments.UseCase, followThreadUC followThreads.UseCase, bookmarkUC bookmarks.UseCase) *UserController {
	return &UserController{
		userUseCase:         userUC,
		threadUseCase:       threadUC,
		commentUseCase:      commentUC,
		followThreadUseCase: followThreadUC,
		bookmarksUseCase:    bookmarkUC,
	}
}

/*
Create
*/

func (userCtrl *UserController) Register(c echo.Context) error {
	var validationErr []helper.ValidationError
	profilePicture, _ := c.FormFile("profilePicture")
	if profilePicture != nil {
		profilePictureExt := filepath.Ext(profilePicture.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if profilePictureExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "profilePicture",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if profilePicture.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "profilePicture",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	userInput := request.Register{}
	c.Bind(&userInput)

	if strings.Contains(userInput.UserName, " ") {
		validationErr = append(validationErr, helper.ValidationError{
			Field:   "userName",
			Message: "This field must not contain any space",
		})
	}

	inputErr := userInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	user, token, err := userCtrl.userUseCase.Register(userInput.ToDomain(), profilePicture)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already") {
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
		Message: "success to register",
		Data: map[string]interface{}{
			"token": token,
			"user":  response.FromDomain(user),
		},
	})
}

/*
Read
*/

func (userCtrl *UserController) Login(c echo.Context) error {
	userInput := request.Login{}
	c.Bind(&userInput)

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	_, token, err := userCtrl.userUseCase.Login(userInput.Key, userInput.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to login",
		Data: map[string]interface{}{
			"token": token,
		},
	})
}

func (userCtrl *UserController) GetManyWithPagination(c echo.Context) error {
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
	} else if !(sort == "_id" || sort == "email" || sort == "userName" || sort == "displayName" || sort == "createdAt" || sort == "updatedAt") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "sort must be _id, email, userName, displayName, createdAt, or updatedAt",
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

	userInputDomain := users.Domain{
		Email:       c.QueryParam("email"),
		UserName:    c.QueryParam("username"),
		DisplayName: c.QueryParam("display-name"),
	}

	pagination := dtoPagination.Request{
		Page:  page,
		Limit: limitNumber,
		Sort:  sort,
		Order: order,
	}

	users, totalPage, totalData, err := userCtrl.userUseCase.GetManyWithPagination(pagination, &userInputDomain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get users",
		Data: map[string]interface{}{
			"users": response.FromDomainArray(users),
		},
		Pagination: helper.Page{
			Size:        limitNumber,
			TotalData:   totalData,
			TotalPage:   totalPage,
			CurrentPage: page,
		},
	})
}

func (userCtrl *UserController) GetByID(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data:    nil,
		})
	}

	user, err := userCtrl.userUseCase.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get user by id",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

func (userCtrl *UserController) GetProfile(c echo.Context) error {
	id, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	user, err := userCtrl.userUseCase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get user profile",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

/*
Update
*/

func (userCtrl *UserController) AdminUpdate(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
			Data:    nil,
		})
	}

	var validationErr []helper.ValidationError
	profilePicture, _ := c.FormFile("profilePicture")
	if profilePicture != nil {
		profilePictureExt := filepath.Ext(profilePicture.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if profilePictureExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "profilePicture",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if profilePicture.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "profilePicture",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	userInput := request.Update{}
	c.Bind(&userInput)

	if strings.Contains(userInput.UserName, " ") {
		validationErr = append(validationErr, helper.ValidationError{
			Field:   "userName",
			Message: "This field must not contain any space",
		})
	}

	inputErr := userInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.Id = userID
	user, err := userCtrl.userUseCase.Update(userDomain, profilePicture)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errors.New("failed to get user") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "already") {
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
		Message: "success to update user",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

func (userCtrl *UserController) UserUpdate(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var validationErr []helper.ValidationError
	profilePicture, _ := c.FormFile("profilePicture")
	if profilePicture != nil {
		profilePictureExt := filepath.Ext(profilePicture.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if profilePictureExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "profilePicture",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if profilePicture.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "profilePicture",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	userInput := request.Update{}
	c.Bind(&userInput)

	if strings.Contains(userInput.UserName, " ") {
		validationErr = append(validationErr, helper.ValidationError{
			Field:   "userName",
			Message: "This field must not contain any space",
		})
	}

	inputErr := userInput.Validate()
	if inputErr != nil {
		validationErr = append(validationErr, inputErr...)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    validationErr,
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.Id = userID
	user, err := userCtrl.userUseCase.Update(userDomain, profilePicture)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errors.New("failed to get user") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "already") {
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
		Message: "success to update user",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

func (userCtrl *UserController) UpdatePassword(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	userInput := request.ChangePassword{}
	c.Bind(&userInput)

	if err := userInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "validation failed",
			Data:    err,
		})
	}

	userDomain := userInput.ToDomain()
	userDomain.Id = userID
	user, err := userCtrl.userUseCase.UpdatePassword(userDomain)

	statusCode := http.StatusInternalServerError
	if err == errors.New("failed to get user") {
		statusCode = http.StatusNotFound
	} else if err == errors.New("password is incorrect") {
		statusCode = http.StatusForbidden
	}

	if err != nil {
		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to change password",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

func (userCtrl *UserController) Suspend(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
			Data:    nil,
		})
	}

	user, err := userCtrl.userUseCase.Suspend(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errors.New("failed to get user") {
			statusCode = http.StatusNotFound
		} else if err == errors.New("user is already suspended") {
			statusCode = http.StatusConflict
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.threadUseCase.SuspendByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.commentUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.followThreadUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.threadUseCase.RemoveUserFromAllLikes(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.bookmarksUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to suspend user",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

func (userCtrl *UserController) Unsuspend(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
			Data:    nil,
		})
	}

	user, err := userCtrl.userUseCase.Unsuspend(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user is not suspended") {
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
		Message: "success to unsuspend user",
		Data: map[string]interface{}{
			"user": response.FromDomain(user),
		},
	})
}

/*
Delete
*/

func (userCtrl *UserController) Delete(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
			Data:    nil,
		})
	}

	deletedUser, err := userCtrl.userUseCase.Delete(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errors.New("failed to get user") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.commentUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.followThreadUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.threadUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.threadUseCase.RemoveUserFromAllLikes(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = userCtrl.bookmarksUseCase.DeleteAllByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to delete user",
		Data: map[string]interface{}{
			"user": response.FromDomain(deletedUser),
		},
	})
}
