package mail

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
)

func TestMailSend(t *testing.T) {
	if !testing.Short() {
		godotenv.Load("../../.env")

		Flags()
		flag.Parse()
		if err := flags.Validate(); err != nil {
			t.Fatalf("Missing SMTP flags: %+v", err)
		}

		message := New()
		message.SetHeader("To", "black@scene-si.org")
		message.SetHeader("Subject", "Hello from Crust tests!")
		message.SetBody("text/html", "Lorem <i>ipsum</i> <u>dolor</u> sit <b>amet</b>!")
		if err := Send(message); err != nil {
			t.Fatalf("E-mail failed to send: %+v", err)
		}
	}
}
