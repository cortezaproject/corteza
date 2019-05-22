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
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Organisation list request parameters
type OrganisationList struct {
	Query string
}

func NewOrganisationList() *OrganisationList {
	return &OrganisationList{}
}

func (r OrganisationList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query

	return out
}

func (r *OrganisationList) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["query"]; ok {
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewOrganisationList()

// Organisation create request parameters
type OrganisationCreate struct {
	Name string
}

func NewOrganisationCreate() *OrganisationCreate {
	return &OrganisationCreate{}
}

func (r OrganisationCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name

	return out
}

func (r *OrganisationCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := post["name"]; ok {
		r.Name = val
	}

	return err
}

var _ RequestFiller = NewOrganisationCreate()

// Organisation update request parameters
type OrganisationUpdate struct {
	ID   uint64 `json:",string"`
	Name string
}

func NewOrganisationUpdate() *OrganisationUpdate {
	return &OrganisationUpdate{}
}

func (r OrganisationUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID
	out["name"] = r.Name

	return out
}

func (r *OrganisationUpdate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ID = parseUInt64(chi.URLParam(req, "id"))
	if val, ok := post["name"]; ok {
		r.Name = val
	}

	return err
}

var _ RequestFiller = NewOrganisationUpdate()

// Organisation delete request parameters
type OrganisationDelete struct {
	ID uint64 `json:",string"`
}

func NewOrganisationDelete() *OrganisationDelete {
	return &OrganisationDelete{}
}

func (r OrganisationDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID

	return out
}

func (r *OrganisationDelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ID = parseUInt64(chi.URLParam(req, "id"))

	return err
}

var _ RequestFiller = NewOrganisationDelete()

// Organisation read request parameters
type OrganisationRead struct {
	ID uint64 `json:",string"`
}

func NewOrganisationRead() *OrganisationRead {
	return &OrganisationRead{}
}

func (r OrganisationRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID

	return out
}

func (r *OrganisationRead) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["id"]; ok {
		r.ID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewOrganisationRead()

// Organisation archive request parameters
type OrganisationArchive struct {
	ID uint64 `json:",string"`
}

func NewOrganisationArchive() *OrganisationArchive {
	return &OrganisationArchive{}
}

func (r OrganisationArchive) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID

	return out
}

func (r *OrganisationArchive) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ID = parseUInt64(chi.URLParam(req, "id"))

	return err
}

var _ RequestFiller = NewOrganisationArchive()
