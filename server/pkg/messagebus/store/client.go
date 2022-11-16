package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
)

type (
	Client interface {
		Add(context.Context, string, []byte) error
		Process(context.Context, uint64, types.QueueMessage) error
	}

	StoreClient interface {
		Client
		Storer
	}

	sClient struct {
		store types.QueueStorer
	}
)

func NewClient(s types.QueueStorer) *sClient {
	return &sClient{store: s}
}

func (c *sClient) Process(ctx context.Context, ID uint64, m types.QueueMessage) (err error) {
	err = c.store.ProcessQueueMessage(ctx, ID, m)

	return
}

func (c *sClient) Add(ctx context.Context, q string, payload []byte) (err error) {
	err = c.store.CreateQueueMessage(ctx, types.QueueMessage{
		Queue:   q,
		Payload: payload,
	})

	return
}

func (c *sClient) GetStore() types.QueueStorer {
	return c.store
}

func (c *sClient) SetStore(s types.QueueStorer) {
	c.store = s
}
