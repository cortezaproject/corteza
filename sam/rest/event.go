package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Event struct{}

func (Event) New() *Event {
	return &Event{}
}

func (ctrl *Event) Edit(ctx context.Context, r *server.EventEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}

func (ctrl *Event) Attach(ctx context.Context, r *server.EventAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}

func (ctrl *Event) Remove(ctx context.Context, r *server.EventRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}

func (ctrl *Event) Read(ctx context.Context, r *server.EventReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}

func (ctrl *Event) Search(ctx context.Context, r *server.EventSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}

func (ctrl *Event) Pin(ctx context.Context, r *server.EventPinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}

func (ctrl *Event) Flag(ctx context.Context, r *server.EventFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Event.....")
}
