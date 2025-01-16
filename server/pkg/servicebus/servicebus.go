package servicebus

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

type (
	client interface {
		CreateQueue(context.Context, string) error
		DeleteQueue(context.Context, string) error
		SendMessage(context.Context, string, []byte) error
		// GetMessage(context.Context, int, string)
	}

	servicebus struct {
		// waitgroup for dispatch
		wg *sync.WaitGroup

		// Read & write locking
		l *sync.RWMutex

		// client for service bus
		client client
	}
)

func Service(logger *zap.Logger, connStr string) (sb *servicebus, err error) {
	sb = &servicebus{
		wg: &sync.WaitGroup{},
		l:  &sync.RWMutex{},
	}

	sb.client, err = newClient(logger, connStr)
	if err != nil {
		return
	}

	return
}

// func (sb *servicebus) CreateQueue(ctx context.Context, q string) (err error) {
//
// }
