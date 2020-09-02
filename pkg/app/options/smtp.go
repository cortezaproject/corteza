package options

type (
	SMTPOpt struct {
		Host string `env:"SMTP_HOST"`
		Port int    `env:"SMTP_PORT"`
		User string `env:"SMTP_USER"`
		Pass string `env:"SMTP_PASS"`
		From string `env:"SMTP_FROM"`

		TlsInsecure   bool   `env:"SMTP_TSL_INSECURE"`
		TlsServerName string `env:"SMTP_TSL_SERVER_NAME"`
	}
)

func SMTP(pfix string) (o *SMTPOpt) {
	o = &SMTPOpt{
		Host: "localhost",
		Port: 25,
		User: "",
		Pass: "",
		From: "",

		TlsInsecure:   false,
		TlsServerName: "",
	}

	fill(o, pfix)

	return
}
