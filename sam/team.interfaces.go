package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `team.go`, `team.util.go` or `team_test.go` to
	implement your API calls, helper functions and tests. The file `team.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type TeamHandlers struct {
	*Team
}

func (TeamHandlers) new() *TeamHandlers {
	return &TeamHandlers{
		Team{}.New(),
	}
}

// Internal API interface
type TeamAPI interface {
	List(*teamListRequest) (interface{}, error)
	Create(*teamCreateRequest) (interface{}, error)
	Edit(*teamEditRequest) (interface{}, error)
	Read(*teamReadRequest) (interface{}, error)
	Remove(*teamRemoveRequest) (interface{}, error)
	Archive(*teamArchiveRequest) (interface{}, error)
	Move(*teamMoveRequest) (interface{}, error)
	Merge(*teamMergeRequest) (interface{}, error)
}

// HTTP API interface
type TeamHandlersAPI interface {
	List(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Remove(http.ResponseWriter, *http.Request)
	Archive(http.ResponseWriter, *http.Request)
	Move(http.ResponseWriter, *http.Request)
	Merge(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ TeamHandlersAPI = &TeamHandlers{}
var _ TeamAPI = &Team{}
