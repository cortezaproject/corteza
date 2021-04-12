package messagebus

import (
	"context"
	"time"
)

const (
	HandlerCorteza HandlerType = "corteza"
	HandlerNoop    HandlerType = "noop"
	HandlerRedis   HandlerType = "redis"
	HandlerSql     HandlerType = "sql"
)

type (
	HandlerType string

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

func HandlerTypes() []HandlerType {
	return []HandlerType{
		HandlerCorteza,
		HandlerRedis,
		HandlerSql,
	}
}
