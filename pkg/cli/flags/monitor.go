package flags

import (
	"time"

	"github.com/spf13/cobra"
)

type (
	MonitorOpt struct {
		Interval time.Duration
	}
)

func Monitor(cmd *cobra.Command, pfix string) (o *MonitorOpt) {
	o = &MonitorOpt{}

	BindDuration(cmd, &o.Interval,
		pFlag(pfix, "monitor-interval"), 300*time.Second,
		"Monitor interval (0 = disable)")

	return
}
