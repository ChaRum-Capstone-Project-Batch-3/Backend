package users_test

import (
	"charum/business/users"
	_userMock "charum/business/users/mocks"
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	_cloudinaryMock "charum/helper/cloudinary/mocks"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	userRepository       _userMock.Repository
	cloudinaryRepository _cloudinaryMock.Function
	userUseCase          users.UseCase
	userDomain           users.Domain
	image                *multipart.FileHeader
)

func TestMain(m *testing.M) {
	userUseCase = users.NewUserUseCase(&userRepository, &cloudinaryRepository)

	userDomain = users.Domain{
		Id:                primitive.NewObjectID(),
		Email:             "test@charum.com",
		Password:          "!Test123",
		UserName:          "tester",
		Biodata:           "",
		SocialMedia:       "",
		ProfilePictureURL: "image",
		Role:              "user",
		IsActive:          true,
		CreatedAt:         primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:         primitive.NewDateTimeFromTime(time.Now()),
	}

	image = &multipart.FileHeader{}

	m.Run()
}

func TestUserRegister(t *testing.T) {
	t.Run("Test Case 1 | Valid Register", func(t *testing.T) {
		userRepository.On("GetByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		userRepository.On("Create", mock.Anything).Return(userDomain, nil).Once()

		actualUser, token, err := userUseCase.Register(&userDomain, image)

		assert.NotNil(t, actualUser)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Register | Email is already registered", func(t *testing.T) {
		expectedErr := errors.New("email is already registered")
		copyDomain := userDomain
		userRepository.On("GetByEmail", copyDomain.Email).Return(copyDomain, nil).Once()

		actualUser, token, err := userUseCase.Register(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test Case 3 | Invalid Register | Username is already used", func(t *testing.T) {
		expectedErr := errors.New("username is already used")
		copyDomain := userDomain
		userRepository.On("GetByEmail", copyDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", copyDomain.UserName).Return(copyDomain, nil).Once()

		actualUser, token, err := userUseCase.Register(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test Case 4 | Invalid Register | Failed to upload profile picture", func(t *testing.T) {
		expectedErr := errors.New("failed to upload profile picture")
		copyDomain := userDomain
		userRepository.On("GetByEmail", copyDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", copyDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("", expectedErr).Once()

		actualUser, token, err := userUseCase.Register(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test Case 4 | Invalid Register | Failed to register user", func(t *testing.T) {
		expectedErr := errors.New("failed to register user")

		userRepository.On("GetByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		userRepository.On("Create", mock.Anything).Return(userDomain, expectedErr).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

		actualUser, token, err := userUseCase.Register(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test Case 5 | Invalid Register | Failed to delete profile picture", func(t *testing.T) {
		expectedErr := errors.New("failed to delete profile picture")

		userRepository.On("GetByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		userRepository.On("Create", mock.Anything).Return(userDomain, errors.New("failed to register user")).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		actualUser, token, err := userUseCase.Register(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Test Case 1 | Valid Login", func(t *testing.T) {
		copyDomain := userDomain
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(copyDomain.Password), bcrypt.DefaultCost)
		copyDomain.Password = string(encryptedPassword)
		userRepository.On("GetByEmail", copyDomain.Email).Return(copyDomain, nil).Once()

		actualUser, token, err := userUseCase.Login(copyDomain.Email, userDomain.Password)

		assert.NotNil(t, actualUser)
		assert.NotEmpty(t, token)
		assert.Nil(t, err)
	})

	t.Run("Test Case 2 | Invalid Login | Wrong Password", func(t *testing.T) {
		expectedErr := errors.New("wrong password")
		copyDomain := userDomain
		copyDomain.Password = "wrong password"
		userRepository.On("GetByEmail", copyDomain.Email).Return(copyDomain, nil).Once()

		actualUser, token, err := userUseCase.Login(userDomain.Email, userDomain.Password)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test Case 3 | Invalid Login | Email or username is not registered", func(t *testing.T) {
		expectedErr := errors.New("email or username is not registered")
		userRepository.On("GetByEmail", userDomain.Email).Return(users.Domain{}, expectedErr).Once()
		userRepository.On("GetByUsername", userDomain.Email).Return(users.Domain{}, expectedErr).Once()

		actualUser, token, err := userUseCase.Login(userDomain.Email, userDomain.Password)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Test Case 4 | Invalid Login | User is suspended", func(t *testing.T) {
		expectedErr := errors.New("user is suspended")
		copyDomain := userDomain
		copyDomain.IsActive = false
		userRepository.On("GetByEmail", userDomain.Email).Return(copyDomain, nil).Once()

		actualUser, token, err := userUseCase.Login(userDomain.Email, userDomain.Password)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Empty(t, token)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetWithSortAndOrder(t *testing.T) {
	t.Run("Test Case 1 | Valid Get Users", func(t *testing.T) {
		query := dtoQuery.Request{
			Skip:  0,
			Limit: 1,
			Sort:  "createdAt",
			Order: -1,
		}
		userRepository.On("GetManyWithPagination", query, &userDomain).Return([]users.Domain{userDomain}, 1, nil).Once()

		pagination := dtoPagination.Request{
			Page:  1,
			Limit: 1,
			Sort:  "createdAt",
			Order: "desc",
		}
		actualUsers, totalPage, totalData, actualErr := userUseCase.GetManyWithPagination(pagination, &userDomain)

		assert.NotZero(t, totalData)
		assert.NotZero(t, totalPage)
		assert.NotNil(t, actualUsers)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Get Users | Error when getting users", func(t *testing.T) {
		expectedErr := errors.New("failed to get users")

		query := dtoQuery.Request{
			Skip:  0,
			Limit: 1,
			Sort:  "createdAt",
			Order: 1,
		}
		userRepository.On("GetManyWithPagination", query, &userDomain).Return([]users.Domain{}, 1, expectedErr).Once()

		pagination := dtoPagination.Request{
			Page:  1,
			Limit: 1,
			Sort:  "createdAt",
			Order: "asc",
		}
		actualUsers, totalPage, totalData, actualErr := userUseCase.GetManyWithPagination(pagination, &userDomain)

		assert.Zero(t, totalData)
		assert.Zero(t, totalPage)
		assert.Equal(t, []users.Domain{}, actualUsers)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("Test Case 1 | Valid Get User", func(t *testing.T) {
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()

		actualUser, actualErr := userUseCase.GetByID(userDomain.Id)

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Get User | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", userDomain.Id).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.GetByID(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Case 1 | Valid Update", func(t *testing.T) {
		copyDomain := userDomain
		copyDomain.UserName = "newUsername"
		copyDomain.Email = "newEmail@charum.com"
		copyDomain.DisplayName = "new display"
		copyDomain.IsActive = false

		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("GetByEmail", copyDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", copyDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		userRepository.On("Update", mock.Anything).Return(copyDomain, nil).Once()

		actualUser, actualErr := userUseCase.Update(&copyDomain, image)

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Update | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", userDomain.Id).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Update(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 3 | Invalid Update | Error when updating user", func(t *testing.T) {
		expectedErr := errors.New("failed to update user")
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", nil).Once()
		userRepository.On("Update", mock.Anything).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Update(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 4 | Invalid Update | Email is already registered", func(t *testing.T) {
		expectedErr := errors.New("email is already registered")
		copyDomain := userDomain
		copyDomain.Email = "existemail@charum.com"
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("GetByEmail", copyDomain.Email).Return(copyDomain, nil).Once()

		actualUser, actualErr := userUseCase.Update(&copyDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 5 | Invalid Update | Username is already used", func(t *testing.T) {
		expectedErr := errors.New("username is already used")
		copyDomain := userDomain
		copyDomain.UserName = "existusername"

		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("GetByUsername", copyDomain.UserName).Return(copyDomain, nil).Once()

		actualUser, actualErr := userUseCase.Update(&copyDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 6 | Invalid Update | Error when upload profile picture", func(t *testing.T) {
		expectedErr := errors.New("failed to upload profile picture")

		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("GetByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		cloudinaryRepository.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return("image", expectedErr).Once()

		actualUser, actualErr := userUseCase.Update(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 7 | Invalid Update | Error when delete old profile picture", func(t *testing.T) {
		expectedErr := errors.New("failed to delete old profile picture")

		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("GetByEmail", userDomain.Email).Return(users.Domain{}, errors.New("not found")).Once()
		userRepository.On("GetByUsername", userDomain.UserName).Return(users.Domain{}, errors.New("not found")).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		actualUser, actualErr := userUseCase.Update(&userDomain, image)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Run("Test Case 1 | Valid Update Password", func(t *testing.T) {
		copyDomain := userDomain
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)
		copyDomain.Password = string(encryptedPassword)
		copyDomain.OldPassword = userDomain.Password
		copyDomain.NewPassword = "newpassword"

		userRepository.On("GetByID", userDomain.Id).Return(copyDomain, nil).Once()
		userRepository.On("UpdatePassword", mock.Anything).Return(copyDomain, nil).Once()
		actualUser, actualErr := userUseCase.UpdatePassword(&copyDomain)

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Update Password | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		copyDomain := userDomain
		copyDomain.OldPassword = userDomain.Password
		copyDomain.NewPassword = "newpassword"

		userRepository.On("GetByID", userDomain.Id).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.UpdatePassword(&copyDomain)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 3 | Invalid Update Password | Error when updating password", func(t *testing.T) {
		expectedErr := errors.New("failed to update password")
		copyDomain := userDomain
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)
		copyDomain.Password = string(encryptedPassword)
		copyDomain.OldPassword = userDomain.Password
		copyDomain.NewPassword = "newpassword"

		userRepository.On("GetByID", userDomain.Id).Return(copyDomain, nil).Once()
		userRepository.On("UpdatePassword", mock.Anything).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.UpdatePassword(&copyDomain)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 4 | Invalid Update Password | Wrong Password", func(t *testing.T) {
		expectedErr := errors.New("wrong password")
		copyDomain := userDomain
		copyDomain.OldPassword = "wrong password"
		userRepository.On("GetByID", copyDomain.Id).Return(copyDomain, nil).Once()

		actualUser, actualErr := userUseCase.UpdatePassword(&copyDomain)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Case 1 | Valid Delete", func(t *testing.T) {
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		userRepository.On("Delete", userDomain.Id).Return(nil).Once()

		actualUser, actualErr := userUseCase.Delete(userDomain.Id)

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Delete | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", userDomain.Id).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Delete(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 3 | Invalid Delete | Error when deleting user", func(t *testing.T) {
		expectedErr := errors.New("failed to delete user")
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
		userRepository.On("Delete", userDomain.Id).Return(expectedErr).Once()

		actualUser, actualErr := userUseCase.Delete(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 4 | Invalid Delete | Error when deleting profile picture", func(t *testing.T) {
		expectedErr := errors.New("failed to delete profile picture")
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		cloudinaryRepository.On("Delete", mock.Anything, mock.Anything).Return(expectedErr).Once()

		actualUser, actualErr := userUseCase.Delete(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestSuspend(t *testing.T) {
	t.Run("Test Case 1 | Valid Suspend", func(t *testing.T) {
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("Update", mock.Anything).Return(userDomain, nil).Once()

		actualUser, actualErr := userUseCase.Suspend(userDomain.Id)

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Suspend | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", userDomain.Id).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Suspend(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 3 | Invalid Suspend | Error when suspending user", func(t *testing.T) {
		expectedErr := errors.New("failed to suspend user")
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()
		userRepository.On("Update", mock.Anything).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Suspend(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 4 | Invalid Suspend | User is already suspended", func(t *testing.T) {
		expectedErr := errors.New("user is already suspended")
		copyDomain := userDomain
		copyDomain.IsActive = false
		userRepository.On("GetByID", userDomain.Id).Return(copyDomain, nil).Once()

		actualUser, actualErr := userUseCase.Suspend(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestUnsuspend(t *testing.T) {
	t.Run("Test Case 1 | Valid Unsuspend", func(t *testing.T) {
		copyDomain := userDomain
		copyDomain.IsActive = false
		userRepository.On("GetByID", userDomain.Id).Return(copyDomain, nil).Once()
		userRepository.On("Update", mock.Anything).Return(userDomain, nil).Once()

		actualUser, actualErr := userUseCase.Unsuspend(userDomain.Id)

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid Unsuspend | Error when getting user", func(t *testing.T) {
		expectedErr := errors.New("failed to get user")
		userRepository.On("GetByID", userDomain.Id).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Unsuspend(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 3 | Invalid Unsuspend | Error when unsuspending user", func(t *testing.T) {
		expectedErr := errors.New("failed to unsuspend user")
		copyDomain := userDomain
		copyDomain.IsActive = false
		userRepository.On("GetByID", userDomain.Id).Return(copyDomain, nil).Once()
		userRepository.On("Update", mock.Anything).Return(users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.Unsuspend(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})

	t.Run("Test Case 4 | Invalid Unsuspend | User is not suspended", func(t *testing.T) {
		expectedErr := errors.New("user is not suspended")
		userRepository.On("GetByID", userDomain.Id).Return(userDomain, nil).Once()

		actualUser, actualErr := userUseCase.Unsuspend(userDomain.Id)

		assert.Equal(t, users.Domain{}, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Test Case 1 | Valid GetAll", func(t *testing.T) {
		userRepository.On("GetAll").Return([]users.Domain{userDomain}, nil).Once()

		actualUser, actualErr := userUseCase.GetAll()

		assert.NotNil(t, actualUser)
		assert.Nil(t, actualErr)
	})

	t.Run("Test Case 2 | Invalid GetAll | Error when getting all users", func(t *testing.T) {
		expectedErr := errors.New("failed to get all users")
		userRepository.On("GetAll").Return([]users.Domain{}, expectedErr).Once()

		actualUser, actualErr := userUseCase.GetAll()

		assert.Equal(t, 0, actualUser)
		assert.Equal(t, expectedErr, actualErr)
	})
}
