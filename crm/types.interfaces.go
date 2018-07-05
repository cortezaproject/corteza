package crm

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type TypesHandlers struct {
	*Types
}

func (TypesHandlers) new() *TypesHandlers {
	return &TypesHandlers{
		Types{}.new(),
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
}

// Compile time check to see if we implement the interfaces
var _ TypesHandlersAPI = &TypesHandlers{}
var _ TypesAPI = &Types{}