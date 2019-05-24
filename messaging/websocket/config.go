package websocket

import (
	"time"
)

type (
	Config struct {
		Timeout     time.Duration
		PingTimeout time.Duration
		PingPeriod  time.Duration
	}
)
