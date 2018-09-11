package repository

import (
	"context"
	"log"
)

type (
	PubSub struct {
		client PubSubClient
	}

	PubSubClient interface {
		Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, message []byte) error) error
		Publish(ctx context.Context, channel string, message string) error
	}

	PubSubPayload struct {
		Channel string
		Message []byte
	}
)

var pubsub PubSubClient

func (PubSub) New() (PubSubClient, error) {
	// return singleton client
	if pubsub != nil {
		return pubsub, nil
	}

	// validate configs and fall back to poll mode on error
	if err := flags.PubSub.Validate(); err != nil {
		log.Printf("[pubsub] An error occured when creating PubSub instance: %+v", err)
		log.Println("[pubsub] Reverting back to 'poll' and trying again")
		flags.PubSub.Mode = "poll"
		if err := flags.PubSub.Validate(); err != nil {
			return nil, err
		}
	}

	// store the singleton instance
	save := func(client PubSubClient, err error) (PubSubClient, error) {
		if err != nil {
			return nil, err
		}
		pubsub = client
		return pubsub, nil
	}

	// create isntances based on mode
	if flags.PubSub.Mode == "redis" {
		return save(PubSubRedis{}.New(flags.PubSub))
	}
	return save(PubSubMemory{}.New(flags.PubSub))
}

func (ps *PubSub) Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, message []byte) error) error {
	return ps.client.Subscribe(ctx, channel, onStart, onMessage)
}

func (ps *PubSub) Publish(ctx context.Context, channel, message string) error {
	return ps.client.Publish(ctx, channel, message)
}
