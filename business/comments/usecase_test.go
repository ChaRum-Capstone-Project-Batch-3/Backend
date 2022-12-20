package comments_test

import (
	"charum/business/comments"
	_commentMock "charum/business/comments/mocks"
	"charum/business/threads"
	_threadMock "charum/business/threads/mocks"
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	_cloudinaryMock "charum/helper/cloudinary/mocks"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	commentRepository    _commentMock.Repository
	threadRepository     _threadMock.Repository
	userRepository       _userMock.Repository
	cloudinaryRepository _cloudinaryMock.Function
	commentUseCase       comments.UseCase
	commentDomain        comments.Domain
	threadDomain         threads.Domain
	userDomain           users.Domain
	image                *multipart.FileHeader
)

func TestMain(m *testing.M) {
	commentUseCase = comments.NewCommentUseCase(&commentRepository, &threadRepository, &userRepository, &cloudinaryRepository)

	userDomain = users.Domain{
		Id:          primitive.NewObjectID(),
		Email:       "test@test.com",
		UserName:    "test",
		DisplayName: "test",
		Password:    "test",
		IsActive:    true,
		Role:        "user",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	threadDomain = threads.Domain{
		Id:            commentDomain.ThreadID,
		TopicID:       primitive.NewObjectID(),
		CreatorID:     userDomain.Id,
		Title:         "test",
		Description:   "test",
		Likes:         []threads.Like{},
		SuspendStatus: "",
		SuspendDetail: "",
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
	}

	commentDomain = comments.Domain{
		Id:        primitive.NewObjectID(),
		ThreadID:  threadDomain.Id,
		UserID:    userDomain.Id,
		ParentID:  primitive.NewObjectID(),
		Comment:   "test",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	image = &multipart.FileHeader{}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test case 1 | Valid create", func(t *testing.T) {
		commentRepository.On("GetByIDAndThreadID", mock.Anything, mock.Anything).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("test", nil).Once()
		commentRepository.On("Create", mock.Anything).Return(commentDomain, nil).Once()

		actualComment, err := commentUseCase.Create(&commentDomain, image)

		assert.Nil(t, err)
		assert.NotEmpty(t, actualComment)
	})

	t.Run("Test case 2 | Invalid create | Failed To Get Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to get comment")
		commentRepository.On("GetByIDAndThreadID", mock.Anything, mock.Anything).Return(comments.Domain{}, expectedErr).Once()

		actualComment, err := commentUseCase.Create(&commentDomain, image)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})

	t.Run("Test case 3 | Invalid create | Failed To Get Thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		commentRepository.On("GetByIDAndThreadID", mock.Anything, mock.Anything).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threads.Domain{}, expectedErr).Once()

		actualComment, err := commentUseCase.Create(&commentDomain, image)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})

	t.Run("Test case 4 | Invalid create | Failed To Upload Image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		commentRepository.On("GetByIDAndThreadID", mock.Anything, mock.Anything).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		actualComment, err := commentUseCase.Create(&commentDomain, image)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})

	t.Run("Test case 5 | Invalid create | Failed To Create Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to create comment")
		commentRepository.On("GetByIDAndThreadID", mock.Anything, mock.Anything).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("test", nil).Once()
		commentRepository.On("Create", mock.Anything).Return(comments.Domain{}, expectedErr).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

		actualComment, err := commentUseCase.Create(&commentDomain, image)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})

	t.Run("Test case 6 | Invalid create | Failed To Delete Image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		commentRepository.On("GetByIDAndThreadID", mock.Anything, mock.Anything).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("test", nil).Once()
		commentRepository.On("Create", mock.Anything).Return(commentDomain, expectedErr).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		actualComment, err := commentUseCase.Create(&commentDomain, image)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})
}

func TestGetByThreadID(t *testing.T) {
	t.Run("Test case 1 | Valid get by thread id", func(t *testing.T) {
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		commentRepository.On("GetByThreadID", commentDomain.ThreadID).Return([]comments.Domain{commentDomain}, nil).Once()

		actualComment, err := commentUseCase.GetByThreadID(commentDomain.ThreadID)

		assert.Nil(t, err)
		assert.NotEmpty(t, actualComment)
	})

	t.Run("Test case 2 | Invalid get by thread id | Failed To Get Thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threads.Domain{}, expectedErr).Once()

		actualComment, err := commentUseCase.GetByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})

	t.Run("Test case 3 | Invalid get by thread id | Failed To Get Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to get comment")
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		commentRepository.On("GetByThreadID", commentDomain.ThreadID).Return([]comments.Domain{}, expectedErr).Once()

		_, err := commentUseCase.GetByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
	})
}

func TestDomainToResponse(t *testing.T) {
	t.Run("Test case 1 | Valid domain to response", func(t *testing.T) {
		userRepository.On("GetByID", commentDomain.UserID).Return(userDomain, nil).Once()

		actualComment, err := commentUseCase.DomainToResponse(commentDomain)

		assert.NotEmpty(t, actualComment)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid domain to response | Failed To Get User", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", commentDomain.UserID).Return(users.Domain{}, expectedErr).Once()

		_, err := commentUseCase.DomainToResponse(commentDomain)
		assert.NotNil(t, err)
	})
}

func TestDomainToResponseArray(t *testing.T) {
	t.Run("Test case 1 | Valid domain to response array", func(t *testing.T) {
		userRepository.On("GetByID", commentDomain.UserID).Return(userDomain, nil).Once()

		actualComment, err := commentUseCase.DomainToResponseArray([]comments.Domain{commentDomain})

		assert.NotEmpty(t, actualComment)
		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid domain to response array | Failed To Get User", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", commentDomain.UserID).Return(users.Domain{}, expectedErr).Once()

		_, err := commentUseCase.DomainToResponseArray([]comments.Domain{commentDomain})
		assert.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test case 1 | Valid update", func(t *testing.T) {
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("test", nil).Once()
		commentRepository.On("Update", mock.Anything).Return(commentDomain, nil).Once()

		actualComment, err := commentUseCase.Update(&commentDomain, image)

		assert.Nil(t, err)
		assert.NotEmpty(t, actualComment)
	})

	t.Run("Test case 2 | Invalid update | Failed To Get Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to get comment")
		commentRepository.On("GetByID", commentDomain.Id).Return(comments.Domain{}, expectedErr).Once()

		_, err := commentUseCase.Update(&commentDomain, image)
		assert.NotNil(t, err)
	})

	t.Run("Test case 3 | Invalid update | Failed To Get Thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threads.Domain{}, expectedErr).Once()

		_, err := commentUseCase.Update(&commentDomain, image)
		assert.NotNil(t, err)
	})

	t.Run("Test case 4 | Invalid update | Failed To Update Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to update comment")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("test", nil).Once()
		commentRepository.On("Update", mock.Anything).Return(comments.Domain{}, expectedErr).Once()

		_, err := commentUseCase.Update(&commentDomain, image)
		assert.NotNil(t, err)
	})

	t.Run("Test case 5 | Invalid update | User Are Not The Owner Of This Comment", func(t *testing.T) {
		copyDomain := commentDomain
		copyDomain.UserID = primitive.NewObjectID()

		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()

		_, err := commentUseCase.Update(&copyDomain, image)
		assert.NotNil(t, err)
	})

	t.Run("Test case 6 | Invalid update | Failed To Delete Image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		_, err := commentUseCase.Update(&commentDomain, image)
		assert.NotNil(t, err)
	})

	t.Run("Test case 7 | Invalid update | Failed To Upload Image", func(t *testing.T) {
		expectedErr := errors.New("failed to upload image")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		_, err := commentUseCase.Update(&commentDomain, image)
		assert.NotNil(t, err)
	})
}

func TestCountByThreadID(t *testing.T) {
	t.Run("Test case 1 | Valid count by thread id", func(t *testing.T) {
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		commentRepository.On("CountByThreadID", commentDomain.ThreadID).Return(1, nil).Once()

		actualComment, err := commentUseCase.CountByThreadID(commentDomain.ThreadID)

		assert.Nil(t, err)
		assert.NotZero(t, actualComment)
	})

	t.Run("Test case 2 | Invalid count by thread id | Failed To Get Thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threads.Domain{}, expectedErr).Once()

		actualComment, err := commentUseCase.CountByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
		assert.Zero(t, actualComment)
	})

	t.Run("Test case 3 | Invalid count by thread id | Failed To Count Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to count comment")
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		commentRepository.On("CountByThreadID", commentDomain.ThreadID).Return(0, expectedErr).Once()

		_, err := commentUseCase.CountByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test case 1 | Valid delete", func(t *testing.T) {
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		commentRepository.On("Delete", commentDomain.Id).Return(nil).Once()

		actaulComment, err := commentUseCase.Delete(commentDomain.Id, commentDomain.UserID)

		assert.Nil(t, err)
		assert.NotEmpty(t, actaulComment)
	})

	t.Run("Test case 2 | Invalid delete | Failed To Get Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to get comment")
		commentRepository.On("GetByID", commentDomain.Id).Return(comments.Domain{}, expectedErr).Once()

		actualComment, err := commentUseCase.Delete(commentDomain.Id, commentDomain.UserID)
		assert.NotNil(t, err)
		assert.Empty(t, actualComment)
	})

	t.Run("Test case 3 | Invalid delete | Failed To Get Thread", func(t *testing.T) {
		expectedErr := errors.New("failed to get thread")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threads.Domain{}, expectedErr).Once()

		_, err := commentUseCase.Delete(commentDomain.Id, commentDomain.UserID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 4 | Invalid delete | Failed To Delete Comment", func(t *testing.T) {
		expectedErr := errors.New("failed to delete comment")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		commentRepository.On("Delete", commentDomain.Id).Return(expectedErr).Once()

		_, err := commentUseCase.Delete(commentDomain.Id, commentDomain.UserID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 5 | Invalid delete | You Are Not The Owner Of This Comment", func(t *testing.T) {
		copyComment := commentDomain
		copyComment.UserID = primitive.NewObjectID()
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()

		_, err := commentUseCase.Delete(commentDomain.Id, copyComment.UserID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 6 | Invalid delete | Failed To Delete Image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		commentRepository.On("GetByID", commentDomain.Id).Return(commentDomain, nil).Once()
		threadRepository.On("GetByID", commentDomain.ThreadID).Return(threadDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		_, err := commentUseCase.Delete(commentDomain.Id, commentDomain.UserID)
		assert.NotNil(t, err)
	})
}

func TestDeleteAllByUserID(t *testing.T) {
	t.Run("Test case 1 | Valid delete all by user id", func(t *testing.T) {
		commentRepository.On("GetAllByUserID", commentDomain.UserID).Return([]comments.Domain{commentDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		commentRepository.On("DeleteAllByUserID", commentDomain.UserID).Return(nil).Once()

		err := commentUseCase.DeleteAllByUserID(commentDomain.UserID)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid delete all by user id | Failed To Delete All Comment By User ID", func(t *testing.T) {
		expectedErr := errors.New("failed to delete all comment by user id")
		commentRepository.On("GetAllByUserID", commentDomain.UserID).Return([]comments.Domain{commentDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		commentRepository.On("DeleteAllByUserID", commentDomain.UserID).Return(expectedErr).Once()

		err := commentUseCase.DeleteAllByUserID(commentDomain.UserID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 3 | Invalid delete all by user id | Failed To Get All Comment By User ID", func(t *testing.T) {
		expectedErr := errors.New("failed to get all comment by user id")
		commentRepository.On("GetAllByUserID", commentDomain.UserID).Return([]comments.Domain{}, expectedErr).Once()

		err := commentUseCase.DeleteAllByUserID(commentDomain.UserID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 4 | Invalid delete all by user id | Failed To Delete Image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		commentRepository.On("GetAllByUserID", commentDomain.UserID).Return([]comments.Domain{commentDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		err := commentUseCase.DeleteAllByUserID(commentDomain.UserID)
		assert.NotNil(t, err)
	})
}

func TestDeleteALlByThreadID(t *testing.T) {
	t.Run("Test case 1 | Valid delete all by thread id", func(t *testing.T) {
		commentRepository.On("GetByThreadID", commentDomain.ThreadID).Return([]comments.Domain{commentDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		commentRepository.On("DeleteAllByThreadID", commentDomain.ThreadID).Return(nil).Once()

		err := commentUseCase.DeleteAllByThreadID(commentDomain.ThreadID)

		assert.Nil(t, err)
	})

	t.Run("Test case 2 | Invalid delete all by thread id | Failed To Delete All Comment By Thread ID", func(t *testing.T) {
		expectedErr := errors.New("failed to delete all comment by thread id")
		commentRepository.On("GetByThreadID", commentDomain.ThreadID).Return([]comments.Domain{commentDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		commentRepository.On("DeleteAllByThreadID", commentDomain.ThreadID).Return(expectedErr).Once()

		err := commentUseCase.DeleteAllByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 3 | Invalid delete all by thread id | Failed To Get All Comment By Thread ID", func(t *testing.T) {
		expectedErr := errors.New("failed to get all comment by thread id")
		commentRepository.On("GetByThreadID", commentDomain.ThreadID).Return([]comments.Domain{}, expectedErr).Once()

		err := commentUseCase.DeleteAllByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
	})

	t.Run("Test case 4 | Invalid delete all by thread id | Failed To Delete Image", func(t *testing.T) {
		expectedErr := errors.New("failed to delete image")
		commentRepository.On("GetByThreadID", commentDomain.ThreadID).Return([]comments.Domain{commentDomain}, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		err := commentUseCase.DeleteAllByThreadID(commentDomain.ThreadID)
		assert.NotNil(t, err)
	})
}
