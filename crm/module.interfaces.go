package crm

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type ModuleHandlers struct {
	*Module
}

func (ModuleHandlers) new() *ModuleHandlers {
	return &ModuleHandlers{
		Module{}.new(),
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
}

// Compile time check to see if we implement the interfaces
var _ ModuleHandlersAPI = &ModuleHandlers{}
var _ ModuleAPI = &Module{}
