package mail

import (
	"crypto/tls"
	"regexp"
	"time"

	gomail "gopkg.in/mail.v2"
)

const (
	hostCheckRE = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
)

// ConfigCheck dials and authenticates to an SMTP server using the same
// gomail dialer used for the regular send path. Returns "" on success
// or a human-readable error string on failure.
func ConfigCheck(host string, port uint, username, password string, tlsConfig *tls.Config) (checkRes string) {
	d := gomail.NewDialer(host, int(port), username, password)
	d.Timeout = 10 * time.Second
	if tlsConfig != nil {
		d.TLSConfig = tlsConfig
	}

	sc, err := d.Dial()
	if err != nil {
		return err.Error()
	}
	_ = sc.Close()
	return ""
}

func IsValidHost(host string) bool {
	hostCheck := regexp.MustCompile(hostCheckRE)
	return hostCheck.MatchString(host)
}
