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
	DataPrivacyRequestList struct {
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

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted requests
		Deleted uint
	}

	DataPrivacyRequestCreate struct {
		// Name POST parameter
		//
		// Request Name
		Name string

		// Kind POST parameter
		//
		// Request Kind
		Kind string
	}

	DataPrivacyRequestUpdate struct {
		// RequestID PATH parameter
		//
		// ID
		RequestID uint64 `json:",string"`

		// Status POST parameter
		//
		// Request Status
		Status string
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

	DataPrivacyRequestRead struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID uint64 `json:",string"`
	}

	DataPrivacyRequestListResponses struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID string
	}

	DataPrivacyRequestCreateResponse struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID string
	}
)

// NewDataPrivacyRequestList request
func NewDataPrivacyRequestList() *DataPrivacyRequestList {
	return &DataPrivacyRequestList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
		"deleted":    r.Deleted,
	}
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

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestList) GetDeleted() uint {
	return r.Deleted
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestList) Fill(req *http.Request) (err error) {

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
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
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
		"name": r.Name,
		"kind": r.Kind,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreate) GetKind() string {
	return r.Kind
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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["kind"]; ok && len(val) > 0 {
				r.Kind, err = val[0], nil
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDataPrivacyRequestUpdate request
func NewDataPrivacyRequestUpdate() *DataPrivacyRequestUpdate {
	return &DataPrivacyRequestUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
		"status":    r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestUpdate) GetRequestID() uint64 {
	return r.RequestID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestUpdate) GetStatus() string {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
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

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
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

// NewDataPrivacyRequestListResponses request
func NewDataPrivacyRequestListResponses() *DataPrivacyRequestListResponses {
	return &DataPrivacyRequestListResponses{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestListResponses) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestListResponses) GetRequestID() string {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestListResponses) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDataPrivacyRequestCreateResponse request
func NewDataPrivacyRequestCreateResponse() *DataPrivacyRequestCreateResponse {
	return &DataPrivacyRequestCreateResponse{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreateResponse) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyRequestCreateResponse) GetRequestID() string {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyRequestCreateResponse) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
