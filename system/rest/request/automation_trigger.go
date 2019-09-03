package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `automation_trigger.go`, `automation_trigger.util.go` or `automation_trigger_test.go` to
	implement your API calls, helper functions and tests. The file `automation_trigger.go`
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

// AutomationTrigger list request parameters
type AutomationTriggerList struct {
	Resource   string
	Event      string
	IncDeleted bool
	Page       uint
	PerPage    uint
	ScriptID   uint64 `json:",string"`
}

func NewAutomationTriggerList() *AutomationTriggerList {
	return &AutomationTriggerList{}
}

func (r AutomationTriggerList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource
	out["event"] = r.Event
	out["incDeleted"] = r.IncDeleted
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["scriptID"] = r.ScriptID

	return out
}

func (r *AutomationTriggerList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["resource"]; ok {
		r.Resource = val
	}
	if val, ok := get["event"]; ok {
		r.Event = val
	}
	if val, ok := get["incDeleted"]; ok {
		r.IncDeleted = parseBool(val)
	}
	if val, ok := get["page"]; ok {
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.PerPage = parseUint(val)
	}
	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))

	return err
}

var _ RequestFiller = NewAutomationTriggerList()

// AutomationTrigger create request parameters
type AutomationTriggerCreate struct {
	Resource  string
	Event     string
	Condition string
	Enabled   bool
	ScriptID  uint64 `json:",string"`
}

func NewAutomationTriggerCreate() *AutomationTriggerCreate {
	return &AutomationTriggerCreate{}
}

func (r AutomationTriggerCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource
	out["event"] = r.Event
	out["condition"] = r.Condition
	out["enabled"] = r.Enabled
	out["scriptID"] = r.ScriptID

	return out
}

func (r *AutomationTriggerCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["resource"]; ok {
		r.Resource = val
	}
	if val, ok := post["event"]; ok {
		r.Event = val
	}
	if val, ok := post["condition"]; ok {
		r.Condition = val
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))

	return err
}

var _ RequestFiller = NewAutomationTriggerCreate()

// AutomationTrigger read request parameters
type AutomationTriggerRead struct {
	TriggerID uint64 `json:",string"`
	ScriptID  uint64 `json:",string"`
}

func NewAutomationTriggerRead() *AutomationTriggerRead {
	return &AutomationTriggerRead{}
}

func (r AutomationTriggerRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["triggerID"] = r.TriggerID
	out["scriptID"] = r.ScriptID

	return out
}

func (r *AutomationTriggerRead) Fill(req *http.Request) (err error) {
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
	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))

	return err
}

var _ RequestFiller = NewAutomationTriggerRead()

// AutomationTrigger update request parameters
type AutomationTriggerUpdate struct {
	TriggerID uint64 `json:",string"`
	ScriptID  uint64 `json:",string"`
	Resource  string
	Event     string
	Condition string
	Enabled   bool
}

func NewAutomationTriggerUpdate() *AutomationTriggerUpdate {
	return &AutomationTriggerUpdate{}
}

func (r AutomationTriggerUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["triggerID"] = r.TriggerID
	out["scriptID"] = r.ScriptID
	out["resource"] = r.Resource
	out["event"] = r.Event
	out["condition"] = r.Condition
	out["enabled"] = r.Enabled

	return out
}

func (r *AutomationTriggerUpdate) Fill(req *http.Request) (err error) {
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
	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))
	if val, ok := post["resource"]; ok {
		r.Resource = val
	}
	if val, ok := post["event"]; ok {
		r.Event = val
	}
	if val, ok := post["condition"]; ok {
		r.Condition = val
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}

	return err
}

var _ RequestFiller = NewAutomationTriggerUpdate()

// AutomationTrigger delete request parameters
type AutomationTriggerDelete struct {
	TriggerID uint64 `json:",string"`
	ScriptID  uint64 `json:",string"`
}

func NewAutomationTriggerDelete() *AutomationTriggerDelete {
	return &AutomationTriggerDelete{}
}

func (r AutomationTriggerDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["triggerID"] = r.TriggerID
	out["scriptID"] = r.ScriptID

	return out
}

func (r *AutomationTriggerDelete) Fill(req *http.Request) (err error) {
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
	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))

	return err
}

var _ RequestFiller = NewAutomationTriggerDelete()
