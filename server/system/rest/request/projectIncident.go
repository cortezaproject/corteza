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
	ProjectIncidentList struct {
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

	ProjectIncidentCreate struct {
		// ProjectID POST parameter
		//
		// Project this incident belongs to
		ProjectID uint64 `json:",string"`

		// RevisionID POST parameter
		//
		// Revision to assign this incident to
		RevisionID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// IncidentType POST parameter
		//
		// Incident type
		IncidentType string

		// GroupSystem POST parameter
		//
		// Affected group/system
		GroupSystem string

		// Status POST parameter
		//
		// Status
		Status string

		// Severity POST parameter
		//
		// Severity
		Severity string

		// Risk POST parameter
		//
		// Risk
		Risk string

		// IssueOwner POST parameter
		//
		// Issue owner user ID
		IssueOwner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// ChangeApprovedBy POST parameter
		//
		// Change approved-by user ID
		ChangeApprovedBy uint64 `json:",string"`

		// RiskIssue POST parameter
		//
		// Risk/issue assessment
		RiskIssue string

		// ChangeRequired POST parameter
		//
		// Change required
		ChangeRequired string

		// RiskChange POST parameter
		//
		// Risk of change
		RiskChange string

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string

		// CompletedDate POST parameter
		//
		// Completed date (ISO)
		CompletedDate string
	}

	ProjectIncidentRead struct {
		// IncidentID PATH parameter
		//
		// Incident ID
		IncidentID uint64 `json:",string"`
	}

	ProjectIncidentUpdate struct {
		// IncidentID PATH parameter
		//
		// Incident ID
		IncidentID uint64 `json:",string"`

		// RevisionID POST parameter
		//
		// Revision to assign this incident to
		RevisionID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// IncidentType POST parameter
		//
		// Incident type
		IncidentType string

		// GroupSystem POST parameter
		//
		// Affected group/system
		GroupSystem string

		// Status POST parameter
		//
		// Status
		Status string

		// Severity POST parameter
		//
		// Severity
		Severity string

		// Risk POST parameter
		//
		// Risk
		Risk string

		// IssueOwner POST parameter
		//
		// Issue owner user ID
		IssueOwner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// ChangeApprovedBy POST parameter
		//
		// Change approved-by user ID
		ChangeApprovedBy uint64 `json:",string"`

		// RiskIssue POST parameter
		//
		// Risk/issue assessment
		RiskIssue string

		// ChangeRequired POST parameter
		//
		// Change required
		ChangeRequired string

		// RiskChange POST parameter
		//
		// Risk of change
		RiskChange string

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string

		// CompletedDate POST parameter
		//
		// Completed date (ISO)
		CompletedDate string
	}

	ProjectIncidentDelete struct {
		// IncidentID PATH parameter
		//
		// Incident ID
		IncidentID uint64 `json:",string"`
	}
)

// NewProjectIncidentList request
func NewProjectIncidentList() *ProjectIncidentList {
	return &ProjectIncidentList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) Auditable() map[string]interface{} {
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
func (r ProjectIncidentList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectIncidentList) Fill(req *http.Request) (err error) {

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

// NewProjectIncidentCreate request
func NewProjectIncidentCreate() *ProjectIncidentCreate {
	return &ProjectIncidentCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":        r.ProjectID,
		"revisionID":       r.RevisionID,
		"title":            r.Title,
		"description":      r.Description,
		"incidentType":     r.IncidentType,
		"groupSystem":      r.GroupSystem,
		"status":           r.Status,
		"severity":         r.Severity,
		"risk":             r.Risk,
		"issueOwner":       r.IssueOwner,
		"changeOwner":      r.ChangeOwner,
		"changeApprovedBy": r.ChangeApprovedBy,
		"riskIssue":        r.RiskIssue,
		"changeRequired":   r.ChangeRequired,
		"riskChange":       r.RiskChange,
		"dateDue":          r.DateDue,
		"completedDate":    r.CompletedDate,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetIncidentType() string {
	return r.IncidentType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetGroupSystem() string {
	return r.GroupSystem
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetIssueOwner() uint64 {
	return r.IssueOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetChangeApprovedBy() uint64 {
	return r.ChangeApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetRiskIssue() string {
	return r.RiskIssue
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetChangeRequired() string {
	return r.ChangeRequired
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetRiskChange() string {
	return r.RiskChange
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetDateDue() string {
	return r.DateDue
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentCreate) GetCompletedDate() string {
	return r.CompletedDate
}

// Fill processes request and fills internal variables
func (r *ProjectIncidentCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["incidentType"]; ok && len(val) > 0 {
				r.IncidentType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["groupSystem"]; ok && len(val) > 0 {
				r.GroupSystem, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["severity"]; ok && len(val) > 0 {
				r.Severity, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["risk"]; ok && len(val) > 0 {
				r.Risk, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["issueOwner"]; ok && len(val) > 0 {
				r.IssueOwner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeOwner"]; ok && len(val) > 0 {
				r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeApprovedBy"]; ok && len(val) > 0 {
				r.ChangeApprovedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["riskIssue"]; ok && len(val) > 0 {
				r.RiskIssue, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeRequired"]; ok && len(val) > 0 {
				r.ChangeRequired, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["riskChange"]; ok && len(val) > 0 {
				r.RiskChange, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["completedDate"]; ok && len(val) > 0 {
				r.CompletedDate, err = val[0], nil
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

		if val, ok := req.Form["incidentType"]; ok && len(val) > 0 {
			r.IncidentType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["groupSystem"]; ok && len(val) > 0 {
			r.GroupSystem, err = val[0], nil
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

		if val, ok := req.Form["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["risk"]; ok && len(val) > 0 {
			r.Risk, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["issueOwner"]; ok && len(val) > 0 {
			r.IssueOwner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeOwner"]; ok && len(val) > 0 {
			r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeApprovedBy"]; ok && len(val) > 0 {
			r.ChangeApprovedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["riskIssue"]; ok && len(val) > 0 {
			r.RiskIssue, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeRequired"]; ok && len(val) > 0 {
			r.ChangeRequired, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["riskChange"]; ok && len(val) > 0 {
			r.RiskChange, err = val[0], nil
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

		if val, ok := req.Form["completedDate"]; ok && len(val) > 0 {
			r.CompletedDate, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewProjectIncidentRead request
func NewProjectIncidentRead() *ProjectIncidentRead {
	return &ProjectIncidentRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"incidentID": r.IncidentID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentRead) GetIncidentID() uint64 {
	return r.IncidentID
}

// Fill processes request and fills internal variables
func (r *ProjectIncidentRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "incidentID")
		r.IncidentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectIncidentUpdate request
func NewProjectIncidentUpdate() *ProjectIncidentUpdate {
	return &ProjectIncidentUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"incidentID":       r.IncidentID,
		"revisionID":       r.RevisionID,
		"title":            r.Title,
		"description":      r.Description,
		"incidentType":     r.IncidentType,
		"groupSystem":      r.GroupSystem,
		"status":           r.Status,
		"severity":         r.Severity,
		"risk":             r.Risk,
		"issueOwner":       r.IssueOwner,
		"changeOwner":      r.ChangeOwner,
		"changeApprovedBy": r.ChangeApprovedBy,
		"riskIssue":        r.RiskIssue,
		"changeRequired":   r.ChangeRequired,
		"riskChange":       r.RiskChange,
		"dateDue":          r.DateDue,
		"completedDate":    r.CompletedDate,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetIncidentID() uint64 {
	return r.IncidentID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetIncidentType() string {
	return r.IncidentType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetGroupSystem() string {
	return r.GroupSystem
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetIssueOwner() uint64 {
	return r.IssueOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetChangeApprovedBy() uint64 {
	return r.ChangeApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetRiskIssue() string {
	return r.RiskIssue
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetChangeRequired() string {
	return r.ChangeRequired
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetRiskChange() string {
	return r.RiskChange
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetDateDue() string {
	return r.DateDue
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentUpdate) GetCompletedDate() string {
	return r.CompletedDate
}

// Fill processes request and fills internal variables
func (r *ProjectIncidentUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["incidentType"]; ok && len(val) > 0 {
				r.IncidentType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["groupSystem"]; ok && len(val) > 0 {
				r.GroupSystem, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["severity"]; ok && len(val) > 0 {
				r.Severity, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["risk"]; ok && len(val) > 0 {
				r.Risk, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["issueOwner"]; ok && len(val) > 0 {
				r.IssueOwner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeOwner"]; ok && len(val) > 0 {
				r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeApprovedBy"]; ok && len(val) > 0 {
				r.ChangeApprovedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["riskIssue"]; ok && len(val) > 0 {
				r.RiskIssue, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeRequired"]; ok && len(val) > 0 {
				r.ChangeRequired, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["riskChange"]; ok && len(val) > 0 {
				r.RiskChange, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["completedDate"]; ok && len(val) > 0 {
				r.CompletedDate, err = val[0], nil
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

		if val, ok := req.Form["incidentType"]; ok && len(val) > 0 {
			r.IncidentType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["groupSystem"]; ok && len(val) > 0 {
			r.GroupSystem, err = val[0], nil
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

		if val, ok := req.Form["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["risk"]; ok && len(val) > 0 {
			r.Risk, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["issueOwner"]; ok && len(val) > 0 {
			r.IssueOwner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeOwner"]; ok && len(val) > 0 {
			r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeApprovedBy"]; ok && len(val) > 0 {
			r.ChangeApprovedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["riskIssue"]; ok && len(val) > 0 {
			r.RiskIssue, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeRequired"]; ok && len(val) > 0 {
			r.ChangeRequired, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["riskChange"]; ok && len(val) > 0 {
			r.RiskChange, err = val[0], nil
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

		if val, ok := req.Form["completedDate"]; ok && len(val) > 0 {
			r.CompletedDate, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "incidentID")
		r.IncidentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectIncidentDelete request
func NewProjectIncidentDelete() *ProjectIncidentDelete {
	return &ProjectIncidentDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"incidentID": r.IncidentID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectIncidentDelete) GetIncidentID() uint64 {
	return r.IncidentID
}

// Fill processes request and fills internal variables
func (r *ProjectIncidentDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "incidentID")
		r.IncidentID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
