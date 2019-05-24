package flags

import (
	"time"

	"github.com/spf13/cobra"
)

type (
	PubSubOpt struct {
		Mode string

		// Mode
		PollingInterval time.Duration

		// Redis
		RedisAddr        string
		RedisTimeout     time.Duration
		RedisPingTimeout time.Duration
		RedisPingPeriod  time.Duration
	}
)

func PubSub(cmd *cobra.Command, pfix string) (o *PubSubOpt) {
	o = &PubSubOpt{}

	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)

	BindString(cmd, &o.Mode,
		pFlag(pfix, "pubsub-mode"), "poll",
		"Pub/Sub mode (poll, redis")

	BindDuration(cmd, &o.RedisPingTimeout,
		pFlag(pfix, "pubsub-polling-interval"), timeout,
		"Sub/Sub polling interval")

	BindString(cmd, &o.RedisAddr,
		pFlag(pfix, "pubsub-redis-addr"), "redis:6379",
		"Pub/Sub mode (poll, redis")

	BindDuration(cmd, &o.RedisTimeout,
		pFlag(pfix, "pubsub-redis-timeout"), timeout,
		"Websocket connection timeout")

	BindDuration(cmd, &o.RedisPingTimeout,
		pFlag(pfix, "pubsub-redis-ping-timeout"), pingTimeout,
		"Pub/Sub connection ping timeout")

	BindDuration(cmd, &o.RedisPingPeriod,
		pFlag(pfix, "pubsub-redis-ping-period"), pingPeriod,
		"Pub/Sub connection ping period (should be lower than timeout)")

	return
}
