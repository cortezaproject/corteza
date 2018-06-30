package crm

import (
	"net/http"
)

// Types list request parameters
type typesListRequest struct {
}

func (typesListRequest) new() *typesListRequest {
	return &typesListRequest{}
}

func (t *typesListRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}
	return nil
}

var _ RequestFiller = typesListRequest{}.new()

// Types type request parameters
type typesTypeRequest struct {
	id string
}

func (typesTypeRequest) new() *typesTypeRequest {
	return &typesTypeRequest{}
}

func (t *typesTypeRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	t.id = path["id"]
	return nil
}

var _ RequestFiller = typesTypeRequest{}.new()
