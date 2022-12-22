package threads

import (
	"charum/business/bookmarks"
	"charum/business/comments"
	followThreads "charum/business/follow_threads"
	"charum/business/reports"
	"charum/business/threads"
	"charum/business/users"
	"charum/controller/threads/request"
	dtoPagination "charum/dto/pagination"
	dtoThread "charum/dto/threads"
	"charum/helper"
	"charum/util"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThreadController struct {
	threadUseCase       threads.UseCase
	commentUseCase      comments.UseCase
	followThreadUseCase followThreads.UseCase
	userUseCase         users.UseCase
	bookmarkUseCase     bookmarks.UseCase
	reportUseCase       reports.UseCase
}

func NewThreadController(threadUC threads.UseCase, commentUC comments.UseCase, followThreadUC followThreads.UseCase, userUC users.UseCase, bookmarkUC bookmarks.UseCase, reportUC reports.UseCase) *ThreadController {
	return &ThreadController{
		threadUseCase:       threadUC,
		commentUseCase:      commentUC,
		followThreadUseCase: followThreadUC,
		userUseCase:         userUC,
		bookmarkUseCase:     bookmarkUC,
		reportUseCase:       reportUC,
	}
}

/*
Create
*/

func (tc *ThreadController) Create(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var validationErr []helper.ValidationError
	image, _ := c.FormFile("image")
	if image != nil {
		imageExt := filepath.Ext(image.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if imageExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if image.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	threadInput := request.Thread{}
	c.Bind(&threadInput)

	inputErr := threadInput.Validate()
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

	threadDomain := threadInput.ToDomain()
	threadDomain.CreatorID = userID
	threadDomain.TopicID, err = primitive.ObjectIDFromHex(threadInput.TopicID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	result, err := tc.threadUseCase.Create(threadDomain, image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThread, err := tc.threadUseCase.DomainToResponse(result, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, helper.BaseResponse{
		Status:  http.StatusCreated,
		Message: "success to create thread",
		Data: map[string]interface{}{
			"thread": responseThread,
		},
	})
}

/*
Read
*/

func (tc *ThreadController) GetManyWithPagination(c echo.Context) error {
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
	} else if !(sort == "_id" || sort == "title" || sort == "createdAt" || sort == "updatedAt" || sort == "likes") {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "sort must be _id, title, likes, createdAt, or updatedAt",
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

	var topicID primitive.ObjectID
	if c.QueryParam("topic-id") != "" {
		topicID, err = primitive.ObjectIDFromHex(c.QueryParam("topic-id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.BaseResponse{
				Status:     http.StatusBadRequest,
				Message:    "invalid topic id",
				Data:       nil,
				Pagination: helper.Page{},
			})
		}
	}

	userInputDomain := threads.Domain{
		TopicID: topicID,
		Title:   c.QueryParam("title"),
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:     http.StatusBadRequest,
			Message:    "invalid topic id",
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	pagination := dtoPagination.Request{
		Page:  page,
		Limit: limitNumber,
		Sort:  sort,
		Order: order,
	}

	threads, totalPage, totalData, err := tc.threadUseCase.GetManyWithPagination(pagination, &userInputDomain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:     http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	var responseThreads []dtoThread.Response
	uid, err := util.GetUIDFromToken(c)
	if err == nil {
		responseThreads, err = tc.threadUseCase.DomainsToResponseArray(threads, uid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		for i, thread := range responseThreads {
			responseThreads[i].IsFollowed, err = tc.followThreadUseCase.CheckFollowedThread(uid, thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:  http.StatusInternalServerError,
					Message: err.Error(),
					Data:    nil,
				})
			}

			responseThreads[i].IsBookmarked, err = tc.bookmarkUseCase.CheckBookmarkedThread(uid, thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:  http.StatusInternalServerError,
					Message: err.Error(),
					Data:    nil,
				})
			}

			responseThreads[i].TotalFollow, err = tc.followThreadUseCase.CountByThreadID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}

			responseThreads[i].TotalComment, err = tc.commentUseCase.CountByThreadID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}

			responseThreads[i].TotalBookmark, err = tc.bookmarkUseCase.CountByThreadID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}

			responseThreads[i].TotalReported, err = tc.reportUseCase.GetByReportedID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}
		}
	} else {
		responseThreads, err := tc.threadUseCase.DomainsToResponseArray(threads, primitive.NilObjectID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		for i, thread := range responseThreads {
			responseThreads[i].TotalFollow, err = tc.followThreadUseCase.CountByThreadID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}

			responseThreads[i].TotalComment, err = tc.commentUseCase.CountByThreadID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}

			responseThreads[i].TotalBookmark, err = tc.bookmarkUseCase.CountByThreadID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}

			responseThreads[i].TotalReported, err = tc.reportUseCase.GetByReportedID(thread.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Data:       nil,
					Pagination: helper.Page{},
				})
			}
		}
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get threads",
		Data: map[string]interface{}{
			"threads": responseThreads,
		},
		Pagination: helper.Page{
			Size:        limitNumber,
			TotalData:   totalData,
			TotalPage:   totalPage,
			CurrentPage: page,
		},
	})
}

func (tc *ThreadController) GetManyByToken(c echo.Context) error {
	uid, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threads, err := tc.threadUseCase.GetAllByUserID(uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThreads, err := tc.threadUseCase.DomainsToResponseArray(threads, uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	for i, thread := range responseThreads {
		responseThreads[i].TotalFollow, err = tc.followThreadUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].TotalComment, err = tc.commentUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].TotalBookmark, err = tc.bookmarkUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}
		responseThreads[i].TotalReported, err = tc.reportUseCase.GetByReportedID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get threads",
		Data: map[string]interface{}{
			"threads": responseThreads,
		},
	})
}

func (tc *ThreadController) GetByID(c echo.Context) error {
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	thread, err := tc.threadUseCase.GetByID(threadID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
	}

	comment, err := tc.commentUseCase.GetByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	totalFollow, err := tc.followThreadUseCase.CountByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	totalBookmark, err := tc.bookmarkUseCase.CountByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	totalReported, err := tc.reportUseCase.GetByReportedID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseComment, err := tc.commentUseCase.DomainToResponseArray(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	var responseThread dtoThread.Response
	uid, err := util.GetUIDFromToken(c)
	if err == nil {
		responseThread, err = tc.threadUseCase.DomainToResponse(thread, uid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThread.IsFollowed, err = tc.followThreadUseCase.CheckFollowedThread(uid, threadID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThread.IsBookmarked, err = tc.bookmarkUseCase.CheckBookmarkedThread(uid, threadID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		err = tc.followThreadUseCase.ResetNotification(threadID, uid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}
	} else {
		responseThread, err = tc.threadUseCase.DomainToResponse(thread, primitive.NilObjectID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}
	}

	responseThread.TotalComment = len(responseComment)
	responseThread.TotalFollow = totalFollow
	responseThread.TotalBookmark = totalBookmark
	responseThread.TotalReported = totalReported

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get thread",
		Data: map[string]interface{}{
			"thread":   responseThread,
			"comments": responseComment,
		},
	})
}

/*
Update
*/

func (tc *ThreadController) UserUpdate(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	var validationErr []helper.ValidationError
	image, _ := c.FormFile("image")
	if image != nil {
		imageExt := filepath.Ext(image.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if imageExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if image.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	threadInput := request.Thread{}
	c.Bind(&threadInput)

	inputErr := threadInput.Validate()
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

	topicID, err := primitive.ObjectIDFromHex(threadInput.TopicID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	threadDomain := threadInput.ToDomain()
	threadDomain.Id = threadID
	threadDomain.TopicID = topicID
	threadDomain.CreatorID = userID

	result, err := tc.threadUseCase.UserUpdate(threadDomain, image)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user are not the thread creator") {
			statusCode = http.StatusForbidden
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThread, err := tc.threadUseCase.DomainToResponse(result, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to update thread",
		Data: map[string]interface{}{
			"thread": responseThread,
		},
	})
}

func (tc *ThreadController) AdminUpdate(c echo.Context) error {
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	var validationErr []helper.ValidationError
	image, _ := c.FormFile("image")
	if image != nil {
		imageExt := filepath.Ext(image.Filename)
		availableExt := []string{".jpg", ".jpeg", ".png"}

		flagExt := false
		for _, ext := range availableExt {
			if imageExt == ext {
				flagExt = true
			}
		}

		if !flagExt {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with .jpg, .jpeg, or .png extension",
			})
		}

		if image.Size > 10000000 {
			validationErr = append(validationErr, helper.ValidationError{
				Field:   "image",
				Message: "This field must be a file with size less than 10 MB",
			})
		}
	}

	threadInput := request.Thread{}
	c.Bind(&threadInput)

	inputErr := threadInput.Validate()
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

	topicID, err := primitive.ObjectIDFromHex(threadInput.TopicID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid topic id",
			Data:    nil,
		})
	}

	threadDomain := threadInput.ToDomain()
	threadDomain.Id = threadID
	threadDomain.TopicID = topicID

	result, err := tc.threadUseCase.AdminUpdate(threadDomain, image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThread, err := tc.threadUseCase.DomainToResponse(result, primitive.NilObjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to update thread",
		Data: map[string]interface{}{
			"thread": responseThread,
		},
	})
}

func (tc *ThreadController) GetLikedThreadByToken(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threads, err := tc.threadUseCase.GetLikedByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThreads, err := tc.threadUseCase.DomainsToResponseArray(threads, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:     http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
			Pagination: helper.Page{},
		})
	}

	for i, thread := range responseThreads {
		responseThreads[i].IsFollowed, err = tc.followThreadUseCase.CheckFollowedThread(userID, thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].IsBookmarked, err = tc.bookmarkUseCase.CheckBookmarkedThread(userID, thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].TotalFollow, err = tc.followThreadUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		responseThreads[i].TotalComment, err = tc.commentUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		responseThreads[i].TotalBookmark, err = tc.bookmarkUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}

		responseThreads[i].TotalReported, err = tc.reportUseCase.GetByReportedID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get liked thread",
		Data: map[string]interface{}{
			"likedThreads": responseThreads,
		},
	})
}

func (tc *ThreadController) GetLikedThreadByUserID(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Param("user-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
			Data:    nil,
		})
	}

	_, err = tc.userUseCase.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "failed to get user",
			Data:    nil,
		})
	}

	result, err := tc.threadUseCase.GetLikedByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThreads, err := tc.threadUseCase.DomainsToResponseArray(result, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	for i, thread := range responseThreads {
		responseThreads[i].IsFollowed, err = tc.followThreadUseCase.CheckFollowedThread(userID, thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].IsBookmarked, err = tc.bookmarkUseCase.CheckBookmarkedThread(userID, thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].TotalFollow, err = tc.followThreadUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].TotalComment, err = tc.commentUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}

		responseThreads[i].TotalBookmark, err = tc.bookmarkUseCase.CountByThreadID(thread.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
				Status:     http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
				Pagination: helper.Page{},
			})
		}
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to get liked thread",
		Data: map[string]interface{}{
			"likedThreads": responseThreads,
		},
	})
}

func (tc *ThreadController) Like(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.BaseResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	err = tc.threadUseCase.Like(userID, threadID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user already") {
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
		Message: "success to like thread",
		Data:    nil,
	})
}

func (tcc *ThreadController) Unlike(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, helper.BaseResponse{
			Status:  http.StatusForbidden,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	err = tcc.threadUseCase.Unlike(userID, threadID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") || strings.Contains(err.Error(), "user not") {
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
		Message: "success to unlike thread",
		Data:    nil,
	})
}

/*
Delete
*/

func (tc *ThreadController) UserDelete(c echo.Context) error {
	userID, err := util.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, helper.BaseResponse{
			Status:  http.StatusForbidden,
			Message: err.Error(),
			Data:    nil,
		})
	}

	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	deletedThread, err := tc.threadUseCase.Delete(userID, threadID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user are not the thread creator") {
			statusCode = http.StatusForbidden
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThread, err := tc.threadUseCase.DomainToResponse(deletedThread, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = tc.commentUseCase.DeleteAllByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = tc.followThreadUseCase.DeleteAllByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = tc.bookmarkUseCase.DeleteAllByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to delete thread",
		Data: map[string]interface{}{
			"thread": responseThread,
		},
	})
}

func (tc *ThreadController) AdminDelete(c echo.Context) error {
	threadID, err := primitive.ObjectIDFromHex(c.Param("thread-id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid thread id",
			Data:    nil,
		})
	}

	deletedThread, err := tc.threadUseCase.AdminDelete(threadID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "failed to get") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "user are not the thread creator") {
			statusCode = http.StatusForbidden
		}

		return c.JSON(statusCode, helper.BaseResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
	}

	responseThread, err := tc.threadUseCase.DomainToResponse(deletedThread, primitive.NilObjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = tc.commentUseCase.DeleteAllByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = tc.followThreadUseCase.DeleteAllByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = tc.bookmarkUseCase.DeleteAllByThreadID(threadID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, helper.BaseResponse{
		Status:  http.StatusOK,
		Message: "success to delete thread",
		Data: map[string]interface{}{
			"thread": responseThread,
		},
	})
}
