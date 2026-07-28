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
	ProjectReviewList struct {
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

	ProjectReviewCreate struct {
		// ProjectID POST parameter
		//
		// Project this review belongs to
		ProjectID uint64 `json:",string"`

		// RevisionID POST parameter
		//
		// Revision to assign this review to
		RevisionID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// ReviewType POST parameter
		//
		// Review type
		ReviewType string

		// ReviewFrequency POST parameter
		//
		// Review frequency
		ReviewFrequency string

		// Scope POST parameter
		//
		// Scope
		Scope string

		// Status POST parameter
		//
		// Status
		Status string

		// Reviewer POST parameter
		//
		// Reviewer user ID
		Reviewer uint64 `json:",string"`

		// ApprovedBy POST parameter
		//
		// Approved-by user ID
		ApprovedBy uint64 `json:",string"`

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string
	}

	ProjectReviewRead struct {
		// ReviewID PATH parameter
		//
		// Review ID
		ReviewID uint64 `json:",string"`
	}

	ProjectReviewUpdate struct {
		// ReviewID PATH parameter
		//
		// Review ID
		ReviewID uint64 `json:",string"`

		// RevisionID POST parameter
		//
		// Revision to assign this review to
		RevisionID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// ReviewType POST parameter
		//
		// Review type
		ReviewType string

		// ReviewFrequency POST parameter
		//
		// Review frequency
		ReviewFrequency string

		// Scope POST parameter
		//
		// Scope
		Scope string

		// Status POST parameter
		//
		// Status
		Status string

		// Reviewer POST parameter
		//
		// Reviewer user ID
		Reviewer uint64 `json:",string"`

		// ApprovedBy POST parameter
		//
		// Approved-by user ID
		ApprovedBy uint64 `json:",string"`

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string
	}

	ProjectReviewDelete struct {
		// ReviewID PATH parameter
		//
		// Review ID
		ReviewID uint64 `json:",string"`
	}
)

// NewProjectReviewList request
func NewProjectReviewList() *ProjectReviewList {
	return &ProjectReviewList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"projectID":  r.ProjectID,
		"revisionID": r.RevisionID,
		"status":     r.Status,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectReviewList) Fill(req *http.Request) (err error) {

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

// NewProjectReviewCreate request
func NewProjectReviewCreate() *ProjectReviewCreate {
	return &ProjectReviewCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":       r.ProjectID,
		"revisionID":      r.RevisionID,
		"title":           r.Title,
		"description":     r.Description,
		"reviewType":      r.ReviewType,
		"reviewFrequency": r.ReviewFrequency,
		"scope":           r.Scope,
		"status":          r.Status,
		"reviewer":        r.Reviewer,
		"approvedBy":      r.ApprovedBy,
		"dateDue":         r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetReviewType() string {
	return r.ReviewType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetReviewFrequency() string {
	return r.ReviewFrequency
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetScope() string {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetReviewer() uint64 {
	return r.Reviewer
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetApprovedBy() uint64 {
	return r.ApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewCreate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectReviewCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["reviewType"]; ok && len(val) > 0 {
				r.ReviewType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["reviewFrequency"]; ok && len(val) > 0 {
				r.ReviewFrequency, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scope"]; ok && len(val) > 0 {
				r.Scope, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["reviewer"]; ok && len(val) > 0 {
				r.Reviewer, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["approvedBy"]; ok && len(val) > 0 {
				r.ApprovedBy, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["reviewType"]; ok && len(val) > 0 {
			r.ReviewType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["reviewFrequency"]; ok && len(val) > 0 {
			r.ReviewFrequency, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope"]; ok && len(val) > 0 {
			r.Scope, err = val[0], nil
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

		if val, ok := req.Form["reviewer"]; ok && len(val) > 0 {
			r.Reviewer, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["approvedBy"]; ok && len(val) > 0 {
			r.ApprovedBy, err = payload.ParseUint64(val[0]), nil
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

// NewProjectReviewRead request
func NewProjectReviewRead() *ProjectReviewRead {
	return &ProjectReviewRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reviewID": r.ReviewID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewRead) GetReviewID() uint64 {
	return r.ReviewID
}

// Fill processes request and fills internal variables
func (r *ProjectReviewRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "reviewID")
		r.ReviewID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectReviewUpdate request
func NewProjectReviewUpdate() *ProjectReviewUpdate {
	return &ProjectReviewUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reviewID":        r.ReviewID,
		"revisionID":      r.RevisionID,
		"title":           r.Title,
		"description":     r.Description,
		"reviewType":      r.ReviewType,
		"reviewFrequency": r.ReviewFrequency,
		"scope":           r.Scope,
		"status":          r.Status,
		"reviewer":        r.Reviewer,
		"approvedBy":      r.ApprovedBy,
		"dateDue":         r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetReviewID() uint64 {
	return r.ReviewID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetReviewType() string {
	return r.ReviewType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetReviewFrequency() string {
	return r.ReviewFrequency
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetScope() string {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetReviewer() uint64 {
	return r.Reviewer
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetApprovedBy() uint64 {
	return r.ApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewUpdate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectReviewUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["reviewType"]; ok && len(val) > 0 {
				r.ReviewType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["reviewFrequency"]; ok && len(val) > 0 {
				r.ReviewFrequency, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scope"]; ok && len(val) > 0 {
				r.Scope, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["reviewer"]; ok && len(val) > 0 {
				r.Reviewer, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["approvedBy"]; ok && len(val) > 0 {
				r.ApprovedBy, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["reviewType"]; ok && len(val) > 0 {
			r.ReviewType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["reviewFrequency"]; ok && len(val) > 0 {
			r.ReviewFrequency, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope"]; ok && len(val) > 0 {
			r.Scope, err = val[0], nil
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

		if val, ok := req.Form["reviewer"]; ok && len(val) > 0 {
			r.Reviewer, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["approvedBy"]; ok && len(val) > 0 {
			r.ApprovedBy, err = payload.ParseUint64(val[0]), nil
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

		val = chi.URLParam(req, "reviewID")
		r.ReviewID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectReviewDelete request
func NewProjectReviewDelete() *ProjectReviewDelete {
	return &ProjectReviewDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reviewID": r.ReviewID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReviewDelete) GetReviewID() uint64 {
	return r.ReviewID
}

// Fill processes request and fills internal variables
func (r *ProjectReviewDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "reviewID")
		r.ReviewID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
