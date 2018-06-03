package sam

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type UserHandlers struct {
	*User
}

func (UserHandlers) new() *UserHandlers {
	return &UserHandlers{
		User{}.new(),
	}
}

// Internal API interface
type UserAPI interface {
	Login(*userLoginRequest) (interface{}, error)
	Search(*userSearchRequest) (interface{}, error)
}

// HTTP API interface
type UserHandlersAPI interface {
	Login(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ UserHandlersAPI = &UserHandlers{}
var _ UserAPI = &User{}
