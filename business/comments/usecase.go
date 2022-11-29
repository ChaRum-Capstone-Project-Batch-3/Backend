package comments

import (
	"charum/business/threads"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentUseCase struct {
	commentRepository Repository
	threadRepository  threads.Repository
}

func NewCommentUseCase(cr Repository, tr threads.Repository) UseCase {
	return &CommentUseCase{
		commentRepository: cr,
		threadRepository:  tr,
	}
}

/*
Create
*/

func (cu *CommentUseCase) Create(domain *Domain) (Domain, error) {
	_, err := cu.threadRepository.GetByID(domain.ThreadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	domain.Id = primitive.NewObjectID()
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	comment, err := cu.commentRepository.Create(domain)
	if err != nil {
		return Domain{}, err
	}

	return comment, nil
}

/*
Read
*/

func (cu *CommentUseCase) GetByThreadID(threadID primitive.ObjectID) ([]Domain, error) {
	_, err := cu.threadRepository.GetByID(threadID)
	if err != nil {
		return []Domain{}, errors.New("failed to get thread")
	}

	comments, err := cu.commentRepository.GetByThreadID(threadID)
	if err != nil {
		return []Domain{}, err
	}

	return comments, nil
}

/*
Update
*/

/*
Delete
*/
