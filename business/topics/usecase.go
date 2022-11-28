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
	_, err := uc.topicsRepository.GetByTopic(domain.Topic)
	if err == nil {
		return Domain{}, errors.New("topic already exist")
	}

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
Get topic
*/

func (uc *TopicUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	result, err := uc.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

func (uc *TopicUseCase) GetAll() ([]Domain, error) {
	result, err := uc.topicsRepository.GetAll()
	if err != nil {
		return []Domain{}, errors.New("failed to get all topic")
	}
	return result, nil
}

// get by topic
func (uc *TopicUseCase) GetByTopic(topic string) (Domain, error) {
	result, err := uc.topicsRepository.GetByTopic(topic)
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

	if domain.Topic != result.Topic {
		_, err := uc.topicsRepository.GetByTopic(domain.Topic)
		if err == nil {
			return Domain{}, errors.New("topic already exist")
		}
	}

	result.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	// update topic and description
	result.Topic = domain.Topic
	result.Description = domain.Description
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
