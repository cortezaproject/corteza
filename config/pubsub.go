package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"time"
)

type (
	PubSub struct {
		Mode            string
		RedisAddr       string
		PollingInterval time.Duration

		Timeout     time.Duration
		PingTimeout time.Duration
		PingPeriod  time.Duration
	}
)

var pubsub *PubSub

func (c *PubSub) Validate() error {
	switch c.Mode {
	case "redis":
		if c.Mode == "redis" && c.RedisAddr == "" {
			return errors.New("No host defined for mode=redis, PubSub.Redis is empty")
		}
	case "poll":
	default:
		return errors.Errorf("Unknown PubSub.Mode: %s", c.Mode)
	}
	return nil
}

func (*PubSub) Init(prefix ...string) *PubSub {
	if pubsub != nil {
		return pubsub
	}

	pubsub = new(PubSub)
	pubsub.Timeout = 15 * time.Second
	pubsub.PingTimeout = 60 * time.Second
	pubsub.PingPeriod = (pubsub.PingTimeout * 10) / 9

	flag.StringVar(&pubsub.Mode, "pubsub", "poll", "Pubsub mode (poll, redis)")
	flag.StringVar(&pubsub.RedisAddr, "pubsub-redis", "", "Redis Pub/Sub hostname")
	flag.DurationVar(&pubsub.PollingInterval, "pubsub-poll-interval", 3*time.Second, "Pub/Sub polling interval (3s, 12m, 3h...)")

	return pubsub
}
