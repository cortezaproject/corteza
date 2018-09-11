package config

import (
	"time"
)

type (
	Websocket struct {
		Timeout     time.Duration
		PingTimeout time.Duration
		PingPeriod  time.Duration
	}
)

var websocket *Websocket

func (c *Websocket) Validate() error {
	return nil
}

func (*Websocket) Init(prefix ...string) *Websocket {
	websocket = new(Websocket)
	websocket.Timeout = 15 * time.Second
	websocket.PingTimeout = 120 * time.Second
	websocket.PingPeriod = (websocket.PingTimeout * 10) / 9
	return websocket
}
