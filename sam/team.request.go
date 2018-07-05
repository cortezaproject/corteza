package sam

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

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

	t.id = parseUInt64(post["id"])

	t.name = post["name"]
	return nil
}

var _ RequestFiller = teamEditRequest{}.new()

// Team remove request parameters
type teamRemoveRequest struct {
	id uint64
}

func (teamRemoveRequest) new() *teamRemoveRequest {
	return &teamRemoveRequest{}
}

func (t *teamRemoveRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = teamRemoveRequest{}.new()

// Team read request parameters
type teamReadRequest struct {
	id uint64
}

func (teamReadRequest) new() *teamReadRequest {
	return &teamReadRequest{}
}

func (t *teamReadRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = teamReadRequest{}.new()

// Team search request parameters
type teamSearchRequest struct {
	query string
}

func (teamSearchRequest) new() *teamSearchRequest {
	return &teamSearchRequest{}
}

func (t *teamSearchRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = teamSearchRequest{}.new()

// Team archive request parameters
type teamArchiveRequest struct {
	id uint64
}

func (teamArchiveRequest) new() *teamArchiveRequest {
	return &teamArchiveRequest{}
}

func (t *teamArchiveRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(post["id"])
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

	t.id = parseUInt64(post["id"])

	t.organisation_id = parseUInt64(post["organisation_id"])
	return nil
}

var _ RequestFiller = teamMoveRequest{}.new()

// Team merge request parameters
type teamMergeRequest struct {
	destination uint64
	source      uint64
}

func (teamMergeRequest) new() *teamMergeRequest {
	return &teamMergeRequest{}
}

func (t *teamMergeRequest) Fill(r *http.Request) error {
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

	t.destination = parseUInt64(post["destination"])

	t.source = parseUInt64(post["source"])
	return nil
}

var _ RequestFiller = teamMergeRequest{}.new()
