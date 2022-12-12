package follow_threads_test

import (
	"charum/business/comments"
	_commentMock "charum/business/comments/mocks"
	followThreads "charum/business/follow_threads"
	_followThreadMock "charum/business/follow_threads/mocks"
	"charum/business/threads"
	_threadMock "charum/business/threads/mocks"
	"charum/business/topics"
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	dtoFollowThread "charum/dto/follow_threads"
	dtoThreads "charum/dto/threads"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	followThreadRepositoryMock _followThreadMock.Repository
	userRepositoryMock         _userMock.Repository
	threadRepositoryMock       _threadMock.Repository
	commentRepositoryMock      _commentMock.Repository
	threadUseCaseMock          _threadMock.UseCase
	followThreadUseCase        followThreads.UseCase
	followThreadDomain         followThreads.Domain
	topicDomain                topics.Domain
	responseThread             dtoThreads.Response
	userDomain                 users.Domain
	threadDomain               threads.Domain
	commentDomain              comments.Domain
)

func TestMain(m *testing.M) {
	followThreadUseCase = followThreads.NewFollowThreadUseCase(&followThreadRepositoryMock, &userRepositoryMock, &threadRepositoryMock, &commentRepositoryMock, &threadUseCaseMock)

	followThreadDomain = followThreads.Domain{
		Id:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ThreadID:  primitive.NewObjectID(),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	userDomain = users.Domain{
		Id:          followThreadDomain.UserID,
		Email:       "test@test.com",
		Password:    "test",
		UserName:    "test",
		DisplayName: "test",
		IsActive:    true,
		Role:        "user",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	threadDomain = threads.Domain{
		Id:            followThreadDomain.ThreadID,
		TopicID:       primitive.NewObjectID(),
		Title:         "test",
		Description:   "test",
		Likes:         []threads.Like{},
		SuspendStatus: "",
		SuspendDetail: "",
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
	}

	commentDomain = comments.Domain{
		Id:        primitive.NewObjectID(),
		ThreadID:  followThreadDomain.ThreadID,
		UserID:    userDomain.Id,
		Comment:   "test",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	topicDomain = topics.Domain{
		Id:          threadDomain.TopicID,
		Topic:       "test",
		Description: "test",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	responseThread = dtoThreads.Response{
		Id:            threadDomain.Id,
		Topic:         topicDomain,
		Creator:       userDomain,
		Title:         threadDomain.Title,
		Description:   threadDomain.Description,
		Likes:         []dtoThreads.Like{},
		TotalFollow:   0,
		TotalComment:  0,
		TotalLike:     len(threadDomain.Likes),
		SuspendStatus: threadDomain.SuspendStatus,
		SuspendDetail: threadDomain.SuspendDetail,
		CreatedAt:     threadDomain.CreatedAt,
		UpdatedAt:     threadDomain.UpdatedAt,
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Case 1 | Valid Create", func(t *testing.T) {
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, errors.New("not found")).Once()
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		followThreadRepositoryMock.On("Create", mock.Anything).Return(followThreadDomain, nil).Once()

		result, err := followThreadUseCase.Create(&followThreadDomain)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Case 2 | Invalid Create | User Not Found", func(t *testing.T) {
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, errors.New("not found")).Once()
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(users.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.Create(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})

	t.Run("Test Case 3 | Invalid Create | Thread Not Found", func(t *testing.T) {
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, errors.New("not found")).Once()
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threads.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.Create(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})

	t.Run("Test Case 4 | Invalid Create | Already Followed", func(t *testing.T) {
		expectedErr := errors.New("user already follow this thread")
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, nil).Once()

		result, err := followThreadUseCase.Create(&followThreadDomain)

		assert.Equal(t, expectedErr, err)
		assert.Equal(t, result, followThreads.Domain{})
	})

	t.Run("Test Case 5 | Invalid Create | Repository Error", func(t *testing.T) {
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, errors.New("not found")).Once()
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		followThreadRepositoryMock.On("Create", mock.Anything).Return(followThreads.Domain{}, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.Create(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})
}

func TestGetAllByUserID(t *testing.T) {
	t.Run("Test Case 1 | Valid Get All By User ID", func(t *testing.T) {
		followThreadRepositoryMock.On("GetAllByUserID", followThreadDomain.UserID).Return([]followThreads.Domain{followThreadDomain}, nil).Once()

		result, err := followThreadUseCase.GetAllByUserID(followThreadDomain.UserID)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Case 2 | Invalid Get All By User ID | Repository Error", func(t *testing.T) {
		followThreadRepositoryMock.On("GetAllByUserID", followThreadDomain.UserID).Return([]followThreads.Domain{}, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.GetAllByUserID(followThreadDomain.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, result, []followThreads.Domain{})
	})
}

func TestCountByThreadID(t *testing.T) {
	t.Run("Test Case 1 | Valid Count By Thread ID", func(t *testing.T) {
		followThreadRepositoryMock.On("CountByThreadID", followThreadDomain.ThreadID).Return(1, nil).Once()

		result, err := followThreadUseCase.CountByThreadID(followThreadDomain.ThreadID)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Case 2 | Invalid Count By Thread ID | Repository Error", func(t *testing.T) {
		followThreadRepositoryMock.On("CountByThreadID", followThreadDomain.ThreadID).Return(0, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.CountByThreadID(followThreadDomain.ThreadID)

		assert.NotNil(t, err)
		assert.Equal(t, result, 0)
	})
}

func TestCheckFollowedThread(t *testing.T) {
	t.Run("Test Case 1 | Valid Check Followed Thread", func(t *testing.T) {
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, nil).Once()

		result, err := followThreadUseCase.CheckFollowedThread(followThreadDomain.UserID, followThreadDomain.ThreadID)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Case 2 | Invalid Check Followed Thread | Repository Error", func(t *testing.T) {
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreads.Domain{}, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.CheckFollowedThread(followThreadDomain.UserID, followThreadDomain.ThreadID)

		assert.Nil(t, err)
		assert.Equal(t, result, false)
	})
}

func TestDomainToResponse(t *testing.T) {
	t.Run("Test Case 1 | Valid Domain To Response", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCaseMock.On("DomainToResponse", mock.Anything, mock.Anything).Return(responseThread, nil).Once()
		commentRepositoryMock.On("CountByThreadID", followThreadDomain.ThreadID).Return(1, nil).Once()

		result, err := followThreadUseCase.DomainToResponse(followThreadDomain, userDomain.Id)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Case 2 | Invalid Domain To Response | User Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(users.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.DomainToResponse(followThreadDomain, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, dtoFollowThread.Response{})
	})

	t.Run("Test Case 3 | Invalid Domain To Response | Thread Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threads.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.DomainToResponse(followThreadDomain, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, dtoFollowThread.Response{})
	})

	t.Run("Test Case 4 | Invalid Domain To Response | Comment Count Error", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCaseMock.On("DomainToResponse", mock.Anything, mock.Anything).Return(responseThread, nil).Once()
		commentRepositoryMock.On("CountByThreadID", followThreadDomain.ThreadID).Return(0, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.DomainToResponse(followThreadDomain, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, dtoFollowThread.Response{})
	})

	t.Run("Test Case 5 | Invalid Domain To Response | Thread Use Case Error", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCaseMock.On("DomainToResponse", mock.Anything, mock.Anything).Return(dtoThreads.Response{}, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.DomainToResponse(followThreadDomain, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, dtoFollowThread.Response{})
	})
}

func TestDomainToResponseArray(t *testing.T) {
	t.Run("Test Case 1 | Valid Domain To Response Array", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCaseMock.On("DomainToResponse", mock.Anything, mock.Anything).Return(responseThread, nil).Once()
		commentRepositoryMock.On("CountByThreadID", followThreadDomain.ThreadID).Return(1, nil).Once()

		result, err := followThreadUseCase.DomainToResponseArray([]followThreads.Domain{followThreadDomain}, userDomain.Id)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Case 2 | Invalid Domain To Response Array | User Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(users.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.DomainToResponseArray([]followThreads.Domain{followThreadDomain}, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, []dtoFollowThread.Response{})
	})

	t.Run("Test Case 3 | Invalid Domain To Response Array | Thread Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threads.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.DomainToResponseArray([]followThreads.Domain{followThreadDomain}, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, []dtoFollowThread.Response{})
	})

	t.Run("Test Case 4 | Invalid Domain To Response Array | Comment Count Error", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCaseMock.On("DomainToResponse", mock.Anything, mock.Anything).Return(responseThread, nil).Once()
		commentRepositoryMock.On("CountByThreadID", followThreadDomain.ThreadID).Return(0, errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.DomainToResponseArray([]followThreads.Domain{followThreadDomain}, userDomain.Id)

		assert.NotNil(t, err)
		assert.Equal(t, result, []dtoFollowThread.Response{})
	})
}

func TestUpdateNotification(t *testing.T) {
	t.Run("Test Case 1 | Valid Update Notification", func(t *testing.T) {
		followThreadRepositoryMock.On("AddOneNotification", followThreadDomain.Id).Return(nil).Once()

		err := followThreadUseCase.UpdateNotification(followThreadDomain.Id)

		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Update Notification", func(t *testing.T) {
		followThreadRepositoryMock.On("AddOneNotification", followThreadDomain.Id).Return(errors.New("unexpected error")).Once()

		err := followThreadUseCase.UpdateNotification(followThreadDomain.Id)

		assert.NotNil(t, err)
	})
}

func TestResetNotification(t *testing.T) {
	t.Run("Test Case 1 | Valid Reset Notification", func(t *testing.T) {
		followThreadRepositoryMock.On("ResetNotification", followThreadDomain.ThreadID, followThreadDomain.UserID).Return(nil).Once()

		err := followThreadUseCase.ResetNotification(followThreadDomain.ThreadID, followThreadDomain.UserID)

		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Reset Notification", func(t *testing.T) {
		followThreadRepositoryMock.On("ResetNotification", followThreadDomain.ThreadID, followThreadDomain.UserID).Return(errors.New("unexpected error")).Once()

		err := followThreadUseCase.ResetNotification(followThreadDomain.ThreadID, followThreadDomain.UserID)

		assert.NotNil(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, nil).Once()
		followThreadRepositoryMock.On("Delete", mock.Anything).Return(nil).Once()

		result, err := followThreadUseCase.Delete(&followThreadDomain)

		assert.Nil(t, err)
		assert.Equal(t, result, followThreadDomain)
	})

	t.Run("Test Case 2 | Invalid Delete | User Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(users.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.Delete(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})

	t.Run("Test Case 3 | Invalid Delete | Thread Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threads.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.Delete(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})

	t.Run("Test Case 4 | Invalid Delete | Follow Thread Not Found", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreads.Domain{}, errors.New("not found")).Once()

		result, err := followThreadUseCase.Delete(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})

	t.Run("Test Case 5 | Invalid Delete | Delete Error", func(t *testing.T) {
		userRepositoryMock.On("GetByID", followThreadDomain.UserID).Return(userDomain, nil).Once()
		threadRepositoryMock.On("GetByID", followThreadDomain.ThreadID).Return(threadDomain, nil).Once()
		followThreadRepositoryMock.On("GetByUserIDAndThreadID", followThreadDomain.UserID, followThreadDomain.ThreadID).Return(followThreadDomain, nil).Once()
		followThreadRepositoryMock.On("Delete", mock.Anything).Return(errors.New("unexpected error")).Once()

		result, err := followThreadUseCase.Delete(&followThreadDomain)

		assert.NotNil(t, err)
		assert.Equal(t, result, followThreads.Domain{})
	})
}

func TestDeleteAllByUserID(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete All By User ID", func(t *testing.T) {
		followThreadRepositoryMock.On("DeleteAllByUserID", followThreadDomain.UserID).Return(nil).Once()

		err := followThreadUseCase.DeleteAllByUserID(followThreadDomain.UserID)

		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Delete All By User ID", func(t *testing.T) {
		followThreadRepositoryMock.On("DeleteAllByUserID", followThreadDomain.UserID).Return(errors.New("unexpected error")).Once()

		err := followThreadUseCase.DeleteAllByUserID(followThreadDomain.UserID)

		assert.NotNil(t, err)
	})
}

func TestDeleteAllByThreadID(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete All By Thread ID", func(t *testing.T) {
		followThreadRepositoryMock.On("DeleteAllByThreadID", followThreadDomain.ThreadID).Return(nil).Once()

		err := followThreadUseCase.DeleteAllByThreadID(followThreadDomain.ThreadID)

		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Delete All By Thread ID", func(t *testing.T) {
		followThreadRepositoryMock.On("DeleteAllByThreadID", followThreadDomain.ThreadID).Return(errors.New("unexpected error")).Once()

		err := followThreadUseCase.DeleteAllByThreadID(followThreadDomain.ThreadID)

		assert.NotNil(t, err)
	})
}
