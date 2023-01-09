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
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/types"
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
	DataPrivacyConnectionList struct {
		// ConnectionID GET parameter
		//
		// Filter by connection ID
		ConnectionID []string

		// Handle GET parameter
		//
		// Search handle to match against connections
		Handle string

		// Type GET parameter
		//
		// Search type to match against connections
		Type string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted connections
		Deleted filter.State
	}

	DataPrivacyRequestList struct {
		// RequestedBy GET parameter
		//
		// Filter by user ID
		RequestedBy []string

		// Query GET parameter
		//
		// Filter requests
		Query string

		// Kind GET parameter
		//
		// Filter by kind: correct, delete, export
		Kind []string

		// Status GET parameter
		//
		// Filter by status: pending, cancel, approve, reject
		Status []string

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

	DataPrivacyRequestCreate struct {
		// Kind POST parameter
		//
		// Request Kind
		Kind string

		// Payload POST parameter
		//
		// Request
		Payload types.DataPrivacyRequestPayloadSet
	}

	DataPrivacyRequestRead struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID uint64 `json:",string"`
	}

	DataPrivacyRequestUpdateStatus struct {
		// RequestID PATH parameter
		//
		// ID
		RequestID uint64 `json:",string"`

		// Status PATH parameter
		//
		// Request Status
		Status string
	}

	DataPrivacyRequestCommentList struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID uint64 `json:",string"`

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

	DataPrivacyRequestCommentCreate struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID uint64 `json:",string"`

		// Comment POST parameter
		//
		// Comment description
		Comment string
	}
)

// NewDataPrivacyConnectionList request
func NewDataPrivacyConnectionList() *DataPrivacyConnectionList {
	return &DataPrivacyConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"handle":       r.Handle,
		"type":         r.Type,
		"deleted":      r.Deleted,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetConnectionID() []string {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetDeleted() filter.State {
	return r.Deleted
}

// Fill processes request and fills internal variables
func (r *DataPrivacyConnectionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["connectionID[]"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["connectionID"]; ok {
			r.ConnectionID, err = val, nil
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
		if val, ok := tmp["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseFilterState(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDataPrivacyRequestList request
func NewDataPrivacyRequestList() *DataPrivacyRequestList {
	return &DataPrivacyRequestList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestedBy": r.RequestedBy,
		"query":       r.Query,
		"kind":        r.Kind,
		"status":      r.Status,
		"limit":       r.Limit,
		"pageCursor":  r.PageCursor,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetRequestedBy() []string {
	return r.RequestedBy
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetKind() []string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetStatus() []string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["requestedBy[]"]; ok {
			r.RequestedBy, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["requestedBy"]; ok {
			r.RequestedBy, err = val, nil
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
		if val, ok := tmp["kind[]"]; ok {
			r.Kind, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["kind"]; ok {
			r.Kind, err = val, nil
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

// NewDataPrivacyRequestCreate request
func NewDataPrivacyRequestCreate() *DataPrivacyRequestCreate {
	return &DataPrivacyRequestCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":    r.Kind,
		"payload": r.Payload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreate) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreate) GetPayload() types.DataPrivacyRequestPayloadSet {
	return r.Payload
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["kind"]; ok && len(val) > 0 {
				r.Kind, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["payload[]"]; ok {
				r.Payload, err = types.ParseDataPrivacyRequestPayload(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["payload"]; ok {
				r.Payload, err = types.ParseDataPrivacyRequestPayload(val)
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

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["payload[]"]; ok {
			r.Payload, err = types.ParseDataPrivacyRequestPayload(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["payload"]; ok {
			r.Payload, err = types.ParseDataPrivacyRequestPayload(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDataPrivacyRequestRead request
func NewDataPrivacyRequestRead() *DataPrivacyRequestRead {
	return &DataPrivacyRequestRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestRead) GetRequestID() uint64 {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDataPrivacyRequestUpdateStatus request
func NewDataPrivacyRequestUpdateStatus() *DataPrivacyRequestUpdateStatus {
	return &DataPrivacyRequestUpdateStatus{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestUpdateStatus) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
		"status":    r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestUpdateStatus) GetRequestID() uint64 {
	return r.RequestID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestUpdateStatus) GetStatus() string {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestUpdateStatus) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "status")
		r.Status, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDataPrivacyRequestCommentList request
func NewDataPrivacyRequestCommentList() *DataPrivacyRequestCommentList {
	return &DataPrivacyRequestCommentList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID":  r.RequestID,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentList) GetRequestID() uint64 {
	return r.RequestID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestCommentList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

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

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDataPrivacyRequestCommentCreate request
func NewDataPrivacyRequestCommentCreate() *DataPrivacyRequestCommentCreate {
	return &DataPrivacyRequestCommentCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
		"comment":   r.Comment,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentCreate) GetRequestID() uint64 {
	return r.RequestID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCommentCreate) GetComment() string {
	return r.Comment
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestCommentCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["comment"]; ok && len(val) > 0 {
				r.Comment, err = val[0], nil
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

		if val, ok := req.Form["comment"]; ok && len(val) > 0 {
			r.Comment, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
