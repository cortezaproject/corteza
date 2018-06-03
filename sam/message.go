package sam

import (
	"github.com/pkg/errors"
)

func (m *Message) Edit(r *messageEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.edit")
}
func (m *Message) Attach(r *messageAttachRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.attach")
}
func (m *Message) Remove(r *messageRemoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.remove")
}
func (m *Message) Read(r *messageReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.read")
}
func (m *Message) Search(r *messageSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.search")
}
func (m *Message) Pin(r *messagePinRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.pin")
}
func (m *Message) Flag(r *messageFlagRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Message.flag")
}
