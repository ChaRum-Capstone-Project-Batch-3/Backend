package forgot_password_test

import (
	"charum/business/forgot_password"
	_forgotPassMock "charum/business/forgot_password/mocks"
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	_mailgunMock "charum/helper/mailgun/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	forgotPasswordRepository _forgotPassMock.Repository
	forgotPasswordUseCase    forgot_password.UseCase
	userRepository           _userMock.Repository
	mailgun                  _mailgunMock.Function
	forgotPasswordDomain     forgot_password.Domain
	userDomain               users.Domain
)

func TestMain(m *testing.M) {
	forgotPasswordUseCase = forgot_password.NewForgotPasswordUseCase(&forgotPasswordRepository, &userRepository, &mailgun)

	userDomain = users.Domain{
		Id:          primitive.NewObjectID(),
		UserName:    "Test",
		DisplayName: "Test",
		Biodata:     "Test",
		SocialMedia: "Test",
		Email:       "Test@mail.com",
		Password:    "Test",
		Role:        "user",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	m.Run()
}

func TestGenerate(t *testing.T) {
	t.Run("Test Case 1 | Valid Generate", func(t *testing.T) {
		userRepository.On("GetByEmail", forgotPasswordDomain.Email).Return(userDomain, nil).Once()
		forgotPasswordRepository.On("Generate", &forgotPasswordDomain).Return(forgotPasswordDomain, nil).Once()
		mailgun.On("SendMail", mock.Anything, mock.Anything).Return("", nil).Once()
		_, err := forgotPasswordUseCase.Generate(&forgotPasswordDomain)

		assert.Nil(t, err)
	})
}

func TestValidate(t *testing.T) {
	forgotPasswordDomain = forgot_password.Domain{
		Id:        primitive.NewObjectID(),
		Email:     "test@mail.com",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		ExpiredAt: primitive.NewDateTimeFromTime(time.Now().Add(time.Minute * 30)),
		IsUsed:    false,
	}
	t.Run("Test Case 1 | Valid Validate", func(t *testing.T) {
		userRepository.On("GetByEmail", forgotPasswordDomain.Email).Return(userDomain, nil).Once()
		forgotPasswordRepository.On("Generate", &forgotPasswordDomain).Return(forgotPasswordDomain, nil).Once()
		forgotPasswordRepository.On("GetByToken", forgotPasswordDomain.Token).Return(forgotPasswordDomain, nil).Once()
		_, err := forgotPasswordUseCase.ValidateToken(forgotPasswordDomain.Token)

		assert.Nil(t, err)
	})

	// create test that token is not valid
	t.Run("Test Case 2 | Token Used ", func(t *testing.T) {
		userRepository.On("GetByEmail", forgotPasswordDomain.Email).Return(users.Domain{}, nil).Once()
		forgotPasswordRepository.On("Generate", &forgotPasswordDomain).Return(forgot_password.Domain{}, nil).Once()
		forgotPasswordRepository.On("GetByToken", forgotPasswordDomain.Token).Return(forgot_password.Domain{}, nil).Once()
		forgotPasswordDomain.IsUsed = true
		_, err := forgotPasswordUseCase.ValidateToken(forgotPasswordDomain.Token)
		assert.NotNil(t, err)
	})

	// create test that token is expired
	t.Run("Test Case 3 | Token Expired ", func(t *testing.T) {
		userRepository.On("GetByEmail", forgotPasswordDomain.Email).Return(users.Domain{}, nil).Once()
		forgotPasswordRepository.On("Generate", &forgotPasswordDomain).Return(forgot_password.Domain{}, nil).Once()
		forgotPasswordRepository.On("GetByToken", forgotPasswordDomain.Token).Return(forgot_password.Domain{}, nil).Once()
		forgotPasswordDomain.ExpiredAt = primitive.NewDateTimeFromTime(time.Now().Add(time.Minute * -50))
		_, err := forgotPasswordUseCase.ValidateToken(forgotPasswordDomain.Token)
		assert.NotNil(t, err)
	})

	// create test that token is not found
	t.Run("Test Case 4 | Token Not Found ", func(t *testing.T) {
		userRepository.On("GetByEmail", forgotPasswordDomain.Email).Return(users.Domain{}, nil).Once()
		forgotPasswordRepository.On("Generate", &forgotPasswordDomain).Return(forgot_password.Domain{}, nil).Once()
		forgotPasswordRepository.On("GetByToken", forgotPasswordDomain.Token).Return(forgot_password.Domain{}, nil).Once()
		forgotPasswordDomain.Token = ""
		_, err := forgotPasswordUseCase.ValidateToken(forgotPasswordDomain.Token)
		assert.NotNil(t, err)
	})
}
