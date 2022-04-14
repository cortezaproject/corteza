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
