package models

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

const (
	DefaultSender = "support@blog.com"
)

type Email struct {
	From string
	To   string
	Subject string
	Plaintext string
	HTML string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

	return &es
}

type EmailService struct {
	DefaultSender string
	dialer        *gomail.Dialer
}

func (es *EmailService) Send(email Email) error {
	msg := gomail.NewMessage()
	es.setFrom(msg, email)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject: "Reset your password",
		To: to,
		Plaintext: "To reset your password, please visit the following link: "+ resetURL,
		HTML: `<p>To reset your password, please visit the following link: <a href="`+ resetURL + `">` +
			resetURL + `</a></p>`,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(msg *gomail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}
