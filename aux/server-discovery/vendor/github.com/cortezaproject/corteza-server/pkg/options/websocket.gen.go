package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/websocket.yaml

import (
	"time"
)

type (
	WebsocketOpt struct {
		LogEnabled  bool          `env:"WEBSOCKET_LOG_ENABLED"`
		Timeout     time.Duration `env:"WEBSOCKET_TIMEOUT"`
		PingTimeout time.Duration `env:"WEBSOCKET_PING_TIMEOUT"`
		PingPeriod  time.Duration `env:"WEBSOCKET_PING_PERIOD"`
	}
)

// Websocket initializes and returns a WebsocketOpt with default values
func Websocket() (o *WebsocketOpt) {
	o = &WebsocketOpt{
		Timeout:     15 * time.Second,
		PingTimeout: 120 * time.Second,
		PingPeriod:  ((120 * time.Second) * 9) / 10,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Websocket) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
