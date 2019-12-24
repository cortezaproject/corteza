package options

import (
	"time"
)

type (
	PubSubOpt struct {
		Mode string `env:"PUBSUB_MODE"`

		// Mode
		PollingInterval time.Duration `env:"PUBSUB_POLLING_INTERVAL"`

		// Redis
		RedisAddr        string        `env:"PUBSUB_REDIS_ADDR"`
		RedisTimeout     time.Duration `env:"PUBSUB_REDIS_TIMEOUT"`
		RedisPingTimeout time.Duration `env:"PUBSUB_REDIS_PING_TIMEOUT"`
		RedisPingPeriod  time.Duration `env:"PUBSUB_REDIS_PING_PERIOD"`
	}
)

func PubSub(pfix string) (o *PubSubOpt) {
	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)

	o = &PubSubOpt{
		Mode:             "poll",
		PollingInterval:  timeout,
		RedisAddr:        "redis:6379",
		RedisTimeout:     timeout,
		RedisPingTimeout: pingTimeout,
		RedisPingPeriod:  pingPeriod,
	}

	fill(o, pfix)

	return
}
