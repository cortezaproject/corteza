package messagebus

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
)

type (
	Client interface {
		Add(ctx context.Context, q string, p []byte) (err error)
		Get(ctx context.Context, q string) (list QueueMessageSet, err error)
		Process(ctx context.Context, m QueueMessage) (err error)
	}

	StoreClient interface {
		Client

		SetStorer(QueueStorer)
		GetStorer() QueueStorer
	}

	sClient struct {
		storer QueueStorer
	}
)

func (c *sClient) Process(ctx context.Context, m QueueMessage) (err error) {
	err = c.GetStorer().UpdateMessagebusQueueMessage(ctx, &QueueMessage{
		ID:        m.ID,
		Queue:     m.Queue,
		Payload:   m.Payload,
		Created:   m.Created,
		Processed: now(),
	})

	return
}

func (c *sClient) Add(ctx context.Context, q string, payload []byte) (err error) {
	err = c.storer.CreateMessagebusQueueMessage(ctx, &QueueMessage{
		ID:      nextID(),
		Queue:   q,
		Created: now(),
		Payload: payload,
	})

	return
}

func (c *sClient) Get(ctx context.Context, q string) (list QueueMessageSet, err error) {
	list, _, err = c.storer.SearchMessagebusQueueMessages(ctx, QueueMessageFilter{
		Queue:     q,
		Processed: filter.StateExcluded})

	return
}

func (c *sClient) GetStorer() QueueStorer {
	return c.storer
}

func (c *sClient) SetStorer(s QueueStorer) {
	c.storer = s
}

func nextID() uint64 {
	return id.Next()
}

func now() *time.Time {
	t := time.Now()
	return &t
}
