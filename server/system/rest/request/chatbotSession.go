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
	"github.com/crusttech/human/server/pkg/payload"
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
	ChatbotSessionList struct {
		// ChatbotID GET parameter
		//
		// Filter by chatbot ID
		ChatbotID uint64 `json:",string"`

		// Status GET parameter
		//
		// Filter by status
		Status []string

		// Source GET parameter
		//
		// Filter by source (widget|preview)
		Source string

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

	ChatbotSessionRead struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string
	}

	ChatbotSessionAdvanceStep struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string
	}

	ChatbotSessionClose struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string
	}

	ChatbotSessionHandoffAccept struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string

		// Operator POST parameter
		//
		// Operator label
		Operator string
	}

	ChatbotSessionOperatorMessage struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string

		// Message POST parameter
		//
		// Message body
		Message string

		// Operator POST parameter
		//
		// Operator label
		Operator string
	}

	ChatbotSessionHandoffComplete struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string
	}

	ChatbotSessionStream struct {
		// SessionID PATH parameter
		//
		// Session ID
		SessionID string
	}
)

// NewChatbotSessionList request
func NewChatbotSessionList() *ChatbotSessionList {
	return &ChatbotSessionList{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID":  r.ChatbotID,
		"status":     r.Status,
		"source":     r.Source,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetChatbotID() uint64 {
	return r.ChatbotID
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetStatus() []string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetSource() string {
	return r.Source
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["chatbotID"]; ok && len(val) > 0 {
			r.ChatbotID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status[]"]; ok {
			r.Status, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["status"]; ok {
			r.Status, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["source"]; ok && len(val) > 0 {
			r.Source, err = val[0], nil
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

// NewChatbotSessionRead request
func NewChatbotSessionRead() *ChatbotSessionRead {
	return &ChatbotSessionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionRead) GetSessionID() string {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotSessionAdvanceStep request
func NewChatbotSessionAdvanceStep() *ChatbotSessionAdvanceStep {
	return &ChatbotSessionAdvanceStep{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionAdvanceStep) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionAdvanceStep) GetSessionID() string {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionAdvanceStep) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotSessionClose request
func NewChatbotSessionClose() *ChatbotSessionClose {
	return &ChatbotSessionClose{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionClose) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionClose) GetSessionID() string {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionClose) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotSessionHandoffAccept request
func NewChatbotSessionHandoffAccept() *ChatbotSessionHandoffAccept {
	return &ChatbotSessionHandoffAccept{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionHandoffAccept) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
		"operator":  r.Operator,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionHandoffAccept) GetSessionID() string {
	return r.SessionID
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionHandoffAccept) GetOperator() string {
	return r.Operator
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionHandoffAccept) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["operator"]; ok && len(val) > 0 {
				r.Operator, err = val[0], nil
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

		if val, ok := req.Form["operator"]; ok && len(val) > 0 {
			r.Operator, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotSessionOperatorMessage request
func NewChatbotSessionOperatorMessage() *ChatbotSessionOperatorMessage {
	return &ChatbotSessionOperatorMessage{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionOperatorMessage) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
		"message":   r.Message,
		"operator":  r.Operator,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionOperatorMessage) GetSessionID() string {
	return r.SessionID
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionOperatorMessage) GetMessage() string {
	return r.Message
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionOperatorMessage) GetOperator() string {
	return r.Operator
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionOperatorMessage) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["message"]; ok && len(val) > 0 {
				r.Message, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["operator"]; ok && len(val) > 0 {
				r.Operator, err = val[0], nil
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

		if val, ok := req.Form["message"]; ok && len(val) > 0 {
			r.Message, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["operator"]; ok && len(val) > 0 {
			r.Operator, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotSessionHandoffComplete request
func NewChatbotSessionHandoffComplete() *ChatbotSessionHandoffComplete {
	return &ChatbotSessionHandoffComplete{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionHandoffComplete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionHandoffComplete) GetSessionID() string {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionHandoffComplete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotSessionStream request
func NewChatbotSessionStream() *ChatbotSessionStream {
	return &ChatbotSessionStream{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionStream) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sessionID": r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotSessionStream) GetSessionID() string {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *ChatbotSessionStream) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
