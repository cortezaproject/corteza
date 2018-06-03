package sam

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type TeamHandlers struct {
	*Team
}

func (TeamHandlers) new() *TeamHandlers {
	return &TeamHandlers{
		Team{}.new(),
	}
}

// Internal API interface
type TeamAPI interface {
	Edit(*TeamEditRequest) (interface{}, error)
	Remove(*TeamRemoveRequest) (interface{}, error)
	Read(*TeamReadRequest) (interface{}, error)
	Search(*TeamSearchRequest) (interface{}, error)
	Archive(*TeamArchiveRequest) (interface{}, error)
	Move(*TeamMoveRequest) (interface{}, error)
}

// HTTP API interface
type TeamHandlersAPI interface {
	Edit(http.ResponseWriter, *http.Request)
	Remove(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Archive(http.ResponseWriter, *http.Request)
	Move(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ TeamHandlersAPI = &TeamHandlers{}
var _ TeamAPI = &Team{}
