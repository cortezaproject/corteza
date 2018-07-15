package sam

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/rest"
	_ "github.com/crusttech/crust/sam/types"
)

type Message struct{}

func (Message) New() *Message {
	return &Message{}
}

var _ = errors.Wrap

func (*Message) Edit(r *rest.MessageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}

func (*Message) Attach(r *rest.MessageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}

func (*Message) Remove(r *rest.MessageRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.remove")
}

func (*Message) Read(r *rest.MessageReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.read")
}

func (*Message) Search(r *rest.MessageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}

func (*Message) Pin(r *rest.MessagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}

func (*Message) Flag(r *rest.MessageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}
