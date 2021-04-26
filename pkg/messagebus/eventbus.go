package messagebus

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Dispatcher interface {
		Dispatch(ctx context.Context, ev eventbus.Event)
	}

	EventbusConsumer struct {
		queue      string
		handle     ConsumerType
		dispatcher Dispatcher
	}
)

func NewEventbusConsumer(settings QueueSettings) *EventbusConsumer {
	h := &EventbusConsumer{
		queue:      settings.Queue,
		handle:     ConsumerEventbus,
		dispatcher: eventbus.Service(),
	}

	return h
}

func (cq *EventbusConsumer) Write(ctx context.Context, p []byte) error {
	cq.dispatcher.Dispatch(ctx, event.QueueOnMessage(makeEvent(cq.queue, p)))
	return nil
}

func (cq *EventbusConsumer) SetStore(s QueueStorer) {}

func makeEvent(q string, p []byte) *types.QueueMessage {
	return &types.QueueMessage{
		Queue:   q,
		Payload: string(p),
	}
}
