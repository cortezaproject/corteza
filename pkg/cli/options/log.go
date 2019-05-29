package options

type (
	// Logger's output leve is configured here, but
	// dev/prod configuration happens earlier
	LogOpt struct {
		Level string
	}
)

func Log(pfix string) (o *LogOpt) {
	o = &LogOpt{
		Level: EnvString(pfix, "LOG_LEVEL", "info"),
	}

	return
}
