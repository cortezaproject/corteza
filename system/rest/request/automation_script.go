package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `automation_script.go`, `automation_script.util.go` or `automation_script_test.go` to
	implement your API calls, helper functions and tests. The file `automation_script.go`
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

	"github.com/cortezaproject/corteza-server/pkg/automation"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// AutomationScript list request parameters
type AutomationScriptList struct {
	Query      string
	Resource   string
	IncDeleted bool
	Page       uint
	PerPage    uint
}

func NewAutomationScriptList() *AutomationScriptList {
	return &AutomationScriptList{}
}

func (r AutomationScriptList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["resource"] = r.Resource
	out["incDeleted"] = r.IncDeleted
	out["page"] = r.Page
	out["perPage"] = r.PerPage

	return out
}

func (r *AutomationScriptList) Fill(req *http.Request) (err error) {
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
	if val, ok := get["resource"]; ok {
		r.Resource = val
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

	return err
}

var _ RequestFiller = NewAutomationScriptList()

// AutomationScript create request parameters
type AutomationScriptCreate struct {
	Name      string
	SourceRef string
	Source    string
	RunAs     uint64 `json:",string"`
	Timeout   uint
	Critical  bool
	Async     bool
	Enabled   bool
	Triggers  automation.TriggerSet
}

func NewAutomationScriptCreate() *AutomationScriptCreate {
	return &AutomationScriptCreate{}
}

func (r AutomationScriptCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["sourceRef"] = r.SourceRef
	out["source"] = r.Source
	out["runAs"] = r.RunAs
	out["timeout"] = r.Timeout
	out["critical"] = r.Critical
	out["async"] = r.Async
	out["enabled"] = r.Enabled
	out["triggers"] = r.Triggers

	return out
}

func (r *AutomationScriptCreate) Fill(req *http.Request) (err error) {
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
	if val, ok := post["sourceRef"]; ok {
		r.SourceRef = val
	}
	if val, ok := post["source"]; ok {
		r.Source = val
	}
	if val, ok := post["runAs"]; ok {
		r.RunAs = parseUInt64(val)
	}
	if val, ok := post["timeout"]; ok {
		r.Timeout = parseUint(val)
	}
	if val, ok := post["critical"]; ok {
		r.Critical = parseBool(val)
	}
	if val, ok := post["async"]; ok {
		r.Async = parseBool(val)
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}

	return err
}

var _ RequestFiller = NewAutomationScriptCreate()

// AutomationScript read request parameters
type AutomationScriptRead struct {
	ScriptID uint64 `json:",string"`
}

func NewAutomationScriptRead() *AutomationScriptRead {
	return &AutomationScriptRead{}
}

func (r AutomationScriptRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["scriptID"] = r.ScriptID

	return out
}

func (r *AutomationScriptRead) Fill(req *http.Request) (err error) {
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

	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))

	return err
}

var _ RequestFiller = NewAutomationScriptRead()

// AutomationScript update request parameters
type AutomationScriptUpdate struct {
	ScriptID  uint64 `json:",string"`
	Name      string
	SourceRef string
	Source    string
	RunAs     uint64 `json:",string"`
	Timeout   uint
	Critical  bool
	Async     bool
	Enabled   bool
	Triggers  automation.TriggerSet
}

func NewAutomationScriptUpdate() *AutomationScriptUpdate {
	return &AutomationScriptUpdate{}
}

func (r AutomationScriptUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["scriptID"] = r.ScriptID
	out["name"] = r.Name
	out["sourceRef"] = r.SourceRef
	out["source"] = r.Source
	out["runAs"] = r.RunAs
	out["timeout"] = r.Timeout
	out["critical"] = r.Critical
	out["async"] = r.Async
	out["enabled"] = r.Enabled
	out["triggers"] = r.Triggers

	return out
}

func (r *AutomationScriptUpdate) Fill(req *http.Request) (err error) {
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

	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["sourceRef"]; ok {
		r.SourceRef = val
	}
	if val, ok := post["source"]; ok {
		r.Source = val
	}
	if val, ok := post["runAs"]; ok {
		r.RunAs = parseUInt64(val)
	}
	if val, ok := post["timeout"]; ok {
		r.Timeout = parseUint(val)
	}
	if val, ok := post["critical"]; ok {
		r.Critical = parseBool(val)
	}
	if val, ok := post["async"]; ok {
		r.Async = parseBool(val)
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}

	return err
}

var _ RequestFiller = NewAutomationScriptUpdate()

// AutomationScript delete request parameters
type AutomationScriptDelete struct {
	ScriptID uint64 `json:",string"`
}

func NewAutomationScriptDelete() *AutomationScriptDelete {
	return &AutomationScriptDelete{}
}

func (r AutomationScriptDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["scriptID"] = r.ScriptID

	return out
}

func (r *AutomationScriptDelete) Fill(req *http.Request) (err error) {
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

	r.ScriptID = parseUInt64(chi.URLParam(req, "scriptID"))

	return err
}

var _ RequestFiller = NewAutomationScriptDelete()

// AutomationScript test request parameters
type AutomationScriptTest struct {
	Source  string
	Payload json.RawMessage
}

func NewAutomationScriptTest() *AutomationScriptTest {
	return &AutomationScriptTest{}
}

func (r AutomationScriptTest) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["source"] = r.Source
	out["payload"] = r.Payload

	return out
}

func (r *AutomationScriptTest) Fill(req *http.Request) (err error) {
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

	if val, ok := post["source"]; ok {
		r.Source = val
	}
	if val, ok := post["payload"]; ok {
		r.Payload = json.RawMessage(val)
	}

	return err
}

var _ RequestFiller = NewAutomationScriptTest()
