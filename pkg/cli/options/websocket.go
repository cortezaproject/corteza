package options

import (
	"time"
)

type (
	WebsocketOpt struct {
		Timeout     time.Duration
		PingTimeout time.Duration
		PingPeriod  time.Duration
	}
)

func Websocket(pfix string) (o *WebsocketOpt) {
	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)

	o = &WebsocketOpt{
		Timeout:     EnvDuration(pfix, "WEBSOCKET_TIMEOUT", timeout),
		PingTimeout: EnvDuration(pfix, "WEBSOCKET_PING_TIMEOUT", pingTimeout),
		PingPeriod:  EnvDuration(pfix, "WEBSOCKET_PING_PERIOD", pingPeriod),
	}

	return
}
