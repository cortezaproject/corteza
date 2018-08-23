package websocket

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"time"
)

type (
	Configuration struct {
		writeTimeout time.Duration
		pingTimeout  time.Duration
		pingPeriod   time.Duration

		pubSubMode     string
		pubSubRedis    string
		pubSubInterval time.Duration
	}
)

// Validate returns error if there is an issue with the config
func (c *Configuration) Validate() error {
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

// Init binds flags to websocket configuration structure
func (c *Configuration) Init() {
	c.writeTimeout = 15 * time.Second
	c.pingTimeout = 120 * time.Second
	c.pingPeriod = (c.pingTimeout * 10) / 9

	flag.StringVar(&c.pubSubMode, "pubsub", "poll", "Pubsub mode (poll, redis)")
	flag.StringVar(&c.pubSubRedis, "pubsub-redis", "", "Redis Pub/Sub hostname")
	flag.DurationVar(&c.pubSubInterval, "pubsub-poll-interval", 3*time.Second, "Pub/Sub polling interval (3s, 12m, 3h...)")
}
