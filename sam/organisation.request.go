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
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Organisation list request parameters
type organisationListRequest struct {
	query string
}

func (organisationListRequest) new() *organisationListRequest {
	return &organisationListRequest{}
}

func (o *organisationListRequest) Fill(r *http.Request) error {
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

	o.query = get["query"]
	return nil
}

var _ RequestFiller = organisationListRequest{}.new()

// Organisation create request parameters
type organisationCreateRequest struct {
	name string
}

func (organisationCreateRequest) new() *organisationCreateRequest {
	return &organisationCreateRequest{}
}

func (o *organisationCreateRequest) Fill(r *http.Request) error {
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

	o.name = post["name"]
	return nil
}

var _ RequestFiller = organisationCreateRequest{}.new()

// Organisation edit request parameters
type organisationEditRequest struct {
	id   uint64
	name string
}

func (organisationEditRequest) new() *organisationEditRequest {
	return &organisationEditRequest{}
}

func (o *organisationEditRequest) Fill(r *http.Request) error {
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

	o.id = chi.URLParam(r, "id")

	o.name = post["name"]
	return nil
}

var _ RequestFiller = organisationEditRequest{}.new()

// Organisation remove request parameters
type organisationRemoveRequest struct {
	id uint64
}

func (organisationRemoveRequest) new() *organisationRemoveRequest {
	return &organisationRemoveRequest{}
}

func (o *organisationRemoveRequest) Fill(r *http.Request) error {
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

	o.id = chi.URLParam(r, "id")
	return nil
}

var _ RequestFiller = organisationRemoveRequest{}.new()

// Organisation read request parameters
type organisationReadRequest struct {
	id uint64
}

func (organisationReadRequest) new() *organisationReadRequest {
	return &organisationReadRequest{}
}

func (o *organisationReadRequest) Fill(r *http.Request) error {
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

	o.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = organisationReadRequest{}.new()

// Organisation archive request parameters
type organisationArchiveRequest struct {
	id uint64
}

func (organisationArchiveRequest) new() *organisationArchiveRequest {
	return &organisationArchiveRequest{}
}

func (o *organisationArchiveRequest) Fill(r *http.Request) error {
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

	o.id = chi.URLParam(r, "id")
	return nil
}

var _ RequestFiller = organisationArchiveRequest{}.new()
