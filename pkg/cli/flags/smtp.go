package flags

import (
	"github.com/spf13/cobra"
)

type (
	SMTPOpt struct {
		Host string
		Port int
		User string
		Pass string
		From string
	}
)

func SMTP(cmd *cobra.Command) (o *SMTPOpt) {
	o = &SMTPOpt{}

	BindString(cmd, &o.Host,
		"smtp-host", "localhost:25",
		"SMTP hostname")

	BindString(cmd, &o.User,
		"smtp-username", "",
		"SMTP server username")

	BindString(cmd, &o.Pass,
		"smtp-pass", "",
		"SMTP server password")

	BindString(cmd, &o.From,
		"smtp-from", "",
		"Sender's email address")

	BindInt(cmd, &o.Port,
		"smtp-port", 25,
		"SMTP port number")

	return
}
