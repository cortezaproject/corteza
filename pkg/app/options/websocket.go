package options

import (
	"time"
)

type (
	WebsocketOpt struct {
		Timeout     time.Duration `env:"WEBSOCKET_TIMEOUT"`
		PingTimeout time.Duration `env:"WEBSOCKET_PING_TIMEOUT"`
		PingPeriod  time.Duration `env:"WEBSOCKET_PING_PERIOD"`
	}
)

func Websocket(pfix string) (o *WebsocketOpt) {
	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)

	o = &WebsocketOpt{
		Timeout:     timeout,
		PingTimeout: pingTimeout,
		PingPeriod:  pingPeriod,
	}

	fill(o, pfix)

	return
}
