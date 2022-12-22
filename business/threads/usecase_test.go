package threads_test

import (
	"charum/business/threads"
	_threadMock "charum/business/threads/mocks"
	"charum/business/topics"
	_topicMock "charum/business/topics/mocks"
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	dtoThread "charum/dto/threads"
	_cloudinaryMock "charum/helper/cloudinary/mocks"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	threadRepository     _threadMock.Repository
	topicRepository      _topicMock.Repository
	userRepository       _userMock.Repository
	cloudinaryRepository _cloudinaryMock.Function
	threadUseCase        threads.UseCase
	topicDomain          topics.Domain
	threadDomain         threads.Domain
	userDomain           users.Domain
	image                *multipart.FileHeader
)

func TestMain(m *testing.M) {
	threadUseCase = threads.NewThreadUseCase(&threadRepository, &topicRepository, &userRepository, &cloudinaryRepository)

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
		Likes: []threads.Like{
			{
				UserID:    primitive.NewObjectID(),
				Timestamp: primitive.NewDateTimeFromTime(time.Now()),
			},
		},
		ImageURL:  "image",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	image = &multipart.FileHeader{}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test case 1 | Valid create thread", func(t *testing.T) {
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("Create", &threadDomain).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()

		result, err := threadUseCase.Create(&threadDomain, image)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid create thread | Topic Not Exist", func(t *testing.T) {
		copyTopic := topicDomain
		copyTopic.Topic = "topic not exist"
		expectedErr := errors.New("failed to get topic")

		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.Create(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test case 3 | Invalid create thread | Error when creating thread", func(t *testing.T) {
		expectedErr := errors.New("failed to create thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Create", &threadDomain).Return(threads.Domain{}, expectedErr).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

		result, err := threadUseCase.Create(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 4 | Invalid create thread | Error when uploading image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		result, err := threadUseCase.Create(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 5 | Invalid create thread | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Create", &threadDomain).Return(threads.Domain{}, errors.New("failed to create")).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		result, err := threadUseCase.Create(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetManyWithPagination(t *testing.T) {
	t.Run("Test case 1 | Valid get thread with sort and order", func(t *testing.T) {
		pagination := dtoPagination.Request{
			Page:  1,
			Limit: 2,
			Sort:  "createdAt",
			Order: "desc",
		}

		query := dtoQuery.Request{
			Skip:  0,
			Limit: 2,
			Sort:  "createdAt",
			Order: -1,
		}
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetManyWithPagination", query, &threadDomain).Return([]threads.Domain{threadDomain}, 1, nil).Once()

		result, totalPage, totalData, err := threadUseCase.GetManyWithPagination(pagination, &threadDomain)

		assert.NotNil(t, result)
		assert.NotZero(t, totalPage)
		assert.NotZero(t, totalData)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get thread with sort and order | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get threads")

		pagination := dtoPagination.Request{
			Page:  1,
			Limit: 2,
			Sort:  "createdAt",
			Order: "asc",
		}

		query := dtoQuery.Request{
			Skip:  0,
			Limit: 2,
			Sort:  "createdAt",
			Order: 1,
		}
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetManyWithPagination", query, &threadDomain).Return([]threads.Domain{}, 0, expectedErr).Once()

		result, totalPage, totalData, err := threadUseCase.GetManyWithPagination(pagination, &threadDomain)

		assert.Equal(t, []threads.Domain{}, result)
		assert.Zero(t, totalPage)
		assert.Zero(t, totalData)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 3 | Invalid get thread with sort and order | Topic Not Exist", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")

		pagination := dtoPagination.Request{
			Page:  1,
			Limit: 2,
			Sort:  "createdAt",
			Order: "asc",
		}

		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, totalPage, totalData, err := threadUseCase.GetManyWithPagination(pagination, &threadDomain)

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

func TestGetAllByTopicID(t *testing.T) {
	t.Run("Test case 1 | Valid get all thread by topic id", func(t *testing.T) {
		threadRepository.On("GetAllByTopicID", threadDomain.TopicID).Return([]threads.Domain{threadDomain}, nil).Once()

		result, err := threadUseCase.GetAllByTopicID(threadDomain.TopicID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get all thread by topic id | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get threads")
		threadRepository.On("GetAllByTopicID", threadDomain.TopicID).Return([]threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.GetAllByTopicID(threadDomain.TopicID)

		assert.Equal(t, []threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetAllByUserID(t *testing.T) {
	t.Run("Test case 1 | Valid get all thread by user id", func(t *testing.T) {
		threadRepository.On("GetAllByUserID", threadDomain.CreatorID).Return([]threads.Domain{threadDomain}, nil).Once()

		result, err := threadUseCase.GetAllByUserID(threadDomain.CreatorID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get all thread by user id | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get threads")
		threadRepository.On("GetAllByUserID", threadDomain.CreatorID).Return([]threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.GetAllByUserID(threadDomain.CreatorID)

		assert.Equal(t, []threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestDomainToResponse(t *testing.T) {
	t.Run("Test case 1 | Valid domain to response", func(t *testing.T) {
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()

		result, actualErr := threadUseCase.DomainToResponse(threadDomain, userDomain.Id)

		assert.NotNil(t, result)
		assert.Nil(t, actualErr)
	})

	t.Run("Test case 2 | Invalid domain to response | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get creator")
		userRepository.On("GetByID", threadDomain.CreatorID).Return(users.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.DomainToResponse(threadDomain, userDomain.Id)

		assert.Equal(t, dtoThread.Response{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test case 3 | Invalid domain to response | Error when getting topic", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.DomainToResponse(threadDomain, primitive.NilObjectID)

		assert.Equal(t, dtoThread.Response{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestGetLikedByUserID(t *testing.T) {
	t.Run("Test case 1 | Valid get liked thread by user id", func(t *testing.T) {
		threadRepository.On("GetLikedByUserID", threadDomain.CreatorID).Return([]threads.Domain{threadDomain}, nil).Once()

		result, err := threadUseCase.GetLikedByUserID(threadDomain.CreatorID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get liked thread by user id | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get liked threads")
		threadRepository.On("GetLikedByUserID", threadDomain.CreatorID).Return([]threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.GetLikedByUserID(threadDomain.CreatorID)

		assert.Equal(t, []threads.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestDomainToResponseArray(t *testing.T) {
	t.Run("Test case 1 | Valid domain to response array", func(t *testing.T) {
		userRepository.On("GetByID", threadDomain.CreatorID).Return(userDomain, nil).Once()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()

		result, actualErr := threadUseCase.DomainsToResponseArray([]threads.Domain{threadDomain}, userDomain.Id)

		assert.NotNil(t, result)
		assert.Nil(t, actualErr)
	})

	t.Run("Test case 2 | Invalid domain to response array | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		userRepository.On("GetByID", threadDomain.CreatorID).Return(users.Domain{}, expectedErr).Once()

		result, actualErr := threadUseCase.DomainsToResponseArray([]threads.Domain{threadDomain}, primitive.NilObjectID)

		assert.Equal(t, []dtoThread.Response{}, result)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestUserUpdate(t *testing.T) {
	t.Run("Test case 1 | Valid user update thread", func(t *testing.T) {
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threadDomain, nil).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid user update thread | Error when updating thread", func(t *testing.T) {
		expectedErr := errors.New("failed to update thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid user update thread | Error when getting topic", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid user update thread | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 5 | Invalid user update thread | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 6 | Invalid user update thread | Error when uploading image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 7 | Invalid user update thread | Error when updating thread", func(t *testing.T) {
		expectedErr := errors.New("failed to update thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 8 | Invalid user update thread | User are not the thread creator", func(t *testing.T) {
		expectedErr := errors.New("user are not the thread creator")
		copyDomain := threadDomain
		copyDomain.CreatorID = primitive.NewObjectID()
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(copyDomain, nil).Once()

		result, err := threadUseCase.UserUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestAdminUpdate(t *testing.T) {
	t.Run("Test case 1 | Valid admin update thread", func(t *testing.T) {
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threadDomain, nil).Once()

		result, err := threadUseCase.AdminUpdate(&threadDomain, image)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid admin update thread | Error when getting topic", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topics.Domain{}, expectedErr).Once()

		result, err := threadUseCase.AdminUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid admin update thread | Error when getting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.AdminUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid admin update thread | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		result, err := threadUseCase.AdminUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 5 | Invalid admin update thread | Error when uploading image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		result, err := threadUseCase.AdminUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 6 | Invalid admin update thread | Error when updating thread", func(t *testing.T) {
		expectedErr := errors.New("failed to update thread")
		topicRepository.On("GetByID", threadDomain.TopicID).Return(topicDomain, nil).Once()
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		threadRepository.On("Update", mock.Anything).Return(threads.Domain{}, expectedErr).Once()

		result, err := threadUseCase.AdminUpdate(&threadDomain, image)

		assert.Equal(t, threads.Domain{}, result)
		assert.Equal(t, expectedErr, err)
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

func TestLike(t *testing.T) {
	t.Run("Test case 1 | Valid like thread", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		threadRepository.On("CheckLikedByUserID", threadDomain.CreatorID, threadDomain.Id).Return(errors.New("not found")).Once()
		threadRepository.On("AppendLike", threadDomain.CreatorID, threadDomain.Id).Return(nil).Once()

		err := threadUseCase.Like(threadDomain.CreatorID, threadDomain.Id)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid like thread | Error when checking liked by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to like thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		threadRepository.On("CheckLikedByUserID", threadDomain.CreatorID, threadDomain.Id).Return(errors.New("not found")).Once()
		threadRepository.On("AppendLike", threadDomain.CreatorID, threadDomain.Id).Return(expectedErr).Once()

		err := threadUseCase.Like(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid like thread | Error when appending like", func(t *testing.T) {
		expectedErr := errors.New("user already like this thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		threadRepository.On("CheckLikedByUserID", threadDomain.CreatorID, threadDomain.Id).Return(nil).Once()

		err := threadUseCase.Like(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid like thread | Error when getting thread by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		err := threadUseCase.Like(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})
}

func TestUnlike(t *testing.T) {
	t.Run("Test case 1 | Valid unlike thread", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		threadRepository.On("CheckLikedByUserID", threadDomain.CreatorID, threadDomain.Id).Return(nil).Once()
		threadRepository.On("RemoveLike", threadDomain.CreatorID, threadDomain.Id).Return(nil).Once()

		err := threadUseCase.Unlike(threadDomain.CreatorID, threadDomain.Id)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid unlike thread | Error when checking liked by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to unlike thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		threadRepository.On("CheckLikedByUserID", threadDomain.CreatorID, threadDomain.Id).Return(nil).Once()
		threadRepository.On("RemoveLike", threadDomain.CreatorID, threadDomain.Id).Return(expectedErr).Once()

		err := threadUseCase.Unlike(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid unlike thread | Error when deleting like", func(t *testing.T) {
		expectedErr := errors.New("user not like this thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		threadRepository.On("CheckLikedByUserID", threadDomain.CreatorID, threadDomain.Id).Return(errors.New("not found")).Once()

		err := threadUseCase.Unlike(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid unlike thread | Error when getting thread by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		err := threadUseCase.Unlike(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})
}

func TestRemoveUserFromAllLikes(t *testing.T) {
	t.Run("Test case 1 | Valid remove user from all likes", func(t *testing.T) {
		threadRepository.On("RemoveUserFromAllLikes", threadDomain.CreatorID).Return(nil).Once()

		err := threadUseCase.RemoveUserFromAllLikes(threadDomain.CreatorID)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid remove user from all likes | Error when removing user from all likes", func(t *testing.T) {
		expectedErr := errors.New("failed to remove user from all likes")
		threadRepository.On("RemoveUserFromAllLikes", threadDomain.CreatorID).Return(expectedErr).Once()

		err := threadUseCase.RemoveUserFromAllLikes(threadDomain.CreatorID)

		assert.Equal(t, expectedErr, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test case 1 | Valid delete thread", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(nil).Once()

		thread, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Nil(t, err)
		assert.Equal(t, threadDomain, thread)
	})

	t.Run("Test case 2 | Invalid delete thread | Error when getting thread by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		_, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid delete thread | Error when deleting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to delete thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(expectedErr).Once()

		_, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 5 | Invalid delete thread | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		_, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 8 | Invalid user delete thread | User are not the thread creator", func(t *testing.T) {
		expectedErr := errors.New("user are not the thread creator")
		copyDomain := threadDomain
		copyDomain.CreatorID = primitive.NewObjectID()

		threadRepository.On("GetByID", threadDomain.Id).Return(copyDomain, nil).Once()

		_, err := threadUseCase.Delete(threadDomain.CreatorID, threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})
}

func TestDeleteAllByUserID(t *testing.T) {
	t.Run("Test case 1 | Valid delete all thread by user id", func(t *testing.T) {
		threadRepository.On("GetAllByUserID", threadDomain.CreatorID).Return([]threads.Domain{threadDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("DeleteAllByUserID", mock.Anything).Return(nil).Once()

		err := threadUseCase.DeleteAllByUserID(threadDomain.CreatorID)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid delete all thread by user id | Error when deleting thread by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to delete user threads")
		threadRepository.On("GetAllByUserID", threadDomain.CreatorID).Return([]threads.Domain{threadDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("DeleteAllByUserID", mock.Anything).Return(expectedErr).Once()

		err := threadUseCase.DeleteAllByUserID(threadDomain.CreatorID)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid delete all thread by user id | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		threadRepository.On("GetAllByUserID", threadDomain.CreatorID).Return([]threads.Domain{threadDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		err := threadUseCase.DeleteAllByUserID(threadDomain.CreatorID)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid delete all thread by user id | Error when getting thread by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to get user threads")
		threadRepository.On("GetAllByUserID", threadDomain.CreatorID).Return([]threads.Domain{}, expectedErr).Once()

		err := threadUseCase.DeleteAllByUserID(threadDomain.CreatorID)

		assert.Equal(t, expectedErr, err)
	})
}

func TestDeleteByThreadID(t *testing.T) {
	t.Run("Test case 1 | Valid delete thread by thread id", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(nil).Once()

		err := threadUseCase.DeleteByThreadID(threadDomain.Id)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid delete thread by thread id | Error when deleting thread by thread id", func(t *testing.T) {
		expectedErr := errors.New("failed to delete thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(expectedErr).Once()

		err := threadUseCase.DeleteByThreadID(threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid delete thread by thread id | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		err := threadUseCase.DeleteByThreadID(threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid delete thread by thread id | Error when getting thread by thread id", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		err := threadUseCase.DeleteByThreadID(threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})
}

func TestAdminDelete(t *testing.T) {
	t.Run("Test case 1 | Valid admin delete thread", func(t *testing.T) {
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(nil).Once()

		thread, err := threadUseCase.AdminDelete(threadDomain.Id)

		assert.Nil(t, err)
		assert.Equal(t, threadDomain, thread)
	})

	t.Run("Test case 2 | Invalid admin delete thread | Error when getting thread by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threads.Domain{}, expectedErr).Once()

		_, err := threadUseCase.AdminDelete(threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 3 | Invalid admin delete thread | Error when deleting thread", func(t *testing.T) {
		expectedErr := errors.New("failed to delete thread")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		threadRepository.On("Delete", mock.Anything).Return(expectedErr).Once()

		_, err := threadUseCase.AdminDelete(threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test case 4 | Invalid admin delete thread | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		threadRepository.On("GetByID", threadDomain.Id).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		_, err := threadUseCase.AdminDelete(threadDomain.Id)

		assert.Equal(t, expectedErr, err)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Test case 1 | Valid admin get all threads", func(t *testing.T) {
		threadRepository.On("GetAll").Return([]threads.Domain{threadDomain}, nil).Once()

		thread, err := threadUseCase.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 1, thread)
	})

	t.Run("Test case 2 | Invalid admin get all threads | Error when getting all threads", func(t *testing.T) {
		expectedErr := errors.New("failed to get threads")
		threadRepository.On("GetAll").Return([]threads.Domain{}, expectedErr).Once()

		_, err := threadUseCase.GetAll()

		assert.Equal(t, expectedErr, err)
	})
}
