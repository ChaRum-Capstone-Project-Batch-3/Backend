package comments

import (
	"charum/business/threads"
	"charum/business/users"
	"charum/dto"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentUseCase struct {
	commentRepository Repository
	threadRepository  threads.Repository
	userRepository    users.Repository
}

func NewCommentUseCase(cr Repository, tr threads.Repository, ur users.Repository) UseCase {
	return &CommentUseCase{
		commentRepository: cr,
		threadRepository:  tr,
		userRepository:    ur,
	}
}

/*
Create
*/

func (cu *CommentUseCase) Create(domain *Domain) (Domain, error) {
	_, err := cu.userRepository.GetByID(domain.UserID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	_, err = cu.threadRepository.GetByID(domain.ThreadID)
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

func (cu *CommentUseCase) DomainToResponse(comment Domain) (dto.ResponseComment, error) {
	responseComment := dto.ResponseComment{}

	user, err := cu.userRepository.GetByID(comment.UserID)
	if err != nil {
		return dto.ResponseComment{}, errors.New("failed to get user")
	}

	responseComment.Id = comment.Id
	responseComment.ThreadID = comment.ThreadID
	responseComment.User = user
	responseComment.Comment = comment.Comment
	responseComment.CreatedAt = comment.CreatedAt
	responseComment.UpdatedAt = comment.UpdatedAt

	return responseComment, nil
}

func (cu *CommentUseCase) DomainToResponseArray(comments []Domain) ([]dto.ResponseComment, error) {
	responseComments := []dto.ResponseComment{}

	for _, comment := range comments {
		responseComment, err := cu.DomainToResponse(comment)
		if err != nil {
			return []dto.ResponseComment{}, err
		}

		responseComments = append(responseComments, responseComment)
	}

	return responseComments, nil
}

/*
Update
*/

func (cu *CommentUseCase) Update(domain *Domain) (Domain, error) {
	comment, err := cu.commentRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, err
	}

	if comment.UserID != domain.UserID {
		return Domain{}, errors.New("you are not the owner of this comment")
	}

	_, err = cu.threadRepository.GetByID(comment.ThreadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	comment.Comment = domain.Comment
	comment.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	comment, err = cu.commentRepository.Update(&comment)
	if err != nil {
		return Domain{}, err
	}

	return comment, nil
}

/*
Delete
*/

func (cu *CommentUseCase) Delete(id primitive.ObjectID, userID primitive.ObjectID) (Domain, error) {
	comment, err := cu.commentRepository.GetByID(id)
	if err != nil {
		return Domain{}, err
	}

	if comment.UserID != userID {
		return Domain{}, errors.New("you are not the owner of this comment")
	}

	_, err = cu.threadRepository.GetByID(comment.ThreadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	err = cu.commentRepository.Delete(id)
	if err != nil {
		return Domain{}, err
	}

	return comment, nil
}
