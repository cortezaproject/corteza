package server

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
	"context"
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type MessageHandlers struct {
	Message MessageAPI
}

// Internal API interface
type MessageAPI interface {
	Create(context.Context, *MessageCreateRequest) (interface{}, error)
	Edit(context.Context, *MessageEditRequest) (interface{}, error)
	Delete(context.Context, *MessageDeleteRequest) (interface{}, error)
	Attach(context.Context, *MessageAttachRequest) (interface{}, error)
	Search(context.Context, *MessageSearchRequest) (interface{}, error)
	Pin(context.Context, *MessagePinRequest) (interface{}, error)
	Unpin(context.Context, *MessageUnpinRequest) (interface{}, error)
	Flag(context.Context, *MessageFlagRequest) (interface{}, error)
	Deflag(context.Context, *MessageDeflagRequest) (interface{}, error)
	React(context.Context, *MessageReactRequest) (interface{}, error)
	Unreact(context.Context, *MessageUnreactRequest) (interface{}, error)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// HTTP API interface
type MessageHandlersAPI interface {
	Create(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Attach(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Pin(http.ResponseWriter, *http.Request)
	Unpin(http.ResponseWriter, *http.Request)
	Flag(http.ResponseWriter, *http.Request)
	Deflag(http.ResponseWriter, *http.Request)
	React(http.ResponseWriter, *http.Request)
	Unreact(http.ResponseWriter, *http.Request)
}
