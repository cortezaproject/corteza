package sam

import (
	"github.com/pkg/errors"
)

var _ = errors.Wrap

func (*Message) Edit(r *messageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}

func (*Message) Attach(r *messageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}

func (*Message) Remove(r *messageRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.remove")
}

func (*Message) Read(r *messageReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.read")
}

func (*Message) Search(r *messageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}

func (*Message) Pin(r *messagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}

func (*Message) Flag(r *messageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}
