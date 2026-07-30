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
	ProjectAiSystemList struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query
		Query string

		// Handle GET parameter
		//
		// Handle filter
		Handle string

		// ProjectAiSystemID GET parameter
		//
		// Filter by AI system IDs
		ProjectAiSystemID []string

		// RiskClass GET parameter
		//
		// EU AI Act risk class filter (prohibited|high|limited|minimal)
		RiskClass string

		// Deleted GET parameter
		//
		// Deleted filter (0=exclude 1=only 2=include)
		Deleted uint

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total count
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort
		Sort string
	}

	ProjectAiSystemCreate struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Name POST parameter
		//
		// Name
		Name string

		// Description POST parameter
		//
		// Description
		Description string

		// IntendedPurpose POST parameter
		//
		// Intended purpose of the AI system
		IntendedPurpose string

		// RiskClass POST parameter
		//
		// EU AI Act risk class (prohibited|high|limited|minimal)
		RiskClass string
	}

	ProjectAiSystemRead struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectAiSystemID PATH parameter
		//
		// AI System ID
		ProjectAiSystemID uint64 `json:",string"`
	}

	ProjectAiSystemUpdate struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectAiSystemID PATH parameter
		//
		// AI System ID
		ProjectAiSystemID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Name POST parameter
		//
		// Name
		Name string

		// Description POST parameter
		//
		// Description
		Description string

		// IntendedPurpose POST parameter
		//
		// Intended purpose of the AI system
		IntendedPurpose string

		// RiskClass POST parameter
		//
		// EU AI Act risk class (prohibited|high|limited|minimal)
		RiskClass string

		// UpdatedAt POST parameter
		//
		// Last update timestamp
		UpdatedAt *time.Time
	}

	ProjectAiSystemDelete struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectAiSystemID PATH parameter
		//
		// AI System ID
		ProjectAiSystemID uint64 `json:",string"`
	}

	ProjectAiSystemEntryAdd struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectAiSystemID PATH parameter
		//
		// AI System ID
		ProjectAiSystemID uint64 `json:",string"`

		// ResourceRef POST parameter
		//
		// Resource ref (corteza::component:type/ID)
		ResourceRef string
	}

	ProjectAiSystemEntryRemove struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectAiSystemID PATH parameter
		//
		// AI System ID
		ProjectAiSystemID uint64 `json:",string"`

		// ResourceRef GET parameter
		//
		// Resource ref (corteza::component:type/ID)
		ResourceRef string
	}
)

// NewProjectAiSystemList request
func NewProjectAiSystemList() *ProjectAiSystemList {
	return &ProjectAiSystemList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":         r.ProjectID,
		"query":             r.Query,
		"handle":            r.Handle,
		"projectAiSystemID": r.ProjectAiSystemID,
		"riskClass":         r.RiskClass,
		"deleted":           r.Deleted,
		"limit":             r.Limit,
		"incTotal":          r.IncTotal,
		"pageCursor":        r.PageCursor,
		"sort":              r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetProjectAiSystemID() []string {
	return r.ProjectAiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetRiskClass() string {
	return r.RiskClass
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemList) Fill(req *http.Request) (err error) {

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
		if val, ok := tmp["projectAiSystemID[]"]; ok {
			r.ProjectAiSystemID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["projectAiSystemID"]; ok {
			r.ProjectAiSystemID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["riskClass"]; ok && len(val) > 0 {
			r.RiskClass, err = val[0], nil
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

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectAiSystemCreate request
func NewProjectAiSystemCreate() *ProjectAiSystemCreate {
	return &ProjectAiSystemCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":       r.ProjectID,
		"handle":          r.Handle,
		"name":            r.Name,
		"description":     r.Description,
		"intendedPurpose": r.IntendedPurpose,
		"riskClass":       r.RiskClass,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) GetIntendedPurpose() string {
	return r.IntendedPurpose
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemCreate) GetRiskClass() string {
	return r.RiskClass
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["intendedPurpose"]; ok && len(val) > 0 {
				r.IntendedPurpose, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["riskClass"]; ok && len(val) > 0 {
				r.RiskClass, err = val[0], nil
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
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

		if val, ok := req.Form["intendedPurpose"]; ok && len(val) > 0 {
			r.IntendedPurpose, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["riskClass"]; ok && len(val) > 0 {
			r.RiskClass, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectAiSystemRead request
func NewProjectAiSystemRead() *ProjectAiSystemRead {
	return &ProjectAiSystemRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":         r.ProjectID,
		"projectAiSystemID": r.ProjectAiSystemID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemRead) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemRead) GetProjectAiSystemID() uint64 {
	return r.ProjectAiSystemID
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectAiSystemID")
		r.ProjectAiSystemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectAiSystemUpdate request
func NewProjectAiSystemUpdate() *ProjectAiSystemUpdate {
	return &ProjectAiSystemUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":         r.ProjectID,
		"projectAiSystemID": r.ProjectAiSystemID,
		"handle":            r.Handle,
		"name":              r.Name,
		"description":       r.Description,
		"intendedPurpose":   r.IntendedPurpose,
		"riskClass":         r.RiskClass,
		"updatedAt":         r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetProjectAiSystemID() uint64 {
	return r.ProjectAiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetIntendedPurpose() string {
	return r.IntendedPurpose
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetRiskClass() string {
	return r.RiskClass
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["intendedPurpose"]; ok && len(val) > 0 {
				r.IntendedPurpose, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["riskClass"]; ok && len(val) > 0 {
				r.RiskClass, err = val[0], nil
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
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

		if val, ok := req.Form["intendedPurpose"]; ok && len(val) > 0 {
			r.IntendedPurpose, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["riskClass"]; ok && len(val) > 0 {
			r.RiskClass, err = val[0], nil
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

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectAiSystemID")
		r.ProjectAiSystemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectAiSystemDelete request
func NewProjectAiSystemDelete() *ProjectAiSystemDelete {
	return &ProjectAiSystemDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":         r.ProjectID,
		"projectAiSystemID": r.ProjectAiSystemID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemDelete) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemDelete) GetProjectAiSystemID() uint64 {
	return r.ProjectAiSystemID
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectAiSystemID")
		r.ProjectAiSystemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectAiSystemEntryAdd request
func NewProjectAiSystemEntryAdd() *ProjectAiSystemEntryAdd {
	return &ProjectAiSystemEntryAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryAdd) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":         r.ProjectID,
		"projectAiSystemID": r.ProjectAiSystemID,
		"resourceRef":       r.ResourceRef,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryAdd) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryAdd) GetProjectAiSystemID() uint64 {
	return r.ProjectAiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryAdd) GetResourceRef() string {
	return r.ResourceRef
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemEntryAdd) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["resourceRef"]; ok && len(val) > 0 {
				r.ResourceRef, err = val[0], nil
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

		if val, ok := req.Form["resourceRef"]; ok && len(val) > 0 {
			r.ResourceRef, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectAiSystemID")
		r.ProjectAiSystemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectAiSystemEntryRemove request
func NewProjectAiSystemEntryRemove() *ProjectAiSystemEntryRemove {
	return &ProjectAiSystemEntryRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":         r.ProjectID,
		"projectAiSystemID": r.ProjectAiSystemID,
		"resourceRef":       r.ResourceRef,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryRemove) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryRemove) GetProjectAiSystemID() uint64 {
	return r.ProjectAiSystemID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAiSystemEntryRemove) GetResourceRef() string {
	return r.ResourceRef
}

// Fill processes request and fills internal variables
func (r *ProjectAiSystemEntryRemove) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["resourceRef"]; ok && len(val) > 0 {
			r.ResourceRef, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectAiSystemID")
		r.ProjectAiSystemID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
