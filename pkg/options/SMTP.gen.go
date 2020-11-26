package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/SMTP.yaml

type (
	SMTPOpt struct {
		Host          string `env:"SMTP_HOST"`
		Port          int    `env:"SMTP_PORT"`
		User          string `env:"SMTP_USER"`
		Pass          string `env:"SMTP_PASS"`
		From          string `env:"SMTP_FROM"`
		TlsInsecure   bool   `env:"SMTP_TLS_INSECURE"`
		TlsServerName string `env:"SMTP_TLS_SERVER_NAME"`
	}
)

// SMTP initializes and returns a SMTPOpt with default values
func SMTP() (o *SMTPOpt) {
	o = &SMTPOpt{
		Host:        "localhost",
		Port:        25,
		TlsInsecure: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *SMTP) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
