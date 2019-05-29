package options

type (
	// Logger's output leve is configured here, but
	// dev/prod configuration happens earlier
	LogOpt struct {
		Level string `env:"LOG_LEVEL"`
	}
)

func Log(pfix string) (o *LogOpt) {
	o = &LogOpt{
		Level: "info",
	}

	fill(o, pfix)

	return
}
