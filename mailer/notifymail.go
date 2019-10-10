package mailer

import (
	"bytes"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
	"os"
)

// Attachable represents an object which could have some attachments.
type Attachable interface {
	GetAttachments() []string
}

// NotifyMail epresents a notification mail.
type NotifyMail struct {
	Title       string       // Title of the mail
	Template    string       // Template used with the mail
	Values      interface{}  // Values to fill the template
	Attachments []Attachable // Absolute paths to the files to attach to the mail
	Level       NotifyLevel  // Level of notification. Filters mails in config file
}

// Send sends a notification mail to the list of notified parties on the config.
func (notification *NotifyMail) Send() error {
	notifyLevel, err := GetNotifyLevel()

	if err == nil && notification.Level >= notifyLevel {
		credentials, err := GetNotifyCredentials()
		if err != nil {
			return err
		}
		toList, err := GetNotifyEmails()
		if err != nil {
			return err
		}

		text := new(bytes.Buffer)
		t, err := template.New("mail").Parse(notification.Template)
		if err != nil {
			return err
		}
		if err := t.Execute(text, notification.Values); err != nil {
			return err
		}
		mail := email.NewEmail()
		mail.From = credentials.GetFormattedAddress()
		mail.To = toList
		mail.Subject = notification.Title
		for _, attachable := range notification.Attachments {
			attachments := attachable.GetAttachments()
			for _, attachment := range attachments {
				// prepareFTP if file exists and its length is greater than 0
				fStat, err := os.Stat(attachment)
				if err != nil {
					continue
				}
				if fStat.Size() > 0 && !fStat.IsDir() {
					_, _ = mail.AttachFile(attachment)
				}
			}
		}
		mail.Text = []byte(text.String())
		auth := smtp.PlainAuth("", credentials.Username, credentials.Password, credentials.Server)
		err = mail.Send(credentials.GetConnString(), auth)
		if err != nil {
			return err
		}
	}
	return nil
}
