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
	AutomationWorkflowList struct {
		// WorkflowID GET parameter
		//
		// Filter by workflow ID
		WorkflowID []string

		// Query GET parameter
		//
		// Filter workflows
		Query string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted workflows
		Deleted uint

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

	AutomationWorkflowCreate struct {
		// Handle POST parameter
		//
		// Workflow name
		Handle string

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Meta POST parameter
		//
		// Workflow meta data
		Meta *types.WorkflowMeta

		// Enabled POST parameter
		//
		// Is workflow enabled
		Enabled bool

		// Trace POST parameter
		//
		// Trace workflow execution
		Trace bool

		// KeepSessions POST parameter
		//
		// Keep old workflow sessions
		KeepSessions int

		// Scope POST parameter
		//
		// Workflow meta data
		Scope wfexec.Variables

		// Steps POST parameter
		//
		// Workflow steps definition
		Steps types.WorkflowStepSet

		// Paths POST parameter
		//
		// Workflow step paths definition
		Paths types.WorkflowPathSet

		// RunAs POST parameter
		//
		// Is workflow enabled
		RunAs uint64 `json:",string"`

		// OwnedBy POST parameter
		//
		// Owner of the workflow
		OwnedBy uint64 `json:",string"`
	}

	AutomationWorkflowUpdate struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Workflow name
		Handle string

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Meta POST parameter
		//
		// Workflow meta data
		Meta *types.WorkflowMeta

		// Enabled POST parameter
		//
		// Is workflow enabled
		Enabled bool

		// Trace POST parameter
		//
		// Trace workflow execution
		Trace bool

		// KeepSessions POST parameter
		//
		// Keep old workflow sessions
		KeepSessions int

		// Scope POST parameter
		//
		// Workflow meta data
		Scope wfexec.Variables

		// Steps POST parameter
		//
		// Workflow steps definition
		Steps types.WorkflowStepSet

		// Paths POST parameter
		//
		// Workflow step paths definition
		Paths types.WorkflowPathSet

		// RunAs POST parameter
		//
		// Is workflow enabled
		RunAs uint64 `json:",string"`

		// OwnedBy POST parameter
		//
		// Owner of the workflow
		OwnedBy uint64 `json:",string"`
	}

	AutomationWorkflowRead struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`
	}

	AutomationWorkflowDelete struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`
	}

	AutomationWorkflowUndelete struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`
	}

	AutomationWorkflowTest struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`

		// Scope POST parameter
		//
		// Workflow meta data
		Scope wfexec.Variables

		// RunAs POST parameter
		//
		// Is workflow enabled
		RunAs bool
	}
)

// NewAutomationWorkflowList request
func NewAutomationWorkflowList() *AutomationWorkflowList {
	return &AutomationWorkflowList{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
		"query":      r.Query,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetWorkflowID() []string {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowList) Fill(req *http.Request) (err error) {
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
		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
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

// NewAutomationWorkflowCreate request
func NewAutomationWorkflowCreate() *AutomationWorkflowCreate {
	return &AutomationWorkflowCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":       r.Handle,
		"labels":       r.Labels,
		"meta":         r.Meta,
		"enabled":      r.Enabled,
		"trace":        r.Trace,
		"keepSessions": r.KeepSessions,
		"scope":        r.Scope,
		"steps":        r.Steps,
		"paths":        r.Paths,
		"runAs":        r.RunAs,
		"ownedBy":      r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetMeta() *types.WorkflowMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetTrace() bool {
	return r.Trace
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetKeepSessions() int {
	return r.KeepSessions
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetScope() wfexec.Variables {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetSteps() types.WorkflowStepSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetPaths() types.WorkflowPathSet {
	return r.Paths
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetRunAs() uint64 {
	return r.RunAs
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowCreate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
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
			r.Meta, err = types.ParseWorkflowMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseWorkflowMeta(val)
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

		if val, ok := req.Form["trace"]; ok && len(val) > 0 {
			r.Trace, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["keepSessions"]; ok && len(val) > 0 {
			r.KeepSessions, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope[]"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scope"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["steps[]"]; ok {
			r.Steps, err = types.ParseWorkflowStepSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["steps"]; ok {
			r.Steps, err = types.ParseWorkflowStepSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["paths[]"]; ok {
			r.Paths, err = types.ParseWorkflowPathSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["paths"]; ok {
			r.Paths, err = types.ParseWorkflowPathSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["runAs"]; ok && len(val) > 0 {
			r.RunAs, err = payload.ParseUint64(val[0]), nil
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

// NewAutomationWorkflowUpdate request
func NewAutomationWorkflowUpdate() *AutomationWorkflowUpdate {
	return &AutomationWorkflowUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID":   r.WorkflowID,
		"handle":       r.Handle,
		"labels":       r.Labels,
		"meta":         r.Meta,
		"enabled":      r.Enabled,
		"trace":        r.Trace,
		"keepSessions": r.KeepSessions,
		"scope":        r.Scope,
		"steps":        r.Steps,
		"paths":        r.Paths,
		"runAs":        r.RunAs,
		"ownedBy":      r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetMeta() *types.WorkflowMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetTrace() bool {
	return r.Trace
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetKeepSessions() int {
	return r.KeepSessions
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetScope() wfexec.Variables {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetSteps() types.WorkflowStepSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetPaths() types.WorkflowPathSet {
	return r.Paths
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetRunAs() uint64 {
	return r.RunAs
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowUpdate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
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
			r.Meta, err = types.ParseWorkflowMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseWorkflowMeta(val)
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

		if val, ok := req.Form["trace"]; ok && len(val) > 0 {
			r.Trace, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["keepSessions"]; ok && len(val) > 0 {
			r.KeepSessions, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope[]"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scope"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["steps[]"]; ok {
			r.Steps, err = types.ParseWorkflowStepSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["steps"]; ok {
			r.Steps, err = types.ParseWorkflowStepSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["paths[]"]; ok {
			r.Paths, err = types.ParseWorkflowPathSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["paths"]; ok {
			r.Paths, err = types.ParseWorkflowPathSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["runAs"]; ok && len(val) > 0 {
			r.RunAs, err = payload.ParseUint64(val[0]), nil
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

		val = chi.URLParam(req, "workflowID")
		r.WorkflowID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationWorkflowRead request
func NewAutomationWorkflowRead() *AutomationWorkflowRead {
	return &AutomationWorkflowRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowRead) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowRead) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "workflowID")
		r.WorkflowID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationWorkflowDelete request
func NewAutomationWorkflowDelete() *AutomationWorkflowDelete {
	return &AutomationWorkflowDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowDelete) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowDelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "workflowID")
		r.WorkflowID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationWorkflowUndelete request
func NewAutomationWorkflowUndelete() *AutomationWorkflowUndelete {
	return &AutomationWorkflowUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowUndelete) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowUndelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "workflowID")
		r.WorkflowID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAutomationWorkflowTest request
func NewAutomationWorkflowTest() *AutomationWorkflowTest {
	return &AutomationWorkflowTest{}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowTest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
		"scope":      r.Scope,
		"runAs":      r.RunAs,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowTest) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowTest) GetScope() wfexec.Variables {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r AutomationWorkflowTest) GetRunAs() bool {
	return r.RunAs
}

// Fill processes request and fills internal variables
func (r *AutomationWorkflowTest) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["scope[]"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scope"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["runAs"]; ok && len(val) > 0 {
			r.RunAs, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "workflowID")
		r.WorkflowID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
