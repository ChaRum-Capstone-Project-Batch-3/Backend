package topics

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TopicsUseCase struct {
	topicsRepository Repository
}

func NewTopicUseCase(tr Repository) UseCase {
	return &TopicsUseCase{
		topicsRepository: tr,
	}
}

/*
Create topic
*/

func (uc *TopicsUseCase) CreateTopic(domain *Domain) (Domain, error) {
	domain.Id = primitive.NewObjectID()
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := uc.topicsRepository.CreateTopic(domain)
	if err != nil {
		return Domain{}, errors.New("failed to create topic")
	}
	return result, nil
}

/*
Get topic by id
*/

func (uc *TopicsUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	result, err := uc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

/*
Update topic
*/
func (uc *TopicsUseCase) UpdateTopic(id primitive.ObjectID, domain *Domain) (Domain, error) {
	result, err := uc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	if domain.Topic != "" {
		result.Topic = domain.Topic
	}
	if domain.Description != "" {
		result.Description = domain.Description
	}

	result.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updatedResult, err := uc.topicsRepository.UpdateTopic(&result)
	if err != nil {
		return Domain{}, errors.New("failed to update data")
	}
	return updatedResult, nil
}

/*
Delete topic
*/

func (uc *TopicsUseCase) DeleteTopic(id primitive.ObjectID) (Domain, error) {
	result, err := uc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	err = uc.topicsRepository.DeleteTopic(id)
	if err != nil {
		return Domain{}, errors.New("failed to delete topic")
	}
	return result, nil
}
