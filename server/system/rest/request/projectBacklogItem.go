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
	ProjectBacklogItemList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// ProjectID GET parameter
		//
		// Filter by project ID
		ProjectID uint64 `json:",string"`

		// RevisionID GET parameter
		//
		// Filter by revision ID
		RevisionID uint64 `json:",string"`

		// EventID GET parameter
		//
		// Filter by parent category item ID
		EventID uint64 `json:",string"`

		// Category GET parameter
		//
		// Filter by category
		Category string

		// Status GET parameter
		//
		// Filter by status
		Status string

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

	ProjectBacklogItemCreate struct {
		// ProjectID POST parameter
		//
		// Project this backlog item belongs to
		ProjectID uint64 `json:",string"`

		// RevisionID POST parameter
		//
		// Revision to assign this backlog item to
		RevisionID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// Category POST parameter
		//
		// Parent category (incident|task|feature|privacy|review)
		Category string

		// EventID POST parameter
		//
		// Parent category item ID
		EventID uint64 `json:",string"`

		// Assignee POST parameter
		//
		// Assignee user ID
		Assignee uint64 `json:",string"`

		// Priority POST parameter
		//
		// Priority
		Priority string

		// Status POST parameter
		//
		// Status
		Status string

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string
	}

	ProjectBacklogItemRead struct {
		// BacklogItemID PATH parameter
		//
		// Backlog item ID
		BacklogItemID uint64 `json:",string"`
	}

	ProjectBacklogItemUpdate struct {
		// BacklogItemID PATH parameter
		//
		// Backlog item ID
		BacklogItemID uint64 `json:",string"`

		// RevisionID POST parameter
		//
		// Revision to assign this backlog item to
		RevisionID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// Category POST parameter
		//
		// Parent category (incident|task|feature|privacy|review)
		Category string

		// EventID POST parameter
		//
		// Parent category item ID
		EventID uint64 `json:",string"`

		// Assignee POST parameter
		//
		// Assignee user ID
		Assignee uint64 `json:",string"`

		// Priority POST parameter
		//
		// Priority
		Priority string

		// Status POST parameter
		//
		// Status
		Status string

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string
	}

	ProjectBacklogItemDelete struct {
		// BacklogItemID PATH parameter
		//
		// Backlog item ID
		BacklogItemID uint64 `json:",string"`
	}
)

// NewProjectBacklogItemList request
func NewProjectBacklogItemList() *ProjectBacklogItemList {
	return &ProjectBacklogItemList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"projectID":  r.ProjectID,
		"revisionID": r.RevisionID,
		"eventID":    r.EventID,
		"category":   r.Category,
		"status":     r.Status,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetEventID() uint64 {
	return r.EventID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetCategory() string {
	return r.Category
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectBacklogItemList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["revisionID"]; ok && len(val) > 0 {
			r.RevisionID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["eventID"]; ok && len(val) > 0 {
			r.EventID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["category"]; ok && len(val) > 0 {
			r.Category, err = val[0], nil
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

// NewProjectBacklogItemCreate request
func NewProjectBacklogItemCreate() *ProjectBacklogItemCreate {
	return &ProjectBacklogItemCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":   r.ProjectID,
		"revisionID":  r.RevisionID,
		"title":       r.Title,
		"description": r.Description,
		"category":    r.Category,
		"eventID":     r.EventID,
		"assignee":    r.Assignee,
		"priority":    r.Priority,
		"status":      r.Status,
		"dateDue":     r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetCategory() string {
	return r.Category
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetEventID() uint64 {
	return r.EventID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetAssignee() uint64 {
	return r.Assignee
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetPriority() string {
	return r.Priority
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemCreate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectBacklogItemCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["projectID"]; ok && len(val) > 0 {
				r.ProjectID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["revisionID"]; ok && len(val) > 0 {
				r.RevisionID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["category"]; ok && len(val) > 0 {
				r.Category, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["eventID"]; ok && len(val) > 0 {
				r.EventID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["assignee"]; ok && len(val) > 0 {
				r.Assignee, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["priority"]; ok && len(val) > 0 {
				r.Priority, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["dateDue"]; ok && len(val) > 0 {
				r.DateDue, err = val[0], nil
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

		if val, ok := req.Form["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["revisionID"]; ok && len(val) > 0 {
			r.RevisionID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["category"]; ok && len(val) > 0 {
			r.Category, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["eventID"]; ok && len(val) > 0 {
			r.EventID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["assignee"]; ok && len(val) > 0 {
			r.Assignee, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["priority"]; ok && len(val) > 0 {
			r.Priority, err = val[0], nil
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

		if val, ok := req.Form["dateDue"]; ok && len(val) > 0 {
			r.DateDue, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewProjectBacklogItemRead request
func NewProjectBacklogItemRead() *ProjectBacklogItemRead {
	return &ProjectBacklogItemRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"backlogItemID": r.BacklogItemID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemRead) GetBacklogItemID() uint64 {
	return r.BacklogItemID
}

// Fill processes request and fills internal variables
func (r *ProjectBacklogItemRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "backlogItemID")
		r.BacklogItemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectBacklogItemUpdate request
func NewProjectBacklogItemUpdate() *ProjectBacklogItemUpdate {
	return &ProjectBacklogItemUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"backlogItemID": r.BacklogItemID,
		"revisionID":    r.RevisionID,
		"title":         r.Title,
		"description":   r.Description,
		"category":      r.Category,
		"eventID":       r.EventID,
		"assignee":      r.Assignee,
		"priority":      r.Priority,
		"status":        r.Status,
		"dateDue":       r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetBacklogItemID() uint64 {
	return r.BacklogItemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetCategory() string {
	return r.Category
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetEventID() uint64 {
	return r.EventID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetAssignee() uint64 {
	return r.Assignee
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetPriority() string {
	return r.Priority
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemUpdate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectBacklogItemUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["revisionID"]; ok && len(val) > 0 {
				r.RevisionID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["category"]; ok && len(val) > 0 {
				r.Category, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["eventID"]; ok && len(val) > 0 {
				r.EventID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["assignee"]; ok && len(val) > 0 {
				r.Assignee, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["priority"]; ok && len(val) > 0 {
				r.Priority, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["dateDue"]; ok && len(val) > 0 {
				r.DateDue, err = val[0], nil
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

		if val, ok := req.Form["revisionID"]; ok && len(val) > 0 {
			r.RevisionID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["category"]; ok && len(val) > 0 {
			r.Category, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["eventID"]; ok && len(val) > 0 {
			r.EventID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["assignee"]; ok && len(val) > 0 {
			r.Assignee, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["priority"]; ok && len(val) > 0 {
			r.Priority, err = val[0], nil
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

		if val, ok := req.Form["dateDue"]; ok && len(val) > 0 {
			r.DateDue, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "backlogItemID")
		r.BacklogItemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectBacklogItemDelete request
func NewProjectBacklogItemDelete() *ProjectBacklogItemDelete {
	return &ProjectBacklogItemDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"backlogItemID": r.BacklogItemID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectBacklogItemDelete) GetBacklogItemID() uint64 {
	return r.BacklogItemID
}

// Fill processes request and fills internal variables
func (r *ProjectBacklogItemDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "backlogItemID")
		r.BacklogItemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
