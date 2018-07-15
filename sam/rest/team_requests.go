package rest

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
type TeamListRequest struct {
	Query string
}

func (TeamListRequest) new() *TeamListRequest {
	return &TeamListRequest{}
}

func (t *TeamListRequest) Fill(r *http.Request) error {
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

	t.Query = get["query"]
	return nil
}

var _ RequestFiller = TeamListRequest{}.new()

// Team create request parameters
type TeamCreateRequest struct {
	Name    string
	Members []uint64
}

func (TeamCreateRequest) new() *TeamCreateRequest {
	return &TeamCreateRequest{}
}

func (t *TeamCreateRequest) Fill(r *http.Request) error {
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

	t.Name = post["name"]
	return nil
}

var _ RequestFiller = TeamCreateRequest{}.new()

// Team edit request parameters
type TeamEditRequest struct {
	ID      uint64
	Name    string
	Members []uint64
}

func (TeamEditRequest) new() *TeamEditRequest {
	return &TeamEditRequest{}
}

func (t *TeamEditRequest) Fill(r *http.Request) error {
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

	t.ID = parseUInt64(chi.URLParam(r, "id"))

	t.Name = post["name"]
	return nil
}

var _ RequestFiller = TeamEditRequest{}.new()

// Team read request parameters
type TeamReadRequest struct {
	ID uint64
}

func (TeamReadRequest) new() *TeamReadRequest {
	return &TeamReadRequest{}
}

func (t *TeamReadRequest) Fill(r *http.Request) error {
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

	t.ID = parseUInt64(chi.URLParam(r, "id"))
	return nil
}

var _ RequestFiller = TeamReadRequest{}.new()

// Team remove request parameters
type TeamRemoveRequest struct {
	ID uint64
}

func (TeamRemoveRequest) new() *TeamRemoveRequest {
	return &TeamRemoveRequest{}
}

func (t *TeamRemoveRequest) Fill(r *http.Request) error {
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

	t.ID = parseUInt64(chi.URLParam(r, "id"))
	return nil
}

var _ RequestFiller = TeamRemoveRequest{}.new()

// Team archive request parameters
type TeamArchiveRequest struct {
	ID uint64
}

func (TeamArchiveRequest) new() *TeamArchiveRequest {
	return &TeamArchiveRequest{}
}

func (t *TeamArchiveRequest) Fill(r *http.Request) error {
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

	t.ID = parseUInt64(chi.URLParam(r, "id"))
	return nil
}

var _ RequestFiller = TeamArchiveRequest{}.new()

// Team move request parameters
type TeamMoveRequest struct {
	ID              uint64
	Organisation_id uint64
}

func (TeamMoveRequest) new() *TeamMoveRequest {
	return &TeamMoveRequest{}
}

func (t *TeamMoveRequest) Fill(r *http.Request) error {
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

	t.ID = parseUInt64(chi.URLParam(r, "id"))

	t.Organisation_id = parseUInt64(post["organisation_id"])
	return nil
}

var _ RequestFiller = TeamMoveRequest{}.new()

// Team merge request parameters
type TeamMergeRequest struct {
	ID          uint64
	Destination uint64
}

func (TeamMergeRequest) new() *TeamMergeRequest {
	return &TeamMergeRequest{}
}

func (t *TeamMergeRequest) Fill(r *http.Request) error {
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

	t.ID = parseUInt64(chi.URLParam(r, "id"))

	t.Destination = parseUInt64(post["destination"])
	return nil
}

var _ RequestFiller = TeamMergeRequest{}.new()
