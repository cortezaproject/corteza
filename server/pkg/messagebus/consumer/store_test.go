package consumer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/messagebus/store"
	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
	"github.com/stretchr/testify/require"
)

type (
	mockClient struct {
		add     func(context.Context, string, []byte) (err error)
		process func(context.Context, uint64, types.QueueMessage) (err error)

		setStore func(types.QueueStorer)
		getStore func() types.QueueStorer
	}
)

var (
	successfulClient = mockClient{
		add: func(c context.Context, q string, p []byte) error {
			return nil
		},
		process: func(c context.Context, u uint64, qm types.QueueMessage) (err error) { return },
	}

	unsuccessfulClient = mockClient{
		add: func(c context.Context, q string, p []byte) error {
			return errors.New("could not write messages")
		},
		process: func(c context.Context, u uint64, qm types.QueueMessage) (err error) {
			return errors.New("could not process messages")
		},
	}
)

func Test_handlerSqlWrite(t *testing.T) {
	var (
		ctx = context.Background()
		tcc = []struct {
			name    string
			err     error
			payload []byte
			client  store.StoreClient
		}{
			{
				name:   "write success",
				err:    nil,
				client: &successfulClient,
			},
			{
				name:   "write error",
				err:    errors.New("could not write messages"),
				client: &unsuccessfulClient,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			h := StoreConsumer{queue: "foobar", handle: types.ConsumerStore, client: tc.client}

			err := h.Write(ctx, []byte("foo bar"))

			req.Equal(tc.err, err)
		})
	}
}

func (mc *mockClient) Add(ctx context.Context, q string, p []byte) error {
	return mc.add(ctx, q, p)
}

func (mc *mockClient) Process(ctx context.Context, ID uint64, m types.QueueMessage) error {
	return mc.process(ctx, ID, m)
}

func (mc *mockClient) GetStore() types.QueueStorer {
	return mc.getStore()
}

func (mc *mockClient) SetStore(s types.QueueStorer) {
	mc.setStore(s)
}

func makeDelay(d time.Duration) *time.Duration {
	return &d
}
