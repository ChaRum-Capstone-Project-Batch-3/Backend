package forgot_password

import (
	"charum/business/users"
	_mailgun "charum/helper/mailgun"
	"charum/util"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordUseCase struct {
	forgotPassword Repository
	userRepository users.Repository
	mailgun        _mailgun.Function
}

func NewForgotPasswordUseCase(fp Repository, ur users.Repository, mg _mailgun.Function) UseCase {
	return &ForgotPasswordUseCase{
		forgotPassword: fp,
		userRepository: ur,
		mailgun:        mg,
	}
}

/*
Create
*/

func (fpu *ForgotPasswordUseCase) Generate(domain *Domain) (Domain, error) {
	_, err := fpu.userRepository.GetByEmail(domain.Email)
	if err != nil {
		return Domain{}, errors.New("email is not registered")
	}

	token := util.GenerateRandomString(80)
	domain.Id = primitive.NewObjectID()
	domain.Token = token
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.ExpiredAt = primitive.NewDateTimeFromTime(time.Now().Add(30 * time.Minute))
	domain.IsUsed = false

	forgotPassword, err := fpu.forgotPassword.Generate(domain)
	if err != nil {
		return Domain{}, errors.New("failed to reset password")
	}

	_, err = fpu.mailgun.SendMail(domain.Email, domain.Token)
	if err != nil {
		delErr := fpu.forgotPassword.Delete(domain.Id)
		if delErr != nil {
			return Domain{}, errors.New("failed to reset password")
		}

		return Domain{}, err
	}

	return forgotPassword, nil
}

/*
Read
*/

func (fpu *ForgotPasswordUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	forgotPassword, err := fpu.forgotPassword.GetByID(id)
	if err != nil {
		return Domain{}, err
	}

	return forgotPassword, nil
}

func (fpu *ForgotPasswordUseCase) GetByToken(token string) (Domain, error) {
	forgotPassword, err := fpu.forgotPassword.GetByToken(token)
	if err != nil {
		return Domain{}, errors.New("failed to get token")
	}

	return forgotPassword, nil
}

func (fpu *ForgotPasswordUseCase) ValidateToken(token string) (Domain, error) {
	tokenData, err := fpu.forgotPassword.GetByToken(token)
	if err != nil {
		return Domain{}, errors.New("failed to get token")
	}

	if tokenData.IsUsed {
		return Domain{}, errors.New("token has been used")
	}

	if tokenData.ExpiredAt.Time().Before(time.Now()) {
		return Domain{}, errors.New("token has expired")
	}

	return tokenData, nil
}

/*
Update
*/

func (fpu *ForgotPasswordUseCase) UpdatePassword(domain *Domain) (Domain, error) {
	tokenData, err := fpu.ValidateToken(domain.Token)
	if err != nil {
		return Domain{}, err
	}

	user, err := fpu.userRepository.GetByEmail(tokenData.Email)
	if err != nil {
		return Domain{}, errors.New("email is not registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)
	if err != nil {
		return Domain{}, errors.New("failed to update password")
	}

	user.Password = string(hashedPassword)
	tokenData.IsUsed = true
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	tokenData.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = fpu.userRepository.UpdatePassword(&user)
	if err != nil {
		return Domain{}, errors.New("failed to update password")
	}

	forgotPassword, err := fpu.forgotPassword.Update(&tokenData)
	if err != nil {
		return Domain{}, errors.New("failed to update token")
	}

	return forgotPassword, nil
}
