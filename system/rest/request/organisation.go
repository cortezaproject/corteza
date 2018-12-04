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
	"github.com/crusttech/crust/internal/rbac"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}
var _ = multipart.FileHeader{}
var _ = rbac.Operation{}

// Organisation list request parameters
type OrganisationList struct {
	Query string
}

func NewOrganisationList() *OrganisationList {
	return &OrganisationList{}
}

func (o *OrganisationList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(o)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

func (o *OrganisationCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(o)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	return err
}

var _ RequestFiller = NewOrganisationCreate()

// Organisation edit request parameters
type OrganisationEdit struct {
	ID   uint64 `json:",string"`
	Name string
}

func NewOrganisationEdit() *OrganisationEdit {
	return &OrganisationEdit{}
}

func (o *OrganisationEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(o)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	return err
}

var _ RequestFiller = NewOrganisationEdit()

// Organisation remove request parameters
type OrganisationRemove struct {
	ID uint64 `json:",string"`
}

func NewOrganisationRemove() *OrganisationRemove {
	return &OrganisationRemove{}
}

func (o *OrganisationRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(o)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	return err
}

var _ RequestFiller = NewOrganisationRemove()

// Organisation read request parameters
type OrganisationRead struct {
	ID uint64 `json:",string"`
}

func NewOrganisationRead() *OrganisationRead {
	return &OrganisationRead{}
}

func (o *OrganisationRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(o)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

func (o *OrganisationArchive) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(o)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	return err
}

var _ RequestFiller = NewOrganisationArchive()
