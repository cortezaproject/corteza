package flags

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// Prefixes flag
func pFlag(pfix, name string) string {
	if pfix != "" {
		name = pfix + "-" + name
	}

	return name
}

// Converts input (flag-name) into ENVIRONMENTAL_VARIABLE_KEY
func envKey(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}

func BindString(cmd *cobra.Command, v *string, flag, def string, desc string) {
	if env, has := os.LookupEnv(envKey(flag)); has {
		def = cast.ToString(env)
	}

	cmd.Flags().StringVar(v, flag, def, desc)
}

func BindBool(cmd *cobra.Command, v *bool, flag string, def bool, desc string) {
	if env, has := os.LookupEnv(envKey(flag)); has {
		def = cast.ToBool(env)
	}

	cmd.Flags().BoolVar(v, flag, def, desc)
}

func BindInt(cmd *cobra.Command, v *int, flag string, def int, desc string) {
	if env, has := os.LookupEnv(envKey(flag)); has {
		def = cast.ToInt(env)
	}

	cmd.Flags().IntVar(v, flag, def, desc)
}

func BindDuration(cmd *cobra.Command, v *time.Duration, flag string, def time.Duration, desc string) {
	if env, has := os.LookupEnv(envKey(flag)); has {
		def = cast.ToDuration(env)
	}

	cmd.Flags().DurationVar(v, flag, def, desc)
}
