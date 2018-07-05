package crm

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Module list request parameters
type moduleListRequest struct {
	id string
}

func (moduleListRequest) new() *moduleListRequest {
	return &moduleListRequest{}
}

func (m *moduleListRequest) Fill(r *http.Request) error {
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

	m.id = get["id"]
	return nil
}

var _ RequestFiller = moduleListRequest{}.new()

// Module edit request parameters
type moduleEditRequest struct {
	id   uint64
	name string
}

func (moduleEditRequest) new() *moduleEditRequest {
	return &moduleEditRequest{}
}

func (m *moduleEditRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(post["id"])

	m.name = post["name"]
	return nil
}

var _ RequestFiller = moduleEditRequest{}.new()

// Module content/list request parameters
type moduleContentListRequest struct {
}

func (moduleContentListRequest) new() *moduleContentListRequest {
	return &moduleContentListRequest{}
}

func (m *moduleContentListRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = moduleContentListRequest{}.new()

// Module content/edit request parameters
type moduleContentEditRequest struct {
	id      uint64
	payload string
}

func (moduleContentEditRequest) new() *moduleContentEditRequest {
	return &moduleContentEditRequest{}
}

func (m *moduleContentEditRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(post["id"])

	m.payload = post["payload"]
	return nil
}

var _ RequestFiller = moduleContentEditRequest{}.new()

// Module content/delete request parameters
type moduleContentDeleteRequest struct {
	id uint64
}

func (moduleContentDeleteRequest) new() *moduleContentDeleteRequest {
	return &moduleContentDeleteRequest{}
}

func (m *moduleContentDeleteRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = moduleContentDeleteRequest{}.new()
