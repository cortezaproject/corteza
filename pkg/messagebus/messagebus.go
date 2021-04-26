package messagebus

import (
	"context"
	"fmt"
	"sync"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
)

var (
	// global service
	gMbus *messageBus
)

type (
	messageBus struct {
		opts   *options.MessagebusOpt
		logger *zap.Logger
		mutex  sync.Mutex
		queues QueueSet

		in     chan message
		quit   chan bool
		reload chan bool
	}

	QueueStorer interface {
		SearchMessagebusQueueSettings(ctx context.Context, f QueueSettingsFilter) (QueueSettingsSet, QueueSettingsFilter, error)
		SearchMessagebusQueueMessages(ctx context.Context, f QueueMessageFilter) (QueueMessageSet, QueueMessageFilter, error)
		CreateMessagebusQueueMessage(ctx context.Context, rr ...*QueueMessage) error
		UpdateMessagebusQueueMessage(ctx context.Context, rr ...*QueueMessage) error
	}
)

func Service() *messageBus {
	return gMbus
}

func Set(m *messageBus) {
	gMbus = m
}

// Setup handles the singleton service
func Setup(opts *options.MessagebusOpt, log *zap.Logger) {
	if gMbus != nil {
		return
	}

	gMbus = New(opts, log)
}

func New(opts *options.MessagebusOpt, logger *zap.Logger) *messageBus {
	return &messageBus{
		opts:   opts,
		logger: logger,
		queues: QueueSet{},

		in:     make(chan message),
		quit:   make(chan bool),
		reload: make(chan bool),
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
	// check in channel and act accordingly
	go func() {
		for {
			select {
			case t := <-mb.in:
				q := mb.queue(t.q)
				log := mb.logger.With(zap.String("queue", t.q))

				if q == nil {
					log.Warn("could not get queue settings")
					break
				}

				// get the consumer
				err := q.consumer.Write(ctx, t.p)

				if err != nil {
					log.Warn("could not add message to queue", zap.Error(err))
					break
				}

				log.Debug("wrote payload to queue")

			case <-mb.quit:
				mb.logger.Debug("quitting from messagebus listener")
				return

			case <-ctx.Done():
				mb.logger.Debug("quitting from messagebus listener")
				return
			}
		}
	}()
}

// Watch checks the channel for restart and loads the queues
// and adds the listeners again (the same process as on boot)
func (mb *messageBus) Watch(ctx context.Context, storer QueueStorer) {
	go func() {
		for {
			select {
			case <-mb.reload:
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
		// quit all the listeners
		mb.quit <- true

		// re-add all the watchers
		mb.reload <- true
	}()
}

func (mb *messageBus) Push(q string, p []byte) {
	if !mb.opts.Enabled {
		mb.logger.Debug("message will not be sent, messagebus disabled", zap.String("queue", q))
		return
	}

	mb.in <- message{p: p, q: q}
}

func (mb *messageBus) Register(ctx context.Context, qs *QueueSettings, consumer Consumer) {
	// associate consumer with the queue
	mb.queues[qs.Queue] = &Queue{
		settings: *qs,
		consumer: consumer,
	}
}

func (mb *messageBus) queue(q string) *Queue {
	if _, ok := mb.queues[q]; ok {
		return mb.queues[q]
	}

	return nil
}

func (mb *messageBus) initQueues(ctx context.Context, storer QueueStorer) error {
	list, _, err := storer.SearchMessagebusQueueSettings(ctx, QueueSettingsFilter{})

	if err != nil {
		return err
	}

	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	// empty first
	mb.queues = make(QueueSet)

	// add to list of consumers
	for _, qs := range list {
		c, err := mb.initConsumer(ctx, *qs)
		mb.logger.Debug("initializing queue", zap.String("queue", qs.Queue))

		if err != nil {
			mb.logger.Warn("could not init consumer for queue", zap.String("queue", qs.Queue), zap.Error(err))
			continue
		}

		if _, is := c.(Storer); is {
			c.SetStore(storer)
		}

		mb.Register(ctx, qs, c)
	}

	return nil
}

// initHandler returns a new instance for a specific handler
func (mb *messageBus) initConsumer(ctx context.Context, settings QueueSettings) (consumer Consumer, err error) {
	handle := settings.Consumer

	switch handle {
	case string(ConsumerEventbus):
		consumer = NewEventbusConsumer(settings)
		return

	case string(ConsumerStore):
		consumer = NewStoreConsumer(settings)
		return

	default:
		err = fmt.Errorf("message queue consumer %s not implemented", handle)
		return
	}
}
