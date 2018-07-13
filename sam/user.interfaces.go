package sam

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
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type UserHandlers struct {
	*User
}

func (UserHandlers) new() *UserHandlers {
	return &UserHandlers{
		User{}.New(),
	}
}

// Internal API interface
type UserAPI interface {
	Login(*userLoginRequest) (interface{}, error)
	Search(*userSearchRequest) (interface{}, error)
}

// HTTP API interface
type UserHandlersAPI interface {
	Login(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// Compile time check to see if we implement the interfaces
var _ UserHandlersAPI = &UserHandlers{}
var _ UserAPI = &User{}
