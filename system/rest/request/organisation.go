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

// OrganisationList request parameters
type OrganisationList struct {
	hasQuery bool
	rawQuery string
	Query    string
}

// NewOrganisationList request
func NewOrganisationList() *OrganisationList {
	return &OrganisationList{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query

	return out
}

// Fill processes request and fills internal variables
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
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewOrganisationList()

// OrganisationCreate request parameters
type OrganisationCreate struct {
	hasName bool
	rawName string
	Name    string
}

// NewOrganisationCreate request
func NewOrganisationCreate() *OrganisationCreate {
	return &OrganisationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name

	return out
}

// Fill processes request and fills internal variables
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
		r.hasName = true
		r.rawName = val
		r.Name = val
	}

	return err
}

var _ RequestFiller = NewOrganisationCreate()

// OrganisationUpdate request parameters
type OrganisationUpdate struct {
	hasID bool
	rawID string
	ID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string
}

// NewOrganisationUpdate request
func NewOrganisationUpdate() *OrganisationUpdate {
	return &OrganisationUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID
	out["name"] = r.Name

	return out
}

// Fill processes request and fills internal variables
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

	r.hasID = true
	r.rawID = chi.URLParam(req, "id")
	r.ID = parseUInt64(chi.URLParam(req, "id"))
	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}

	return err
}

var _ RequestFiller = NewOrganisationUpdate()

// OrganisationDelete request parameters
type OrganisationDelete struct {
	hasID bool
	rawID string
	ID    uint64 `json:",string"`
}

// NewOrganisationDelete request
func NewOrganisationDelete() *OrganisationDelete {
	return &OrganisationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasID = true
	r.rawID = chi.URLParam(req, "id")
	r.ID = parseUInt64(chi.URLParam(req, "id"))

	return err
}

var _ RequestFiller = NewOrganisationDelete()

// OrganisationRead request parameters
type OrganisationRead struct {
	hasID bool
	rawID string
	ID    uint64 `json:",string"`
}

// NewOrganisationRead request
func NewOrganisationRead() *OrganisationRead {
	return &OrganisationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasID = true
	r.rawID = chi.URLParam(req, "id")
	r.ID = parseUInt64(chi.URLParam(req, "id"))

	return err
}

var _ RequestFiller = NewOrganisationRead()

// OrganisationArchive request parameters
type OrganisationArchive struct {
	hasID bool
	rawID string
	ID    uint64 `json:",string"`
}

// NewOrganisationArchive request
func NewOrganisationArchive() *OrganisationArchive {
	return &OrganisationArchive{}
}

// Auditable returns all auditable/loggable parameters
func (r OrganisationArchive) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["id"] = r.ID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasID = true
	r.rawID = chi.URLParam(req, "id")
	r.ID = parseUInt64(chi.URLParam(req, "id"))

	return err
}

var _ RequestFiller = NewOrganisationArchive()

// HasQuery returns true if query was set
func (r *OrganisationList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *OrganisationList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *OrganisationList) GetQuery() string {
	return r.Query
}

// HasName returns true if name was set
func (r *OrganisationCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *OrganisationCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *OrganisationCreate) GetName() string {
	return r.Name
}

// HasID returns true if id was set
func (r *OrganisationUpdate) HasID() bool {
	return r.hasID
}

// RawID returns raw value of id parameter
func (r *OrganisationUpdate) RawID() string {
	return r.rawID
}

// GetID returns casted value of  id parameter
func (r *OrganisationUpdate) GetID() uint64 {
	return r.ID
}

// HasName returns true if name was set
func (r *OrganisationUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *OrganisationUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *OrganisationUpdate) GetName() string {
	return r.Name
}

// HasID returns true if id was set
func (r *OrganisationDelete) HasID() bool {
	return r.hasID
}

// RawID returns raw value of id parameter
func (r *OrganisationDelete) RawID() string {
	return r.rawID
}

// GetID returns casted value of  id parameter
func (r *OrganisationDelete) GetID() uint64 {
	return r.ID
}

// HasID returns true if id was set
func (r *OrganisationRead) HasID() bool {
	return r.hasID
}

// RawID returns raw value of id parameter
func (r *OrganisationRead) RawID() string {
	return r.rawID
}

// GetID returns casted value of  id parameter
func (r *OrganisationRead) GetID() uint64 {
	return r.ID
}

// HasID returns true if id was set
func (r *OrganisationArchive) HasID() bool {
	return r.hasID
}

// RawID returns raw value of id parameter
func (r *OrganisationArchive) RawID() string {
	return r.rawID
}

// GetID returns casted value of  id parameter
func (r *OrganisationArchive) GetID() uint64 {
	return r.ID
}
