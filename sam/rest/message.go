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

func (ctrl *Message) Create(ctx context.Context, r *server.MessageCreateRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.create")
}

func (ctrl *Message) Edit(ctx context.Context, r *server.MessageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}

func (ctrl *Message) Delete(ctx context.Context, r *server.MessageDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.delete")
}

func (ctrl *Message) Attach(ctx context.Context, r *server.MessageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}

func (ctrl *Message) Search(ctx context.Context, r *server.MessageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}

func (ctrl *Message) Pin(ctx context.Context, r *server.MessagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}

func (ctrl *Message) Unpin(ctx context.Context, r *server.MessageUnpinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.unpin")
}

func (ctrl *Message) Flag(ctx context.Context, r *server.MessageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}

func (ctrl *Message) Deflag(ctx context.Context, r *server.MessageDeflagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.deflag")
}

func (ctrl *Message) React(ctx context.Context, r *server.MessageReactRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.react")
}

func (ctrl *Message) Unreact(ctx context.Context, r *server.MessageUnreactRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.unreact")
}
