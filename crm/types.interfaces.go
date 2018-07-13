package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `types.go`, `types.util.go` or `types_test.go` to
	implement your API calls, helper functions and tests. The file `types.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type TypesHandlers struct {
	*Types
}

func (TypesHandlers) new() *TypesHandlers {
	return &TypesHandlers{
		Types{}.New(),
	}
}

// Internal API interface
type TypesAPI interface {
	List(*typesListRequest) (interface{}, error)
	Type(*typesTypeRequest) (interface{}, error)
}

// HTTP API interface
type TypesHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Type(http.ResponseWriter, *http.Request)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// Compile time check to see if we implement the interfaces
var _ TypesHandlersAPI = &TypesHandlers{}
var _ TypesAPI = &Types{}
