package options

type (
	SMTPOpt struct {
		Host string `env:"SMTP_HOST"`
		Port int    `env:"SMTP_PORT"`
		User string `env:"SMTP_USER"`
		Pass string `env:"SMTP_PASS"`
		From string `env:"SMTP_FROM"`
	}
)

func SMTP(pfix string) (o *SMTPOpt) {
	o = &SMTPOpt{
		Host: "localhost:25",
		Port: 25,
		User: "",
		Pass: "",
		From: "",
	}

	fill(o, pfix)

	return
}
