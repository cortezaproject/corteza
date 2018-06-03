package sam

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type MessageHandlers struct {
	*Message
}

func (MessageHandlers) new() *MessageHandlers {
	return &MessageHandlers{
		Message{}.new(),
	}
}

// Internal API interface
type MessageAPI interface {
	Edit(*MessageEditRequest) (interface{}, error)
	Attach(*MessageAttachRequest) (interface{}, error)
	Remove(*MessageRemoveRequest) (interface{}, error)
	Read(*MessageReadRequest) (interface{}, error)
	Search(*MessageSearchRequest) (interface{}, error)
	Pin(*MessagePinRequest) (interface{}, error)
	Flag(*MessageFlagRequest) (interface{}, error)
}

// HTTP API interface
type MessageHandlersAPI interface {
	Edit(http.ResponseWriter, *http.Request)
	Attach(http.ResponseWriter, *http.Request)
	Remove(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Pin(http.ResponseWriter, *http.Request)
	Flag(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ MessageHandlersAPI = &MessageHandlers{}
var _ MessageAPI = &Message{}
