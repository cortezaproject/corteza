package websocket

import (
	"time"
)

type (
	Config struct {
		LogEnabled  bool
		Timeout     time.Duration
		PingTimeout time.Duration
		PingPeriod  time.Duration
	}
)
