package messagebus

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type (
	mockClient struct {
		get     func(context.Context, string) (QueueMessageSet, error)
		add     func(context.Context, string, []byte) (err error)
		process func(context.Context, QueueMessage) (err error)

		setStorer func(QueueStorer)
		getStorer func() QueueStorer
	}
)

var (
	fooMessageSet = QueueMessageSet{
		{ID: 1},
		{ID: 2},
	}

	successfulClient = mockClient{
		get: func(c context.Context, q string) (QueueMessageSet, error) {
			return fooMessageSet, nil
		},
		add: func(c context.Context, q string, p []byte) error {
			return nil
		},
		process: func(c context.Context, m QueueMessage) error {
			return nil
		},
	}

	unsuccessfulClient = mockClient{
		get: func(c context.Context, s string) (QueueMessageSet, error) {
			return QueueMessageSet{}, errors.New("could not get messages")
		},
		add: func(c context.Context, q string, p []byte) error {
			return errors.New("could not write messages")
		},
		process: func(c context.Context, m QueueMessage) error {
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
			expect  QueueMessageSet
			client  StoreClient
		}{
			{
				name:   "write success",
				err:    nil,
				expect: fooMessageSet,
				client: &successfulClient,
			},
			{
				name:   "write error",
				err:    errors.New("could not write messages"),
				expect: QueueMessageSet{},
				client: &unsuccessfulClient,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			h := StoreConsumer{queue: "foobar", handle: ConsumerStore, client: tc.client}

			err := h.Write(ctx, []byte("foo bar"))

			req.Equal(tc.err, err)
		})
	}
}

func (mc *mockClient) Get(ctx context.Context, q string) (set QueueMessageSet, err error) {
	return mc.get(ctx, q)
}

func (mc *mockClient) Add(ctx context.Context, q string, p []byte) error {
	return mc.add(ctx, q, p)
}

func (mc *mockClient) Process(ctx context.Context, m QueueMessage) error {
	return mc.process(ctx, m)
}

func (mc *mockClient) GetStorer() QueueStorer {
	return mc.getStorer()
}

func (mc *mockClient) SetStorer(s QueueStorer) {
	mc.setStorer(s)
}

func makeDelay(d time.Duration) *time.Duration {
	return &d
}
