package follow_threads

import (
	"charum/business/threads"
	"charum/business/users"
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

func (ftu *FollowThreadUseCase) Delete(id primitive.ObjectID) (Domain, error) {
	result, err := ftu.followThreadRepository.GetByID(id)
	if err != nil {
		return Domain{}, err
	}

	err = ftu.followThreadRepository.Delete(id)
	if err != nil {
		return Domain{}, err
	}

	return result, nil
}
