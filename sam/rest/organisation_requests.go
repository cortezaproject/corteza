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
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Organisation list request parameters
type OrganisationListRequest struct {
	Query string
}

func (OrganisationListRequest) new() *OrganisationListRequest {
	return &OrganisationListRequest{}
}

func (o *OrganisationListRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	o.Query = get["query"]
	return nil
}

var _ RequestFiller = OrganisationListRequest{}.new()

// Organisation create request parameters
type OrganisationCreateRequest struct {
	Name string
}

func (OrganisationCreateRequest) new() *OrganisationCreateRequest {
	return &OrganisationCreateRequest{}
}

func (o *OrganisationCreateRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	o.Name = post["name"]
	return nil
}

var _ RequestFiller = OrganisationCreateRequest{}.new()

// Organisation edit request parameters
type OrganisationEditRequest struct {
	ID   uint64
	Name string
}

func (OrganisationEditRequest) new() *OrganisationEditRequest {
	return &OrganisationEditRequest{}
}

func (o *OrganisationEditRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	o.ID = parseUInt64(chi.URLParam(r, "id"))

	o.Name = post["name"]
	return nil
}

var _ RequestFiller = OrganisationEditRequest{}.new()

// Organisation remove request parameters
type OrganisationRemoveRequest struct {
	ID uint64
}

func (OrganisationRemoveRequest) new() *OrganisationRemoveRequest {
	return &OrganisationRemoveRequest{}
}

func (o *OrganisationRemoveRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	o.ID = parseUInt64(chi.URLParam(r, "id"))
	return nil
}

var _ RequestFiller = OrganisationRemoveRequest{}.new()

// Organisation read request parameters
type OrganisationReadRequest struct {
	ID uint64
}

func (OrganisationReadRequest) new() *OrganisationReadRequest {
	return &OrganisationReadRequest{}
}

func (o *OrganisationReadRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	o.ID = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = OrganisationReadRequest{}.new()

// Organisation archive request parameters
type OrganisationArchiveRequest struct {
	ID uint64
}

func (OrganisationArchiveRequest) new() *OrganisationArchiveRequest {
	return &OrganisationArchiveRequest{}
}

func (o *OrganisationArchiveRequest) Fill(r *http.Request) error {
	r.ParseForm()
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

	o.ID = parseUInt64(chi.URLParam(r, "id"))
	return nil
}

var _ RequestFiller = OrganisationArchiveRequest{}.new()
