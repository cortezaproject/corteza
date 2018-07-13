package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type ModuleHandlers struct {
	*Module
}

func (ModuleHandlers) new() *ModuleHandlers {
	return &ModuleHandlers{
		Module{}.New(),
	}
}

// Internal API interface
type ModuleAPI interface {
	List(*moduleListRequest) (interface{}, error)
	Edit(*moduleEditRequest) (interface{}, error)
	ContentList(*moduleContentListRequest) (interface{}, error)
	ContentEdit(*moduleContentEditRequest) (interface{}, error)
	ContentDelete(*moduleContentDeleteRequest) (interface{}, error)
}

// HTTP API interface
type ModuleHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	ContentList(http.ResponseWriter, *http.Request)
	ContentEdit(http.ResponseWriter, *http.Request)
	ContentDelete(http.ResponseWriter, *http.Request)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// Compile time check to see if we implement the interfaces
var _ ModuleHandlersAPI = &ModuleHandlers{}
var _ ModuleAPI = &Module{}
