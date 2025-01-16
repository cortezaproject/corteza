package consumer

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/servicebus"
	"go.uber.org/zap"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
)

type (
	// make this common type
	// Dispatcher interface {
	// 	Dispatch(ctx context.Context, ev eventbus.Event)
	// 	WaitFor(ctx context.Context, ev eventbus.Event) (err error)
	// }

	QueueEventBuilder interface {
		CreateQueueEvent(string, []byte) eventbus.Event
	}

	ServicebusConsumer struct {
		logger *zap.Logger

		// @fixme:not sure about this
		// Azure service bus connection string
		connStr string

		queue  string
		handle types.ConsumerType
		// explorer servicebus.Client
		poll *time.Ticker

		// @fixme: not sure about this weather to keep it here or not
		servicer QueueEventBuilder

		s Dispatcher
	}
)

func NewServicebusConsumer(ctx context.Context, logger *zap.Logger, connStr string, q string, servicer types.QueueEventBuilder) (sb *ServicebusConsumer, err error) {
	sb = &ServicebusConsumer{
		logger: logger,

		queue:    q,
		handle:   types.ConsumerServicebus,
		s:        servicebus.Service(logger, connStr),
		servicer: servicer,
	}

	return
}

func (sb *ServicebusConsumer) Write(ctx context.Context, p []byte) (err error) {
	err = sb.client.SendMessage(ctx, sb.queue, p)
	if err != nil {
		return
	}

	sb.dispatcher.Dispatch(ctx, sb.servicer.CreateQueueEvent(sb.queue, p))

	// @fixme: invoke workflow events directly
	_ = sb.servicer.CreateQueueEvent(sb.queue, p)

	// sb.dispatcher.WaitFor(ctx, )

	// if err = sb.dispatcher.WaitFor(ctx, event(upd, old, m, ns, rve, nil)); err != nil {
	// 	return
	// }

	return
}

func (sb *ServicebusConsumer) GetConsumerType() string {
	if sb.handle == "" {
		return ""
	}

	return string(sb.handle)
}
