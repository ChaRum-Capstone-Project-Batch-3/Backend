package users

import (
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	"charum/helper/cloudinary"
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
	cloudinary     cloudinary.Function
}

func NewUserUseCase(ur Repository, cld cloudinary.Function) UseCase {
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
		cloudinaryURL, err := uu.cloudinary.Upload("profilePicture", profilePicture, util.GenerateUUID())
		if err != nil {
			return Domain{}, "", errors.New("failed to upload profile picture")
		}

		domain.ProfilePictureURL = cloudinaryURL
	}

	domain.Id = primitive.NewObjectID()
	domain.Password = string(encryptedPassword)
	domain.Role = "user"
	domain.IsActive = true
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	user, err := uu.userRepository.Create(domain)
	if err != nil {
		if domain.ProfilePictureURL != "" {
			err = uu.cloudinary.Delete("profilePicture", util.GetFilenameWithoutExtension(domain.ProfilePictureURL))
			if err != nil {
				return Domain{}, "", errors.New("failed to delete profile picture")
			}
		}

		return Domain{}, "", errors.New("failed to register user")
	}

	token := util.GenerateToken(user.Id.Hex(), user.Role)
	return user, token, nil
}

/*
Read
*/

func (uu *UserUseCase) Login(key string, password string) (Domain, string, error) {
	var user Domain

	user, err := uu.userRepository.GetByEmail(key)
	if err != nil {
		user, err = uu.userRepository.GetByUsername(key)
		if err != nil {
			return Domain{}, "", errors.New("email or username is not registered")
		}
	}

	if !user.IsActive {
		return Domain{}, "", errors.New("user is suspended")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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

func (uu *UserUseCase) GetAll() (int, error) {
	users, err := uu.userRepository.GetAll()
	if err != nil {
		return 0, errors.New("failed to get all users")
	}

	return len(users), nil
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
		if user.ProfilePictureURL != "" {
			err := uu.cloudinary.Delete("profilePicture", util.GetFilenameWithoutExtension(user.ProfilePictureURL))
			if err != nil {
				return Domain{}, errors.New("failed to delete old profile picture")
			}
		}

		cloudinaryURL, err := uu.cloudinary.Upload("profilePicture", profilePicture, util.GenerateUUID())
		if err != nil {
			return Domain{}, errors.New("failed to upload profile picture")
		}

		user.ProfilePictureURL = cloudinaryURL
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

func (uu *UserUseCase) UpdatePassword(domain *Domain) (Domain, error) {
	user, err := uu.userRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(domain.OldPassword))
	if err != nil {
		return Domain{}, errors.New("wrong password")
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(domain.NewPassword), bcrypt.DefaultCost)
	user.Password = string(encryptedPassword)
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updatedUser, err := uu.userRepository.UpdatePassword(&user)
	if err != nil {
		return Domain{}, errors.New("failed to update password")
	}

	return updatedUser, nil
}

func (uu *UserUseCase) Suspend(id primitive.ObjectID) (Domain, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get user")
	}

	if user.Role == "admin" {
		return Domain{}, errors.New("admin cannot be suspended")
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

	if deletedUser.Role == "admin" {
		return Domain{}, errors.New("admin cannot be deleted")
	}

	if deletedUser.ProfilePictureURL != "" {
		err = uu.cloudinary.Delete("profilePicture", util.GetFilenameWithoutExtension(deletedUser.ProfilePictureURL))
		if err != nil {
			return Domain{}, errors.New("failed to delete profile picture")
		}
	}

	err = uu.userRepository.Delete(id)
	if err != nil {
		return Domain{}, errors.New("failed to delete user")
	}

	return deletedUser, nil
}
