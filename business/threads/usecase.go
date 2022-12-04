package threads

import (
	"charum/business/topics"
	"charum/business/users"
	"charum/dto"
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

func (tu *ThreadUseCase) Create(domain *Domain) (Domain, error) {
	_, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	domain.Id = primitive.NewObjectID()
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

func (tu *ThreadUseCase) GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, int, error) {
	skip := limit * (page - 1)
	var orderInMongo int

	if order == "asc" {
		orderInMongo = 1
	} else {
		orderInMongo = -1
	}

	threads, totalData, err := tu.threadRepository.GetWithSortAndOrder(skip, limit, sort, orderInMongo)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("failed to get threads")
	}

	totalPage := math.Ceil(float64(totalData) / float64(limit))

	return threads, int(totalPage), totalData, nil
}

func (tu *ThreadUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	thread, err := tu.threadRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	return thread, nil
}

func (tu *ThreadUseCase) DomainToResponse(domain Domain) (dto.ResponseThread, error) {
	creator, err := tu.userRepository.GetByID(domain.CreatorID)
	if err != nil {
		return dto.ResponseThread{}, errors.New("failed to get creator")
	}

	likes := []dto.Like{}
	for _, like := range domain.Likes {
		user, err := tu.userRepository.GetByID(like.UserID)
		if err != nil {
			return dto.ResponseThread{}, errors.New("failed to get user who like thread")
		}

		likes = append(likes, dto.Like{
			User:      user,
			CreatedAt: domain.CreatedAt,
		})
	}

	topic, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return dto.ResponseThread{}, errors.New("failed to get topic")
	}

	return dto.ResponseThread{
		Id:            domain.Id,
		Topic:         topic,
		Creator:       creator,
		Title:         domain.Title,
		Description:   domain.Description,
		Likes:         likes,
		TotalLike:     len(domain.Likes),
		SuspendStatus: domain.SuspendStatus,
		SuspendDetail: domain.SuspendDetail,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}, nil
}

func (tu *ThreadUseCase) DomainsToResponseArray(domains []Domain) ([]dto.ResponseThread, error) {
	var responses []dto.ResponseThread
	for _, domain := range domains {
		response, err := tu.DomainToResponse(domain)
		if err != nil {
			return []dto.ResponseThread{}, errors.New("failed to get thread")
		}

		responses = append(responses, response)
	}

	return responses, nil
}

/*
Update
*/

func (tu *ThreadUseCase) Update(domain *Domain) (Domain, error) {
	_, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	thread, err := tu.threadRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	user, err := tu.userRepository.GetByID(domain.CreatorID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if user.Role != "admin" && thread.CreatorID != domain.CreatorID {
		return Domain{}, errors.New("user are not the thread creator")
	}

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
		return Domain{}, errors.New("user are not the thread creator")
	}

	err = tu.threadRepository.Delete(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to delete thread")
	}

	return thread, nil
}
