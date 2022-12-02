package bookmarks

import (
	"charum/business/threads"
	"charum/business/users"
	"errors"
	"fmt"
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

	user, err := bu.userRepository.GetByID(userID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	bookmark, err := bu.bookmarkRepository.GetByID(userID)
	if err == nil {
		// loop check if bookmark already exist or not, if not append thread id
		for _, v := range bookmark.Threads {
			if v == thread.Id {
				return Domain{}, errors.New("thread already bookmarked")
			} else {
				return bu.UpdateBookmark(userID, threadID, domain)
			}
		}
	}
	domain.Id = user.BookmarkID
	domain.UserID = userID
	domain.Threads = append(domain.Threads, thread.Id)
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	bookmark, err = bu.bookmarkRepository.AddBookmark(domain)
	fmt.Println("bookmark", bookmark, err)
	if err != nil {
		return Domain{}, err
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
		return Domain{}, err
	}
	return result, nil
}

// update bookmark
func (bu *BookmarkUseCase) UpdateBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error) {
	// check user already bookmarked or not
	thread, err := bu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	user, err := bu.userRepository.GetByID(userID)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	bookmark, err := bu.bookmarkRepository.GetByID(userID)
	if err != nil {
		return Domain{}, err
	}

	bookmark.Id = user.BookmarkID
	bookmark.UserID = userID
	// append bookmark threads without remove the previous threads
	bookmark.Threads = append(bookmark.Threads, thread.Id)
	bookmark.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	bookmark, err = bu.bookmarkRepository.UpdateBookmark(&bookmark)
	if err != nil {
		return Domain{}, err
	}

	return bookmark, nil
}

// return all bookmarked threads by user id
func (bu *BookmarkUseCase) GetAllBookmark(userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	bookmark, err := bu.bookmarkRepository.GetByID(userID)
	if err != nil {
		return []primitive.ObjectID{}, err
	}
	var result []primitive.ObjectID
	for _, v := range bookmark.Threads {
		thread, err := bu.threadRepository.GetByID(v)
		if err != nil {
			return []primitive.ObjectID{}, err
		}
		result = append(result, thread.Id)

	}
	return result, nil
}
