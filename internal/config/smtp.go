package config

import (
	"errors"
	"strconv"
	"strings"

	"github.com/namsral/flag"
)

type (
	SMTP struct {
		Host string
		Port int
		User string
		Pass string
		From string
	}
)

const defaultSMTPPort = 25

var smtp *SMTP

// Validate
//
// No actual validation here for SMTP, we allow mis/un-configured
func (c *SMTP) Validate() error {
	if strings.Contains(c.Host, ":") {
		parts := strings.SplitN(c.Host, ":", 2)
		c.Port, _ = strconv.Atoi(parts[1])
		c.Host = parts[0]
	}

	return nil
}

// RuntimeValidation is used for run-time configuration validation
func (c *SMTP) RuntimeValidation() error {
	if c == nil {
		return errors.New("SMTP config missing")
	}

	if c.Host == "" {
		return errors.New("No hostname provided for SMTP")
	}
	if c.Port == 0 {
		return errors.New("No port provided for SMTP")
	}
	if c.From == "" {
		return errors.New("Sender for SMTP is not set")
	}

	return nil
}

func (*SMTP) Init(prefix ...string) *SMTP {
	if smtp != nil {
		return smtp
	}
	smtp = new(SMTP)
	flag.StringVar(&smtp.Host, "smtp-host", "", "SMTP hostname (may be host:port)")
	flag.IntVar(&smtp.Port, "smtp-port", defaultSMTPPort, "SMTP port number")
	flag.StringVar(&smtp.User, "smtp-user", "", "SMTP server username")
	flag.StringVar(&smtp.Pass, "smtp-pass", "", "SMTP server password")
	flag.StringVar(&smtp.From, "smtp-from", "", "SMTP sender header")

	return smtp
}
