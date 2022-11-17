package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"time"
)

const (
	hostCheckRE = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
)

// ConfigCheck dials and authenticates to an SMTP server.
func ConfigCheck(host string, port uint, username string, password string, tlsConfig *tls.Config) (checkRes string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	// Dial the tcp connection
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return err.Error()
	}

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Sprintf("SMTP connection to %s timeout!", addr)
	}
	defer c.Close()

	// Check whether STARTTLS extension is supported
	ok, _ := c.Extension("STARTTLS")
	if ok {
		err = c.StartTLS(tlsConfig)
		if err != nil {
			return err.Error()
		}
	}

	// Authentication
	auth := smtp.PlainAuth("", username, password, host)

	err = c.Auth(auth)
	if err != nil {
		return err.Error()
	}

	return
}

func IsValidHost(host string) bool {
	hostCheck := regexp.MustCompile(hostCheckRE)
	return hostCheck.MatchString(host)
}
