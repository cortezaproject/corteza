package rest

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `channel.go`, `channel.util.go` or `channel_test.go` to
	implement your API calls, helper functions and tests. The file `channel.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type ChannelHandlers struct {
	Channel ChannelAPI
}

// Internal API interface
type ChannelAPI interface {
	List(*ChannelListRequest) (interface{}, error)
	Create(*ChannelCreateRequest) (interface{}, error)
	Edit(*ChannelEditRequest) (interface{}, error)
	Read(*ChannelReadRequest) (interface{}, error)
	Delete(*ChannelDeleteRequest) (interface{}, error)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// HTTP API interface
type ChannelHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}
