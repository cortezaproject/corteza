package options

type (
	ActionLogOpt struct {
		Enabled bool `env:"ACTIONLOG_ENABLED"`
		Debug   bool `env:"ACTIONLOG_DEBUG"`
	}
)

func ActionLog() (o *ActionLogOpt) {
	o = &ActionLogOpt{
		Enabled: true,
		Debug:   false,
	}

	fill(o)

	return
}
