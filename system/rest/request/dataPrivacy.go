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
	DataPrivacyListRequests struct {
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

	DataPrivacyCreateRequest struct {
		// Name POST parameter
		//
		// Request Name
		Name string

		// RequestType POST parameter
		//
		// Request Type
		RequestType uint
	}

	DataPrivacyUpdateRequest struct {
		// RequestID PATH parameter
		//
		// ID
		RequestID uint64 `json:",string"`

		// Status POST parameter
		//
		// Request Status
		Status uint
	}

	DataPrivacyUpdateRequestStatus struct {
		// RequestID PATH parameter
		//
		// ID
		RequestID uint64 `json:",string"`

		// Status PATH parameter
		//
		// Request Status
		Status uint
	}

	DataPrivacyReadRequest struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID uint64 `json:",string"`
	}

	DataPrivacyListResponsesOfRequest struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID string
	}

	DataPrivacyCreateResponseForRequest struct {
		// RequestID PATH parameter
		//
		// Request ID
		RequestID string
	}
)

// NewDataPrivacyListRequests request
func NewDataPrivacyListRequests() *DataPrivacyListRequests {
	return &DataPrivacyListRequests{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListRequests) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
		"deleted":    r.Deleted,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListRequests) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListRequests) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListRequests) GetSort() string {
	return r.Sort
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListRequests) GetDeleted() uint {
	return r.Deleted
}

// Fill processes request and fills internal variables
func (r *DataPrivacyListRequests) Fill(req *http.Request) (err error) {

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

// NewDataPrivacyCreateRequest request
func NewDataPrivacyCreateRequest() *DataPrivacyCreateRequest {
	return &DataPrivacyCreateRequest{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyCreateRequest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":        r.Name,
		"requestType": r.RequestType,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyCreateRequest) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyCreateRequest) GetRequestType() uint {
	return r.RequestType
}

// Fill processes request and fills internal variables
func (r *DataPrivacyCreateRequest) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["requestType"]; ok && len(val) > 0 {
				r.RequestType, err = payload.ParseUint(val[0]), nil
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

		if val, ok := req.Form["requestType"]; ok && len(val) > 0 {
			r.RequestType, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDataPrivacyUpdateRequest request
func NewDataPrivacyUpdateRequest() *DataPrivacyUpdateRequest {
	return &DataPrivacyUpdateRequest{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyUpdateRequest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
		"status":    r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyUpdateRequest) GetRequestID() uint64 {
	return r.RequestID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyUpdateRequest) GetStatus() uint {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *DataPrivacyUpdateRequest) Fill(req *http.Request) (err error) {

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
				r.Status, err = payload.ParseUint(val[0]), nil
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
			r.Status, err = payload.ParseUint(val[0]), nil
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

// NewDataPrivacyUpdateRequestStatus request
func NewDataPrivacyUpdateRequestStatus() *DataPrivacyUpdateRequestStatus {
	return &DataPrivacyUpdateRequestStatus{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyUpdateRequestStatus) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
		"status":    r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyUpdateRequestStatus) GetRequestID() uint64 {
	return r.RequestID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyUpdateRequestStatus) GetStatus() uint {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *DataPrivacyUpdateRequestStatus) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "requestID")
		r.RequestID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "status")
		r.Status, err = payload.ParseUint(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDataPrivacyReadRequest request
func NewDataPrivacyReadRequest() *DataPrivacyReadRequest {
	return &DataPrivacyReadRequest{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyReadRequest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyReadRequest) GetRequestID() uint64 {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyReadRequest) Fill(req *http.Request) (err error) {

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

// NewDataPrivacyListResponsesOfRequest request
func NewDataPrivacyListResponsesOfRequest() *DataPrivacyListResponsesOfRequest {
	return &DataPrivacyListResponsesOfRequest{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListResponsesOfRequest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyListResponsesOfRequest) GetRequestID() string {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyListResponsesOfRequest) Fill(req *http.Request) (err error) {

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

// NewDataPrivacyCreateResponseForRequest request
func NewDataPrivacyCreateResponseForRequest() *DataPrivacyCreateResponseForRequest {
	return &DataPrivacyCreateResponseForRequest{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyCreateResponseForRequest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"requestID": r.RequestID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyCreateResponseForRequest) GetRequestID() string {
	return r.RequestID
}

// Fill processes request and fills internal variables
func (r *DataPrivacyCreateResponseForRequest) Fill(req *http.Request) (err error) {

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
