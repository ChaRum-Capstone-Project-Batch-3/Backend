package topics

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

func (uc *TopicUseCase) CreateTopic(domain *Domain) (Domain, error) {
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

func (uc *TopicUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	result, err := uc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

/*
Update topic
*/
func (uc *TopicUseCase) UpdateTopic(id primitive.ObjectID, domain *Domain) (Domain, error) {
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
		return Domain{}, errors.New("failed to update topic")
	}
	return updatedResult, nil
}

/*
Delete topic
*/

func (uc *TopicUseCase) DeleteTopic(id primitive.ObjectID) (Domain, error) {
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
