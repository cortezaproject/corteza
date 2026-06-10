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
	ProjectGroupList struct {
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

		// ProjectGroupID GET parameter
		//
		// Filter by group IDs
		ProjectGroupID []string

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

	ProjectGroupCreate struct {
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
	}

	ProjectGroupRead struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectGroupID PATH parameter
		//
		// Project Group ID
		ProjectGroupID uint64 `json:",string"`
	}

	ProjectGroupUpdate struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectGroupID PATH parameter
		//
		// Project Group ID
		ProjectGroupID uint64 `json:",string"`

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

		// UpdatedAt POST parameter
		//
		// Last update timestamp
		UpdatedAt *time.Time
	}

	ProjectGroupDelete struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectGroupID PATH parameter
		//
		// Project Group ID
		ProjectGroupID uint64 `json:",string"`
	}

	ProjectGroupEntryAdd struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectGroupID PATH parameter
		//
		// Project Group ID
		ProjectGroupID uint64 `json:",string"`

		// ResourceRef PATH parameter
		//
		// Resource ref (corteza::component:type/ID)
		ResourceRef string
	}

	ProjectGroupEntryRemove struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// ProjectGroupID PATH parameter
		//
		// Project Group ID
		ProjectGroupID uint64 `json:",string"`

		// ResourceRef PATH parameter
		//
		// Resource ref (corteza::component:type/ID)
		ResourceRef string
	}
)

// NewProjectGroupList request
func NewProjectGroupList() *ProjectGroupList {
	return &ProjectGroupList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":      r.ProjectID,
		"query":          r.Query,
		"handle":         r.Handle,
		"projectGroupID": r.ProjectGroupID,
		"deleted":        r.Deleted,
		"limit":          r.Limit,
		"incTotal":       r.IncTotal,
		"pageCursor":     r.PageCursor,
		"sort":           r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetProjectGroupID() []string {
	return r.ProjectGroupID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectGroupList) Fill(req *http.Request) (err error) {

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
		if val, ok := tmp["projectGroupID[]"]; ok {
			r.ProjectGroupID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["projectGroupID"]; ok {
			r.ProjectGroupID, err = val, nil
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

// NewProjectGroupCreate request
func NewProjectGroupCreate() *ProjectGroupCreate {
	return &ProjectGroupCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":   r.ProjectID,
		"handle":      r.Handle,
		"name":        r.Name,
		"description": r.Description,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupCreate) GetDescription() string {
	return r.Description
}

// Fill processes request and fills internal variables
func (r *ProjectGroupCreate) Fill(req *http.Request) (err error) {

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

// NewProjectGroupRead request
func NewProjectGroupRead() *ProjectGroupRead {
	return &ProjectGroupRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":      r.ProjectID,
		"projectGroupID": r.ProjectGroupID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupRead) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupRead) GetProjectGroupID() uint64 {
	return r.ProjectGroupID
}

// Fill processes request and fills internal variables
func (r *ProjectGroupRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectGroupID")
		r.ProjectGroupID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectGroupUpdate request
func NewProjectGroupUpdate() *ProjectGroupUpdate {
	return &ProjectGroupUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":      r.ProjectID,
		"projectGroupID": r.ProjectGroupID,
		"handle":         r.Handle,
		"name":           r.Name,
		"description":    r.Description,
		"updatedAt":      r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) GetProjectGroupID() uint64 {
	return r.ProjectGroupID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *ProjectGroupUpdate) Fill(req *http.Request) (err error) {

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

		val = chi.URLParam(req, "projectGroupID")
		r.ProjectGroupID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectGroupDelete request
func NewProjectGroupDelete() *ProjectGroupDelete {
	return &ProjectGroupDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":      r.ProjectID,
		"projectGroupID": r.ProjectGroupID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupDelete) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupDelete) GetProjectGroupID() uint64 {
	return r.ProjectGroupID
}

// Fill processes request and fills internal variables
func (r *ProjectGroupDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectGroupID")
		r.ProjectGroupID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectGroupEntryAdd request
func NewProjectGroupEntryAdd() *ProjectGroupEntryAdd {
	return &ProjectGroupEntryAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryAdd) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":      r.ProjectID,
		"projectGroupID": r.ProjectGroupID,
		"resourceRef":    r.ResourceRef,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryAdd) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryAdd) GetProjectGroupID() uint64 {
	return r.ProjectGroupID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryAdd) GetResourceRef() string {
	return r.ResourceRef
}

// Fill processes request and fills internal variables
func (r *ProjectGroupEntryAdd) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectGroupID")
		r.ProjectGroupID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "resourceRef")
		r.ResourceRef, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectGroupEntryRemove request
func NewProjectGroupEntryRemove() *ProjectGroupEntryRemove {
	return &ProjectGroupEntryRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":      r.ProjectID,
		"projectGroupID": r.ProjectGroupID,
		"resourceRef":    r.ResourceRef,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryRemove) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryRemove) GetProjectGroupID() uint64 {
	return r.ProjectGroupID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGroupEntryRemove) GetResourceRef() string {
	return r.ResourceRef
}

// Fill processes request and fills internal variables
func (r *ProjectGroupEntryRemove) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "projectGroupID")
		r.ProjectGroupID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "resourceRef")
		r.ResourceRef, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
