package threads_test

import (
	"charum/business/threads"
	_threadMock "charum/business/threads/mocks"
	"charum/business/topics"
	_topicMock "charum/business/topics/mocks"
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	dtoThread "charum/dto/threads"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	threadRepository _threadMock.Repository
	topicRepository  _topicMock.Repository
	userRepository   _userMock.Repository
	threadUseCase    threads.UseCase
	topicDomain      topics.Domain
	threadDomain     threads.Domain
	userDomain       users.Domain
)

func TestMain(m *testing.M) {
	threadUseCase = threads.NewThreadUseCase(&threadRepository, &topicRepository, &userRepository)

	userDomain = users.Domain{
		Id:          primitive.NewObjectID(),
		Email:       "email",
		UserName:    "username",
		DisplayName: "displayname",
		Password:    "password",
		IsActive:    true,
		Role:        "user",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	topicDomain = topics.Domain{
		Id:          primitive.NewObjectID(),
		Topic:       "topic",
		Description: "description",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	threadDomain = threads.Domain{
		Id:          primitive.NewObjectID(),
		TopicID:     primitive.NewObjectID(),
		CreatorID:   primitive.NewObjectID(),
		Title:       "Test Thread",
		Description: "Test Thread Description",
		Likes:       []threads.Like{},
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test case 1 | Valid create thread", func(t *testing.T) {
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("Create", &threadDomain).Return(threadDomain, nil).Once()

		result, err := threadUseCase.Create(&threadDomain)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid create thread | Topic Not Exist", func(t *testing.T) {
		copyTopic := topicDomain
		copyTopic.Topic = "topic not exist"
		expectedErr := errors.New("failed to get topic")

		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.Create(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test case 3 | Invalid create thread | Error when creating thread", func(t *testing.T) {
		expectedErr := errors.New("failed to create thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, nil).Once()
		threadRepository.On("Create", &threadDomain).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.Create(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetWithSortAndOrder(t *testing.T) {
	t.Run("Test case 1 | Valid get thread with sort and order", func(t *testing.T) {
		threadRepository.On("GetWithSortAndOrder", 0, 2, "createdAt", -1).Return([]threads.Domain{threadDomain}, 1, nil).Once()

		result, totalPage, totalData, err := threadUseCase.GetWithSortAndOrder(1, 2, "createdAt", "desc")

		assert.NotNil(t, result)
		assert.NotZero(t, totalPage)
		assert.NotZero(t, totalData)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get thread with sort and order | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get threads")
		threadRepository.On("GetWithSortAndOrder", 0, 2, "createdAt", 1).Return([]threads.Domain{}, 0, expectedErr).Once()

		result, totalPage, totalData, err := threadUseCase.GetWithSortAndOrder(1, 2, "createdAt", "asc")

		assert.Equal(t, []threads.Domain{}, result)
		assert.Zero(t, totalPage)
		assert.Zero(t, totalData)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("Test case 1 | Valid get thread by id", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()

		result, err := threadUseCase.GetByID(threadDomain.Id)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get thread by id | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.GetByID(threadDomain.Id)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestDomainToResponse(t *testing.T) {
	t.Run("Test case 1 | Valid domain to response", func(t *testing.T) {
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()

		result, actualErr := threadUseCase.DomainToResponse(threadDomain)

		assert.NotNil(t, result)
		assert.Nil(t, actualErr)
	})

	t.Run("Test case 2 | Invalid domain to response | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get creator")
		userRepository.On("GetByID", threadDomain.CreatorID).Return(users.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.DomainToResponse(threadDomain)

		assert.Equal(t, dtoThread.Response{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test case 3 | Invalid domain to response | Error when getting topic", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.DomainToResponse(threadDomain)

		assert.Equal(t, dtoThread.Response{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestDomainToResponseArray(t *testing.T) {
	t.Run("Test case 1 | Valid domain to response array", func(t *testing.T) {
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()

		result, actualErr := threadUseCase.DomainsToResponseArray([]threads.Domain{threadDomain})

		assert.NotNil(t, result)
		assert.Nil(t, actualErr)
	})

	t.Run("Test case 2 | Invalid domain to response array | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		userRepository.On("GetByID", threadDomain.CreatorID).Return(users.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.DomainsToResponseArray([]threads.Domain{threadDomain})

		assert.Equal(t, []dtoThread.Response{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test case 1 | Valid update thread", func(t *testing.T) {
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threadDomain, nil).Once()

		result, err := threadUseCase.Update(&threadDomain)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid update thread | Error when updating thread", func(t *testing.T) {
		expectedErr := errors.New("failed to update thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.Update(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 3 | Invalid update thread | Topic Not Exist", func(t *testing.T) {
		copyTopic := topicDomain
		copyTopic.Topic = "topic not exist"
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, err := threadUseCase.Update(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 4 | Invalid update thread | Thread Not Exist", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.Update(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 5 | Invalid update thread | User Not Exist", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")

		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(users.Domain{}, expectedErr).Once()

		result, err := threadUseCase.Update(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 6 | Invalid update thread | User Not Creator", func(t *testing.T) {
		copyThread := threadDomain
		copyThread.CreatorID = primitive.NewObjectID()
		expectedErr := errors.New("user are not the thread creator")

		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(copyThread, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()

		result, err := threadUseCase.Update(&threadDomain)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestSuspendByUserID(t *testing.T) {
	t.Run("Test case 1 | Valid suspend thread by user id", func(t *testing.T) {
		threadRepository.On("SuspendByUserID", mock.Anything).Return(nil).Once()

		err := threadUseCase.SuspendByUserID(threadDomain.CreatorID)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid suspend thread by user id | Error when getting thread by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to suspend user threads")
		threadRepository.On("SuspendByUserID", mock.Anything).Return(expectedErr).Once()

		err := threadUseCase.SuspendByUserID(threadDomain.CreatorID)

		assert.Equal(t, expectedErr, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test case 1 | Valid delete thread", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(nil).Once()

		thread, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Nil(t, err)
		assert.Equal(t, threadDomain, thread)
	})

	t.Run("Test case 2 | Invalid delete thread | Error when deleting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to delete thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(expectedErr).Once()

		thread, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, err, expectedErr)
		assert.Equal(t, threads.Domain{}, thread)
	})

	t.Run("Test case 3 | Invalid delete thread | Thread Not Exist", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		thread, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, err, expectedErr)
		assert.Equal(t, threads.Domain{}, thread)
	})

	t.Run("Test case 4 | Invalid delete thread | User Not Exist", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", threadDomain.CreatorID).Return(users.Domain{}, expectedErr).Once()

		thread, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, err, expectedErr)
		assert.Equal(t, threads.Domain{}, thread)
	})

	t.Run("Test case 5 | Invalid delete thread | User Not Creator", func(t *testing.T) {
		expectedErr := errors.New("user are not the thread creator")

		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()

		thread, err := threadUseCase.Delete(userDomain.Id, threadDomain.Id)

		assert.Equal(t, err, expectedErr)
		assert.Equal(t, threads.Domain{}, thread)
	})
}
