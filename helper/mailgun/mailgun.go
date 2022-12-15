package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

type Function interface {
	SendMail(email string, token string) (string, error)
}

type Mailgun struct {
	Mailgun     *mailgun.MailgunImpl
	EmailDomain string
}

func Init(emailDomain string, mailgunKey string) Function {
	mg := mailgun.NewMailgun(emailDomain, mailgunKey)

	return &Mailgun{
		Mailgun:     mg,
		EmailDomain: emailDomain,
	}
}

func (mg *Mailgun) SendMail(email string, token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	m := mg.Mailgun.NewMessage(fmt.Sprintf("Charum No-Reply <noreply@%s>", mg.EmailDomain), "Your Reset Password Link", "")
	m.SetTemplate("charum")
	if err := m.AddRecipient(email); err != nil {
		return "", err
	}

	vars, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return "", err
	}
	m.AddHeader("X-Mailgun-Template-Variables", string(vars))

	_, id, err := mg.Mailgun.Send(ctx, m)
	return id, err
}
