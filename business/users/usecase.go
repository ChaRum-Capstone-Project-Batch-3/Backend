package users

import (
	"charum/util"
	"errors"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepository Repository
}

func NewUserUseCase(ur Repository) UseCase {
	return &UserUseCase{
		userRepository: ur,
	}
}

/*
Create
*/

func (uu *UserUseCase) Register(domain *Domain) (Domain, string, error) {
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
	user, err := uu.userRepository.GetByEmail(domain.Email)
	if err != nil {
		return Domain{}, "", errors.New("email is not registered")
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

func (uu *UserUseCase) GetWithSortAndOrder(page int, limit int, sort string, order string) ([]Domain, int, int, error) {
	skip := limit * (page - 1)
	var orderInMongo int

	if order == "asc" {
		orderInMongo = 1
	} else {
		orderInMongo = -1
	}

	users, totalData, err := uu.userRepository.GetWithSortAndOrder(skip, limit, sort, orderInMongo)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("failed to get users")
	}

	totalPage := math.Ceil(float64(totalData) / float64(limit))
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

func (uu *UserUseCase) Update(domain *Domain) (Domain, error) {
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

	user.Email = domain.Email
	user.UserName = domain.UserName
	user.DisplayName = domain.DisplayName
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
