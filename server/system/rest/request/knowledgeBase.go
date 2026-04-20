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
	"github.com/crusttech/human/server/system/types"
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
	KnowledgeBaseList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// Handle GET parameter
		//
		// Search by handle
		Handle string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted knowledge bases
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

	KnowledgeBaseCreate struct {
		// Handle POST parameter
		//
		// Handle
		Handle string

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// Context POST parameter
		//
		// Compose context
		Context types.KnowledgeBaseContext
	}

	KnowledgeBaseRead struct {
		// KnowledgeBaseID PATH parameter
		//
		// Knowledge Base ID
		KnowledgeBaseID uint64 `json:",string"`
	}

	KnowledgeBaseUpdate struct {
		// KnowledgeBaseID PATH parameter
		//
		// Knowledge Base ID
		KnowledgeBaseID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// Context POST parameter
		//
		// Compose context
		Context types.KnowledgeBaseContext

		// UpdatedAt POST parameter
		//
		// Last update (used for stale data check)
		UpdatedAt *time.Time
	}

	KnowledgeBaseDelete struct {
		// KnowledgeBaseID PATH parameter
		//
		// Knowledge Base ID
		KnowledgeBaseID uint64 `json:",string"`
	}

	KnowledgeBaseUndelete struct {
		// KnowledgeBaseID PATH parameter
		//
		// Knowledge Base ID
		KnowledgeBaseID uint64 `json:",string"`
	}
)

// NewKnowledgeBaseList request
func NewKnowledgeBaseList() *KnowledgeBaseList {
	return &KnowledgeBaseList{}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"handle":     r.Handle,
		"deleted":    r.Deleted,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *KnowledgeBaseList) Fill(req *http.Request) (err error) {

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

// NewKnowledgeBaseCreate request
func NewKnowledgeBaseCreate() *KnowledgeBaseCreate {
	return &KnowledgeBaseCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":      r.Handle,
		"title":       r.Title,
		"description": r.Description,
		"context":     r.Context,
	}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseCreate) GetContext() types.KnowledgeBaseContext {
	return r.Context
}

// Fill processes request and fills internal variables
func (r *KnowledgeBaseCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["context[]"]; ok {
				r.Context, err = types.ParseKnowledgeBaseContext(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["context"]; ok {
				r.Context, err = types.ParseKnowledgeBaseContext(val)
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

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["context[]"]; ok {
			r.Context, err = types.ParseKnowledgeBaseContext(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["context"]; ok {
			r.Context, err = types.ParseKnowledgeBaseContext(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewKnowledgeBaseRead request
func NewKnowledgeBaseRead() *KnowledgeBaseRead {
	return &KnowledgeBaseRead{}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"knowledgeBaseID": r.KnowledgeBaseID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseRead) GetKnowledgeBaseID() uint64 {
	return r.KnowledgeBaseID
}

// Fill processes request and fills internal variables
func (r *KnowledgeBaseRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "knowledgeBaseID")
		r.KnowledgeBaseID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewKnowledgeBaseUpdate request
func NewKnowledgeBaseUpdate() *KnowledgeBaseUpdate {
	return &KnowledgeBaseUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"knowledgeBaseID": r.KnowledgeBaseID,
		"handle":          r.Handle,
		"title":           r.Title,
		"description":     r.Description,
		"context":         r.Context,
		"updatedAt":       r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) GetKnowledgeBaseID() uint64 {
	return r.KnowledgeBaseID
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) GetContext() types.KnowledgeBaseContext {
	return r.Context
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *KnowledgeBaseUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["context[]"]; ok {
				r.Context, err = types.ParseKnowledgeBaseContext(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["context"]; ok {
				r.Context, err = types.ParseKnowledgeBaseContext(val)
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

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["context[]"]; ok {
			r.Context, err = types.ParseKnowledgeBaseContext(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["context"]; ok {
			r.Context, err = types.ParseKnowledgeBaseContext(val)
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

		val = chi.URLParam(req, "knowledgeBaseID")
		r.KnowledgeBaseID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewKnowledgeBaseDelete request
func NewKnowledgeBaseDelete() *KnowledgeBaseDelete {
	return &KnowledgeBaseDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"knowledgeBaseID": r.KnowledgeBaseID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseDelete) GetKnowledgeBaseID() uint64 {
	return r.KnowledgeBaseID
}

// Fill processes request and fills internal variables
func (r *KnowledgeBaseDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "knowledgeBaseID")
		r.KnowledgeBaseID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewKnowledgeBaseUndelete request
func NewKnowledgeBaseUndelete() *KnowledgeBaseUndelete {
	return &KnowledgeBaseUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"knowledgeBaseID": r.KnowledgeBaseID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r KnowledgeBaseUndelete) GetKnowledgeBaseID() uint64 {
	return r.KnowledgeBaseID
}

// Fill processes request and fills internal variables
func (r *KnowledgeBaseUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "knowledgeBaseID")
		r.KnowledgeBaseID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
