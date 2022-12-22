package reports_test

import (
	"charum/business/reports"
	_ReportMock "charum/business/reports/mocks"
	"charum/business/threads"
	_threadMock "charum/business/threads/mocks"
	"charum/business/topics"
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
	ReportUseCase    reports.UseCase
	ReportRepository _ReportMock.Repository
	ThreadRepository _threadMock.Repository
	UserRepository   _userMock.Repository
	userDomain       users.Domain
	threadDomain     threads.Domain
	reportDomain     reports.Domain
	threadResponse   dtoThread.Response
	topicDomain      topics.Domain
)

func TestMain(m *testing.M) {
	ReportUseCase = reports.NewReportUseCase(&ReportRepository, &UserRepository, &ThreadRepository)

	reportDomain = reports.Domain{
		Id:           primitive.NewObjectID(),
		UserID:       primitive.NewObjectID(),
		ReportedID:   primitive.NewObjectID(),
		ReportedType: "",
		ReportDetail: "Inappropriate content or behavior",
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
	}

	userDomain = users.Domain{
		Id:          reportDomain.UserID,
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
		Id:          primitive.NewObjectID(),
		Topic:       "Test",
		Description: "Test",
		ImageURL:    "url",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	threadDomain = threads.Domain{
		Id:            reportDomain.ReportedID,
		TopicID:       topicDomain.Id,
		CreatorID:     userDomain.Id,
		Title:         "Thread Title",
		Description:   "Thread Description",
		Likes:         []threads.Like{},
		SuspendStatus: "Test",
		SuspendDetail: "Test",
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
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

// get all report
func TestCreate(t *testing.T) {

	t.Run("Test Case 1 | Valid Create Report Thread", func(t *testing.T) {
		reportDomain.ReportedType = "thread"
		ReportRepository.On("Create", &reportDomain).Return(reportDomain, nil).Once()
		ReportRepository.On("GetByID", reportDomain.Id).Return(reportDomain, nil).Once()
		ReportRepository.On("GetByReportedID", reportDomain.ReportedID).Return(reportDomain, nil).Once()
		ThreadRepository.On("GetByID", mock.Anything).Return(threadDomain, nil).Once()
		UserRepository.On("GetByID", mock.Anything).Return(userDomain, nil).Once()
		ReportRepository.On("CheckByUserID", mock.Anything, mock.Anything).Return(reports.Domain{}, errors.New("not reported")).Once()

		_, err := ReportUseCase.Create(&reportDomain)

		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Valid Create Report User", func(t *testing.T) {
		reportDomain.ReportedType = "user"
		ReportRepository.On("Create", &reportDomain).Return(reportDomain, nil).Once()
		ReportRepository.On("GetByID", reportDomain.Id).Return(reportDomain, nil).Once()
		ReportRepository.On("GetByReportedID", reportDomain.ReportedID).Return(reportDomain, nil).Once()
		ThreadRepository.On("GetByID", mock.Anything).Return(threadDomain, nil).Once()
		UserRepository.On("GetByID", mock.Anything).Return(userDomain, nil).Once()
		ReportRepository.On("CheckByUserID", mock.Anything, mock.Anything).Return(reports.Domain{}, errors.New("not reported")).Once()

		_, err := ReportUseCase.Create(&reportDomain)

		assert.Nil(t, err)
	})

	t.Run("Test Case 3 | Invalid Create Report - already reported", func(t *testing.T) {
		expectedErr := errors.New("already reported")
		// check if the id is user id or thread id, if user then return "user", and if thread then return "thread"
		ReportRepository.On("Create", &reportDomain).Return(reportDomain, nil).Once()
		ReportRepository.On("GetByID", reportDomain.Id).Return(reportDomain, nil).Once()
		ReportRepository.On("GetByReportedID", reportDomain.ReportedID).Return(reportDomain, nil).Once()
		ThreadRepository.On("GetByID", mock.Anything).Return(threadDomain, nil).Once()
		UserRepository.On("GetByID", mock.Anything).Return(userDomain, nil).Once()
		ReportRepository.On("CheckByUserID", mock.Anything, mock.Anything).Return(reportDomain, nil).Once()

		_, err := ReportUseCase.Create(&reportDomain)

		assert.Equal(t, expectedErr, err)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Test Case 1 | Valid Get All Report", func(t *testing.T) {
		ReportRepository.On("GetAll").Return([]reports.Domain{reportDomain}, nil).Once()
		_, err := ReportUseCase.GetAll()

		assert.Nil(t, err)
	})
	t.Run("Test Case 2 | Invalid Get All Report", func(t *testing.T) {
		expectedErr := errors.New("failed to get reports")
		ReportRepository.On("GetAll").Return([]reports.Domain{}, expectedErr).Once()
		res, err := ReportUseCase.GetAll()

		assert.Equal(t, 0, res)
		assert.Equal(t, err, expectedErr)
	})
}

// get all reported threads
func TestGetAllReportedThreads(t *testing.T) {
	t.Run("Test Case 1 | Valid Get All Reported Threads", func(t *testing.T) {
		ReportRepository.On("GetAllReportedThreads").Return([]reports.Domain{reportDomain}, nil).Once()
		_, err := ReportUseCase.GetAllReportedThreads()

		assert.Nil(t, err)
	})
	t.Run("Test Case 2 | Invalid Get All Reported Threads", func(t *testing.T) {
		expectedErr := errors.New("failed to get reports")
		ReportRepository.On("GetAllReportedThreads").Return([]reports.Domain{}, expectedErr).Once()
		res, err := ReportUseCase.GetAllReportedThreads()

		assert.Equal(t, 0, res)
		assert.Equal(t, err, expectedErr)
	})
}

// get all reported users
func TestGetAllReportedUsers(t *testing.T) {
	t.Run("Test Case 1 | Valid Get All Reported Users", func(t *testing.T) {
		ReportRepository.On("GetAllReportedUsers").Return([]reports.Domain{reportDomain}, nil).Once()
		_, err := ReportUseCase.GetAllReportedUsers()

		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Get All Reported Users", func(t *testing.T) {
		expectedErr := errors.New("failed to get reports")
		ReportRepository.On("GetAllReportedUsers").Return([]reports.Domain{}, expectedErr).Once()
		res, err := ReportUseCase.GetAllReportedUsers()

		assert.Equal(t, 0, res)
		assert.Equal(t, err, expectedErr)
	})
}
