package repository

import (
	"context"
)

type (
	PubSub struct {
		client pubSubModule
	}

	PubSubClient interface {
		pubSubModule
		Event(ctx context.Context, message string) error
	}

	pubSubModule interface {
		Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, message []byte) error) error
		Publish(ctx context.Context, channel string, message string) error
	}

	PubSubPayload struct {
		Channel string
		Message []byte
	}
)

// var pubsub *PubSub

func (PubSub) New() *PubSub {
	panic("pending reimplementation")
	// @todo should be configured much earlier, do not depend on flags here
	// // return singleton client
	// if pubsub != nil {
	// 	return pubsub
	// }
	//
	// // store the singleton instance
	// save := func(client pubSubModule) *PubSub {
	// 	pubsub = &PubSub{client}
	// 	return pubsub
	// }

	// // create isntances based on mode
	// if flags != nil && flags.PubSub.Mode == "redis" {
	// 	return save(PubSubRedis{}.New(flags.PubSub))
	// }
	// return save(PubSubMemory{}.New(flags.PubSub))
}

func (ps *PubSub) Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, message []byte) error) error {
	return ps.client.Subscribe(ctx, channel, onStart, onMessage)
}

func (ps *PubSub) Event(ctx context.Context, message string) error {
	return ps.Publish(ctx, "events", message)
}

func (ps *PubSub) Publish(ctx context.Context, channel, message string) error {
	return ps.client.Publish(ctx, channel, message)
}
