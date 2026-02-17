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
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
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
	AgentList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// Handle GET parameter
		//
		// Search by handle
		Handle string

		// Status GET parameter
		//
		// Filter by status
		Status string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted agents
		Deleted uint

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total counter
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	AgentCreate struct {
		// Handle POST parameter
		//
		// Agent handle
		Handle string

		// Status POST parameter
		//
		// Agent status
		Status string

		// Meta POST parameter
		//
		// Agent meta
		Meta types.AgentMeta

		// Behavior POST parameter
		//
		// Agent behavior
		Behavior types.AgentBehavior

		// Execution POST parameter
		//
		// Agent execution
		Execution types.AgentExecution

		// Access POST parameter
		//
		// Agent access
		Access types.AgentAccess

		// Invocation POST parameter
		//
		// Agent invocation
		Invocation types.AgentInvocation
	}

	AgentRead struct {
		// AgentID PATH parameter
		//
		// Agent ID
		AgentID uint64 `json:",string"`
	}

	AgentUpdate struct {
		// AgentID PATH parameter
		//
		// Agent ID
		AgentID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Agent handle
		Handle string

		// Status POST parameter
		//
		// Agent status
		Status string

		// Meta POST parameter
		//
		// Agent meta
		Meta types.AgentMeta

		// Behavior POST parameter
		//
		// Agent behavior
		Behavior types.AgentBehavior

		// Execution POST parameter
		//
		// Agent execution
		Execution types.AgentExecution

		// Access POST parameter
		//
		// Agent access
		Access types.AgentAccess

		// Invocation POST parameter
		//
		// Agent invocation
		Invocation types.AgentInvocation

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	AgentDelete struct {
		// AgentID PATH parameter
		//
		// Agent ID
		AgentID uint64 `json:",string"`
	}

	AgentUndelete struct {
		// AgentID PATH parameter
		//
		// Agent ID
		AgentID uint64 `json:",string"`
	}

	AgentExec struct {
		// AgentID PATH parameter
		//
		// Agent ID
		AgentID uint64 `json:",string"`

		// Input POST parameter
		//
		// User input message
		Input string

		// ConversationID POST parameter
		//
		// Conversation ID for multi-turn
		ConversationID uint64 `json:",string"`
	}
)

// NewAgentList request
func NewAgentList() *AgentList {
	return &AgentList{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"handle":     r.Handle,
		"status":     r.Status,
		"deleted":    r.Deleted,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r AgentList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *AgentList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
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
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
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

// NewAgentCreate request
func NewAgentCreate() *AgentCreate {
	return &AgentCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":     r.Handle,
		"status":     r.Status,
		"meta":       r.Meta,
		"behavior":   r.Behavior,
		"execution":  r.Execution,
		"access":     r.Access,
		"invocation": r.Invocation,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetMeta() types.AgentMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetBehavior() types.AgentBehavior {
	return r.Behavior
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetExecution() types.AgentExecution {
	return r.Execution
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetAccess() types.AgentAccess {
	return r.Access
}

// Auditable returns all auditable/loggable parameters
func (r AgentCreate) GetInvocation() types.AgentInvocation {
	return r.Invocation
}

// Fill processes request and fills internal variables
func (r *AgentCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseAgentMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseAgentMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["behavior[]"]; ok {
				r.Behavior, err = types.ParseAgentBehavior(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["behavior"]; ok {
				r.Behavior, err = types.ParseAgentBehavior(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["execution[]"]; ok {
				r.Execution, err = types.ParseAgentExecution(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["execution"]; ok {
				r.Execution, err = types.ParseAgentExecution(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["access[]"]; ok {
				r.Access, err = types.ParseAgentAccess(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["access"]; ok {
				r.Access, err = types.ParseAgentAccess(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["invocation[]"]; ok {
				r.Invocation, err = types.ParseAgentInvocation(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["invocation"]; ok {
				r.Invocation, err = types.ParseAgentInvocation(val)
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseAgentMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseAgentMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["behavior[]"]; ok {
			r.Behavior, err = types.ParseAgentBehavior(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["behavior"]; ok {
			r.Behavior, err = types.ParseAgentBehavior(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["execution[]"]; ok {
			r.Execution, err = types.ParseAgentExecution(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["execution"]; ok {
			r.Execution, err = types.ParseAgentExecution(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["access[]"]; ok {
			r.Access, err = types.ParseAgentAccess(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["access"]; ok {
			r.Access, err = types.ParseAgentAccess(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["invocation[]"]; ok {
			r.Invocation, err = types.ParseAgentInvocation(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["invocation"]; ok {
			r.Invocation, err = types.ParseAgentInvocation(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAgentRead request
func NewAgentRead() *AgentRead {
	return &AgentRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"agentID": r.AgentID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentRead) GetAgentID() uint64 {
	return r.AgentID
}

// Fill processes request and fills internal variables
func (r *AgentRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "agentID")
		r.AgentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAgentUpdate request
func NewAgentUpdate() *AgentUpdate {
	return &AgentUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"agentID":    r.AgentID,
		"handle":     r.Handle,
		"status":     r.Status,
		"meta":       r.Meta,
		"behavior":   r.Behavior,
		"execution":  r.Execution,
		"access":     r.Access,
		"invocation": r.Invocation,
		"updatedAt":  r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetAgentID() uint64 {
	return r.AgentID
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetMeta() types.AgentMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetBehavior() types.AgentBehavior {
	return r.Behavior
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetExecution() types.AgentExecution {
	return r.Execution
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetAccess() types.AgentAccess {
	return r.Access
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetInvocation() types.AgentInvocation {
	return r.Invocation
}

// Auditable returns all auditable/loggable parameters
func (r AgentUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *AgentUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseAgentMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseAgentMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["behavior[]"]; ok {
				r.Behavior, err = types.ParseAgentBehavior(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["behavior"]; ok {
				r.Behavior, err = types.ParseAgentBehavior(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["execution[]"]; ok {
				r.Execution, err = types.ParseAgentExecution(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["execution"]; ok {
				r.Execution, err = types.ParseAgentExecution(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["access[]"]; ok {
				r.Access, err = types.ParseAgentAccess(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["access"]; ok {
				r.Access, err = types.ParseAgentAccess(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["invocation[]"]; ok {
				r.Invocation, err = types.ParseAgentInvocation(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["invocation"]; ok {
				r.Invocation, err = types.ParseAgentInvocation(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["updatedAt"]; ok && len(val) > 0 {
				r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseAgentMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseAgentMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["behavior[]"]; ok {
			r.Behavior, err = types.ParseAgentBehavior(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["behavior"]; ok {
			r.Behavior, err = types.ParseAgentBehavior(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["execution[]"]; ok {
			r.Execution, err = types.ParseAgentExecution(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["execution"]; ok {
			r.Execution, err = types.ParseAgentExecution(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["access[]"]; ok {
			r.Access, err = types.ParseAgentAccess(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["access"]; ok {
			r.Access, err = types.ParseAgentAccess(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["invocation[]"]; ok {
			r.Invocation, err = types.ParseAgentInvocation(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["invocation"]; ok {
			r.Invocation, err = types.ParseAgentInvocation(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["updatedAt"]; ok && len(val) > 0 {
			r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "agentID")
		r.AgentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAgentDelete request
func NewAgentDelete() *AgentDelete {
	return &AgentDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"agentID": r.AgentID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentDelete) GetAgentID() uint64 {
	return r.AgentID
}

// Fill processes request and fills internal variables
func (r *AgentDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "agentID")
		r.AgentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAgentUndelete request
func NewAgentUndelete() *AgentUndelete {
	return &AgentUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"agentID": r.AgentID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentUndelete) GetAgentID() uint64 {
	return r.AgentID
}

// Fill processes request and fills internal variables
func (r *AgentUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "agentID")
		r.AgentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAgentExec request
func NewAgentExec() *AgentExec {
	return &AgentExec{}
}

// Auditable returns all auditable/loggable parameters
func (r AgentExec) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"agentID":        r.AgentID,
		"input":          r.Input,
		"conversationID": r.ConversationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AgentExec) GetAgentID() uint64 {
	return r.AgentID
}

// Auditable returns all auditable/loggable parameters
func (r AgentExec) GetInput() string {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r AgentExec) GetConversationID() uint64 {
	return r.ConversationID
}

// Fill processes request and fills internal variables
func (r *AgentExec) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
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

			if val, ok := req.MultipartForm.Value["input"]; ok && len(val) > 0 {
				r.Input, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["conversationID"]; ok && len(val) > 0 {
				r.ConversationID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["input"]; ok && len(val) > 0 {
			r.Input, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["conversationID"]; ok && len(val) > 0 {
			r.ConversationID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "agentID")
		r.AgentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
