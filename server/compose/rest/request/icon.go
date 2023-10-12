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
	IconList struct {
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

	IconUpload struct {
		// Icon POST parameter
		//
		// Icon to upload
		Icon *multipart.FileHeader
	}

	IconDelete struct {
		// IconID PATH parameter
		//
		// Icon ID
		IconID uint64 `json:",string"`
	}
)

// NewIconList request
func NewIconList() *IconList {
	return &IconList{}
}

// Auditable returns all auditable/loggable parameters
func (r IconList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r IconList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r IconList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r IconList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r IconList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *IconList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

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

// NewIconUpload request
func NewIconUpload() *IconUpload {
	return &IconUpload{}
}

// Auditable returns all auditable/loggable parameters
func (r IconUpload) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"icon": r.Icon,
	}
}

// Auditable returns all auditable/loggable parameters
func (r IconUpload) GetIcon() *multipart.FileHeader {
	return r.Icon
}

// Fill processes request and fills internal variables
func (r *IconUpload) Fill(req *http.Request) (err error) {

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

			// Ignoring icon as its handled in the POST params section
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if _, r.Icon, err = req.FormFile("icon"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

	}

	return err
}

// NewIconDelete request
func NewIconDelete() *IconDelete {
	return &IconDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r IconDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"iconID": r.IconID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r IconDelete) GetIconID() uint64 {
	return r.IconID
}

// Fill processes request and fills internal variables
func (r *IconDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "iconID")
		r.IconID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
