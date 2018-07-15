package rest

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `message.go`, `message.util.go` or `message_test.go` to
	implement your API calls, helper functions and tests. The file `message.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type MessageHandlers struct {
	Message MessageAPI
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

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
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
