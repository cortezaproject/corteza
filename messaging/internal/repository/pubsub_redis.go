package repository

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

type PubSubRedis struct {
	addr string

	timeout     time.Duration
	pingTimeout time.Duration
	pingPeriod  time.Duration
}

func (PubSubRedis) New(addr string, to, pt, pp time.Duration) *PubSubRedis {
	return &PubSubRedis{
		addr:        addr,
		timeout:     to,
		pingTimeout: pt,
		pingPeriod:  pp,
	}
}

func (ps *PubSubRedis) dial() (redis.Conn, error) {
	return redis.Dial(
		"tcp",
		ps.addr,
		redis.DialReadTimeout(ps.pingTimeout+ps.timeout),
		redis.DialWriteTimeout(ps.timeout),
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
		case <-time.After(ps.pingPeriod):
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
