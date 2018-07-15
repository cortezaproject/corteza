package rest

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/rest/server"
	_ "github.com/crusttech/crust/sam/types"
)

type Message struct{}

func (Message) New() *Message {
	return &Message{}
}

var _ = errors.Wrap

func (*Message) Edit(r *server.MessageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}

func (*Message) Attach(r *server.MessageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}

func (*Message) Remove(r *server.MessageRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.remove")
}

func (*Message) Read(r *server.MessageReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.read")
}

func (*Message) Search(r *server.MessageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}

func (*Message) Pin(r *server.MessagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}

func (*Message) Flag(r *server.MessageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}
