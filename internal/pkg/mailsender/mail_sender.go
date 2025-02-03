package mailsender

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	dialer *gomail.Dialer
}

func NewEmailSender(smtpServer string, smtpPort int, from, password string) *EmailSender {
	// SMTP server configuration
	dialer := gomail.NewDialer(smtpServer, smtpPort, from, password)
	return &EmailSender{dialer: dialer}
}

func (e *EmailSender) SendEmail(to, subject, body string) error {

	// Email message setup
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", e.dialer.Username)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)
	// Send the email
	err := e.dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	return nil
}
