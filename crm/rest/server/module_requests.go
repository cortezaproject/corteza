package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Module list request parameters
type ModuleListRequest struct {
	ID uint64
}

func (ModuleListRequest) new() *ModuleListRequest {
	return &ModuleListRequest{}
}

func (m *ModuleListRequest) Fill(r *http.Request) error {
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

	m.ID = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = ModuleListRequest{}.new()

// Module edit request parameters
type ModuleEditRequest struct {
	ID   uint64
	Name string
}

func (ModuleEditRequest) new() *ModuleEditRequest {
	return &ModuleEditRequest{}
}

func (m *ModuleEditRequest) Fill(r *http.Request) error {
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

	m.ID = parseUInt64(post["id"])

	m.Name = post["name"]
	return nil
}

var _ RequestFiller = ModuleEditRequest{}.new()

// Module content/list request parameters
type ModuleContentListRequest struct {
	ID uint64
}

func (ModuleContentListRequest) new() *ModuleContentListRequest {
	return &ModuleContentListRequest{}
}

func (m *ModuleContentListRequest) Fill(r *http.Request) error {
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

	m.ID = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = ModuleContentListRequest{}.new()

// Module content/edit request parameters
type ModuleContentEditRequest struct {
	ID      uint64
	Payload string
}

func (ModuleContentEditRequest) new() *ModuleContentEditRequest {
	return &ModuleContentEditRequest{}
}

func (m *ModuleContentEditRequest) Fill(r *http.Request) error {
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

	m.ID = parseUInt64(post["id"])

	m.Payload = post["payload"]
	return nil
}

var _ RequestFiller = ModuleContentEditRequest{}.new()

// Module content/delete request parameters
type ModuleContentDeleteRequest struct {
	ID uint64
}

func (ModuleContentDeleteRequest) new() *ModuleContentDeleteRequest {
	return &ModuleContentDeleteRequest{}
}

func (m *ModuleContentDeleteRequest) Fill(r *http.Request) error {
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

	m.ID = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = ModuleContentDeleteRequest{}.new()
