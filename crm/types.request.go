package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `types.go`, `types.util.go` or `types_test.go` to
	implement your API calls, helper functions and tests. The file `types.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

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

	t.id = chi.URLParam(r, "id")
	return nil
}

var _ RequestFiller = typesTypeRequest{}.new()
