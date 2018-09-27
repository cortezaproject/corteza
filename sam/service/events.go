package service

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	event struct {
		ctx context.Context

		events repository.EventsRepository
	}

	EventService interface {
		With(ctx context.Context) EventService
		Message(m *types.Message) error
		Channel(m *types.Channel) error
	}
)

func Event() EventService {
	return (&event{events: repository.Events()}).With(context.Background())
}

func (svc *event) With(ctx context.Context) EventService {
	return &event{
		ctx:    ctx,
		events: svc.events,
	}
}

func (svc *event) Message(m *types.Message) error {
	return svc.push(payload.Message(m), m.ChannelID)
}

func (svc *event) Channel(m *types.Channel) error {
	return svc.push(payload.Channel(m), m.ID)
}

func (svc *event) push(m outgoing.MessageEncoder, sub uint64) error {
	var enc, err = m.EncodeMessage()
	if err != nil {
		return err
	}

	println("PUSHING EVENT !!", string(enc))

	item := &types.EventQueueItem{Payload: enc}

	if sub > 0 {
		item.Subscriber = payload.Uint64toa(sub)
	}

	return svc.events.Push(svc.ctx, item)
}
