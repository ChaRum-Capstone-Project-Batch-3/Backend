package bookmarks

import (
	"charum/business/threads"
	"charum/business/users"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookmarkUseCase struct {
	bookmarkRepository Repository
	threadRepository   threads.Repository
	userRepository     users.Repository
}

func NewBookmarkUseCase(br Repository, tr threads.Repository, ur users.Repository) UseCase {
	return &BookmarkUseCase{
		bookmarkRepository: br,
		threadRepository:   tr,
		userRepository:     ur,
	}
}

/*
Add Bookmark
*/
func (bu *BookmarkUseCase) AddBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error) {
	// check user already bookmarked or not
	thread, err := bu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	domain.Id = primitive.NewObjectID()
	domain.UserID = userID
	domain.Threads = append(domain.Threads, Thread{
		Id:    thread.Id,
		Title: thread.Title,
	})
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	bookmark, err := bu.bookmarkRepository.AddBookmark(domain)
	if err != nil {
		return Domain{}, errors.New("failed to add thread to bookmark")
	}

	return bookmark, nil
}

/*
Read
*/

func (bu *BookmarkUseCase) CheckBookmark(userID primitive.ObjectID) (bool, error) {
	_, err := bu.bookmarkRepository.GetByID(userID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// get bookmark by user id
func (bu *BookmarkUseCase) GetByID(userID primitive.ObjectID) (Domain, error) {
	result, err := bu.bookmarkRepository.GetByID(userID)
	if err != nil {
		return Domain{}, errors.New("failed to get bookmark")
	}
	return result, nil
}
