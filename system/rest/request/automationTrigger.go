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
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
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
)

type (
	// Internal API interface
	AutomationTriggerList struct {
		// WorkflowID GET parameter
		//
		// Filter by workflow ID
		WorkflowID []string

		// TriggerID GET parameter
		//
		// Filter by trigger ID
		TriggerID []string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted triggers
		Deleted uint

		// EventType GET parameter
		//
		// Filter triggers by event type
		EventType string

		// ResourceType GET parameter
		//
		// Filter triggers by resoure type
		ResourceType string

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	AutomationTriggerCreate struct {
		// EventType POST parameter
		//
		// Event type
		EventType string

		// ResourceType POST parameter
		//
		// Resource type
		ResourceType string

		// Enabled POST parameter
		//
		// Is trigger enabled
		Enabled bool

		// WorkflowID POST parameter
		//
		// Workflow to be triggered
		WorkflowID uint64 `json:",string"`

		// WorkflowStepID POST parameter
		//
		// Start workflow in a specific step
		WorkflowStepID uint64 `json:",string"`

		// Input POST parameter
		//
		// Workflow meta data
		Input wfexec.Variables

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Constraints POST parameter
		//
		// Workflow steps definition
		Constraints types.TriggerConstraintSet

		// OwnedBy POST parameter
		//
		// Owner of the trigger
		OwnedBy uint64 `json:",string"`
	}

	AutomationTriggerUpdate struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`

		// EventType POST parameter
		//
		// Event type
		EventType string

		// ResourceType POST parameter
		//
		// Resource type
		ResourceType string

		// Enabled POST parameter
		//
		// Is trigger enabled
		Enabled bool

		// WorkflowID POST parameter
		//
		// Workflow to be triggered
		WorkflowID uint64 `json:",string"`

		// WorkflowStepID POST parameter
		//
		// Start workflow in a specific step
		WorkflowStepID uint64 `json:",string"`

		// Input POST parameter
		//
		// Workflow meta data
		Input wfexec.Variables

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Constraints POST parameter
		//
		// Workflow steps definition
		Constraints types.TriggerConstraintSet

		// OwnedBy POST parameter
		//
		// Owner of the trigger
		OwnedBy uint64 `json:",string"`
	}

	AutomationTriggerRead struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`
	}

	AutomationTriggerDelete struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`
	}

	AutomationTriggerUndelete struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`
	}
)

// NewAutomationTriggerList request
func NewAutomationTriggerList() *AutomationTriggerList {
	return &AutomationTriggerList{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID":   r.WorkflowID,
		"triggerID":    r.TriggerID,
		"deleted":      r.Deleted,
		"eventType":    r.EventType,
		"resourceType": r.ResourceType,
		"labels":       r.Labels,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetWorkflowID() []string {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetTriggerID() []string {
	return r.TriggerID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerList) Fill(req *http.Request) (err error) {
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
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["workflowID[]"]; ok {
			r.WorkflowID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["workflowID"]; ok {
			r.WorkflowID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["triggerID[]"]; ok {
			r.TriggerID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["triggerID"]; ok {
			r.TriggerID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["eventType"]; ok && len(val) > 0 {
			r.EventType, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := tmp["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAutomationTriggerCreate request
func NewAutomationTriggerCreate() *AutomationTriggerCreate {
	return &AutomationTriggerCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"eventType":      r.EventType,
		"resourceType":   r.ResourceType,
		"enabled":        r.Enabled,
		"workflowID":     r.WorkflowID,
		"workflowStepID": r.WorkflowStepID,
		"input":          r.Input,
		"labels":         r.Labels,
		"constraints":    r.Constraints,
		"ownedBy":        r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetWorkflowStepID() uint64 {
	return r.WorkflowStepID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetInput() wfexec.Variables {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetConstraints() types.TriggerConstraintSet {
	return r.Constraints
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerCreate) Fill(req *http.Request) (err error) {
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
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["eventType"]; ok && len(val) > 0 {
			r.EventType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["workflowID"]; ok && len(val) > 0 {
			r.WorkflowID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["workflowStepID"]; ok && len(val) > 0 {
			r.WorkflowStepID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["input[]"]; ok {
			r.Input, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["input"]; ok {
			r.Input, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["constraints[]"]; ok {
			r.Constraints, err = types.ParseTriggerConstraintSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["constraints"]; ok {
			r.Constraints, err = types.ParseTriggerConstraintSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAutomationTriggerUpdate request
func NewAutomationTriggerUpdate() *AutomationTriggerUpdate {
	return &AutomationTriggerUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID":      r.TriggerID,
		"eventType":      r.EventType,
		"resourceType":   r.ResourceType,
		"enabled":        r.Enabled,
		"workflowID":     r.WorkflowID,
		"workflowStepID": r.WorkflowStepID,
		"input":          r.Input,
		"labels":         r.Labels,
		"constraints":    r.Constraints,
		"ownedBy":        r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetTriggerID() uint64 {
	return r.TriggerID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetWorkflowStepID() uint64 {
	return r.WorkflowStepID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetInput() wfexec.Variables {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetConstraints() types.TriggerConstraintSet {
	return r.Constraints
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerUpdate) Fill(req *http.Request) (err error) {
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
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["eventType"]; ok && len(val) > 0 {
			r.EventType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["workflowID"]; ok && len(val) > 0 {
			r.WorkflowID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["workflowStepID"]; ok && len(val) > 0 {
			r.WorkflowStepID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["input[]"]; ok {
			r.Input, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["input"]; ok {
			r.Input, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["constraints[]"]; ok {
			r.Constraints, err = types.ParseTriggerConstraintSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["constraints"]; ok {
			r.Constraints, err = types.ParseTriggerConstraintSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "triggerID")
		r.TriggerID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationTriggerRead request
func NewAutomationTriggerRead() *AutomationTriggerRead {
	return &AutomationTriggerRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID": r.TriggerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerRead) GetTriggerID() uint64 {
	return r.TriggerID
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerRead) Fill(req *http.Request) (err error) {
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
		var val string
		// path params

		val = chi.URLParam(req, "triggerID")
		r.TriggerID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationTriggerDelete request
func NewAutomationTriggerDelete() *AutomationTriggerDelete {
	return &AutomationTriggerDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID": r.TriggerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerDelete) GetTriggerID() uint64 {
	return r.TriggerID
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerDelete) Fill(req *http.Request) (err error) {
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
		var val string
		// path params

		val = chi.URLParam(req, "triggerID")
		r.TriggerID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationTriggerUndelete request
func NewAutomationTriggerUndelete() *AutomationTriggerUndelete {
	return &AutomationTriggerUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID": r.TriggerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationTriggerUndelete) GetTriggerID() uint64 {
	return r.TriggerID
}

// Fill processes request and fills internal variables
func (r *AutomationTriggerUndelete) Fill(req *http.Request) (err error) {
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
		var val string
		// path params

		val = chi.URLParam(req, "triggerID")
		r.TriggerID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
