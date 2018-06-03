package sam

import (
	"github.com/pkg/errors"
)

func (m *Message) Edit(r *MessageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}
func (m *Message) Attach(r *MessageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}
func (m *Message) Remove(r *MessageRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.remove")
}
func (m *Message) Read(r *MessageReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.read")
}
func (m *Message) Search(r *MessageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}
func (m *Message) Pin(r *MessagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}
func (m *Message) Flag(r *MessageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}
