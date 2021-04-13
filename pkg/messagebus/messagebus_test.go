package messagebus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
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
		setStorer    func(qs QueueStorer)
		process      func(ctx context.Context, qm QueueMessage) error
	}
)

var (
	sqlQueueSettings    = QueueSettings{ID: 1, Handler: "sql", Queue: "sql"}
	foobarQueueSettings = QueueSettings{ID: 1, Handler: "foobar", Queue: "foobar"}
	foobarQueueMessage  = QueueMessage{ID: 1, Queue: "foobar", Payload: []byte(`{}`), Created: now()}
	logger              = zap.NewNop()
)

func Test_messageBusRegister(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	mb := New(logger, mockDispatcher{})
	mb.Register(ctx, &QueueSettings{Handler: "foobar", Queue: "foobar"}, &mockQueueHandler{})

	req.NotEmpty(mb.queues)
	req.NotEmpty(mb.queues["foobar"])
	req.Empty(mb.queues["non_existing_queue"])
}

func Test_dispatchEvents(t *testing.T) {
	var (
		ctx = context.Background()
		tcc = []struct {
			name   string
			expect string
			fn     func(mb *messageBus)
		}{
			{
				name:   "dispatch enabled",
				expect: "dispatched",
				fn: func(mb *messageBus) {
					foobarQueueSettings.Meta.DispatchEvents = true
				},
			},
			{
				name:   "dispatch disabled",
				expect: "",
				fn: func(mb *messageBus) {
					foobarQueueSettings.Meta.DispatchEvents = false
				},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			tst := ""

			mb := New(logger, mockDispatcher{dispatch: func(ctx context.Context, ev eventbus.Event) {
				tst = "dispatched"
			}})

			// prepare
			tc.fn(mb)

			mb.queues["foobar"] = &Queue{
				settings: foobarQueueSettings,
				dispatch: make(chan []byte),
			}

			mb.dispatchEvents(ctx)
			mb.queues["foobar"].dispatch <- []byte("trigger this chan")

			req.Equal(0, len(mb.queues["foobar"].dispatch))
			req.Equal(tc.expect, tst)
		})
	}

}

func Test_updateProcessedMessages(t *testing.T) {

	req := require.New(t)
	ctx := context.Background()

	mb := New(logger, mockDispatcher{})

	numProcessed := 0

	mb.queues["foobar"] = &Queue{
		settings: foobarQueueSettings,
		dispatch: make(chan []byte),
	}

	mb.updateProcessedMessages(ctx)

	mb.queues["foobar"].processed <- &foobarQueueMessage

	req.Equal(1, numProcessed)
}

func Test_readListenerPoll(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	mb := New(logger, mockDispatcher{dispatch: func(ctx context.Context, ev eventbus.Event) {}})

	ticker := time.NewTicker(time.Millisecond)

	go func() {
		for ; true; <-ticker.C {
		}
	}()

	mb.queues["foobar"] = &Queue{
		dispatch:  make(chan []byte),
		out:       make(chan []byte),
		settings:  foobarQueueSettings,
		processed: make(chan *QueueMessage),
		handler: &mockQueueHandler{
			read: func(ctx context.Context) (QueueMessageSet, error) {
				ticker.Stop()
				return QueueMessageSet{&foobarQueueMessage}, nil
			},
			ticker: func(ctx context.Context) <-chan time.Time {
				return ticker.C
			},
			notification: func(ctx context.Context) <-chan interface{} {
				return make(<-chan interface{})
			},
		},
	}

	mb.Listen(ctx)

	req.Equal(foobarQueueMessage.Payload, <-mb.queues["foobar"].out)
	req.Equal(&foobarQueueMessage, <-mb.queues["foobar"].processed)
	req.Equal(foobarQueueMessage.Payload, <-mb.queues["foobar"].dispatch)
}

func Test_writeListener(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()
	w := sync.WaitGroup{}

	mb := New(logger, mockDispatcher{dispatch: func(ctx context.Context, ev eventbus.Event) {}})

	mockDb := [][]byte{}
	expectedDb := [][]byte{
		[]byte("mock payload"),
		[]byte("second mock payload"),
	}

	mb.queues["foobar"] = &Queue{
		in:       make(chan []byte),
		settings: foobarQueueSettings,
		handler: &mockQueueHandler{
			write: func(ctx context.Context, p []byte) error {
				mockDb = append(mockDb, p)
				w.Done()
				return nil
			},
			ticker: func(ctx context.Context) <-chan time.Time {
				return make(<-chan time.Time)
			},
			notification: func(ctx context.Context) <-chan interface{} {
				return make(<-chan interface{})
			},
		},
	}

	mb.Listen(ctx)

	w.Add(2)
	mb.queues["foobar"].in <- expectedDb[0]
	mb.queues["foobar"].in <- expectedDb[1]
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
func (mh *mockQueueHandler) SetStorer(qs QueueStorer) {
	mh.setStorer(qs)
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
