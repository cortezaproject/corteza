package server

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
	"context"
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type ModuleHandlers struct {
	Module ModuleAPI
}

// Internal API interface
type ModuleAPI interface {
	List(context.Context, *ModuleListRequest) (interface{}, error)
	Create(context.Context, *ModuleCreateRequest) (interface{}, error)
	Read(context.Context, *ModuleReadRequest) (interface{}, error)
	Edit(context.Context, *ModuleEditRequest) (interface{}, error)
	Delete(context.Context, *ModuleDeleteRequest) (interface{}, error)
	ContentList(context.Context, *ModuleContentListRequest) (interface{}, error)
	ContentCreate(context.Context, *ModuleContentCreateRequest) (interface{}, error)
	ContentRead(context.Context, *ModuleContentReadRequest) (interface{}, error)
	ContentEdit(context.Context, *ModuleContentEditRequest) (interface{}, error)
	ContentDelete(context.Context, *ModuleContentDeleteRequest) (interface{}, error)
}

// HTTP API interface
type ModuleHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	ContentList(http.ResponseWriter, *http.Request)
	ContentCreate(http.ResponseWriter, *http.Request)
	ContentRead(http.ResponseWriter, *http.Request)
	ContentEdit(http.ResponseWriter, *http.Request)
	ContentDelete(http.ResponseWriter, *http.Request)
}
