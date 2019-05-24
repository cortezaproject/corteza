package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type PubSubMemory struct {
	pollingInterval time.Duration
	input           chan *PubSubPayload
}

func (PubSubMemory) New(pollingInterval time.Duration) *PubSubMemory {
	return &PubSubMemory{
		pollingInterval: pollingInterval,
		input:           make(chan *PubSubPayload, 512),
	}
}

func (ps *PubSubMemory) Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, payload []byte) error) error {
	polling := func() error {
		if err := onStart(); err != nil {
			return err
		}
		for {
			select {
			// context cancelled
			case <-ctx.Done():
				return ctx.Err()
			// triggered local event
			case msg := <-ps.input:
				if msg.Channel == channel {
					onMessage(msg.Channel, msg.Message)
				}
			// polling event
			case <-time.After(ps.pollingInterval):
				onMessage(channel, []byte("pubsub tick event"))
			}
		}
	}
	defer func() {
		close(ps.input)
	}()
	return polling()
}

func (ps *PubSubMemory) Publish(ctx context.Context, channel string, message string) (err error) {
	defer func() {
		// trying to send on closed channel panic, recover
		if r := recover(); r != nil {
			err = errors.Errorf("PubSubMemory.Publish: %+v", r)
		}
	}()

	payload := &PubSubPayload{
		channel,
		[]byte(message),
	}
	select {
	case ps.input <- payload:
	case <-time.After(50 * time.Millisecond):
		return errors.New("PubSubMemory.Publish: send timeout")
	}
	return nil
}
