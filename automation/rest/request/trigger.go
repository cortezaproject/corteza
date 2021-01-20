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
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/payload"
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
	TriggerList struct {
		// TriggerID GET parameter
		//
		// Filter by trigger ID
		TriggerID []string

		// WorkflowID GET parameter
		//
		// Filter by workflow ID
		WorkflowID []string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted triggers
		Deleted uint

		// Disabled GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) disabled triggers
		Disabled uint

		// EventType GET parameter
		//
		// Filter triggers by event type
		EventType string

		// ResourceType GET parameter
		//
		// Filter triggers by resource type
		ResourceType string

		// Query GET parameter
		//
		// Filter workflows,
		Query string

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

	TriggerCreate struct {
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
		Input *expr.Vars

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Meta POST parameter
		//
		// Trigger meta data
		Meta *types.TriggerMeta

		// Constraints POST parameter
		//
		// Workflow steps definition
		Constraints types.TriggerConstraintSet

		// OwnedBy POST parameter
		//
		// Owner of the trigger
		OwnedBy uint64 `json:",string"`
	}

	TriggerUpdate struct {
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
		Input *expr.Vars

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Meta POST parameter
		//
		// Trigger meta data
		Meta *types.TriggerMeta

		// Constraints POST parameter
		//
		// Workflow steps definition
		Constraints types.TriggerConstraintSet

		// OwnedBy POST parameter
		//
		// Owner of the trigger
		OwnedBy uint64 `json:",string"`
	}

	TriggerRead struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`
	}

	TriggerDelete struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`
	}

	TriggerUndelete struct {
		// TriggerID PATH parameter
		//
		// Trigger ID
		TriggerID uint64 `json:",string"`
	}
)

// NewTriggerList request
func NewTriggerList() *TriggerList {
	return &TriggerList{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID":    r.TriggerID,
		"workflowID":   r.WorkflowID,
		"deleted":      r.Deleted,
		"disabled":     r.Disabled,
		"eventType":    r.EventType,
		"resourceType": r.ResourceType,
		"query":        r.Query,
		"labels":       r.Labels,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetTriggerID() []string {
	return r.TriggerID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetWorkflowID() []string {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetDisabled() uint {
	return r.Disabled
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r TriggerList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *TriggerList) Fill(req *http.Request) (err error) {
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
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["disabled"]; ok && len(val) > 0 {
			r.Disabled, err = payload.ParseUint(val[0]), nil
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
		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
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

// NewTriggerCreate request
func NewTriggerCreate() *TriggerCreate {
	return &TriggerCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"eventType":      r.EventType,
		"resourceType":   r.ResourceType,
		"enabled":        r.Enabled,
		"workflowID":     r.WorkflowID,
		"workflowStepID": r.WorkflowStepID,
		"input":          r.Input,
		"labels":         r.Labels,
		"meta":           r.Meta,
		"constraints":    r.Constraints,
		"ownedBy":        r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetWorkflowStepID() uint64 {
	return r.WorkflowStepID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetInput() *expr.Vars {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetMeta() *types.TriggerMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetConstraints() types.TriggerConstraintSet {
	return r.Constraints
}

// Auditable returns all auditable/loggable parameters
func (r TriggerCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *TriggerCreate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTriggerMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTriggerMeta(val)
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

// NewTriggerUpdate request
func NewTriggerUpdate() *TriggerUpdate {
	return &TriggerUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID":      r.TriggerID,
		"eventType":      r.EventType,
		"resourceType":   r.ResourceType,
		"enabled":        r.Enabled,
		"workflowID":     r.WorkflowID,
		"workflowStepID": r.WorkflowStepID,
		"input":          r.Input,
		"labels":         r.Labels,
		"meta":           r.Meta,
		"constraints":    r.Constraints,
		"ownedBy":        r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetTriggerID() uint64 {
	return r.TriggerID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetWorkflowStepID() uint64 {
	return r.WorkflowStepID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetInput() *expr.Vars {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetMeta() *types.TriggerMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetConstraints() types.TriggerConstraintSet {
	return r.Constraints
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *TriggerUpdate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTriggerMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTriggerMeta(val)
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

// NewTriggerRead request
func NewTriggerRead() *TriggerRead {
	return &TriggerRead{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID": r.TriggerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerRead) GetTriggerID() uint64 {
	return r.TriggerID
}

// Fill processes request and fills internal variables
func (r *TriggerRead) Fill(req *http.Request) (err error) {
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

// NewTriggerDelete request
func NewTriggerDelete() *TriggerDelete {
	return &TriggerDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID": r.TriggerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDelete) GetTriggerID() uint64 {
	return r.TriggerID
}

// Fill processes request and fills internal variables
func (r *TriggerDelete) Fill(req *http.Request) (err error) {
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

// NewTriggerUndelete request
func NewTriggerUndelete() *TriggerUndelete {
	return &TriggerUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerID": r.TriggerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerUndelete) GetTriggerID() uint64 {
	return r.TriggerID
}

// Fill processes request and fills internal variables
func (r *TriggerUndelete) Fill(req *http.Request) (err error) {
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
