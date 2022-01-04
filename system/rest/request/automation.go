package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	AutomationList struct {
		// ResourceTypePrefixes GET parameter
		//
		// Filter by resource prefix
		ResourceTypePrefixes []string

		// ResourceTypes GET parameter
		//
		// Filter by resource type
		ResourceTypes []string

		// EventTypes GET parameter
		//
		// Filter by event type
		EventTypes []string

		// ExcludeInvalid GET parameter
		//
		// Exclude scripts that cannot be used (errors)
		ExcludeInvalid bool

		// ExcludeClientScripts GET parameter
		//
		// Do not include client scripts
		ExcludeClientScripts bool

		// ExcludeServerScripts GET parameter
		//
		// Do not include server scripts
		ExcludeServerScripts bool
	}

	AutomationBundle struct {
		// Bundle PATH parameter
		//
		// Name of the bundle
		Bundle string

		// Type PATH parameter
		//
		// Bundle type
		Type string

		// Ext PATH parameter
		//
		// Bundle extension
		Ext string
	}

	AutomationTriggerScript struct {
		// Script POST parameter
		//
		// Script to execute
		Script string

		// Args POST parameter
		//
		// Arguments to pass to the script
		Args map[string]interface{}
	}
)

// NewAutomationList request
func NewAutomationList() *AutomationList {
	return &AutomationList{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resourceTypePrefixes": r.ResourceTypePrefixes,
		"resourceTypes":        r.ResourceTypes,
		"eventTypes":           r.EventTypes,
		"excludeInvalid":       r.ExcludeInvalid,
		"excludeClientScripts": r.ExcludeClientScripts,
		"excludeServerScripts": r.ExcludeServerScripts,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) GetResourceTypePrefixes() []string {
	return r.ResourceTypePrefixes
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) GetResourceTypes() []string {
	return r.ResourceTypes
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) GetEventTypes() []string {
	return r.EventTypes
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) GetExcludeInvalid() bool {
	return r.ExcludeInvalid
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) GetExcludeClientScripts() bool {
	return r.ExcludeClientScripts
}

// Auditable returns all auditable/loggable parameters
func (r AutomationList) GetExcludeServerScripts() bool {
	return r.ExcludeServerScripts
}

// Fill processes request and fills internal variables
func (r *AutomationList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["resourceTypePrefixes[]"]; ok {
			r.ResourceTypePrefixes, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["resourceTypePrefixes"]; ok {
			r.ResourceTypePrefixes, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resourceTypes[]"]; ok {
			r.ResourceTypes, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["resourceTypes"]; ok {
			r.ResourceTypes, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["eventTypes[]"]; ok {
			r.EventTypes, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["eventTypes"]; ok {
			r.EventTypes, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["excludeInvalid"]; ok && len(val) > 0 {
			r.ExcludeInvalid, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["excludeClientScripts"]; ok && len(val) > 0 {
			r.ExcludeClientScripts, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["excludeServerScripts"]; ok && len(val) > 0 {
			r.ExcludeServerScripts, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAutomationBundle request
func NewAutomationBundle() *AutomationBundle {
	return &AutomationBundle{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationBundle) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"bundle": r.Bundle,
		"type":   r.Type,
		"ext":    r.Ext,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationBundle) GetBundle() string {
	return r.Bundle
}

// Auditable returns all auditable/loggable parameters
func (r AutomationBundle) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r AutomationBundle) GetExt() string {
	return r.Ext
}

// Fill processes request and fills internal variables
func (r *AutomationBundle) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "bundle")
		r.Bundle, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "type")
		r.Type, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "ext")
		r.Ext, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationTriggerScript request
func NewAutomationTriggerScript() *AutomationTriggerScript {
	return &AutomationTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"script": r.Script,
		"args":   r.Args,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerScript) GetScript() string {
	return r.Script
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerScript) GetArgs() map[string]interface{} {
	return r.Args
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerScript) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["script"]; ok && len(val) > 0 {
				r.Script, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["script"]; ok && len(val) > 0 {
			r.Script, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["args[]"]; ok {
			r.Args, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["args"]; ok {
			r.Args, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}
