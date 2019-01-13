package mail

import (
	"regexp"

	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"

	"github.com/crusttech/crust/internal/config"
)

type (
	Dialer interface {
		DialAndSend(...*gomail.Message) error
	}
)

const (
	addressCheckRE = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

var (
	addressCheck       *regexp.Regexp
	defaultDialer      Dialer
	defaultFrom        string
	defaultDialerError error
)

func init() {
	addressCheck = regexp.MustCompile(addressCheckRE)
}

func SetupDialer(config *config.SMTP) {
	defaultDialerError = config.RuntimeValidation()

	if defaultDialerError != nil {
		return
	}

	defaultFrom = config.From
	defaultDialer = gomail.NewDialer(
		config.Host,
		config.Port,
		config.User,
		config.Pass,
	)
}

func New() *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", defaultFrom)
	return message
}

// Sends message with SMTP dialer
func Send(message *gomail.Message, dd ...Dialer) (err error) {

	for _, d := range append(dd, defaultDialer) {
		if d == nil {
			continue
		}

		return errors.WithStack(d.DialAndSend(message))
	}

	// At this point, none of the dialer could be used,
	// is there an error with default dialer?
	if defaultDialerError != nil {
		return errors.Wrap(defaultDialerError, "could not send email")
	}

	return errors.New("unable to find configured and working SMTP dialer")

}

func IsValidAddress(addr string) bool {
	return addressCheck.MatchString(addr)
}
