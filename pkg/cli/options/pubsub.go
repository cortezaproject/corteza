package options

import (
	"time"
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

func PubSub(pfix string) (o *PubSubOpt) {
	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)

	o = &PubSubOpt{
		Mode:             EnvString(pfix, "PUBSUB_MODE", "poll"),
		PollingInterval:  EnvDuration(pfix, "PUBSUB_POLLING_INTERVAL", timeout),
		RedisAddr:        EnvString(pfix, "PUBSUB_REDIS_ADDR", "redis:6379"),
		RedisTimeout:     EnvDuration(pfix, "PUBSUB_REDIS_TIMEOUT", timeout),
		RedisPingTimeout: EnvDuration(pfix, "PUBSUB_REDIS_PING_TIMEOUT", pingTimeout),
		RedisPingPeriod:  EnvDuration(pfix, "PUBSUB_REDIS_PING_PERIOD", pingPeriod),
	}

	return
}
