package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `automation.go`, `automation.util.go` or `automation_test.go` to
	implement your API calls, helper functions and tests. The file `automation.go`
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

// Automation list request parameters
type AutomationList struct {
	ResourceTypePrefixes []string
	ResourceTypes        []string
	EventTypes           []string
	ExcludeInvalid       bool
	ExcludeClientScripts bool
	ExcludeServerScripts bool
}

func NewAutomationList() *AutomationList {
	return &AutomationList{}
}

func (r AutomationList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resourceTypePrefixes"] = r.ResourceTypePrefixes
	out["resourceTypes"] = r.ResourceTypes
	out["eventTypes"] = r.EventTypes
	out["excludeInvalid"] = r.ExcludeInvalid
	out["excludeClientScripts"] = r.ExcludeClientScripts
	out["excludeServerScripts"] = r.ExcludeServerScripts

	return out
}

func (r *AutomationList) Fill(req *http.Request) (err error) {
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

	if val, ok := urlQuery["resourceTypePrefixes[]"]; ok {
		r.ResourceTypePrefixes = parseStrings(val)
	} else if val, ok = urlQuery["resourceTypePrefixes"]; ok {
		r.ResourceTypePrefixes = parseStrings(val)
	}

	if val, ok := urlQuery["resourceTypes[]"]; ok {
		r.ResourceTypes = parseStrings(val)
	} else if val, ok = urlQuery["resourceTypes"]; ok {
		r.ResourceTypes = parseStrings(val)
	}

	if val, ok := urlQuery["eventTypes[]"]; ok {
		r.EventTypes = parseStrings(val)
	} else if val, ok = urlQuery["eventTypes"]; ok {
		r.EventTypes = parseStrings(val)
	}

	if val, ok := get["excludeInvalid"]; ok {
		r.ExcludeInvalid = parseBool(val)
	}
	if val, ok := get["excludeClientScripts"]; ok {
		r.ExcludeClientScripts = parseBool(val)
	}
	if val, ok := get["excludeServerScripts"]; ok {
		r.ExcludeServerScripts = parseBool(val)
	}

	return err
}

var _ RequestFiller = NewAutomationList()

// Automation bundle request parameters
type AutomationBundle struct {
	Bundle string
	Type   string
	Ext    string
}

func NewAutomationBundle() *AutomationBundle {
	return &AutomationBundle{}
}

func (r AutomationBundle) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["bundle"] = r.Bundle
	out["type"] = r.Type
	out["ext"] = r.Ext

	return out
}

func (r *AutomationBundle) Fill(req *http.Request) (err error) {
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

	r.Bundle = chi.URLParam(req, "bundle")
	r.Type = chi.URLParam(req, "type")
	r.Ext = chi.URLParam(req, "ext")

	return err
}

var _ RequestFiller = NewAutomationBundle()

// Automation triggerScript request parameters
type AutomationTriggerScript struct {
	Script string
}

func NewAutomationTriggerScript() *AutomationTriggerScript {
	return &AutomationTriggerScript{}
}

func (r AutomationTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["script"] = r.Script

	return out
}

func (r *AutomationTriggerScript) Fill(req *http.Request) (err error) {
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

	if val, ok := post["script"]; ok {
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewAutomationTriggerScript()
