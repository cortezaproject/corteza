package flags

import (
	"github.com/spf13/cobra"
)

type (
	LogOpt struct {
		Level string
		JSON  bool
	}
)

func Log(cmd *cobra.Command) (o *LogOpt) {
	o = &LogOpt{}

	bindString(cmd, &o.Level,
		"log-level", "info",
		"Log level (debug, info, warn, error, panic, fatal)")

	bindBool(cmd, &o.JSON,
		"log-json", true,
		"Log in JSON format")

	return
}
