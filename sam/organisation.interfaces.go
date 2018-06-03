package sam

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type OrganisationHandlers struct {
	*Organisation
}

func (OrganisationHandlers) new() *OrganisationHandlers {
	return &OrganisationHandlers{
		Organisation{}.new(),
	}
}

// Internal API interface
type OrganisationAPI interface {
	Edit(*OrganisationEditRequest) (interface{}, error)
	Remove(*OrganisationRemoveRequest) (interface{}, error)
	Read(*OrganisationReadRequest) (interface{}, error)
	Search(*OrganisationSearchRequest) (interface{}, error)
	Archive(*OrganisationArchiveRequest) (interface{}, error)
}

// HTTP API interface
type OrganisationHandlersAPI interface {
	Edit(http.ResponseWriter, *http.Request)
	Remove(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Archive(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ OrganisationHandlersAPI = &OrganisationHandlers{}
var _ OrganisationAPI = &Organisation{}
