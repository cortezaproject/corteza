package request

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
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Organisation list request parameters
type OrganisationList struct {
	Query string
}

func NewOrganisationList() *OrganisationList {
	return &OrganisationList{}
}

func (o *OrganisationList) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(o)

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
		o.Query = val
	}

	return nil
}

var _ RequestFiller = NewOrganisationList()

// Organisation create request parameters
type OrganisationCreate struct {
	Name string
}

func NewOrganisationCreate() *OrganisationCreate {
	return &OrganisationCreate{}
}

func (o *OrganisationCreate) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(o)

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
		o.Name = val
	}

	return nil
}

var _ RequestFiller = NewOrganisationCreate()

// Organisation edit request parameters
type OrganisationEdit struct {
	ID   uint64
	Name string
}

func NewOrganisationEdit() *OrganisationEdit {
	return &OrganisationEdit{}
}

func (o *OrganisationEdit) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(o)

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
	if val, ok := post["name"]; ok {
		o.Name = val
	}

	return nil
}

var _ RequestFiller = NewOrganisationEdit()

// Organisation remove request parameters
type OrganisationRemove struct {
	ID uint64
}

func NewOrganisationRemove() *OrganisationRemove {
	return &OrganisationRemove{}
}

func (o *OrganisationRemove) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(o)

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

var _ RequestFiller = NewOrganisationRemove()

// Organisation read request parameters
type OrganisationRead struct {
	ID uint64
}

func NewOrganisationRead() *OrganisationRead {
	return &OrganisationRead{}
}

func (o *OrganisationRead) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(o)

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

	if val, ok := get["id"]; ok {
		o.ID = parseUInt64(val)
	}

	return nil
}

var _ RequestFiller = NewOrganisationRead()

// Organisation archive request parameters
type OrganisationArchive struct {
	ID uint64
}

func NewOrganisationArchive() *OrganisationArchive {
	return &OrganisationArchive{}
}

func (o *OrganisationArchive) Fill(r *http.Request) error {
	json.NewDecoder(r.Body).Decode(o)

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

var _ RequestFiller = NewOrganisationArchive()
