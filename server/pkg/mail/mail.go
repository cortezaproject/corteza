package mail

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gomail "gopkg.in/mail.v2"
)

type (
	Dialer interface {
		DialAndSend(...*gomail.Message) error
	}

	applyCfg func(*gomail.Dialer)
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
func SetupDialer(host string, port int, user, pass, from string, ff ...applyCfg) {
	if host == "" {
		defaultDialerError = fmt.Errorf("no hostname provided for SMTP")
		return
	}

	if strings.Contains(host, ":") {
		parts := strings.SplitN(host, ":", 2)
		port, _ = strconv.Atoi(parts[1])
		host = parts[0]
	}

	if port == 0 {
		defaultDialerError = fmt.Errorf("no port provided for SMTP")
		return
	}

	if from == "" {
		from = "placeholder@example.net"
	}

	if defaultDialerError != nil {
		return
	}

	defaultFrom = from
	dialer := gomail.NewDialer(
		host,
		port,
		user,
		pass,
	)

	for _, fn := range ff {
		fn(dialer)
	}

	defaultDialer = dialer
}

func New() *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", defaultFrom)
	return message
}

// Send message with SMTP dialer
func Send(message *gomail.Message, dd ...Dialer) (err error) {
	for _, d := range append(dd, defaultDialer) {
		if d == nil {
			continue
		}

		if err = d.DialAndSend(message); err != nil {
			return fmt.Errorf("could not send email: %w", err)
		} else {
			return nil
		}
	}

	// At this point, none of the dialer could be used,
	// is there an error with default dialer?
	if defaultDialerError != nil {
		return fmt.Errorf("could not send email: %w", defaultDialerError)
	}

	return fmt.Errorf("unable to find configured and working SMTP dialer")

}

func IsValidAddress(addr string) bool {
	return addressCheck.MatchString(addr)
}
