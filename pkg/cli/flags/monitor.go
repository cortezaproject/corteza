package flags

import (
	"github.com/spf13/cobra"
)

type (
	MonitorOpt struct {
		Interval int
	}
)

func Monitor(cmd *cobra.Command, pfix string) (o *MonitorOpt) {
	o = &MonitorOpt{}

	bindInt(cmd, &o.Interval,
		pFlag(pfix, "monitor-interval"), 300,
		"Monitor interval (seconds, 0 = disable)")

	return
}
