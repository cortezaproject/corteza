package mail

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"
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

// SetupDialer setups SMTP dialer
//
// Host variable can contain "<host>:<port>" that will override port value
func SetupDialer(host string, port int, user, pass, from string) {
	if host == "" {
		defaultDialerError = errors.New("No hostname provided for SMTP")
		return
	}

	if strings.Contains(host, ":") {
		parts := strings.SplitN(host, ":", 2)
		port, _ = strconv.Atoi(parts[1])
		host = parts[0]
	}

	if port == 0 {
		defaultDialerError = errors.New("No port provided for SMTP")
		return
	}
	if from == "" {
		defaultDialerError = errors.New("Sender for SMTP is not set")
		return
	}

	if defaultDialerError != nil {
		return
	}

	defaultFrom = from
	defaultDialer = gomail.NewDialer(
		host,
		port,
		user,
		pass,
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
