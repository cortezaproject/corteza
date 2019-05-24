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

	bindString(cmd, &o.Host,
		"smtp-host", "localhost:25",
		"SMTP hostname")

	bindString(cmd, &o.User,
		"smtp-username", "",
		"SMTP server username")

	bindString(cmd, &o.Pass,
		"smtp-pass", "",
		"SMTP server password")

	bindString(cmd, &o.From,
		"smtp-from", "",
		"Sender's email address")

	bindInt(cmd, &o.Port,
		"smtp-port", 25,
		"SMTP port number")

	return
}
