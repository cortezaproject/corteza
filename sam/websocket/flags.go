package websocket

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"time"
)

type (
	configuration struct {
		writeTimeout time.Duration
		pingTimeout  time.Duration
		pingPeriod   time.Duration

		pubSubMode     string
		pubSubRedis    string
		pubSubInterval time.Duration
	}
)

var config configuration

func (c configuration) validate() error {
	switch c.pubSubMode {
	case "redis", "poll":
	default:
		return errors.Errorf("Unknown pubSubMode: %s", c.pubSubMode)
	}
	if c.pubSubMode == "redis" && c.pubSubRedis == "" {
		return errors.New("No host defined for mode=redis, pubSubRedis is empty")
	}
	return nil
}

// Flags should be called from main to register flags
func Flags() {
	config.writeTimeout = 15 * time.Second
	config.pingTimeout = 120 * time.Second
	config.pingPeriod = (config.pingTimeout * 10) / 9

	flag.StringVar(&config.pubSubMode, "pubsub", "poll", "Pubsub mode (poll, redis)")
	flag.StringVar(&config.pubSubRedis, "pubsub-redis", "", "Redis Pub/Sub hostname")
	flag.DurationVar(&config.pubSubInterval, "pubsub-poll-interval", 3*time.Second, "Pub/Sub polling interval (3s, 12m, 3h...)")
}
