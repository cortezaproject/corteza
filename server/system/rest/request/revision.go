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
	"github.com/cortezaproject/corteza/server/pkg/revisions"
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
	RevisionList struct {
		// ResourceType GET parameter
		//
		// Filter by resource type
		ResourceType string

		// Status GET parameter
		//
		// Filter by status (draft or empty for published)
		Status string

		// ResourceID GET parameter
		//
		// Filter by resource ID
		ResourceID uint64 `json:",string"`

		// DeletedOnly GET parameter
		//
		// Return only soft-deleted revisions
		DeletedOnly bool

		// Since GET parameter
		//
		// Return revisions after this timestamp
		Since *time.Time

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

	RevisionCreate struct {
		// ResourceType POST parameter
		//
		// Resource type
		ResourceType string

		// ResourceID POST parameter
		//
		// Resource ID
		ResourceID uint64 `json:",string"`

		// Status POST parameter
		//
		// Revision status (draft or empty)
		Status string

		// Changes POST parameter
		//
		// Changes (delta)
		Changes revisions.ChangeSet

		// Comment POST parameter
		//
		// Optional comment
		Comment string
	}

	RevisionDelete struct {
		// RevisionID PATH parameter
		//
		// Revision ID
		RevisionID uint64 `json:",string"`
	}
)

// NewRevisionList request
func NewRevisionList() *RevisionList {
	return &RevisionList{}
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resourceType": r.ResourceType,
		"status":       r.Status,
		"resourceID":   r.ResourceID,
		"deletedOnly":  r.DeletedOnly,
		"since":        r.Since,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetResourceID() uint64 {
	return r.ResourceID
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetDeletedOnly() bool {
	return r.DeletedOnly
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetSince() *time.Time {
	return r.Since
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r RevisionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *RevisionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
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
		if val, ok := tmp["resourceID"]; ok && len(val) > 0 {
			r.ResourceID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deletedOnly"]; ok && len(val) > 0 {
			r.DeletedOnly, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["since"]; ok && len(val) > 0 {
			r.Since, err = payload.ParseISODatePtrWithErr(val[0])
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

// NewRevisionCreate request
func NewRevisionCreate() *RevisionCreate {
	return &RevisionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r RevisionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resourceType": r.ResourceType,
		"resourceID":   r.ResourceID,
		"status":       r.Status,
		"changes":      r.Changes,
		"comment":      r.Comment,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RevisionCreate) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r RevisionCreate) GetResourceID() uint64 {
	return r.ResourceID
}

// Auditable returns all auditable/loggable parameters
func (r RevisionCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r RevisionCreate) GetChanges() revisions.ChangeSet {
	return r.Changes
}

// Auditable returns all auditable/loggable parameters
func (r RevisionCreate) GetComment() string {
	return r.Comment
}

// Fill processes request and fills internal variables
func (r *RevisionCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["resourceType"]; ok && len(val) > 0 {
				r.ResourceType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["resourceID"]; ok && len(val) > 0 {
				r.ResourceID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resourceID"]; ok && len(val) > 0 {
			r.ResourceID, err = payload.ParseUint64(val[0]), nil
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

		//if val, ok := req.Form["changes[]"]; ok && len(val) > 0  {
		//    r.Changes, err = revisions.ChangeSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["comment"]; ok && len(val) > 0 {
			r.Comment, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewRevisionDelete request
func NewRevisionDelete() *RevisionDelete {
	return &RevisionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RevisionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"revisionID": r.RevisionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RevisionDelete) GetRevisionID() uint64 {
	return r.RevisionID
}

// Fill processes request and fills internal variables
func (r *RevisionDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "revisionID")
		r.RevisionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
