package bookmarks

import (
	"charum/business/threads"
	"charum/business/topics"
	"charum/business/users"
	"charum/dto"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarkUseCase struct {
	bookmarkRepository Repository
	threadRepository   threads.Repository
	userRepository     users.Repository
	topicRepository    topics.Repository
}

func NewBookmarkUseCase(br Repository, tr threads.Repository, ur users.Repository, tc topics.Repository) UseCase {
	return &BookmarkUseCase{
		bookmarkRepository: br,
		threadRepository:   tr,
		userRepository:     ur,
		topicRepository:    tc,
	}
}

/*
Add Bookmark
*/
func (bu *BookmarkUseCase) AddBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error) {
	// check thread is exist or not
	thread, err := bu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	// check the thread already bookmarked or not
	checkBookmark, err := bu.bookmarkRepository.GetByID(userID, threadID)
	if err != nil {
		return Domain{}, err
	}
	if checkBookmark.ThreadID == thread.Id && checkBookmark.UserID == userID {
		return Domain{}, errors.New("thread already bookmarked")
	}

	domain.Id = primitive.NewObjectID()
	domain.UserID = userID
	domain.ThreadID = thread.Id
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	bookmark, err := bu.bookmarkRepository.AddBookmark(domain)
	fmt.Println("bookmark error:", err)
	if err != nil {
		return Domain{}, errors.New("failed to add thread to bookmark")
	}

	return bookmark, nil
}

/*
Read
*/

// get bookmark by user id
func (bu *BookmarkUseCase) GetByID(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error) {
	result, err := bu.bookmarkRepository.GetByID(userID, threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get bookmark")
	}
	return result, nil
}

// get all bookmark by user id
func (bu *BookmarkUseCase) GetAllBookmark(userID primitive.ObjectID) ([]Domain, error) {
	// get bookmark data by user id, and get thread data by thread id and return it with array
	result, err := bu.bookmarkRepository.GetAllBookmark(userID)
	if err != nil {
		return []Domain{}, errors.New("failed to get all bookmark")
	}

	return result, nil
}

func (bu *BookmarkUseCase) DomainToResponse(domain Domain) (dto.ResponseBookmark, error) {
	// get thread data by thread id
	thread, err := bu.threadRepository.GetByID(domain.ThreadID)
	if err != nil {
		return dto.ResponseBookmark{}, errors.New("failed to get thread")
	}
	creator, err := bu.userRepository.GetByID(thread.CreatorID)
	if err != nil {
		return dto.ResponseBookmark{}, errors.New("failed to get creator")
	}

	likes := []dto.Like{}
	for _, like := range thread.Likes {
		user, err := bu.userRepository.GetByID(like.UserID)
		if err != nil {
			return dto.ResponseBookmark{}, err
		}

		likes = append(likes, dto.Like{
			User:      user,
			CreatedAt: domain.CreatedAt,
		})
	}

	topic, err := bu.topicRepository.GetByID(thread.TopicID)
	if err != nil {
		return dto.ResponseBookmark{}, errors.New("failed to get topic")
	}

	return dto.ResponseBookmark{
		Id:            domain.Id,
		ThreadId:      thread.Id,
		Topic:         topic,
		Creator:       creator,
		Title:         thread.Title,
		Description:   thread.Description,
		Likes:         likes,
		TotalLike:     len(thread.Likes),
		SuspendStatus: thread.SuspendStatus,
		SuspendDetail: thread.SuspendDetail,
		CreatedAt:     thread.CreatedAt,
		UpdatedAt:     thread.UpdatedAt,
	}, nil
}

func (bu *BookmarkUseCase) DomainsToResponseArray(domains []Domain) ([]dto.ResponseBookmark, error) {
	var responses []dto.ResponseBookmark
	for _, domain := range domains {
		response, err := bu.DomainToResponse(domain)
		if err != nil {
			return []dto.ResponseBookmark{}, errors.New("failed to get thread")
		}

		responses = append(responses, response)
	}

	return responses, nil
}

/*
Update
*/

func (bu *BookmarkUseCase) UpdateBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *Domain) (Domain, error) {
	// check user already bookmarked or not
	thread, err := bu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	bookmark, err := bu.bookmarkRepository.GetByID(userID, threadID)
	if err != nil {
		return Domain{}, err
	}

	if bookmark.ThreadID == thread.Id && bookmark.UserID == userID {
		return Domain{}, errors.New("thread already bookmarked")
	}

	bookmark.Id = domain.Id
	// append bookmark threads without remove the previous threads
	bookmark.ThreadID = thread.Id
	bookmark.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	bookmark, err = bu.bookmarkRepository.UpdateBookmark(userID, threadID, &bookmark)
	if err != nil {
		return Domain{}, err
	}

	return bookmark, nil
}

/*
Delete
*/

// delete bookmark by bookmark id
func (bu *BookmarkUseCase) DeleteBookmark(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error) {
	bookmark, err := bu.bookmarkRepository.GetByID(userID, threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get bookmark")
	}

	result := bu.bookmarkRepository.DeleteBookmark(userID, bookmark.ThreadID)
	if result != nil {
		return Domain{}, errors.New("failed to delete bookmark")
	}
	return bookmark, nil
}
