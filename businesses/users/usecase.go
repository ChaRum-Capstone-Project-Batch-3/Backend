package users

import (
	"charum/util"
	"errors"
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

func (uu *UserUseCase) UserRegister(domain *Domain) (Domain, string, error) {
	domain.UserName = strings.ToLower(domain.UserName)
	_, err := uu.userRepository.GetUserByEmail(domain.Email)
	if err == nil {
		return Domain{}, "", errors.New("email is already registered")
	}

	_, err = uu.userRepository.GetUserByUsername(domain.UserName)
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
		return Domain{}, "", err
	}

	token := util.GenerateToken(user.Id.Hex(), user.Role)
	return user, token, nil
}

/*
Read
*/

func (uu *UserUseCase) Login(domain *Domain) (string, error) {
	user, err := uu.userRepository.GetUserByEmail(domain.Email)
	if err != nil {
		return "", errors.New("email is not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(domain.Password))
	if err != nil {
		return "", errors.New("wrong password")
	}

	token := util.GenerateToken(user.Id.Hex(), user.Role)
	return token, nil
}

/*
Update
*/

/*
Delete
*/
