package messagebus

import (
	"context"
	"fmt"
	"sync"

	"github.com/cortezaproject/corteza/server/pkg/messagebus/consumer"
	"github.com/cortezaproject/corteza/server/pkg/messagebus/store"
	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"go.uber.org/zap"
)

var (
	// global service
	gMbus *messageBus
)

type (
	messageBus struct {
		opts      *options.MessagebusOpt
		qservicer types.QueueServicer
		logger    *zap.Logger
		mutex     sync.Mutex
		queues    types.QueueSet

		in     chan types.Message
		quit   chan bool
		reload chan bool
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
		queues: types.QueueSet{},

		in:     make(chan types.Message),
		quit:   make(chan bool),
		reload: make(chan bool),
	}
}

// Init takes care of preloading the queue and creating
// a connection to the store of their choice
func (mb *messageBus) Init(ctx context.Context, storer types.QueueServicer) {
	// set store now, on New() we do not have store yet
	mb.qservicer = storer

	mb.initQueues(ctx)
}

// Listen sets read and write channel listeners,
// used on boot
func (mb *messageBus) Listen(ctx context.Context) {
	// check in channel and act accordingly
	go func() {
		for {
			select {
			case t := <-mb.in:
				q := mb.queue(t.Q)
				log := mb.logger.With(zap.String("queue", t.Q))

				if q == nil {
					log.Warn("could not get queue settings")
					break
				}

				// get the consumer
				err := q.Consumer.Write(ctx, t.P)

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
func (mb *messageBus) Watch(ctx context.Context, storer types.QueueServicer) {
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

	mb.logger.Debug("pushing message", zap.String("queue", q))

	mb.in <- types.Message{P: p, Q: q}
}

func (mb *messageBus) Register(ctx context.Context, qs *types.Queue) {
	mb.queues[qs.Name] = qs
}

func (mb *messageBus) queue(q string) *types.Queue {
	if _, ok := mb.queues[q]; ok {
		return mb.queues[q]
	}

	return nil
}

func (mb *messageBus) initQueues(ctx context.Context) error {
	list, _, err := mb.qservicer.SearchQueues(ctx, types.QueueFilter{})

	if err != nil {
		return err
	}

	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	// empty first
	mb.queues = make(types.QueueSet)

	// add to list of consumers
	for _, q := range list {
		mb.logger.Debug("initializing queue", zap.String("queue", q.Queue))

		c, err := mb.initConsumer(ctx, q.Queue, q.Consumer)

		if err != nil {
			mb.logger.Warn("could not init consumer for queue", zap.String("queue", q.Queue), zap.Error(err))
			continue
		}

		if _, is := c.(store.Storer); is {
			c.(store.Storer).SetStore(mb.qservicer)
		}

		mb.Register(ctx, &types.Queue{
			Consumer: c,
			Name:     q.Queue,
			Meta:     types.QueueMeta(q.Meta),
		})
	}

	return nil
}

// initHandler returns a new instance for a specific handler
func (mb *messageBus) initConsumer(ctx context.Context, q string, c string) (cns types.Consumer, err error) {
	switch c {
	case string(types.ConsumerEventbus):
		cns = consumer.NewEventbusConsumer(q, mb.qservicer)
		return

	case string(types.ConsumerStore):
		cns = consumer.NewStoreConsumer(q, mb.qservicer)
		return

	default:
		err = fmt.Errorf("message queue consumer %s not implemented", c)
		return
	}
}
