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

// AutomationList request parameters
type AutomationList struct {
	hasResourceTypePrefixes bool
	rawResourceTypePrefixes []string
	ResourceTypePrefixes    []string

	hasResourceTypes bool
	rawResourceTypes []string
	ResourceTypes    []string

	hasEventTypes bool
	rawEventTypes []string
	EventTypes    []string

	hasExcludeInvalid bool
	rawExcludeInvalid string
	ExcludeInvalid    bool

	hasExcludeClientScripts bool
	rawExcludeClientScripts string
	ExcludeClientScripts    bool

	hasExcludeServerScripts bool
	rawExcludeServerScripts string
	ExcludeServerScripts    bool
}

// NewAutomationList request
func NewAutomationList() *AutomationList {
	return &AutomationList{}
}

// Auditable returns all auditable/loggable parameters
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

// Fill processes request and fills internal variables
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
		r.hasResourceTypePrefixes = true
		r.rawResourceTypePrefixes = val
		r.ResourceTypePrefixes = parseStrings(val)
	} else if val, ok = urlQuery["resourceTypePrefixes"]; ok {
		r.hasResourceTypePrefixes = true
		r.rawResourceTypePrefixes = val
		r.ResourceTypePrefixes = parseStrings(val)
	}

	if val, ok := urlQuery["resourceTypes[]"]; ok {
		r.hasResourceTypes = true
		r.rawResourceTypes = val
		r.ResourceTypes = parseStrings(val)
	} else if val, ok = urlQuery["resourceTypes"]; ok {
		r.hasResourceTypes = true
		r.rawResourceTypes = val
		r.ResourceTypes = parseStrings(val)
	}

	if val, ok := urlQuery["eventTypes[]"]; ok {
		r.hasEventTypes = true
		r.rawEventTypes = val
		r.EventTypes = parseStrings(val)
	} else if val, ok = urlQuery["eventTypes"]; ok {
		r.hasEventTypes = true
		r.rawEventTypes = val
		r.EventTypes = parseStrings(val)
	}

	if val, ok := get["excludeInvalid"]; ok {
		r.hasExcludeInvalid = true
		r.rawExcludeInvalid = val
		r.ExcludeInvalid = parseBool(val)
	}
	if val, ok := get["excludeClientScripts"]; ok {
		r.hasExcludeClientScripts = true
		r.rawExcludeClientScripts = val
		r.ExcludeClientScripts = parseBool(val)
	}
	if val, ok := get["excludeServerScripts"]; ok {
		r.hasExcludeServerScripts = true
		r.rawExcludeServerScripts = val
		r.ExcludeServerScripts = parseBool(val)
	}

	return err
}

var _ RequestFiller = NewAutomationList()

// AutomationBundle request parameters
type AutomationBundle struct {
	hasBundle bool
	rawBundle string
	Bundle    string

	hasType bool
	rawType string
	Type    string

	hasExt bool
	rawExt string
	Ext    string
}

// NewAutomationBundle request
func NewAutomationBundle() *AutomationBundle {
	return &AutomationBundle{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationBundle) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["bundle"] = r.Bundle
	out["type"] = r.Type
	out["ext"] = r.Ext

	return out
}

// Fill processes request and fills internal variables
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

	r.hasBundle = true
	r.rawBundle = chi.URLParam(req, "bundle")
	r.Bundle = chi.URLParam(req, "bundle")
	r.hasType = true
	r.rawType = chi.URLParam(req, "type")
	r.Type = chi.URLParam(req, "type")
	r.hasExt = true
	r.rawExt = chi.URLParam(req, "ext")
	r.Ext = chi.URLParam(req, "ext")

	return err
}

var _ RequestFiller = NewAutomationBundle()

// AutomationTriggerScript request parameters
type AutomationTriggerScript struct {
	hasScript bool
	rawScript string
	Script    string
}

// NewAutomationTriggerScript request
func NewAutomationTriggerScript() *AutomationTriggerScript {
	return &AutomationTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
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
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewAutomationTriggerScript()

// HasResourceTypePrefixes returns true if resourceTypePrefixes was set
func (r *AutomationList) HasResourceTypePrefixes() bool {
	return r.hasResourceTypePrefixes
}

// RawResourceTypePrefixes returns raw value of resourceTypePrefixes parameter
func (r *AutomationList) RawResourceTypePrefixes() []string {
	return r.rawResourceTypePrefixes
}

// GetResourceTypePrefixes returns casted value of  resourceTypePrefixes parameter
func (r *AutomationList) GetResourceTypePrefixes() []string {
	return r.ResourceTypePrefixes
}

// HasResourceTypes returns true if resourceTypes was set
func (r *AutomationList) HasResourceTypes() bool {
	return r.hasResourceTypes
}

// RawResourceTypes returns raw value of resourceTypes parameter
func (r *AutomationList) RawResourceTypes() []string {
	return r.rawResourceTypes
}

// GetResourceTypes returns casted value of  resourceTypes parameter
func (r *AutomationList) GetResourceTypes() []string {
	return r.ResourceTypes
}

// HasEventTypes returns true if eventTypes was set
func (r *AutomationList) HasEventTypes() bool {
	return r.hasEventTypes
}

// RawEventTypes returns raw value of eventTypes parameter
func (r *AutomationList) RawEventTypes() []string {
	return r.rawEventTypes
}

// GetEventTypes returns casted value of  eventTypes parameter
func (r *AutomationList) GetEventTypes() []string {
	return r.EventTypes
}

// HasExcludeInvalid returns true if excludeInvalid was set
func (r *AutomationList) HasExcludeInvalid() bool {
	return r.hasExcludeInvalid
}

// RawExcludeInvalid returns raw value of excludeInvalid parameter
func (r *AutomationList) RawExcludeInvalid() string {
	return r.rawExcludeInvalid
}

// GetExcludeInvalid returns casted value of  excludeInvalid parameter
func (r *AutomationList) GetExcludeInvalid() bool {
	return r.ExcludeInvalid
}

// HasExcludeClientScripts returns true if excludeClientScripts was set
func (r *AutomationList) HasExcludeClientScripts() bool {
	return r.hasExcludeClientScripts
}

// RawExcludeClientScripts returns raw value of excludeClientScripts parameter
func (r *AutomationList) RawExcludeClientScripts() string {
	return r.rawExcludeClientScripts
}

// GetExcludeClientScripts returns casted value of  excludeClientScripts parameter
func (r *AutomationList) GetExcludeClientScripts() bool {
	return r.ExcludeClientScripts
}

// HasExcludeServerScripts returns true if excludeServerScripts was set
func (r *AutomationList) HasExcludeServerScripts() bool {
	return r.hasExcludeServerScripts
}

// RawExcludeServerScripts returns raw value of excludeServerScripts parameter
func (r *AutomationList) RawExcludeServerScripts() string {
	return r.rawExcludeServerScripts
}

// GetExcludeServerScripts returns casted value of  excludeServerScripts parameter
func (r *AutomationList) GetExcludeServerScripts() bool {
	return r.ExcludeServerScripts
}

// HasBundle returns true if bundle was set
func (r *AutomationBundle) HasBundle() bool {
	return r.hasBundle
}

// RawBundle returns raw value of bundle parameter
func (r *AutomationBundle) RawBundle() string {
	return r.rawBundle
}

// GetBundle returns casted value of  bundle parameter
func (r *AutomationBundle) GetBundle() string {
	return r.Bundle
}

// HasType returns true if type was set
func (r *AutomationBundle) HasType() bool {
	return r.hasType
}

// RawType returns raw value of type parameter
func (r *AutomationBundle) RawType() string {
	return r.rawType
}

// GetType returns casted value of  type parameter
func (r *AutomationBundle) GetType() string {
	return r.Type
}

// HasExt returns true if ext was set
func (r *AutomationBundle) HasExt() bool {
	return r.hasExt
}

// RawExt returns raw value of ext parameter
func (r *AutomationBundle) RawExt() string {
	return r.rawExt
}

// GetExt returns casted value of  ext parameter
func (r *AutomationBundle) GetExt() string {
	return r.Ext
}

// HasScript returns true if script was set
func (r *AutomationTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *AutomationTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *AutomationTriggerScript) GetScript() string {
	return r.Script
}
