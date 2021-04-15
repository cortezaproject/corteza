package messagebus

import (
	"context"
	"time"
)

type (
	StoreConsumer struct {
		queue  string
		handle ConsumerType
		client StoreClient
		poll   *time.Ticker
	}
)

func NewStoreConsumer(settings QueueSettings) *StoreConsumer {
	h := &StoreConsumer{
		queue:  settings.Queue,
		handle: ConsumerStore,
		client: &sClient{},
	}

	return h
}

func (cq *StoreConsumer) Write(ctx context.Context, p []byte) error {
	return cq.client.Add(ctx, cq.queue, p)
}

func (cq *StoreConsumer) SetStore(s QueueStorer) {
	cq.client.SetStorer(s)
}
