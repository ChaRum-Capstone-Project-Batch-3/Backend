package threads

import (
	"charum/business/topics"
	"charum/business/users"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThreadUseCase struct {
	threadRepository Repository
	topicRepository  topics.Repository
	userRepository   users.Repository
}

func NewThreadUseCase(thr Repository, tor topics.Repository, ur users.Repository) UseCase {
	return &ThreadUseCase{
		threadRepository: thr,
		topicRepository:  tor,
		userRepository:   ur,
	}
}

/*
Create
*/

func (tu *ThreadUseCase) Create(creatorID primitive.ObjectID, topicName string, domain *Domain) (Domain, error) {
	topic, err := tu.topicRepository.GetByTopic(topicName)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	domain.Id = primitive.NewObjectID()
	domain.TopicID = topic.Id
	domain.CreatorID = creatorID
	domain.Likes = []Like{}
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	thread, err := tu.threadRepository.Create(domain)
	if err != nil {
		return Domain{}, errors.New("failed to create thread")
	}

	return thread, nil
}

/*
Read
*/

func (tu *ThreadUseCase) GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, error) {
	skip := limit * (page - 1)
	var orderInMongo int

	if order == "asc" {
		orderInMongo = 1
	} else {
		orderInMongo = -1
	}

	users, totalData, err := tu.threadRepository.GetWithSortAndOrder(skip, limit, sort, orderInMongo)
	if err != nil {
		return []Domain{}, 0, errors.New("failed to get threads")
	}

	totalPage := math.Ceil(float64(totalData) / float64(limit))
	return users, int(totalPage), nil
}

func (tu *ThreadUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	thread, err := tu.threadRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	return thread, nil
}

/*
Update
*/

func (tu *ThreadUseCase) Update(userID primitive.ObjectID, threadID primitive.ObjectID, topicName string, domain *Domain) (Domain, error) {
	topic, err := tu.topicRepository.GetByTopic(topicName)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	thread, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	user, err := tu.userRepository.GetByID(userID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if user.Role != "admin" && thread.CreatorID != userID {
		return Domain{}, errors.New("you are not the thread creator")
	}

	thread.TopicID = topic.Id
	thread.Title = domain.Title
	thread.Description = domain.Description
	thread.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updatedThread, err := tu.threadRepository.Update(&thread)
	if err != nil {
		return Domain{}, errors.New("failed to update thread")
	}

	return updatedThread, nil
}

/*
Delete
*/

func (tu *ThreadUseCase) Delete(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error) {
	thread, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	user, err := tu.userRepository.GetByID(userID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if user.Role != "admin" && thread.CreatorID != userID {
		return Domain{}, errors.New("you are not the thread creator")
	}

	err = tu.threadRepository.Delete(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to delete thread")
	}

	return thread, nil
}
