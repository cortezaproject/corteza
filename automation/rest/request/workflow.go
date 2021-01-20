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
	WorkflowList struct {
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

		// Disabled GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) disabled workflows
		Disabled uint

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

	WorkflowCreate struct {
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
		Scope *expr.Vars

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

	WorkflowUpdate struct {
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
		Scope *expr.Vars

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

	WorkflowRead struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`
	}

	WorkflowDelete struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`
	}

	WorkflowUndelete struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`
	}

	WorkflowTest struct {
		// WorkflowID PATH parameter
		//
		// Workflow ID
		WorkflowID uint64 `json:",string"`

		// Scope POST parameter
		//
		// Workflow meta data
		Scope *expr.Vars

		// RunAs POST parameter
		//
		// Is workflow enabled
		RunAs bool
	}
)

// NewWorkflowList request
func NewWorkflowList() *WorkflowList {
	return &WorkflowList{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
		"query":      r.Query,
		"deleted":    r.Deleted,
		"disabled":   r.Disabled,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetWorkflowID() []string {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetDisabled() uint {
	return r.Disabled
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *WorkflowList) Fill(req *http.Request) (err error) {
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
		if val, ok := tmp["disabled"]; ok && len(val) > 0 {
			r.Disabled, err = payload.ParseUint(val[0]), nil
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

// NewWorkflowCreate request
func NewWorkflowCreate() *WorkflowCreate {
	return &WorkflowCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) Auditable() map[string]interface{} {
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
func (r WorkflowCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetMeta() *types.WorkflowMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetTrace() bool {
	return r.Trace
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetKeepSessions() int {
	return r.KeepSessions
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetScope() *expr.Vars {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetSteps() types.WorkflowStepSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetPaths() types.WorkflowPathSet {
	return r.Paths
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetRunAs() uint64 {
	return r.RunAs
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *WorkflowCreate) Fill(req *http.Request) (err error) {
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

// NewWorkflowUpdate request
func NewWorkflowUpdate() *WorkflowUpdate {
	return &WorkflowUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) Auditable() map[string]interface{} {
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
func (r WorkflowUpdate) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetMeta() *types.WorkflowMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetTrace() bool {
	return r.Trace
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetKeepSessions() int {
	return r.KeepSessions
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetScope() *expr.Vars {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetSteps() types.WorkflowStepSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetPaths() types.WorkflowPathSet {
	return r.Paths
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetRunAs() uint64 {
	return r.RunAs
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *WorkflowUpdate) Fill(req *http.Request) (err error) {
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

// NewWorkflowRead request
func NewWorkflowRead() *WorkflowRead {
	return &WorkflowRead{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowRead) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Fill processes request and fills internal variables
func (r *WorkflowRead) Fill(req *http.Request) (err error) {
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

// NewWorkflowDelete request
func NewWorkflowDelete() *WorkflowDelete {
	return &WorkflowDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowDelete) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Fill processes request and fills internal variables
func (r *WorkflowDelete) Fill(req *http.Request) (err error) {
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

// NewWorkflowUndelete request
func NewWorkflowUndelete() *WorkflowUndelete {
	return &WorkflowUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowUndelete) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Fill processes request and fills internal variables
func (r *WorkflowUndelete) Fill(req *http.Request) (err error) {
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

// NewWorkflowTest request
func NewWorkflowTest() *WorkflowTest {
	return &WorkflowTest{}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowTest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"workflowID": r.WorkflowID,
		"scope":      r.Scope,
		"runAs":      r.RunAs,
	}
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowTest) GetWorkflowID() uint64 {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowTest) GetScope() *expr.Vars {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r WorkflowTest) GetRunAs() bool {
	return r.RunAs
}

// Fill processes request and fills internal variables
func (r *WorkflowTest) Fill(req *http.Request) (err error) {
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
