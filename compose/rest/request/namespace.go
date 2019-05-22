package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `namespace.go`, `namespace.util.go` or `namespace_test.go` to
	implement your API calls, helper functions and tests. The file `namespace.go`
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

	sqlxTypes "github.com/jmoiron/sqlx/types"
	"time"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Namespace list request parameters
type NamespaceList struct {
	Query   string
	Page    uint
	PerPage uint
}

func NewNamespaceList() *NamespaceList {
	return &NamespaceList{}
}

func (r NamespaceList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["page"] = r.Page
	out["perPage"] = r.PerPage

	return out
}

func (r *NamespaceList) Fill(req *http.Request) (err error) {
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
	if val, ok := get["page"]; ok {
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.PerPage = parseUint(val)
	}

	return err
}

var _ RequestFiller = NewNamespaceList()

// Namespace create request parameters
type NamespaceCreate struct {
	Name    string
	Slug    string
	Enabled bool
	Meta    sqlxTypes.JSONText
}

func NewNamespaceCreate() *NamespaceCreate {
	return &NamespaceCreate{}
}

func (r NamespaceCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["slug"] = r.Slug
	out["enabled"] = r.Enabled
	out["meta"] = r.Meta

	return out
}

func (r *NamespaceCreate) Fill(req *http.Request) (err error) {
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
	if val, ok := post["slug"]; ok {
		r.Slug = val
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	if val, ok := post["meta"]; ok {

		if r.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewNamespaceCreate()

// Namespace read request parameters
type NamespaceRead struct {
	NamespaceID uint64 `json:",string"`
}

func NewNamespaceRead() *NamespaceRead {
	return &NamespaceRead{}
}

func (r NamespaceRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID

	return out
}

func (r *NamespaceRead) Fill(req *http.Request) (err error) {
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

	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewNamespaceRead()

// Namespace update request parameters
type NamespaceUpdate struct {
	NamespaceID uint64 `json:",string"`
	Name        string
	Slug        string
	Enabled     bool
	Meta        sqlxTypes.JSONText
	UpdatedAt   *time.Time
}

func NewNamespaceUpdate() *NamespaceUpdate {
	return &NamespaceUpdate{}
}

func (r NamespaceUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID
	out["name"] = r.Name
	out["slug"] = r.Slug
	out["enabled"] = r.Enabled
	out["meta"] = r.Meta
	out["updatedAt"] = r.UpdatedAt

	return out
}

func (r *NamespaceUpdate) Fill(req *http.Request) (err error) {
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

	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["slug"]; ok {
		r.Slug = val
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	if val, ok := post["meta"]; ok {

		if r.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["updatedAt"]; ok {

		if r.UpdatedAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewNamespaceUpdate()

// Namespace delete request parameters
type NamespaceDelete struct {
	NamespaceID uint64 `json:",string"`
}

func NewNamespaceDelete() *NamespaceDelete {
	return &NamespaceDelete{}
}

func (r NamespaceDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID

	return out
}

func (r *NamespaceDelete) Fill(req *http.Request) (err error) {
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

	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewNamespaceDelete()
