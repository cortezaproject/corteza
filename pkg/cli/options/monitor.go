package options

import (
	"time"
)

type (
	MonitorOpt struct {
		Interval time.Duration
	}
)

func Monitor(pfix string) (o *MonitorOpt) {
	o = &MonitorOpt{
		Interval: EnvDuration(pfix, "MONITOR_INTERVAL", 300*time.Second),
	}

	return
}
