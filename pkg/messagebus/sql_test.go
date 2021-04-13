package messagebus

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type (
	mockClient struct {
		_get     func(context.Context, string) (QueueMessageSet, error)
		_add     func(context.Context, string, []byte) (err error)
		_process func(context.Context, QueueMessage) (err error)

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
		_get: func(c context.Context, q string) (QueueMessageSet, error) {
			return fooMessageSet, nil
		},
		_add: func(c context.Context, q string, p []byte) error {
			return nil
		},
		_process: func(c context.Context, m QueueMessage) error {
			return nil
		},
	}

	unsuccessfulClient = mockClient{
		_get: func(c context.Context, s string) (QueueMessageSet, error) {
			return QueueMessageSet{}, errors.New("could not get messages")
		},
		_add: func(c context.Context, q string, p []byte) error {
			return errors.New("could not write messages")
		},
		_process: func(c context.Context, m QueueMessage) error {
			return errors.New("could not process messages")
		},
	}
)

func Test_handlerSqlRead(t *testing.T) {
	var (
		ctx = context.Background()
		tcc = []struct {
			name   string
			err    error
			expect QueueMessageSet
			client SqlClient
		}{
			{
				name:   "read success",
				err:    nil,
				expect: fooMessageSet,
				client: &successfulClient,
			},
			{
				name:   "read error",
				err:    errors.New("could not get messages"),
				expect: QueueMessageSet{},
				client: &unsuccessfulClient,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			h := SqlQueueHandler{queue: "foobar", handle: HandlerSql, client: tc.client}

			set, err := h.Read(ctx)

			req.Equal(tc.err, err)
			req.Equal(tc.expect, set)
		})
	}
}

func Test_handlerSqlWrite(t *testing.T) {
	var (
		ctx = context.Background()
		tcc = []struct {
			name    string
			err     error
			payload []byte
			expect  QueueMessageSet
			client  SqlClient
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
			h := SqlQueueHandler{queue: "foobar", handle: HandlerSql, client: tc.client}

			err := h.Write(ctx, []byte("foo bar"))

			req.Equal(tc.err, err)
		})
	}
}

func Test_handlerSqlProcess(t *testing.T) {
	var (
		ctx = context.Background()
		tcc = []struct {
			name    string
			err     error
			payload []byte
			expect  QueueMessageSet
			client  SqlClient
		}{
			{
				name:   "process success",
				err:    nil,
				expect: fooMessageSet,
				client: &successfulClient,
			},
			{
				name:   "process error",
				err:    errors.New("could not process messages"),
				expect: QueueMessageSet{},
				client: &unsuccessfulClient,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			h := SqlQueueHandler{queue: "foobar", handle: HandlerSql, client: tc.client}

			err := h.Process(ctx, QueueMessage{})

			req.Equal(tc.err, err)
		})
	}
}

func Test_handlerSqlTicker(t *testing.T) {
	var (
		tcc = []struct {
			name     string
			expect   interface{}
			settings QueueSettings
		}{
			{
				name:     "ticker created",
				expect:   1,
				settings: QueueSettings{Queue: "foobar", Meta: QueueSettingsMeta{PollDelay: makeDelay(time.Second)}},
			},
			{
				name:     "ticker empty",
				expect:   0,
				settings: QueueSettings{Queue: "foobar", Meta: QueueSettingsMeta{}},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			h := NewSqlHandler(tc.settings)
			c := h.Ticker(context.Background())

			req.Equal(reflect.TypeOf(c), reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(time.Time{})))
			req.Equal(tc.expect, cap(c))
		})
	}
}

func (mc *mockClient) get(ctx context.Context, q string) (set QueueMessageSet, err error) {
	return mc._get(ctx, q)
}

func (mc *mockClient) add(ctx context.Context, q string, p []byte) error {
	return mc._add(ctx, q, p)
}

func (mc *mockClient) process(ctx context.Context, m QueueMessage) error {
	return mc._process(ctx, m)
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
