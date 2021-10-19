package types

import (
	"context"
	"time"
)

type (
	Producer interface {
		Reader
		Subscriber
		Poller
	}

	Subscriber interface {
		Subscribe(ctx context.Context) <-chan interface{}
	}

	Poller interface {
		Poll(ctx context.Context) <-chan time.Time
	}

	Reader interface {
		Read(ctx context.Context) ([]QueueMessage, error)
	}
)
