package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `websocket.go`, `websocket.util.go` or `websocket_test.go` to
	implement your API calls, helper functions and tests. The file `websocket.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

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
