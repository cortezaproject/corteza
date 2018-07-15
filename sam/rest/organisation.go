package rest

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
	Organisation OrganisationAPI
}

// Internal API interface
type OrganisationAPI interface {
	List(*OrganisationListRequest) (interface{}, error)
	Create(*OrganisationCreateRequest) (interface{}, error)
	Edit(*OrganisationEditRequest) (interface{}, error)
	Remove(*OrganisationRemoveRequest) (interface{}, error)
	Read(*OrganisationReadRequest) (interface{}, error)
	Archive(*OrganisationArchiveRequest) (interface{}, error)

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
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
