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
	Edit(*messageEditRequest) (interface{}, error)
	Attach(*messageAttachRequest) (interface{}, error)
	Remove(*messageRemoveRequest) (interface{}, error)
	Read(*messageReadRequest) (interface{}, error)
	Search(*messageSearchRequest) (interface{}, error)
	Pin(*messagePinRequest) (interface{}, error)
	Flag(*messageFlagRequest) (interface{}, error)
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
