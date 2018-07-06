package sam

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type WebsocketHandlers struct {
	*Websocket
}

func (WebsocketHandlers) new() *WebsocketHandlers {
	return &WebsocketHandlers{
		Websocket{}.New(),
	}
}

// Internal API interface
type WebsocketAPI interface {
	Client(*websocketClientRequest) (interface{}, error)
}

// HTTP API interface
type WebsocketHandlersAPI interface {
	Client(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ WebsocketHandlersAPI = &WebsocketHandlers{}
var _ WebsocketAPI = &Websocket{}
