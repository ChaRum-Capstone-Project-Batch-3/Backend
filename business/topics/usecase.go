package topics

import (
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicUseCase struct {
	topicsRepository Repository
}

func NewTopicUseCase(tr Repository) UseCase {
	return &TopicUseCase{
		topicsRepository: tr,
	}
}

/*
Create topic
*/

func (tc *TopicUseCase) CreateTopic(domain *Domain) (Domain, error) {
	_, err := tc.topicsRepository.GetByTopic(domain.Topic)
	if err == nil {
		return Domain{}, errors.New("topic already exist")
	}

	domain.Id = primitive.NewObjectID()
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := tc.topicsRepository.CreateTopic(domain)
	if err != nil {
		return Domain{}, errors.New("failed to create topic")
	}
	return result, nil
}

/*
Get topic
*/

func (tc *TopicUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	result, err := tc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

func (tc *TopicUseCase) GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error) {
	skip := pagination.Limit * (pagination.Page - 1)
	var orderInMongo int

	if pagination.Order == "asc" {
		orderInMongo = 1
	} else {
		orderInMongo = -1
	}

	query := dtoQuery.Request{
		Skip:  skip,
		Limit: pagination.Limit,
		Order: orderInMongo,
		Sort:  pagination.Sort,
	}

	users, totalData, err := tc.topicsRepository.GetManyWithPagination(query, domain)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("failed to get topics")
	}

	totalPage := math.Ceil(float64(totalData) / float64(pagination.Limit))
	return users, int(totalPage), totalData, nil
}

// get by topic
func (tc *TopicUseCase) GetByTopic(topic string) (Domain, error) {
	result, err := tc.topicsRepository.GetByTopic(topic)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

/*
Update topic
*/
func (tc *TopicUseCase) UpdateTopic(id primitive.ObjectID, domain *Domain) (Domain, error) {
	result, err := tc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	if domain.Topic != result.Topic {
		_, err := tc.topicsRepository.GetByTopic(domain.Topic)
		if err == nil {
			return Domain{}, errors.New("topic already exist")
		}
	}

	result.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	// update topic and description
	result.Topic = domain.Topic
	result.Description = domain.Description
	updatedResult, err := tc.topicsRepository.UpdateTopic(&result)
	if err != nil {
		return Domain{}, errors.New("failed to update topic")
	}
	return updatedResult, nil
}

/*
Delete topic
*/

func (tc *TopicUseCase) DeleteTopic(id primitive.ObjectID) (Domain, error) {
	result, err := tc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	err = tc.topicsRepository.DeleteTopic(id)
	if err != nil {
		return Domain{}, errors.New("failed to delete topic")
	}
	return result, nil
}
