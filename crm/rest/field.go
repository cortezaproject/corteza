package rest

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `field.go`, `field.util.go` or `field_test.go` to
	implement your API calls, helper functions and tests. The file `field.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type FieldHandlers struct {
	Field FieldAPI
}

// Internal API interface
type FieldAPI interface {
	List(*FieldListRequest) (interface{}, error)
	Type(*FieldTypeRequest) (interface{}, error)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// HTTP API interface
type FieldHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Type(http.ResponseWriter, *http.Request)
}
