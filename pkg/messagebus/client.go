package messagebus

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
)

type (
	Client interface {
		add(ctx context.Context, q string, p []byte) (err error)
		get(ctx context.Context, q string) (list QueueMessageSet, err error)
		process(ctx context.Context, m QueueMessage) (err error)
	}

	SqlClient interface {
		Client

		SetStorer(QueueStorer)
		GetStorer() QueueStorer
	}

	sClient struct {
		storer QueueStorer
	}
)

func (c *sClient) process(ctx context.Context, m QueueMessage) (err error) {
	err = c.GetStorer().UpdateMessagebusQueuemessage(ctx, &QueueMessage{
		ID:        m.ID,
		Queue:     m.Queue,
		Payload:   m.Payload,
		Created:   m.Created,
		Processed: now(),
	})

	return
}

func (c *sClient) add(ctx context.Context, q string, payload []byte) (err error) {
	err = c.storer.CreateMessagebusQueuemessage(ctx, &QueueMessage{
		ID:      nextID(),
		Queue:   q,
		Created: now(),
		Payload: payload,
	})

	return
}

func (c *sClient) get(ctx context.Context, q string) (list QueueMessageSet, err error) {
	list, _, err = c.storer.SearchMessagebusQueuemessages(ctx, QueueMessageFilter{
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
