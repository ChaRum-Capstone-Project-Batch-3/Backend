package bookmarks

import (
	"charum/business/threads"
	"charum/business/topics"
	"charum/business/users"
	dtoBookmark "charum/dto/bookmarks"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarkUseCase struct {
	bookmarkRepository Repository
	threadRepository   threads.Repository
	userRepository     users.Repository
	topicRepository    topics.Repository
	threadUseCase      threads.UseCase
}

func NewBookmarkUseCase(br Repository, tr threads.Repository, ur users.Repository, tc topics.Repository, tuc threads.UseCase) UseCase {
	return &BookmarkUseCase{
		bookmarkRepository: br,
		threadRepository:   tr,
		userRepository:     ur,
		topicRepository:    tc,
		threadUseCase:      tuc,
	}
}

/*
Add Bookmark
*/
func (bu *BookmarkUseCase) Create(domain *Domain) (Domain, error) {
	_, err := bu.threadRepository.GetByID(domain.ThreadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	_, err = bu.bookmarkRepository.GetByUserIDAndThreadID(domain.UserID, domain.ThreadID)
	if err == nil {
		return Domain{}, errors.New("thread already bookmarked")
	}

	domain.Id = primitive.NewObjectID()
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	bookmark, err := bu.bookmarkRepository.Create(domain)
	if err != nil {
		return Domain{}, errors.New("failed to add thread to bookmark")
	}

	return bookmark, nil
}

/*
Read
*/

func (bu *BookmarkUseCase) GetAllByUserID(userID primitive.ObjectID) ([]Domain, error) {
	result, err := bu.bookmarkRepository.GetAllByUserID(userID)
	if err != nil {
		return []Domain{}, errors.New("failed to get all bookmark")
	}

	return result, nil
}

func (bu *BookmarkUseCase) CountByThreadID(threadID primitive.ObjectID) (int, error) {
	_, err := bu.threadRepository.GetByID(threadID)
	if err != nil {
		return 0, errors.New("failed to get thread")
	}

	result, err := bu.bookmarkRepository.CountByThreadID(threadID)
	if err != nil {
		return 0, errors.New("failed to count bookmark")
	}

	return result, nil
}

func (bu *BookmarkUseCase) CheckBookmarkedThread(userID primitive.ObjectID, threadID primitive.ObjectID) (bool, error) {
	_, err := bu.threadRepository.GetByID(threadID)
	if err != nil {
		return false, errors.New("failed to get thread")
	}

	_, err = bu.bookmarkRepository.GetByUserIDAndThreadID(userID, threadID)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (bu *BookmarkUseCase) DomainToResponse(domain Domain, userID primitive.ObjectID) (dtoBookmark.Response, error) {
	thread, err := bu.threadRepository.GetByID(domain.ThreadID)
	if err != nil {
		return dtoBookmark.Response{}, errors.New("failed to get bookmark")
	}

	responseThread, err := bu.threadUseCase.DomainToResponse(thread, userID)
	if err != nil {
		return dtoBookmark.Response{}, errors.New("failed to get bookmark")
	}

	responseBookmark := dtoBookmark.Response{
		Id:     domain.Id,
		UserID: domain.UserID,
		Thread: responseThread,
	}

	return responseBookmark, nil
}

func (bu *BookmarkUseCase) DomainsToResponseArray(domains []Domain, userID primitive.ObjectID) ([]dtoBookmark.Response, error) {
	var responses []dtoBookmark.Response
	for _, domain := range domains {
		response, err := bu.DomainToResponse(domain, userID)
		if err != nil {
			return []dtoBookmark.Response{}, errors.New("failed to get bookmark")
		}

		responses = append(responses, response)
	}

	return responses, nil
}

/*
Delete
*/

func (bu *BookmarkUseCase) Delete(domain *Domain) (Domain, error) {
	bookmark, err := bu.bookmarkRepository.GetByUserIDAndThreadID(domain.UserID, domain.ThreadID)
	if err != nil {
		return Domain{}, errors.New("failed to get bookmark")
	}

	result := bu.bookmarkRepository.Delete(domain)
	if result != nil {
		return Domain{}, errors.New("failed to delete bookmark")
	}

	return bookmark, nil
}

func (bu *BookmarkUseCase) DeleteAllByUserID(userID primitive.ObjectID) error {
	result := bu.bookmarkRepository.DeleteAllByUserID(userID)
	if result != nil {
		return errors.New("failed to delete bookmark")
	}

	return nil
}

func (bu *BookmarkUseCase) DeleteAllByThreadID(threadID primitive.ObjectID) error {
	result := bu.bookmarkRepository.DeleteAllByThreadID(threadID)
	if result != nil {
		return errors.New("failed to delete bookmark")
	}

	return nil
}
