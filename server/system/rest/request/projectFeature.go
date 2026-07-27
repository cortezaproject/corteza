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
	ProjectFeatureList struct {
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

	ProjectFeatureCreate struct {
		// ProjectID POST parameter
		//
		// Project this feature belongs to
		ProjectID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// FeatureType POST parameter
		//
		// Feature type
		FeatureType string

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

		// FeatureOwner POST parameter
		//
		// Feature owner user ID
		FeatureOwner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// ChangeApprovedBy POST parameter
		//
		// Change approved-by user ID
		ChangeApprovedBy uint64 `json:",string"`

		// RiskFeature POST parameter
		//
		// Risk of feature
		RiskFeature string

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

	ProjectFeatureRead struct {
		// FeatureID PATH parameter
		//
		// Feature ID
		FeatureID uint64 `json:",string"`
	}

	ProjectFeatureUpdate struct {
		// FeatureID PATH parameter
		//
		// Feature ID
		FeatureID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// FeatureType POST parameter
		//
		// Feature type
		FeatureType string

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

		// FeatureOwner POST parameter
		//
		// Feature owner user ID
		FeatureOwner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// ChangeApprovedBy POST parameter
		//
		// Change approved-by user ID
		ChangeApprovedBy uint64 `json:",string"`

		// RiskFeature POST parameter
		//
		// Risk of feature
		RiskFeature string

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

	ProjectFeatureDelete struct {
		// FeatureID PATH parameter
		//
		// Feature ID
		FeatureID uint64 `json:",string"`
	}
)

// NewProjectFeatureList request
func NewProjectFeatureList() *ProjectFeatureList {
	return &ProjectFeatureList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) Auditable() map[string]interface{} {
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
func (r ProjectFeatureList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectFeatureList) Fill(req *http.Request) (err error) {

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

// NewProjectFeatureCreate request
func NewProjectFeatureCreate() *ProjectFeatureCreate {
	return &ProjectFeatureCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":        r.ProjectID,
		"title":            r.Title,
		"description":      r.Description,
		"featureType":      r.FeatureType,
		"status":           r.Status,
		"severity":         r.Severity,
		"risk":             r.Risk,
		"featureOwner":     r.FeatureOwner,
		"changeOwner":      r.ChangeOwner,
		"changeApprovedBy": r.ChangeApprovedBy,
		"riskFeature":      r.RiskFeature,
		"changeRequired":   r.ChangeRequired,
		"riskChange":       r.RiskChange,
		"dateDue":          r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetFeatureType() string {
	return r.FeatureType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetFeatureOwner() uint64 {
	return r.FeatureOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetChangeApprovedBy() uint64 {
	return r.ChangeApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetRiskFeature() string {
	return r.RiskFeature
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetChangeRequired() string {
	return r.ChangeRequired
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetRiskChange() string {
	return r.RiskChange
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureCreate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectFeatureCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["featureType"]; ok && len(val) > 0 {
				r.FeatureType, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["featureOwner"]; ok && len(val) > 0 {
				r.FeatureOwner, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["riskFeature"]; ok && len(val) > 0 {
				r.RiskFeature, err = val[0], nil
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

		if val, ok := req.Form["featureType"]; ok && len(val) > 0 {
			r.FeatureType, err = val[0], nil
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

		if val, ok := req.Form["featureOwner"]; ok && len(val) > 0 {
			r.FeatureOwner, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["riskFeature"]; ok && len(val) > 0 {
			r.RiskFeature, err = val[0], nil
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

// NewProjectFeatureRead request
func NewProjectFeatureRead() *ProjectFeatureRead {
	return &ProjectFeatureRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"featureID": r.FeatureID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureRead) GetFeatureID() uint64 {
	return r.FeatureID
}

// Fill processes request and fills internal variables
func (r *ProjectFeatureRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "featureID")
		r.FeatureID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectFeatureUpdate request
func NewProjectFeatureUpdate() *ProjectFeatureUpdate {
	return &ProjectFeatureUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"featureID":        r.FeatureID,
		"title":            r.Title,
		"description":      r.Description,
		"featureType":      r.FeatureType,
		"status":           r.Status,
		"severity":         r.Severity,
		"risk":             r.Risk,
		"featureOwner":     r.FeatureOwner,
		"changeOwner":      r.ChangeOwner,
		"changeApprovedBy": r.ChangeApprovedBy,
		"riskFeature":      r.RiskFeature,
		"changeRequired":   r.ChangeRequired,
		"riskChange":       r.RiskChange,
		"dateDue":          r.DateDue,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetFeatureID() uint64 {
	return r.FeatureID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetFeatureType() string {
	return r.FeatureType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetFeatureOwner() uint64 {
	return r.FeatureOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetChangeApprovedBy() uint64 {
	return r.ChangeApprovedBy
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetRiskFeature() string {
	return r.RiskFeature
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetChangeRequired() string {
	return r.ChangeRequired
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetRiskChange() string {
	return r.RiskChange
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureUpdate) GetDateDue() string {
	return r.DateDue
}

// Fill processes request and fills internal variables
func (r *ProjectFeatureUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["featureType"]; ok && len(val) > 0 {
				r.FeatureType, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["featureOwner"]; ok && len(val) > 0 {
				r.FeatureOwner, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["riskFeature"]; ok && len(val) > 0 {
				r.RiskFeature, err = val[0], nil
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

		if val, ok := req.Form["featureType"]; ok && len(val) > 0 {
			r.FeatureType, err = val[0], nil
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

		if val, ok := req.Form["featureOwner"]; ok && len(val) > 0 {
			r.FeatureOwner, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["riskFeature"]; ok && len(val) > 0 {
			r.RiskFeature, err = val[0], nil
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

		val = chi.URLParam(req, "featureID")
		r.FeatureID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectFeatureDelete request
func NewProjectFeatureDelete() *ProjectFeatureDelete {
	return &ProjectFeatureDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"featureID": r.FeatureID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectFeatureDelete) GetFeatureID() uint64 {
	return r.FeatureID
}

// Fill processes request and fills internal variables
func (r *ProjectFeatureDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "featureID")
		r.FeatureID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
