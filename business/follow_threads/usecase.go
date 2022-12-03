package follow_threads

import (
	"charum/business/threads"
	"charum/business/users"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowThreadUseCase struct {
	followThreadRepository Repository
	userRepository         users.Repository
	threadRepository       threads.Repository
}

func NewFollowThreadUseCase(ftr Repository, ur users.Repository, tr threads.Repository) UseCase {
	return &FollowThreadUseCase{
		followThreadRepository: ftr,
		userRepository:         ur,
		threadRepository:       tr,
	}
}

/*
Create
*/

func (ftu *FollowThreadUseCase) Create(domain *Domain) (Domain, error) {
	_, err := ftu.userRepository.GetByID(domain.UserID)
	if err != nil {
		return Domain{}, err
	}

	_, err = ftu.threadRepository.GetByID(domain.ThreadID)
	if err != nil {
		return Domain{}, err
	}

	domain.Id = primitive.NewObjectID()
	domain.Notification = 0
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := ftu.followThreadRepository.Create(domain)
	if err != nil {
		return Domain{}, err
	}

	return result, nil
}

/*
Read
*/

/*
Update
*/

/*
Delete
*/

func (ftu *FollowThreadUseCase) Delete(domain *Domain) (Domain, error) {
	_, err := ftu.userRepository.GetByID(domain.UserID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	_, err = ftu.threadRepository.GetByID(domain.ThreadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	result, err := ftu.followThreadRepository.GetByUserIDAndThreadID(domain.UserID, domain.ThreadID)
	if err != nil {
		return Domain{}, errors.New("follow thread not found")
	}

	if result.UserID != domain.UserID {
		return Domain{}, errors.New("you are not the owner of this follow thread")
	}

	err = ftu.followThreadRepository.Delete(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to unfollow thread")
	}

	return result, nil
}
