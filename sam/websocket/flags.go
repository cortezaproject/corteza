package websocket

import (
	"time"
)

type (
	configuration struct {
		writeTimeout time.Duration
		pingTimeout  time.Duration
		pingPeriod   time.Duration
	}
)

var config configuration

func (c configuration) validate() error {
	return nil
}

// Flags should be called from main to register flags
func Flags() {
	config.writeTimeout = 15 * time.Second
	config.pingTimeout = 120 * time.Second
	config.pingPeriod = (config.pingTimeout * 10) / 9
}
