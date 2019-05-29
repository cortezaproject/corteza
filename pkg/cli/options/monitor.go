package options

import (
	"time"
)

type (
	MonitorOpt struct {
		Interval time.Duration `env:"MONITOR_INTERVAL"`
	}
)

func Monitor(pfix string) (o *MonitorOpt) {
	o = &MonitorOpt{
		Interval: 300 * time.Second,
	}
	fill(o, pfix)

	return
}
