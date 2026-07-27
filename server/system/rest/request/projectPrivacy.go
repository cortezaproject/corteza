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
	ProjectPrivacyList struct {
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

	ProjectPrivacyCreate struct {
		// ProjectID POST parameter
		//
		// Project this item belongs to
		ProjectID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// RequestType POST parameter
		//
		// Request type
		RequestType string

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

		// RequestOwner POST parameter
		//
		// Request owner user ID
		RequestOwner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// ChangeApprovedBy POST parameter
		//
		// Change approved-by user ID
		ChangeApprovedBy uint64 `json:",string"`

		// RiskAssessment POST parameter
		//
		// Risk assessment
		RiskAssessment string

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
	}

	ProjectPrivacyRead struct {
		// PrivacyID PATH parameter
		//
		// Privacy item ID
		PrivacyID uint64 `json:",string"`
	}

	ProjectPrivacyUpdate struct {
		// PrivacyID PATH parameter
		//
		// Privacy item ID
		PrivacyID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// RequestType POST parameter
		//
		// Request type
		RequestType string

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

		// RequestOwner POST parameter
		//
		// Request owner user ID
		RequestOwner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// ChangeApprovedBy POST parameter
		//
		// Change approved-by user ID
		ChangeApprovedBy uint64 `json:",string"`

		// RiskAssessment POST parameter
		//
		// Risk assessment
		RiskAssessment string

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
	}

	ProjectPrivacyDelete struct {
		// PrivacyID PATH parameter
		//
		// Privacy item ID
		PrivacyID uint64 `json:",string"`
	}
)

// NewProjectPrivacyList request
func NewProjectPrivacyList() *ProjectPrivacyList {
	return &ProjectPrivacyList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) Auditable() map[string]interface{} {
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
func (r ProjectPrivacyList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectPrivacyList) Fill(req *http.Request) (err error) {

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

// NewProjectPrivacyCreate request
func NewProjectPrivacyCreate() *ProjectPrivacyCreate {
	return &ProjectPrivacyCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":        r.ProjectID,
		"title":            r.Title,
		"description":      r.Description,
		"requestType":      r.RequestType,
		"status":           r.Status,
		"severity":         r.Severity,
		"risk":             r.Risk,
		"requestOwner":     r.RequestOwner,
		"changeOwner":      r.ChangeOwner,
		"changeApprovedBy": r.ChangeApprovedBy,
		"riskAssessment":   r.RiskAssessment,
		"changeRequired":   r.ChangeRequired,
		"riskChange":       r.RiskChange,
		"dateDue":          r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetRequestType() string {
	return r.RequestType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetRequestOwner() uint64 {
	return r.RequestOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetChangeApprovedBy() uint64 {
	return r.ChangeApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetRiskAssessment() string {
	return r.RiskAssessment
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetChangeRequired() string {
	return r.ChangeRequired
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetRiskChange() string {
	return r.RiskChange
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyCreate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectPrivacyCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["requestType"]; ok && len(val) > 0 {
				r.RequestType, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["requestOwner"]; ok && len(val) > 0 {
				r.RequestOwner, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["riskAssessment"]; ok && len(val) > 0 {
				r.RiskAssessment, err = val[0], nil
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

		if val, ok := req.Form["requestType"]; ok && len(val) > 0 {
			r.RequestType, err = val[0], nil
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

		if val, ok := req.Form["requestOwner"]; ok && len(val) > 0 {
			r.RequestOwner, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["riskAssessment"]; ok && len(val) > 0 {
			r.RiskAssessment, err = val[0], nil
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
	}

	return err
}

// NewProjectPrivacyRead request
func NewProjectPrivacyRead() *ProjectPrivacyRead {
	return &ProjectPrivacyRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"privacyID": r.PrivacyID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyRead) GetPrivacyID() uint64 {
	return r.PrivacyID
}

// Fill processes request and fills internal variables
func (r *ProjectPrivacyRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "privacyID")
		r.PrivacyID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectPrivacyUpdate request
func NewProjectPrivacyUpdate() *ProjectPrivacyUpdate {
	return &ProjectPrivacyUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"privacyID":        r.PrivacyID,
		"title":            r.Title,
		"description":      r.Description,
		"requestType":      r.RequestType,
		"status":           r.Status,
		"severity":         r.Severity,
		"risk":             r.Risk,
		"requestOwner":     r.RequestOwner,
		"changeOwner":      r.ChangeOwner,
		"changeApprovedBy": r.ChangeApprovedBy,
		"riskAssessment":   r.RiskAssessment,
		"changeRequired":   r.ChangeRequired,
		"riskChange":       r.RiskChange,
		"dateDue":          r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetPrivacyID() uint64 {
	return r.PrivacyID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetRequestType() string {
	return r.RequestType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetRequestOwner() uint64 {
	return r.RequestOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetChangeApprovedBy() uint64 {
	return r.ChangeApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetRiskAssessment() string {
	return r.RiskAssessment
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetChangeRequired() string {
	return r.ChangeRequired
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetRiskChange() string {
	return r.RiskChange
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyUpdate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectPrivacyUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["requestType"]; ok && len(val) > 0 {
				r.RequestType, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["requestOwner"]; ok && len(val) > 0 {
				r.RequestOwner, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["riskAssessment"]; ok && len(val) > 0 {
				r.RiskAssessment, err = val[0], nil
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
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

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

		if val, ok := req.Form["requestType"]; ok && len(val) > 0 {
			r.RequestType, err = val[0], nil
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

		if val, ok := req.Form["requestOwner"]; ok && len(val) > 0 {
			r.RequestOwner, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["riskAssessment"]; ok && len(val) > 0 {
			r.RiskAssessment, err = val[0], nil
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
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "privacyID")
		r.PrivacyID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectPrivacyDelete request
func NewProjectPrivacyDelete() *ProjectPrivacyDelete {
	return &ProjectPrivacyDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"privacyID": r.PrivacyID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPrivacyDelete) GetPrivacyID() uint64 {
	return r.PrivacyID
}

// Fill processes request and fills internal variables
func (r *ProjectPrivacyDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "privacyID")
		r.PrivacyID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
