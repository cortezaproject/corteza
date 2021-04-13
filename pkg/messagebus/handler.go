package messagebus

import (
	"context"
	"time"
)

const (
	HandlerCorteza handler = "corteza"
	HandlerNoop    handler = "noop"
	HandlerRedis   handler = "redis"
	HandlerSql     handler = "sql"
)

type (
	handler string

	Handler interface {
		ReadHandler
		WriteHandler
	}

	ReadHandler interface {
		Reader
		Subscriber
		Poller
		Storer
	}

	WriteHandler interface {
		Writer
		Storer
	}

	Subscriber interface {
		Notification(ctx context.Context) <-chan interface{}
	}

	Poller interface {
		Ticker(ctx context.Context) <-chan time.Time
	}

	Reader interface {
		Read(ctx context.Context) (QueueMessageSet, error)
	}
	Writer interface {
		Write(ctx context.Context, p []byte) error
	}

	Storer interface {
		// todo
		SetStorer(QueueStorer)
		Process(context.Context, QueueMessage) error
	}
)
