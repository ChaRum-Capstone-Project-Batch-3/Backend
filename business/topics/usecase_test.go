package topics_test

import (
	"charum/business/topics"
	_topicMock "charum/business/topics/mocks"
	dtoPagination "charum/dto/pagination"
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
	topicRepository      _topicMock.Repository
	cloudinaryRepository _cloudinaryMock.Function
	topicUseCase         topics.UseCase
	topicDomain          topics.Domain
	image                *multipart.FileHeader
)

func TestMain(m *testing.M) {
	topicUseCase = topics.NewTopicUseCase(&topicRepository, &cloudinaryRepository)

	topicDomain = topics.Domain{
		Id:          primitive.NewObjectID(),
		Topic:       "Test Topic",
		Description: "Test Topic Description",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	image = &multipart.FileHeader{}

	m.Run()

}

func TestCreate(t *testing.T) {
	t.Run("Test case 1 | Valid create topic", func(t *testing.T) {
		topicRepository.On("GetByTopic", topicDomain.Topic).Return(topics.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		topicRepository.On("Create", &topicDomain).Return(topicDomain, nil).Once()

		result, err := topicUseCase.Create(&topicDomain, image)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid create topic | Error when getting topic by topic", func(t *testing.T) {
		expectedErr := errors.New("topic already exist")
		topicRepository.On("GetByTopic", topicDomain.Topic).Return(topicDomain, nil).Once()

		result, err := topicUseCase.Create(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 3 | Invalid create topic | Error when uploading image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		topicRepository.On("GetByTopic", topicDomain.Topic).Return(topics.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		result, err := topicUseCase.Create(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 4 | Invalid create topic | Error when creating topic", func(t *testing.T) {
		expectedErr := errors.New("failed to create topic")
		topicRepository.On("GetByTopic", topicDomain.Topic).Return(topics.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		topicRepository.On("Create", &topicDomain).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.Create(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("Test case 1 | Valid get topic by id", func(t *testing.T) {
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()

		result, err := topicUseCase.GetByID(topicDomain.Id)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get topic by id | Error when getting topic by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.GetByID(topicDomain.Id)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestGetByTopic(t *testing.T) {
	t.Run("Test case 1 | Valid get topic by topic", func(t *testing.T) {
		topicRepository.On("GetByTopic", topicDomain.Topic).Return(topicDomain, nil).Once()

		result, err := topicUseCase.GetByTopic(topicDomain.Topic)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get topic by topic | Error when getting topic by topic", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByTopic", topicDomain.Topic).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.GetByTopic(topicDomain.Topic)

		assert.Equal(t, topics.Domain{}, result)
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
		topicRepository.On("GetManyWithPagination", mock.Anything, mock.Anything).Return([]topics.Domain{topicDomain}, 1, nil).Once()

		result, totalPage, totalData, err := topicUseCase.GetManyWithPagination(pagination, &topicDomain)

		assert.NotNil(t, result)
		assert.NotZero(t, totalPage)
		assert.NotZero(t, totalData)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid get thread with sort and order | Error when getting thread with sort and order", func(t *testing.T) {
		pagination := dtoPagination.Request{
			Page:  1,
			Limit: 2,
			Sort:  "createdAt",
			Order: "asc",
		}

		expectedErr := errors.New("failed to get topics")
		topicRepository.On("GetManyWithPagination", mock.Anything, mock.Anything).Return([]topics.Domain{}, 0, expectedErr).Once()

		result, totalPage, totalData, err := topicUseCase.GetManyWithPagination(pagination, &topicDomain)

		assert.Equal(t, []topics.Domain{}, result)
		assert.Zero(t, totalPage)
		assert.Zero(t, totalData)
		assert.Equal(t, err, expectedErr)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test case 1 | Valid update topic", func(t *testing.T) {
		copyDomain := topicDomain
		copyDomain.Topic = "Updated Topic"
		copyDomain.Description = "Updated Description"

		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		topicRepository.On("Update", mock.Anything).Return(topicDomain, nil).Once()

		result, err := topicUseCase.Update(&topicDomain, image)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid update topic | Error when getting topic by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.Update(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 3 | Invalid update topic | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		result, err := topicUseCase.Update(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 4 | Invalid update topic | Error when uploading image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		result, err := topicUseCase.Update(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 5 | Invalid update topic | Error when updating topic", func(t *testing.T) {
		expectedErr := errors.New("failed to update topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		topicRepository.On("Update", mock.Anything).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.Update(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 5 | Invalid update topic | Topic already exist", func(t *testing.T) {
		expectedErr := errors.New("topic already exist")
		copyDomain := topicDomain
		copyDomain.Topic = "Updated Topic"
		topicRepository.On("GetByID", copyDomain.Id).Return(copyDomain, nil).Once()
		topicRepository.On("GetByTopic", mock.Anything).Return(copyDomain, nil).Once()

		result, err := topicUseCase.Update(&topicDomain, image)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test case 1 | Valid delete topic", func(t *testing.T) {
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		topicRepository.On("Delete", topicDomain.Id).Return(nil).Once()

		result, err := topicUseCase.Delete(topicDomain.Id)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid delete topic | Error when getting topic by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.Delete(topicDomain.Id)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 3 | Invalid delete topic | Error when deleting image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		result, err := topicUseCase.Delete(topicDomain.Id)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test case 4 | Invalid delete topic | Error when deleting topic", func(t *testing.T) {
		expectedErr := errors.New("failed to delete topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		topicRepository.On("Delete", topicDomain.Id).Return(expectedErr).Once()

		result, err := topicUseCase.Delete(topicDomain.Id)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}
