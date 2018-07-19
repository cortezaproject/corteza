package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Message struct{}

func (Message) New() *Message {
	return &Message{}
}

func (*Message) Edit(ctx context.Context, r *server.MessageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}

func (*Message) Attach(ctx context.Context, r *server.MessageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}

func (*Message) Remove(ctx context.Context, r *server.MessageRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.remove")
}

func (*Message) Read(ctx context.Context, r *server.MessageReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.read")
}

func (*Message) Search(ctx context.Context, r *server.MessageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}

func (*Message) Pin(ctx context.Context, r *server.MessagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}

func (*Message) Flag(ctx context.Context, r *server.MessageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}
