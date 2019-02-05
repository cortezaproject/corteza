package repository

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/crusttech/crust/internal/config"
)

type PubSubRedis struct {
	config *config.PubSub
}

func (PubSubRedis) New(config *config.PubSub) *PubSubRedis {
	return &PubSubRedis{config}
}

func (ps *PubSubRedis) dial() (redis.Conn, error) {
	return redis.Dial(
		"tcp",
		ps.config.RedisAddr,
		redis.DialReadTimeout(ps.config.PingTimeout+ps.config.Timeout),
		redis.DialWriteTimeout(ps.config.Timeout),
	)
}

func (ps *PubSubRedis) Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, payload []byte) error) error {
	// main redis connection
	conn, err := ps.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	// pubsub object
	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(redis.Args{}.Add(channel)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case 1:
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	cleanup := func(err error) error {
		psc.Unsubscribe()
		return err
	}

	for {
		select {
		case <-time.After(ps.config.PingPeriod):
			if err := psc.Ping(""); err != nil {
				return cleanup(err)
			}
		case <-ctx.Done():
			return cleanup(ctx.Err())
		case err := <-done:
			return err
		}
	}
}

func (ps *PubSubRedis) Publish(ctx context.Context, channel, message string) error {
	conn, err := ps.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("PUBLISH", channel, message)
	return err
}
