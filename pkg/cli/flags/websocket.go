package flags

import (
	"time"

	"github.com/spf13/cobra"
)

type (
	WebsocketOpt struct {
		Timeout     time.Duration
		PingTimeout time.Duration
		PingPeriod  time.Duration
	}
)

func Websocket(cmd *cobra.Command, pfix string) (o *WebsocketOpt) {
	o = &WebsocketOpt{}

	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)

	bindDuration(cmd, &o.Timeout,
		pFlag(pfix, "websocket-timeout"), timeout,
		"Websocket connection timeout")

	bindDuration(cmd, &o.PingTimeout,
		pFlag(pfix, "websocket-ping-timeout"), pingTimeout,
		"Websocket connection ping timeout")

	bindDuration(cmd, &o.PingPeriod,
		pFlag(pfix, "websocket-ping-period"), pingPeriod,
		"Websocket connection ping period (should be lower than timeout)")

	return
}
