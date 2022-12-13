package forgot_password

import (
	"charum/business/users"
	"charum/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/mailgun/mailgun-go/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type ForgotPasswordUseCase struct {
	forgotPassword Repository
	userRepository users.Repository
}

func NewForgotPasswordUseCase(fp Repository, ur users.Repository) UseCase {
	return &ForgotPasswordUseCase{
		forgotPassword: fp,
		userRepository: ur,
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
	// generate random string
	token := util.GenerateRandomString(80)
	domain.Id = primitive.NewObjectID()
	domain.Token = token
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.ExpiredAt = primitive.NewDateTimeFromTime(time.Now().Add(time.Minute * 30))
	domain.IsUsed = false

	forgotPassword, err := fpu.forgotPassword.Generate(domain)
	if err != nil {
		return Domain{}, errors.New("failed to reset password")
	}

	// sendmail from driver.mail
	_, err = fpu.SendMail(domain)
	if err != nil {
		return Domain{}, errors.New("failed to send email")
	}

	return forgotPassword, nil
}

func (fpu *ForgotPasswordUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	forgotPassword, err := fpu.forgotPassword.GetByID(id)
	if err != nil {
		return Domain{}, err
	}

	return forgotPassword, nil
}

// get by token
func (fpu *ForgotPasswordUseCase) GetByToken(token string) (Domain, error) {
	forgotPassword, err := fpu.forgotPassword.GetByToken(token)
	if err != nil {
		return Domain{}, err
	}

	return forgotPassword, nil
}

func (fpu *ForgotPasswordUseCase) ValidateToken(token string) (Domain, error) {
	tokenData, err := fpu.forgotPassword.GetByToken(token)
	if err != nil {
		return Domain{}, errors.New("token invalid/not found")
	}

	if tokenData.IsUsed {
		return Domain{}, errors.New("token has been used")
	}

	if tokenData.ExpiredAt.Time().Before(time.Now()) {
		return Domain{}, errors.New("token has expired")
	}

	return tokenData, nil
}

// update password
func (fpu *ForgotPasswordUseCase) UpdatePassword(domain *Domain) (Domain, error) {
	// validate token
	tokenData, err := fpu.ValidateToken(domain.Token)
	if err != nil {
		return Domain{}, err
	}

	user, err := fpu.userRepository.GetByEmail(tokenData.Email)
	if err != nil {
		return Domain{}, errors.New("email is not registered")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)
	if err != nil {
		return Domain{}, errors.New("failed to update password")
	}

	// update password
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

func (fpu *ForgotPasswordUseCase) SendMail(domain *Domain) (string, error) {
	// get mailgun_api_key from env
	mailgunKey := util.GetConfig("MAILGUN_API_KEY")
	mg := mailgun.NewMailgun("charum.razanfawwaz.xyz", mailgunKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	m := mg.NewMessage("Charum No-Reply <noreply@charum.razanfawwaz.xyz>", "Your Reset Password Link", "")
	m.SetTemplate("charum")
	if err := m.AddRecipient(domain.Email); err != nil {
		return "", err
	}

	vars, err := json.Marshal(map[string]string{
		"token": domain.Token,
	})
	if err != nil {
		return "", err
	}
	m.AddHeader("X-Mailgun-Template-Variables", string(vars))

	_, id, err := mg.Send(ctx, m)

	return id, err
}
