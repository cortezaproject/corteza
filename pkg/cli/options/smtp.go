package options

type (
	SMTPOpt struct {
		Host string
		Port int
		User string
		Pass string
		From string
	}
)

func SMTP(pfix string) (o *SMTPOpt) {
	o = &SMTPOpt{
		Host: EnvString(pfix, "SMTP_HOST", "localhost:25"),
		Port: EnvInt(pfix, "SMTP_PORT", 25),
		User: EnvString(pfix, "SMTP_USERNAME", ""),
		Pass: EnvString(pfix, "SMTP_PASS", ""),
		From: EnvString(pfix, "SMTP_FROM", ""),
	}

	return
}
