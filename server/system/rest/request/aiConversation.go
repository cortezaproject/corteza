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
	AiConversationList struct {
		// AgentID GET parameter
		//
		// Filter by agent ID
		AgentID uint64 `json:",string"`

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted conversations
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

	AiConversationRead struct {
		// AiConversationID PATH parameter
		//
		// AI Conversation ID
		AiConversationID uint64 `json:",string"`
	}

	AiConversationDelete struct {
		// AiConversationID PATH parameter
		//
		// AI Conversation ID
		AiConversationID uint64 `json:",string"`
	}

	AiConversationUndelete struct {
		// AiConversationID PATH parameter
		//
		// AI Conversation ID
		AiConversationID uint64 `json:",string"`
	}
)

// NewAiConversationList request
func NewAiConversationList() *AiConversationList {
	return &AiConversationList{}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"agentID":    r.AgentID,
		"deleted":    r.Deleted,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) GetAgentID() uint64 {
	return r.AgentID
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *AiConversationList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["agentID"]; ok && len(val) > 0 {
			r.AgentID, err = payload.ParseUint64(val[0]), nil
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

// NewAiConversationRead request
func NewAiConversationRead() *AiConversationRead {
	return &AiConversationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"aiConversationID": r.AiConversationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationRead) GetAiConversationID() uint64 {
	return r.AiConversationID
}

// Fill processes request and fills internal variables
func (r *AiConversationRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "aiConversationID")
		r.AiConversationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAiConversationDelete request
func NewAiConversationDelete() *AiConversationDelete {
	return &AiConversationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"aiConversationID": r.AiConversationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationDelete) GetAiConversationID() uint64 {
	return r.AiConversationID
}

// Fill processes request and fills internal variables
func (r *AiConversationDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "aiConversationID")
		r.AiConversationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewAiConversationUndelete request
func NewAiConversationUndelete() *AiConversationUndelete {
	return &AiConversationUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"aiConversationID": r.AiConversationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AiConversationUndelete) GetAiConversationID() uint64 {
	return r.AiConversationID
}

// Fill processes request and fills internal variables
func (r *AiConversationUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "aiConversationID")
		r.AiConversationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
