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
	ProjectReportReport struct {
		// Resource GET parameter
		//
		// Category to report on (incident, task, feature, privacy, review)
		Resource string

		// ProjectID GET parameter
		//
		// Project scope
		ProjectID uint64 `json:",string"`

		// RevisionID GET parameter
		//
		// Scope the report to a single revision (default aggregates the whole chain)
		RevisionID uint64 `json:",string"`

		// Dimensions GET parameter
		//
		// Group rows by one or more dimensions
		Dimensions []string

		// Metrics GET parameter
		//
		// Metrics to compute per group
		Metrics []string

		// From GET parameter
		//
		// Created-at window start
		From *time.Time

		// To GET parameter
		//
		// Created-at window end
		To *time.Time
	}
)

// NewProjectReportReport request
func NewProjectReportReport() *ProjectReportReport {
	return &ProjectReportReport{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource":   r.Resource,
		"projectID":  r.ProjectID,
		"revisionID": r.RevisionID,
		"dimensions": r.Dimensions,
		"metrics":    r.Metrics,
		"from":       r.From,
		"to":         r.To,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetRevisionID() uint64 {
	return r.RevisionID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetDimensions() []string {
	return r.Dimensions
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetMetrics() []string {
	return r.Metrics
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r ProjectReportReport) GetTo() *time.Time {
	return r.To
}

// Fill processes request and fills internal variables
func (r *ProjectReportReport) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
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
		if val, ok := tmp["dimensions[]"]; ok {
			r.Dimensions, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["dimensions"]; ok {
			r.Dimensions, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["metrics[]"]; ok {
			r.Metrics, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["metrics"]; ok {
			r.Metrics, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["from"]; ok && len(val) > 0 {
			r.From, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["to"]; ok && len(val) > 0 {
			r.To, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	return err
}
