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
	ResourceTypes        []string
	EventTypes           []string
	ExcludeClientScripts bool
	ExcludeServerScripts bool
}

func NewAutomationList() *AutomationList {
	return &AutomationList{}
}

func (r AutomationList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resourceTypes"] = r.ResourceTypes
	out["eventTypes"] = r.EventTypes
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

	if val, ok := get["excludeClientScripts"]; ok {
		r.ExcludeClientScripts = parseBool(val)
	}
	if val, ok := get["excludeServerScripts"]; ok {
		r.ExcludeServerScripts = parseBool(val)
	}

	return err
}

var _ RequestFiller = NewAutomationList()

// Automation trigger request parameters
type AutomationTrigger struct {
	Script string
}

func NewAutomationTrigger() *AutomationTrigger {
	return &AutomationTrigger{}
}

func (r AutomationTrigger) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["script"] = r.Script

	return out
}

func (r *AutomationTrigger) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewAutomationTrigger()
