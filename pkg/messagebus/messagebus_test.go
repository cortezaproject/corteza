package messagebus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	mockStorer struct{}

	// only use this one to get .messages
	// from mockQueueHandler
	messageTesterHandler interface {
		GetAllMessages() [][]byte
	}

	mockDispatcher struct {
		register func(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) (p uintptr)
		dispatch func(ctx context.Context, ev eventbus.Event)
		waitFor  func(ctx context.Context, ev eventbus.Event) (err error)
	}

	mockQueueHandler struct {
		messages     [][]byte
		notification func(ctx context.Context) <-chan interface{}
		ticker       func(ctx context.Context) <-chan time.Time
		read         func(ctx context.Context) (QueueMessageSet, error)
		write        func(ctx context.Context, p []byte) error
		setStore     func(qs QueueStorer)
		process      func(ctx context.Context, qm QueueMessage) error
	}
)

var (
	mOptions            = &options.MessagebusOpt{Enabled: true, LogEnabled: true}
	sqlQueueSettings    = QueueSettings{ID: 1, Consumer: "sql", Queue: "sql"}
	foobarQueueSettings = QueueSettings{ID: 1, Consumer: "foobar", Queue: "foobar"}
	foobarQueueMessage  = QueueMessage{ID: 1, Queue: "foobar", Payload: []byte(`{}`), Created: now()}
	logger              = zap.NewNop()
)

func Test_messageBusRegister(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	mb := New(mOptions, logger)
	mb.Register(ctx, &QueueSettings{Consumer: "foobar", Queue: "foobar"}, &mockQueueHandler{})

	req.NotEmpty(mb.queues)
	req.NotEmpty(mb.queues["foobar"])
	req.Empty(mb.queues["non_existing_queue"])
}

func Test_consume(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()
	w := sync.WaitGroup{}

	mb := New(mOptions, logger)

	mockDb := [][]byte{}
	expectedDb := [][]byte{
		[]byte("mock payload"),
		[]byte("second mock payload"),
	}

	mb.queues[foobarQueueSettings.Queue] = &Queue{
		settings: foobarQueueSettings,
		consumer: &mockQueueHandler{
			write: func(ctx context.Context, p []byte) error {
				mockDb = append(mockDb, p)
				w.Done()
				return nil
			},
		},
	}

	mb.Listen(ctx)

	w.Add(2)
	mb.Push("foobar", expectedDb[0])
	mb.Push("foobar", expectedDb[1])
	w.Wait()

	req.Equal(expectedDb, mockDb)
}

func (mh *mockQueueHandler) Notification(ctx context.Context) <-chan interface{} {
	return mh.notification(ctx)
}
func (mh *mockQueueHandler) Ticker(ctx context.Context) <-chan time.Time {
	return mh.ticker(ctx)
}
func (mh *mockQueueHandler) Read(ctx context.Context) (QueueMessageSet, error) {
	return mh.read(ctx)
}
func (mh *mockQueueHandler) Write(ctx context.Context, p []byte) error {
	return mh.write(ctx, p)
}
func (mh *mockQueueHandler) SetStore(qs QueueStorer) {
	mh.setStore(qs)
}
func (mh *mockQueueHandler) Process(ctx context.Context, qm QueueMessage) error {
	return mh.process(ctx, qm)
}

func (md mockDispatcher) Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) (p uintptr) {
	return md.register(h, ops...)
}
func (md mockDispatcher) Dispatch(ctx context.Context, ev eventbus.Event) {
	md.dispatch(ctx, ev)
}
func (md mockDispatcher) WaitFor(ctx context.Context, ev eventbus.Event) (err error) {
	return md.waitFor(ctx, ev)
}

func tickOnce(tt time.Ticker) {
	go func() {
		for ; true; <-tt.C {
		}
	}()
}
