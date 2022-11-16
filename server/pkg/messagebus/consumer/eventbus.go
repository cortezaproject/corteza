package consumer

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
	st "github.com/cortezaproject/corteza/server/system/types"
)

type (
	Dispatcher interface {
		Dispatch(ctx context.Context, ev eventbus.Event)
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
	}

	EventbusConsumer struct {
		queue      string
		handle     types.ConsumerType
		dispatcher Dispatcher
		servicer   types.QueueEventBuilder
	}
)

func NewEventbusConsumer(q string, servicer types.QueueEventBuilder) *EventbusConsumer {
	h := &EventbusConsumer{
		queue:      q,
		handle:     types.ConsumerEventbus,
		dispatcher: eventbus.Service(),
		servicer:   servicer,
	}

	return h
}

func (cq *EventbusConsumer) Write(ctx context.Context, p []byte) error {
	cq.dispatcher.Dispatch(ctx, cq.servicer.CreateQueueEvent(cq.queue, p))
	return nil
}

func makeEvent(q string, p []byte) *st.QueueMessage {
	return &st.QueueMessage{
		Queue:   q,
		Payload: p,
	}
}
