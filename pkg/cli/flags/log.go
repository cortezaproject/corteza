package flags

import (
	"github.com/spf13/cobra"
)

type (
	// Logger's output leve is configured here, but
	// dev/prod configuration happens earlier
	LogOpt struct {
		Level string
	}
)

func Log(cmd *cobra.Command) (o *LogOpt) {
	o = &LogOpt{}

	BindString(cmd, &o.Level,
		"log-level", "info",
		"Log level (debug, info, warn, error, panic, fatal)")

	return
}
