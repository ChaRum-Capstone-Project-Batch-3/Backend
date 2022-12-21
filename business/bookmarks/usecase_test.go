package bookmarks_test

import (
	"charum/business/bookmarks"
	_bookmarkMock "charum/business/bookmarks/mocks"
	"charum/business/threads"
	_threadMock "charum/business/threads/mocks"
	"charum/business/topics"
	_topicMock "charum/business/topics/mocks"
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	dtoThread "charum/dto/threads"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	bookmarkRepository _bookmarkMock.Repository
	threadRepository   _threadMock.Repository
	userRepository     _userMock.Repository
	topicRepository    _topicMock.Repository
	BookmarkUseCase    bookmarks.UseCase
	threadUseCase      _threadMock.UseCase
	bookmarkDomain     bookmarks.Domain
	threadDomain       threads.Domain
	userDomain         users.Domain
	topicDomain        topics.Domain
	threadResponse     dtoThread.Response
)

func TestMain(m *testing.M) {
	BookmarkUseCase = bookmarks.NewBookmarkUseCase(&bookmarkRepository, &threadRepository, &userRepository, &topicRepository, &threadUseCase)

	bookmarkDomain = bookmarks.Domain{
		Id:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ThreadID:  primitive.NewObjectID(),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	threadDomain = threads.Domain{
		Id:            bookmarkDomain.ThreadID,
		TopicID:       primitive.NewObjectID(),
		CreatorID:     bookmarkDomain.UserID,
		Title:         "Thread Title",
		Description:   "Thread Description",
		Likes:         []threads.Like{},
		SuspendStatus: "Test",
		SuspendDetail: "Test",
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
	}

	userDomain = users.Domain{
		Id:          bookmarkDomain.UserID,
		UserName:    "Test",
		DisplayName: "Test",
		Biodata:     "Test",
		SocialMedia: "Test",
		Email:       "Test",
		Password:    "Test",
		Role:        "user",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	topicDomain = topics.Domain{
		Id:          threadDomain.TopicID,
		Topic:       "Test",
		Description: "Test",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	threadResponse = dtoThread.Response{
		Id:            threadDomain.Id,
		Topic:         topicDomain,
		Creator:       userDomain,
		Title:         threadDomain.Title,
		Description:   threadDomain.Description,
		Likes:         []dtoThread.Like{},
		SuspendStatus: threadDomain.SuspendStatus,
		SuspendDetail: threadDomain.SuspendDetail,
		TotalLike:     len(threadDomain.Likes),
		TotalFollow:   0,
		TotalComment:  0,
		TotalBookmark: 0,
		CreatedAt:     threadDomain.CreatedAt,
		UpdatedAt:     threadDomain.UpdatedAt,
	}

	m.Run()
}

func TestAddBookmark(t *testing.T) {
	t.Run("Test Case 1 | Valid Add Bookmark", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("GetByUserIDAndThreadID", bookmarkDomain.UserID, bookmarkDomain.ThreadID).Return(bookmarkDomain, errors.New("not found")).Once()
		bookmarkRepository.On("Create", mock.Anything).Return(bookmarkDomain, nil).Once()

		_, err := BookmarkUseCase.Create(&bookmarkDomain)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Add Bookmark | Thread Not Found", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.Create(&bookmarkDomain)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 3 | Invalid Add Bookmark | Bookmark Already Exist", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("GetByUserIDAndThreadID", bookmarkDomain.UserID, bookmarkDomain.ThreadID).Return(bookmarkDomain, nil).Once()

		_, err := BookmarkUseCase.Create(&bookmarkDomain)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 4 | Invalid Add Bookmark | Repository Error", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("GetByUserIDAndThreadID", bookmarkDomain.UserID, bookmarkDomain.ThreadID).Return(bookmarkDomain, errors.New("not found")).Once()
		bookmarkRepository.On("Create", mock.Anything).Return(bookmarkDomain, errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.Create(&bookmarkDomain)
		assert.NotNil(t, err)
	})
}

func TestGetAllByUserID(t *testing.T) {
	t.Run("Test Case 1 | Valid Get All By User ID", func(t *testing.T) {
		bookmarkRepository.On("GetAllByUserID", bookmarkDomain.UserID).Return([]bookmarks.Domain{bookmarkDomain}, nil).Once()

		_, err := BookmarkUseCase.GetAllByUserID(bookmarkDomain.UserID)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Get All By User ID | Repository Error", func(t *testing.T) {
		bookmarkRepository.On("GetAllByUserID", bookmarkDomain.UserID).Return([]bookmarks.Domain{bookmarkDomain}, errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.GetAllByUserID(bookmarkDomain.UserID)
		assert.NotNil(t, err)
	})
}

func TestCountByThreadID(t *testing.T) {
	t.Run("Test Case 1 | Valid Count By Thread ID", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("CountByThreadID", bookmarkDomain.ThreadID).Return(1, nil).Once()

		_, err := BookmarkUseCase.CountByThreadID(bookmarkDomain.ThreadID)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Count By Thread ID | Thread Not Found", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.CountByThreadID(bookmarkDomain.ThreadID)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 3 | Invalid Count By Thread ID | Repository Error", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("CountByThreadID", bookmarkDomain.ThreadID).Return(1, errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.CountByThreadID(bookmarkDomain.ThreadID)
		assert.NotNil(t, err)
	})
}

func TestCheckBookmarkedThread(t *testing.T) {
	t.Run("Test Case 1 | Valid Check Bookmarked Thread", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("GetByUserIDAndThreadID", bookmarkDomain.UserID, bookmarkDomain.ThreadID).Return(bookmarkDomain, nil).Once()

		_, err := BookmarkUseCase.CheckBookmarkedThread(bookmarkDomain.UserID, bookmarkDomain.ThreadID)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Check Bookmarked Thread | Thread Not Found", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.CheckBookmarkedThread(bookmarkDomain.UserID, bookmarkDomain.ThreadID)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 3 | Invalid Check Bookmarked Thread | Bookmark Not Found", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("GetByUserIDAndThreadID", bookmarkDomain.UserID, bookmarkDomain.ThreadID).Return(bookmarkDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.CheckBookmarkedThread(bookmarkDomain.UserID, bookmarkDomain.ThreadID)
		assert.Nil(t, err)
	})

	t.Run("Test Case 4 | Invalid Check Bookmarked Thread | Repository Error", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		bookmarkRepository.On("GetByUserIDAndThreadID", bookmarkDomain.UserID, bookmarkDomain.ThreadID).Return(bookmarkDomain, errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.CheckBookmarkedThread(bookmarkDomain.UserID, bookmarkDomain.ThreadID)
		assert.Nil(t, err)
	})
}

func TestDomainToResponse(t *testing.T) {
	t.Run("Test Case 1 | Valid Domain To Response", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCase.On("DomainToResponse", threadDomain, bookmarkDomain.UserID).Return(threadResponse, nil).Once()

		_, err := BookmarkUseCase.DomainToResponse(bookmarkDomain, userDomain.Id)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Domain To Response | Thread Not Found", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.DomainToResponse(bookmarkDomain, userDomain.Id)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 3 | Invalid Domain To Response | Repository Error", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCase.On("DomainToResponse", threadDomain, primitive.NilObjectID).Return(threadResponse, errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.DomainToResponse(bookmarkDomain, primitive.NilObjectID)
		assert.NotNil(t, err)
	})
}

func TestDomainsToResponseArray(t *testing.T) {
	t.Run("Test Case 1 | Valid Domains To Response Array", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCase.On("DomainToResponse", threadDomain, bookmarkDomain.UserID).Return(threadResponse, nil).Once()

		_, err := BookmarkUseCase.DomainsToResponseArray([]bookmarks.Domain{bookmarkDomain}, userDomain.Id)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Domains To Response Array | Thread Not Found", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.DomainsToResponseArray([]bookmarks.Domain{bookmarkDomain}, primitive.NilObjectID)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 3 | Invalid Domains To Response Array | Repository Error", func(t *testing.T) {
		threadRepository.On("GetByID", bookmarkDomain.ThreadID).Return(threadDomain, nil).Once()
		threadUseCase.On("DomainToResponse", threadDomain, primitive.NilObjectID).Return(threadResponse, errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.DomainsToResponseArray([]bookmarks.Domain{bookmarkDomain}, primitive.NilObjectID)
		assert.NotNil(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete", func(t *testing.T) {
		bookmarkRepository.On("GetByUserIDAndThreadID", mock.Anything, mock.Anything).Return(bookmarkDomain, nil).Once()
		bookmarkRepository.On("Delete", mock.Anything).Return(nil).Once()

		result, err := BookmarkUseCase.Delete(&bookmarkDomain)
		assert.Nil(t, err)
		assert.Equal(t, result, bookmarkDomain)
	})

	t.Run("Test Case 2 | Invalid Delete | Bookmark Not Found", func(t *testing.T) {
		bookmarkRepository.On("GetByUserIDAndThreadID", mock.Anything, mock.Anything).Return(bookmarkDomain, errors.New("not found")).Once()

		_, err := BookmarkUseCase.Delete(&bookmarkDomain)
		assert.NotNil(t, err)
	})

	t.Run("Test Case 3 | Invalid Delete | Repository Error", func(t *testing.T) {
		bookmarkRepository.On("GetByUserIDAndThreadID", mock.Anything, mock.Anything).Return(bookmarkDomain, nil).Once()
		bookmarkRepository.On("Delete", mock.Anything).Return(errors.New("unexpected error")).Once()

		_, err := BookmarkUseCase.Delete(&bookmarkDomain)
		assert.NotNil(t, err)
	})
}

func TestDeleteAllByUserID(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete All By User ID", func(t *testing.T) {
		bookmarkRepository.On("DeleteAllByUserID", mock.Anything).Return(nil).Once()

		err := BookmarkUseCase.DeleteAllByUserID(bookmarkDomain.UserID)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Delete All By User ID | Repository Error", func(t *testing.T) {
		bookmarkRepository.On("DeleteAllByUserID", mock.Anything).Return(errors.New("unexpected error")).Once()

		err := BookmarkUseCase.DeleteAllByUserID(bookmarkDomain.UserID)
		assert.NotNil(t, err)
	})
}

func TestDeleteAllByThreadID(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete All By Thread ID", func(t *testing.T) {
		bookmarkRepository.On("DeleteAllByThreadID", mock.Anything).Return(nil).Once()

		err := BookmarkUseCase.DeleteAllByThreadID(bookmarkDomain.ThreadID)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Delete All By Thread ID | Repository Error", func(t *testing.T) {
		bookmarkRepository.On("DeleteAllByThreadID", mock.Anything).Return(errors.New("unexpected error")).Once()

		err := BookmarkUseCase.DeleteAllByThreadID(bookmarkDomain.ThreadID)
		assert.NotNil(t, err)
	})
}
