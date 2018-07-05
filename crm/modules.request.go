package crm

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Modules list request parameters
type modulesListRequest struct {
	id string
}

func (modulesListRequest) new() *modulesListRequest {
	return &modulesListRequest{}
}

func (m *modulesListRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = modulesListRequest{}.new()

// Modules edit request parameters
type modulesEditRequest struct {
	id   uint64
	name string
}

func (modulesEditRequest) new() *modulesEditRequest {
	return &modulesEditRequest{}
}

func (m *modulesEditRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = modulesEditRequest{}.new()

// Modules content/list request parameters
type modulesContentListRequest struct {
}

func (modulesContentListRequest) new() *modulesContentListRequest {
	return &modulesContentListRequest{}
}

func (m *modulesContentListRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = modulesContentListRequest{}.new()

// Modules content/edit request parameters
type modulesContentEditRequest struct {
	id      uint64
	payload string
}

func (modulesContentEditRequest) new() *modulesContentEditRequest {
	return &modulesContentEditRequest{}
}

func (m *modulesContentEditRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = modulesContentEditRequest{}.new()

// Modules content/delete request parameters
type modulesContentDeleteRequest struct {
	id uint64
}

func (modulesContentDeleteRequest) new() *modulesContentDeleteRequest {
	return &modulesContentDeleteRequest{}
}

func (m *modulesContentDeleteRequest) Fill(r *http.Request) error {
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

var _ RequestFiller = modulesContentDeleteRequest{}.new()
