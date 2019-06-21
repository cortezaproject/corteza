package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `trigger.go`, `trigger.util.go` or `trigger_test.go` to
	implement your API calls, helper functions and tests. The file `trigger.go`
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

	"time"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Trigger list request parameters
type TriggerList struct {
	ModuleID    uint64 `json:",string"`
	Query       string
	Page        uint
	PerPage     uint
	NamespaceID uint64 `json:",string"`
}

func NewTriggerList() *TriggerList {
	return &TriggerList{}
}

func (r TriggerList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["moduleID"] = r.ModuleID
	out["query"] = r.Query
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["namespaceID"] = r.NamespaceID

	return out
}

func (r *TriggerList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["moduleID"]; ok {
		r.ModuleID = parseUInt64(val)
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
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewTriggerList()

// Trigger create request parameters
type TriggerCreate struct {
	ModuleID    uint64 `json:",string"`
	Name        string
	Actions     []string
	Enabled     bool
	Source      string
	UpdatedAt   *time.Time
	NamespaceID uint64 `json:",string"`
}

func NewTriggerCreate() *TriggerCreate {
	return &TriggerCreate{}
}

func (r TriggerCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["moduleID"] = r.ModuleID
	out["name"] = r.Name
	out["actions"] = r.Actions
	out["enabled"] = r.Enabled
	out["source"] = r.Source
	out["updatedAt"] = r.UpdatedAt
	out["namespaceID"] = r.NamespaceID

	return out
}

func (r *TriggerCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["moduleID"]; ok {
		r.ModuleID = parseUInt64(val)
	}
	if val, ok := post["name"]; ok {
		r.Name = val
	}

	if val, ok := req.Form["actions"]; ok {
		r.Actions = parseStrings(val)
	}

	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	if val, ok := post["source"]; ok {
		r.Source = val
	}
	if val, ok := post["updatedAt"]; ok {

		if r.UpdatedAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewTriggerCreate()

// Trigger read request parameters
type TriggerRead struct {
	TriggerID   uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewTriggerRead() *TriggerRead {
	return &TriggerRead{}
}

func (r TriggerRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["triggerID"] = r.TriggerID
	out["namespaceID"] = r.NamespaceID

	return out
}

func (r *TriggerRead) Fill(req *http.Request) (err error) {
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

	r.TriggerID = parseUInt64(chi.URLParam(req, "triggerID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewTriggerRead()

// Trigger update request parameters
type TriggerUpdate struct {
	TriggerID   uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	Name        string
	Actions     []string
	Enabled     bool
	Source      string
}

func NewTriggerUpdate() *TriggerUpdate {
	return &TriggerUpdate{}
}

func (r TriggerUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["triggerID"] = r.TriggerID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["name"] = r.Name
	out["actions"] = r.Actions
	out["enabled"] = r.Enabled
	out["source"] = r.Source

	return out
}

func (r *TriggerUpdate) Fill(req *http.Request) (err error) {
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

	r.TriggerID = parseUInt64(chi.URLParam(req, "triggerID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["moduleID"]; ok {
		r.ModuleID = parseUInt64(val)
	}
	if val, ok := post["name"]; ok {
		r.Name = val
	}

	if val, ok := req.Form["actions"]; ok {
		r.Actions = parseStrings(val)
	}

	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	if val, ok := post["source"]; ok {
		r.Source = val
	}

	return err
}

var _ RequestFiller = NewTriggerUpdate()

// Trigger delete request parameters
type TriggerDelete struct {
	TriggerID   uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewTriggerDelete() *TriggerDelete {
	return &TriggerDelete{}
}

func (r TriggerDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["triggerID"] = r.TriggerID
	out["namespaceID"] = r.NamespaceID

	return out
}

func (r *TriggerDelete) Fill(req *http.Request) (err error) {
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

	r.TriggerID = parseUInt64(chi.URLParam(req, "triggerID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewTriggerDelete()
