package users

import (
	_cloudinary "charum/driver/cloudinary"
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	"charum/helper"
	"charum/util"
	"errors"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepository Repository
	cloudinary     _cloudinary.Function
}

func NewUserUseCase(ur Repository, cld _cloudinary.Function) UseCase {
	return &UserUseCase{
		userRepository: ur,
		cloudinary:     cld,
	}
}

/*
Create
*/

func (uu *UserUseCase) Register(domain *Domain, profilePicture *multipart.FileHeader) (Domain, string, error) {
	domain.UserName = strings.ToLower(domain.UserName)
	_, err := uu.userRepository.GetByEmail(domain.Email)
	if err == nil {
		return Domain{}, "", errors.New("email is already registered")
	}

	_, err = uu.userRepository.GetByUsername(domain.UserName)
	if err == nil {
		return Domain{}, "", errors.New("username is already used")
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)

	if profilePicture != nil {
		uploadResult, err := uu.cloudinary.Upload("profilePicture", profilePicture, helper.GenerateUUID())
		if err != nil {
			return Domain{}, "", errors.New("failed to upload profile picture")
		}

		domain.ProfilePictureURL = uploadResult
	}

	domain.Id = primitive.NewObjectID()
	domain.Password = string(encryptedPassword)
	domain.Role = "user"
	domain.IsActive = true
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	user, err := uu.userRepository.Create(domain)

	if err != nil {
		return Domain{}, "", errors.New("failed to register user")
	}

	token := util.GenerateToken(user.Id.Hex(), user.Role)
	return user, token, nil
}

/*
Read
*/

func (uu *UserUseCase) Login(domain *Domain) (Domain, string, error) {
	var user Domain

	user, err := uu.userRepository.GetByEmail(domain.Email)
	if err != nil {
		user, err = uu.userRepository.GetByUsername(domain.Email)
		if err != nil {
			return Domain{}, "", errors.New("email or username is not registered")
		}
	}

	if !user.IsActive {
		return Domain{}, "", errors.New("user is suspended")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(domain.Password))
	if err != nil {
		return Domain{}, "", errors.New("wrong password")
	}

	token := util.GenerateToken(user.Id.Hex(), user.Role)
	return user, token, nil
}

func (uu *UserUseCase) GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error) {
	skip := pagination.Limit * (pagination.Page - 1)
	var orderInMongo int

	if pagination.Order == "asc" {
		orderInMongo = 1
	} else {
		orderInMongo = -1
	}

	query := dtoQuery.Request{
		Skip:  skip,
		Limit: pagination.Limit,
		Order: orderInMongo,
		Sort:  pagination.Sort,
	}

	users, totalData, err := uu.userRepository.GetManyWithPagination(query, domain)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("failed to get users")
	}

	totalPage := math.Ceil(float64(totalData) / float64(pagination.Limit))
	return users, int(totalPage), totalData, nil
}

func (uu *UserUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	return user, nil
}

/*
Update
*/

func (uu *UserUseCase) Update(domain *Domain, profilePicture *multipart.FileHeader) (Domain, error) {
	user, err := uu.userRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if domain.Email != user.Email {
		_, err = uu.userRepository.GetByEmail(domain.Email)
		if err == nil {
			return Domain{}, errors.New("email is already registered")
		}
	}

	if domain.UserName != user.UserName {
		_, err = uu.userRepository.GetByUsername(domain.UserName)
		if err == nil {
			return Domain{}, errors.New("username is already used")
		}
	}

	if profilePicture != nil {
		err := uu.cloudinary.Delete("profilePicture", helper.GetFilenameWithoutExtension(user.ProfilePictureURL))
		if err != nil {
			return Domain{}, errors.New("failed to delete old profile picture")
		}

		uploadResult, err := uu.cloudinary.Upload("profilePicture", profilePicture, helper.GenerateUUID())
		if err != nil {
			return Domain{}, errors.New("failed to upload profile picture")
		}

		user.ProfilePictureURL = uploadResult
	}

	user.Email = domain.Email
	user.UserName = domain.UserName
	user.DisplayName = domain.DisplayName
	user.Biodata = domain.Biodata
	user.SocialMedia = domain.SocialMedia
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updatedUser, err := uu.userRepository.Update(&user)
	if err != nil {
		return Domain{}, errors.New("failed to update user")
	}

	return updatedUser, nil
}

func (uu *UserUseCase) Suspend(id primitive.ObjectID) (Domain, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if !user.IsActive {
		return Domain{}, errors.New("user is already suspended")
	}

	user.IsActive = false
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	suspendedUser, err := uu.userRepository.Update(&user)
	if err != nil {
		return Domain{}, errors.New("failed to suspend user")
	}

	return suspendedUser, nil
}

func (uu *UserUseCase) Unsuspend(id primitive.ObjectID) (Domain, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if user.IsActive {
		return Domain{}, errors.New("user is not suspended")
	}

	user.IsActive = true
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	unsuspendedUser, err := uu.userRepository.Update(&user)
	if err != nil {
		return Domain{}, errors.New("failed to unsuspend user")
	}

	return unsuspendedUser, nil
}

/*
Delete
*/

func (uu *UserUseCase) Delete(id primitive.ObjectID) (Domain, error) {
	deletedUser, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	err = uu.userRepository.Delete(id)
	if err != nil {
		return Domain{}, errors.New("failed to delete user")
	}

	return deletedUser, nil
}
