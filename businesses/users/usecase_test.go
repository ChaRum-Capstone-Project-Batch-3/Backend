package users_test

import (
	"charum/businesses/users"
	_userMock "charum/businesses/users/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	userRepository _userMock.Repository
	userUseCase    users.UseCase
	userDomain     users.Domain
)

func TestMain(m *testing.M) {
	userUseCase = users.NewUserUseCase(&userRepository)

	userDomain = users.Domain{
		Id:        primitive.NewObjectID(),
		Email:     "test@gmail.com",
		Password:  "test123",
		UserName:  "tester",
		Role:      "user",
		IsActive:  true,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	m.Run()
}

func TestUserRegister(t *testing.T) {
	t.Run("Test Case 1 | Valid Register", func(t *testing.T) {
		userRepository.On("GetUserByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetUserByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("Create", mock.Anything).Return(userDomain, nil).Once()

		user, token, err := userUseCase.UserRegister(&userDomain)

		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Register | Email already registered", func(t *testing.T) {
		expectedErr := errors.New("email is already registered")
		userRepository.On("GetUserByEmail", userDomain.Email).Return(userDomain, nil).Once()

		user, token, err := userUseCase.UserRegister(&userDomain)

		assert.Equal(t, users.Domain{}, user)
		assert.Empty(t, token)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test Case 3 | Invalid Register | Username already used", func(t *testing.T) {
		expectedErr := errors.New("username is already used")
		userRepository.On("GetUserByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetUserByUsername", userDomain.UserName).Return(userDomain, nil).Once()

		user, token, err := userUseCase.UserRegister(&userDomain)

		assert.Equal(t, users.Domain{}, user)
		assert.Empty(t, token)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test Case 4 | Invalid Register | Error when creating user", func(t *testing.T) {
		expectedErr := errors.New("error when creating user")
		userRepository.On("GetUserByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetUserByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("Create", mock.Anything).Return(users.Domain{}, expectedErr).Once()

		user, token, err := userUseCase.UserRegister(&userDomain)

		assert.Equal(t, users.Domain{}, user)
		assert.Empty(t, token)
		assert.Equal(t, err, expectedErr)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Test Case 1 | Valid Login", func(t *testing.T) {
		copyDomain := userDomain
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(copyDomain.Password), bcrypt.DefaultCost)
		copyDomain.Password = string(encryptedPassword)

		userRepository.On("GetUserByEmail", userDomain.Email).Return(copyDomain, nil).Once()

		token, err := userUseCase.Login(&userDomain)

		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Login | Email not found", func(t *testing.T) {
		expectedErr := errors.New("email not found")
		userRepository.On("GetUserByEmail", userDomain.Email).Return(users.Domain{}, expectedErr).Once()

		token, err := userUseCase.Login(&userDomain)

		assert.Empty(t, token)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("Test Case 3 | Invalid Login | Wrong password", func(t *testing.T) {
		expectedErr := errors.New("wrong password")
		userRepository.On("GetUserByEmail", userDomain.Email).Return(userDomain, nil).Once()

		token, err := userUseCase.Login(&userDomain)

		assert.Empty(t, token)
		assert.Equal(t, err, expectedErr)
	})
}
