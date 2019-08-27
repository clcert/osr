package mailer

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

// SendTestMail sends a test mail.
func SendTestMail() error {
	credentials, err := GetNotifyCredentials()
	if err != nil {
		return err
	}
	toList, err := GetNotifyEmails()
	if err != nil {
		return err
	}
	testMail := email.NewEmail()
	testMail.Subject = "Correo de Prueba"
	testMail.Text = []byte("Este es un correo de prueba del OSR.")
	testMail.From = credentials.GetFormattedAddress()
	testMail.To = toList
	auth := smtp.PlainAuth("", credentials.Username, credentials.Password, credentials.Server)
	err = testMail.Send(credentials.GetConnString(), auth)
	if err != nil {
		return err
	}
	return nil
}
