package mail

import (
	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"
)

func New() *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", flags.From)
	return message
}

func Send(message *gomail.Message) error {
	dialer := gomail.NewDialer(
		flags.Host,
		flags.Port,
		flags.User,
		flags.Pass,
	)
	return errors.WithStack(dialer.DialAndSend(message))
}
