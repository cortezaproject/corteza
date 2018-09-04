package request

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
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}

// Module list request parameters
type ModuleList struct {
	Query string
}

func NewModuleList() *ModuleList {
	return &ModuleList{}
}

func (m *ModuleList) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	if val, ok := get["query"]; ok {

		m.Query = val
	}

	return err
}

var _ RequestFiller = NewModuleList()

// Module create request parameters
type ModuleCreate struct {
	Name   string
	Fields types.JSONText
}

func NewModuleCreate() *ModuleCreate {
	return &ModuleCreate{}
}

func (m *ModuleCreate) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	if val, ok := post["name"]; ok {

		m.Name = val
	}
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONText(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleCreate()

// Module read request parameters
type ModuleRead struct {
	ID uint64
}

func NewModuleRead() *ModuleRead {
	return &ModuleRead{}
}

func (m *ModuleRead) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.ID = parseUInt64(chi.URLParam(r, "id"))

	return err
}

var _ RequestFiller = NewModuleRead()

// Module edit request parameters
type ModuleEdit struct {
	ID     uint64
	Name   string
	Fields types.JSONText
}

func NewModuleEdit() *ModuleEdit {
	return &ModuleEdit{}
}

func (m *ModuleEdit) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.ID = parseUInt64(chi.URLParam(r, "id"))
	if val, ok := post["name"]; ok {

		m.Name = val
	}
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONText(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleEdit()

// Module delete request parameters
type ModuleDelete struct {
	ID uint64
}

func NewModuleDelete() *ModuleDelete {
	return &ModuleDelete{}
}

func (m *ModuleDelete) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.ID = parseUInt64(chi.URLParam(r, "id"))

	return err
}

var _ RequestFiller = NewModuleDelete()

// Module content/list request parameters
type ModuleContentList struct {
	Module uint64
}

func NewModuleContentList() *ModuleContentList {
	return &ModuleContentList{}
}

func (m *ModuleContentList) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.Module = parseUInt64(chi.URLParam(r, "module"))

	return err
}

var _ RequestFiller = NewModuleContentList()

// Module content/create request parameters
type ModuleContentCreate struct {
	Module uint64
	Fields types.JSONText
}

func NewModuleContentCreate() *ModuleContentCreate {
	return &ModuleContentCreate{}
}

func (m *ModuleContentCreate) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.Module = parseUInt64(chi.URLParam(r, "module"))
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONText(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleContentCreate()

// Module content/read request parameters
type ModuleContentRead struct {
	Module uint64
	ID     uint64
}

func NewModuleContentRead() *ModuleContentRead {
	return &ModuleContentRead{}
}

func (m *ModuleContentRead) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.Module = parseUInt64(chi.URLParam(r, "module"))
	m.ID = parseUInt64(chi.URLParam(r, "id"))

	return err
}

var _ RequestFiller = NewModuleContentRead()

// Module content/edit request parameters
type ModuleContentEdit struct {
	Module uint64
	ID     uint64
	Fields types.JSONText
}

func NewModuleContentEdit() *ModuleContentEdit {
	return &ModuleContentEdit{}
}

func (m *ModuleContentEdit) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.Module = parseUInt64(chi.URLParam(r, "module"))
	m.ID = parseUInt64(chi.URLParam(r, "id"))
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONText(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleContentEdit()

// Module content/delete request parameters
type ModuleContentDelete struct {
	Module uint64
	ID     uint64
}

func NewModuleContentDelete() *ModuleContentDelete {
	return &ModuleContentDelete{}
}

func (m *ModuleContentDelete) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}
	}

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

	m.Module = parseUInt64(chi.URLParam(r, "module"))
	m.ID = parseUInt64(chi.URLParam(r, "id"))

	return err
}

var _ RequestFiller = NewModuleContentDelete()
