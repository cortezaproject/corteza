package crm

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type ModulesHandlers struct {
	*Modules
}

func (ModulesHandlers) new() *ModulesHandlers {
	return &ModulesHandlers{
		Modules{}.new(),
	}
}

// Internal API interface
type ModulesAPI interface {
	List(*modulesListRequest) (interface{}, error)
	Edit(*modulesEditRequest) (interface{}, error)
	ContentList(*modulesContentListRequest) (interface{}, error)
	ContentEdit(*modulesContentEditRequest) (interface{}, error)
	ContentDelete(*modulesContentDeleteRequest) (interface{}, error)
}

// HTTP API interface
type ModulesHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	ContentList(http.ResponseWriter, *http.Request)
	ContentEdit(http.ResponseWriter, *http.Request)
	ContentDelete(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ ModulesHandlersAPI = &ModulesHandlers{}
var _ ModulesAPI = &Modules{}
