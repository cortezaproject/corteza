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
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Team list request parameters
type teamListRequest struct {
	query string
}

func (teamListRequest) new() *teamListRequest {
	return &teamListRequest{}
}

func (t *teamListRequest) Fill(r *http.Request) error {
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

	t.query = get["query"]
	return nil
}

var _ RequestFiller = teamListRequest{}.new()

// Team create request parameters
type teamCreateRequest struct {
	name    string
	members []uint64
}

func (teamCreateRequest) new() *teamCreateRequest {
	return &teamCreateRequest{}
}

func (t *teamCreateRequest) Fill(r *http.Request) error {
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

	t.name = post["name"]
	return nil
}

var _ RequestFiller = teamCreateRequest{}.new()

// Team edit request parameters
type teamEditRequest struct {
	id      uint64
	name    string
	members []uint64
}

func (teamEditRequest) new() *teamEditRequest {
	return &teamEditRequest{}
}

func (t *teamEditRequest) Fill(r *http.Request) error {
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

	t.id = chi.URLParam(r, "id")

	t.name = post["name"]
	return nil
}

var _ RequestFiller = teamEditRequest{}.new()

// Team read request parameters
type teamReadRequest struct {
	id uint64
}

func (teamReadRequest) new() *teamReadRequest {
	return &teamReadRequest{}
}

func (t *teamReadRequest) Fill(r *http.Request) error {
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

	t.id = chi.URLParam(r, "id")
	return nil
}

var _ RequestFiller = teamReadRequest{}.new()

// Team remove request parameters
type teamRemoveRequest struct {
	id uint64
}

func (teamRemoveRequest) new() *teamRemoveRequest {
	return &teamRemoveRequest{}
}

func (t *teamRemoveRequest) Fill(r *http.Request) error {
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

	t.id = chi.URLParam(r, "id")
	return nil
}

var _ RequestFiller = teamRemoveRequest{}.new()

// Team archive request parameters
type teamArchiveRequest struct {
	id uint64
}

func (teamArchiveRequest) new() *teamArchiveRequest {
	return &teamArchiveRequest{}
}

func (t *teamArchiveRequest) Fill(r *http.Request) error {
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

	t.id = chi.URLParam(r, "id")
	return nil
}

var _ RequestFiller = teamArchiveRequest{}.new()

// Team move request parameters
type teamMoveRequest struct {
	id              uint64
	organisation_id uint64
}

func (teamMoveRequest) new() *teamMoveRequest {
	return &teamMoveRequest{}
}

func (t *teamMoveRequest) Fill(r *http.Request) error {
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

	t.id = chi.URLParam(r, "id")

	t.organisation_id = parseUInt64(post["organisation_id"])
	return nil
}

var _ RequestFiller = teamMoveRequest{}.new()

// Team merge request parameters
type teamMergeRequest struct {
	id          uint64
	destination uint64
}

func (teamMergeRequest) new() *teamMergeRequest {
	return &teamMergeRequest{}
}

func (t *teamMergeRequest) Fill(r *http.Request) error {
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

	t.id = chi.URLParam(r, "id")

	t.destination = parseUInt64(post["destination"])
	return nil
}

var _ RequestFiller = teamMergeRequest{}.new()
