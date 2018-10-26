package config

import (
	"github.com/namsral/flag"
)

type (
	Monitor struct {
		Interval int
	}
)

var monitor *Monitor

func (c *Monitor) Validate() error {
	return nil
}

func (*Monitor) Init(prefix ...string) *Monitor {
	if monitor != nil {
		return monitor
	}

	monitor = new(Monitor)
	flag.IntVar(&monitor.Interval, "monitor-interval", 300, "Monitor interval (seconds, 0 = disable)")
	return monitor
}
