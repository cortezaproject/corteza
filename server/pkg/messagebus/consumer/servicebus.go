package consumer

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/servicebus"
	stypes "github.com/cortezaproject/corteza/server/system/types"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
)

type (
	// make this common type
	// Dispatcher interface {
	// 	Dispatch(ctx context.Context, ev eventbus.Event)
	// 	WaitFor(ctx context.Context, ev eventbus.Event) (err error)
	// }

	ServicebusConsumer struct {
		queue  string
		handle types.ConsumerType
		client servicebus.Client
		poll   *time.Ticker

		// @fixme: not sure about this weather to keep it here or not
		servicer types.QueueEventBuilder
	}
)

func NewServicebusConsumer(q string, servicer types.QueueEventBuilder) *ServicebusConsumer {
	// fixme: implement error
	h := &ServicebusConsumer{
		queue:    q,
		handle:   types.ConsumerEventbus,
		client:   servicebus.NewClient(),
		servicer: servicer,
	}

	return h
}

func (sb *ServicebusConsumer) Write(ctx context.Context, p []byte) (err error) {
	err = sb.client.SendMessage(ctx, sb.servicer.CreateQueueEvent(sb.queue, p))

	return
}
