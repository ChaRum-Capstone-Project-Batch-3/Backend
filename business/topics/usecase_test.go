package topics_test

import (
	"charum/business/topics"
	_topicMock "charum/business/topics/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var (
	topicRepository _topicMock.Repository
	topicUseCase    topics.UseCase
	topicDomain     topics.Domain
)

func TestMain(m *testing.M) {
	topicUseCase = topics.NewTopicUseCase(&topicRepository)

	topicDomain = topics.Domain{
		Id:          primitive.NewObjectID(),
		Topic:       "Test Topic",
		Description: "Test Topic Description",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	m.Run()

}

func TestCreateTopic(t *testing.T) {
	t.Run("Test case 1 | Valid create topic", func(t *testing.T) {
		topicRepository.On("CreateTopic", &topicDomain).Return(topicDomain, nil).Once()

		result, err := topicUseCase.CreateTopic(&topicDomain)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})
	t.Run("Test case 2 | Invalid create topic | Invalid Register | Error when creating topic", func(t *testing.T) {
		expectedErr := errors.New("failed to create topic")
		topicRepository.On("CreateTopic", &topicDomain).Return(topics.Domain{}, errors.New("not found")).Once()
		topicRepository.On("CreateTopic", mock.Anything).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.CreateTopic(&topicDomain)

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

func TestUpdateTopic(t *testing.T) {
	t.Run("Test case 1 | Valid update topic", func(t *testing.T) {
		copyDomain := topicDomain
		copyDomain.Topic = "Updated Topic"
		copyDomain.Description = "Updated Description"

		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		topicRepository.On("UpdateTopic", mock.Anything).Return(topicDomain, nil).Once()

		result, err := topicUseCase.UpdateTopic(topicDomain.Id, &topicDomain)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})
	t.Run("Test case 2 | Invalid update topic | Error when getting topic by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.UpdateTopic(topicDomain.Id, &topicDomain)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
	t.Run("Test case 3 | Invalid update topic | Error when updating topic", func(t *testing.T) {
		expectedErr := errors.New("failed to update topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		topicRepository.On("UpdateTopic", mock.Anything).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.UpdateTopic(topicDomain.Id, &topicDomain)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}

func TestDeleteTopic(t *testing.T) {
	t.Run("Test case 1 | Valid delete topic", func(t *testing.T) {
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		topicRepository.On("DeleteTopic", topicDomain.Id).Return(nil).Once()

		result, err := topicUseCase.DeleteTopic(topicDomain.Id)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})
	t.Run("Test case 2 | Invalid delete topic | Error when getting topic by id", func(t *testing.T) {
		expectedErr := errors.New("failed to get topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topics.Domain{}, expectedErr).Once()

		result, err := topicUseCase.DeleteTopic(topicDomain.Id)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
	t.Run("Test case 3 | Invalid delete topic | Error when deleting topic", func(t *testing.T) {
		expectedErr := errors.New("failed to delete topic")
		topicRepository.On("GetByID", topicDomain.Id).Return(topicDomain, nil).Once()
		topicRepository.On("DeleteTopic", topicDomain.Id).Return(expectedErr).Once()

		result, err := topicUseCase.DeleteTopic(topicDomain.Id)

		assert.Equal(t, topics.Domain{}, result)
		assert.Equal(t, err, expectedErr)
	})
}
