package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `user.go`, `user.util.go` or `user_test.go` to
	implement your API calls, helper functions and tests. The file `user.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type UserHandlers struct {
	User UserAPI
}

// Internal API interface
type UserAPI interface {
	Search(context.Context, *UserSearchRequest) (interface{}, error)
	Message(context.Context, *UserMessageRequest) (interface{}, error)
}

// HTTP API interface
type UserHandlersAPI interface {
	Search(http.ResponseWriter, *http.Request)
	Message(http.ResponseWriter, *http.Request)
}
