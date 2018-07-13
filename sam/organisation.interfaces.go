package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `organisation.go`, `organisation.util.go` or `organisation_test.go` to
	implement your API calls, helper functions and tests. The file `organisation.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type OrganisationHandlers struct {
	*Organisation
}

func (OrganisationHandlers) new() *OrganisationHandlers {
	return &OrganisationHandlers{
		Organisation{}.New(),
	}
}

// Internal API interface
type OrganisationAPI interface {
	List(*organisationListRequest) (interface{}, error)
	Create(*organisationCreateRequest) (interface{}, error)
	Edit(*organisationEditRequest) (interface{}, error)
	Remove(*organisationRemoveRequest) (interface{}, error)
	Read(*organisationReadRequest) (interface{}, error)
	Archive(*organisationArchiveRequest) (interface{}, error)
}

// HTTP API interface
type OrganisationHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Remove(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Archive(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ OrganisationHandlersAPI = &OrganisationHandlers{}
var _ OrganisationAPI = &Organisation{}
