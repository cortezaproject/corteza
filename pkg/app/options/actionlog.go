package options

type (
	ActionLogOpt struct {
		Debug bool `env:"ACTIONLOG_DEBUG"`
	}
)

func ActionLog() (o *ActionLogOpt) {
	o = &ActionLogOpt{
		Debug: false,
	}

	fill(o, "")

	return
}
