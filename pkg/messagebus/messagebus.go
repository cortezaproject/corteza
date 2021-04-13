package messagebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"go.uber.org/zap"
)

var (
	// global service
	gMbus *messageBus
)

type (
	messageBus struct {
		logger   *zap.Logger
		eventbus dispatcher
		mutex    sync.Mutex

		queues QueueSet
		w      chan bool
	}

	QueueStorer interface {
		SearchMessagebusQueuesettings(ctx context.Context, f QueueSettingsFilter) (QueueSettingsSet, QueueSettingsFilter, error)
		SearchMessagebusQueuemessages(ctx context.Context, f QueueMessageFilter) (QueueMessageSet, QueueMessageFilter, error)
		CreateMessagebusQueuemessage(ctx context.Context, rr ...*QueueMessage) error
		UpdateMessagebusQueuemessage(ctx context.Context, rr ...*QueueMessage) error
	}

	dispatcher interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

func Service() *messageBus {
	return gMbus
}

func Set(m *messageBus) {
	gMbus = m
}

// Setup handles the singleton service
func Setup(logger *zap.Logger, d dispatcher) {
	if gMbus != nil {
		return
	}

	gMbus = New(logger, d)
}

func New(logger *zap.Logger, d dispatcher) *messageBus {
	return &messageBus{
		eventbus: d,
		logger:   logger,
		queues:   QueueSet{},
		w:        make(chan bool),
	}
}

// Init takes care of preloading the queue and creating
// a connection to the store of their choice
func (mb *messageBus) Init(ctx context.Context, storer QueueStorer) {
	mb.initQueues(ctx, storer)
}

// Listen sets read and write channel listeners,
// used on boot
func (mb *messageBus) Listen(ctx context.Context) {
	for _, qq := range mb.queues {
		// add data listeners
		mb.readListener(ctx, qq)
		mb.writeListener(ctx, qq)

		// listen on out channel and call registered events
		mb.dispatchEvents(ctx, qq)

		// listen on processed channel and update message status in store
		mb.updateProcessedMessages(ctx, qq)
	}
}

// Watch checks the channel for restart and loads the queues
// and adds the listeners again (the same process as on boot)
func (mb *messageBus) Watch(ctx context.Context, storer QueueStorer) {
	go func() {
		for {
			select {
			case <-mb.w:
				mb.logger.Debug("reloading queues")

				// refresh queue list
				mb.Init(ctx, storer)

				// reload listeners, we removed them in ReloadQueues()
				mb.Listen(ctx)

			case <-ctx.Done():
				return
			}
		}
	}()
}

// ReloadQueues sends an error signal to all of the read and
// write watchers, dispatch and processed signal watchers
//
// Once the watchers are finished, the Init() and Listen()
// reruns with the new list of queues
func (mb *messageBus) ReloadQueues() {
	go func() {
		// first, exit all the dispatchers and processing watchers
		for _, qq := range mb.queues {
			mb.logger.Debug("stopping all listeners for queue", zap.String("queue", qq.settings.Queue))
			qq.err <- errors.New("stop")
			qq.err <- errors.New("stop")
			qq.err <- errors.New("stop")
			qq.err <- errors.New("stop")
			mb.logger.Debug("stopped all listeners for queue", zap.String("queue", qq.settings.Queue))
		}

		// re-add all the watchers
		mb.w <- true
	}()
}

func (mb *messageBus) Handler(q string) Handler {
	if _, ok := mb.queues[q]; ok {
		return mb.queues[q].handler
	}

	return nil
}

func (mb *messageBus) Register(ctx context.Context, qs *QueueSettings, handler Handler) {
	// associate handler with the queue
	mb.queues[qs.Queue] = &Queue{
		settings:  *qs,
		handler:   handler,
		in:        make(chan []byte, 10),
		out:       make(chan []byte),
		dispatch:  make(chan []byte),
		err:       make(chan error),
		processed: make(chan *QueueMessage),
	}
}

// dispatchEvents listens on dispatch channel and dispatches
// the event
func (mb *messageBus) dispatchEvents(ctx context.Context, queues ...*Queue) {
	for _, qq := range queues {
		go func(ctx context.Context, q *Queue) {
			for {
				select {
				case p := <-q.dispatch:
					if !q.settings.CanDispatch() {
						break
					}

					mb.logger.Debug("dispatching queue with payload", zap.String("queue", q.settings.Queue), zap.String("payload", string(p)))
					pp, _ := expr.NewString(p)
					mb.eventbus.Dispatch(ctx, event.QueueOnMessage(pp))

				case err := <-q.err:
					mb.logger.Debug("closing dispatch messages watcher", zap.String("queue", q.settings.Queue), zap.Error(err))
					return

				case <-ctx.Done():
					mb.logger.Debug("exiting dispatch messages watcher", zap.String("queue", q.settings.Queue))
					return
				}
			}
		}(ctx, qq)
	}
}

// updateProcessedMessages listens on processed channel and
// marks the message as read, works with Storer
func (mb *messageBus) updateProcessedMessages(ctx context.Context, queues ...*Queue) {
	for _, qq := range queues {
		go func(ctx context.Context, q *Queue) {
			for {
				select {
				case message := <-q.processed:
					err := q.handler.Process(ctx, *message)

					if err != nil {
						mb.logger.Debug("could not mark message as processed", zap.String("queue", message.Queue), zap.Error(err))
					} else {
						mb.logger.Debug("message processed", zap.String("queue", message.Queue))
					}

				case err := <-q.err:
					mb.logger.Debug("closing processed messages watcher", zap.String("queue", q.settings.Queue), zap.Error(err))
					return

				case <-ctx.Done():
					mb.logger.Debug("exiting processed messages watcher", zap.String("queue", q.settings.Queue))
					return
				}
			}
		}(ctx, qq)
	}
}

// Read returns the read channel for the messages
func (mb *messageBus) Read(q string) chan []byte {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	if _, ok := mb.queues[q]; ok {
		return mb.queues[q].out
	}

	return nil
}

// Write returns the write channel for the messages
func (mb *messageBus) Write(q string) chan []byte {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	if _, ok := mb.queues[q]; ok {
		return mb.queues[q].in
	}

	return nil
}

// writeListener is the main watcher for messagebus on pushing to queue
func (mb *messageBus) writeListener(ctx context.Context, queues ...*Queue) {
	// concurrent per-queue
	for _, qq := range queues {
		mb.logger.Debug("adding write listener for queue", zap.String("queue", qq.settings.Queue))
		go mb.writeToQueue(ctx, qq)
	}
}

// readListener is the main watcher for messagebus on reading from queue
func (mb *messageBus) readListener(ctx context.Context, queues ...*Queue) {
	// concurrent per-queue
	for _, qq := range queues {
		mb.logger.Debug("adding read listener for queue", zap.String("queue", qq.settings.Queue))
		go mb.readFromQueue(ctx, qq)
	}
}

// writeToQueue sends the payloads from the <-in channel
func (mb *messageBus) writeToQueue(ctx context.Context, q *Queue) {
	for {
		select {
		case p := <-q.in:
			err := q.handler.Write(ctx, p)

			if err != nil {
				mb.logger.Info("could not add message to queue", zap.String("queue", q.settings.Queue), zap.Error(err))
				break
			}

			mb.logger.Debug("wrote payload to queue", zap.String("queue", q.settings.Queue))

		case <-q.err:
			mb.logger.Info("closing message queue listener", zap.String("queue", q.settings.Queue), zap.Int("processed", len(q.processed)))
			return

		case <-ctx.Done():
			mb.logger.Info("closing message queue listener", zap.String("queue", q.settings.Queue))
			return
		}
	}
}

// readFromQueue gets the data from store via handler and writes to out
func (mb *messageBus) readFromQueue(ctx context.Context, q *Queue) {
	for {
		select {
		case <-mb.getNotification(ctx, q.handler):
			// todo
			break

		case <-mb.getTicker(ctx, q.handler):
			list, err := q.handler.Read(ctx)

			if err != nil {
				mb.logger.Info("could not read message from queue", zap.String("queue", q.settings.Queue), zap.Error(err))
				break
			}

			for _, message := range list {
				q.out <- []byte(message.Payload)
				mb.logger.Debug("added payload to out channel", zap.String("queue", q.settings.Queue))

				// processed = sent to out channel + dispatched (if enabled)
				// this should block dispatch until done
				q.processed <- message
				mb.logger.Debug("added processed signal channel", zap.String("queue", q.settings.Queue))

				// once message is consumed and processed,
				// dispatch the info about it
				q.dispatch <- []byte(message.Payload)
				mb.logger.Debug("added dispatch signal channel", zap.String("queue", q.settings.Queue))
			}

		case <-q.err:
			mb.logger.Info("closing message queue listener", zap.String("queue", q.settings.Queue), zap.Int("processed", len(q.processed)))
			return

		case <-ctx.Done():
			mb.logger.Info("closing message queue listener", zap.String("queue", q.settings.Queue), zap.Int("processed", len(q.processed)))
			return
		}
	}
}

func (mb *messageBus) initQueues(ctx context.Context, storer QueueStorer) error {
	list, _, err := storer.SearchMessagebusQueuesettings(ctx, QueueSettingsFilter{})

	if err != nil {
		return err
	}

	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	// add to list of handlers
	for _, qs := range list {
		handler, err := mb.initHandler(ctx, *qs)
		mb.logger.Debug("initializing queue", zap.String("queue", qs.Queue))

		if err != nil {
			mb.logger.Info("could not init handler for queue", zap.String("queue", qs.Queue), zap.Error(err))
			continue
		}

		// todo
		handler.SetStorer(storer)

		mb.Register(ctx, qs, handler)
	}

	return nil
}

// initHandler returns a new instance for a specific handler
func (mb *messageBus) initHandler(ctx context.Context, settings QueueSettings) (handler Handler, err error) {
	handle := settings.Handler

	switch handle {
	// case string(handlerCorteza):
	// 	handler = NewCortezaHandler(settings)
	// 	return

	case string(HandlerSql):
		handler = NewSqlHandler(settings)
		return

	// case string(HandlerRedis):
	// 	handler = NewRedisHandler(settings)
	// 	return

	default:
		err = fmt.Errorf("message queue handler %s not implemented", handle)
		return
	}
}

func (mb *messageBus) getNotification(ctx context.Context, sub Subscriber) <-chan interface{} {
	return sub.Notification(ctx)
}

func (mb *messageBus) getTicker(ctx context.Context, p Poller) <-chan time.Time {
	return p.Ticker(ctx)
}
