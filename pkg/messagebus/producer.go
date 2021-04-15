package messagebus

import (
	"context"
	"time"
)

type (
	Producer interface {
		Reader
		Subscriber
		Poller
		Storer
	}

	Subscriber interface {
		Subscribe(ctx context.Context) <-chan interface{}
	}

	Poller interface {
		Poll(ctx context.Context) <-chan time.Time
	}

	Reader interface {
		Read(ctx context.Context) (QueueMessageSet, error)
	}
)
