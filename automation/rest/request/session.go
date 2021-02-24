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
	SessionList struct {
		// SessionID GET parameter
		//
		// Filter by session ID
		SessionID []string

		// WorkflowID GET parameter
		//
		// Filter by workflow ID
		WorkflowID []string

		// Completed GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) completed sessions
		Completed uint

		// Status GET parameter
		//
		// Filter by status: started (0), prompted (1), suspended (2), failed (3) and completed (4)
		Status []uint

		// EventType GET parameter
		//
		// Filter event type
		EventType string

		// ResourceType GET parameter
		//
		// Filter resource type
		ResourceType string

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

	SessionRead struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID uint64 `json:",string"`
	}

	SessionTrace struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID uint64 `json:",string"`
	}

	SessionDelete struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID uint64 `json:",string"`
	}

	SessionListPrompts struct {
	}

	SessionResumeState struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID uint64 `json:",string"`

		// StateID PATH parameter
		//
		// State ID
		StateID uint64 `json:",string"`

		// Input POST parameter
		//
		// Prompt variables
		Input *expr.Vars
	}

	SessionDeleteState struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID uint64 `json:",string"`

		// StateID PATH parameter
		//
		// State ID
		StateID uint64 `json:",string"`
	}
)

// NewSessionList request
func NewSessionList() *SessionList {
	return &SessionList{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID":    r.SessionID,
		"workflowID":   r.WorkflowID,
		"completed":    r.Completed,
		"status":       r.Status,
		"eventType":    r.EventType,
		"resourceType": r.ResourceType,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetSessionID() []string {
	return r.SessionID
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetWorkflowID() []string {
	return r.WorkflowID
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetCompleted() uint {
	return r.Completed
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetStatus() []uint {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r SessionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *SessionList) Fill(req *http.Request) (err error) {
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

		if val, ok := tmp["sessionID[]"]; ok {
			r.SessionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["sessionID"]; ok {
			r.SessionID, err = val, nil
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
		if val, ok := tmp["completed"]; ok && len(val) > 0 {
			r.Completed, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status[]"]; ok {
			r.Status, err = payload.ParseUints(val), nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["status"]; ok {
			r.Status, err = payload.ParseUints(val), nil
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

// NewSessionRead request
func NewSessionRead() *SessionRead {
	return &SessionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SessionRead) GetSessionID() uint64 {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *SessionRead) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSessionTrace request
func NewSessionTrace() *SessionTrace {
	return &SessionTrace{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionTrace) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SessionTrace) GetSessionID() uint64 {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *SessionTrace) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSessionDelete request
func NewSessionDelete() *SessionDelete {
	return &SessionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SessionDelete) GetSessionID() uint64 {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *SessionDelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSessionListPrompts request
func NewSessionListPrompts() *SessionListPrompts {
	return &SessionListPrompts{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionListPrompts) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *SessionListPrompts) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	return err
}

// NewSessionResumeState request
func NewSessionResumeState() *SessionResumeState {
	return &SessionResumeState{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionResumeState) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
		"stateID":   r.StateID,
		"input":     r.Input,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SessionResumeState) GetSessionID() uint64 {
	return r.SessionID
}

// Auditable returns all auditable/loggable parameters
func (r SessionResumeState) GetStateID() uint64 {
	return r.StateID
}

// Auditable returns all auditable/loggable parameters
func (r SessionResumeState) GetInput() *expr.Vars {
	return r.Input
}

// Fill processes request and fills internal variables
func (r *SessionResumeState) Fill(req *http.Request) (err error) {
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
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "stateID")
		r.StateID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSessionDeleteState request
func NewSessionDeleteState() *SessionDeleteState {
	return &SessionDeleteState{}
}

// Auditable returns all auditable/loggable parameters
func (r SessionDeleteState) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
		"stateID":   r.StateID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SessionDeleteState) GetSessionID() uint64 {
	return r.SessionID
}

// Auditable returns all auditable/loggable parameters
func (r SessionDeleteState) GetStateID() uint64 {
	return r.StateID
}

// Fill processes request and fills internal variables
func (r *SessionDeleteState) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "stateID")
		r.StateID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
