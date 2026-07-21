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
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
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
	ProjectList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// Handle GET parameter
		//
		// Search by handle
		Handle string

		// Status GET parameter
		//
		// Filter by status
		Status string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted projects
		Deleted uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

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

	ProjectCreate struct {
		// Handle POST parameter
		//
		// Project handle
		Handle string

		// Status POST parameter
		//
		// Project status
		Status string

		// Config POST parameter
		//
		// Project config
		Config types.ProjectConfig

		// Meta POST parameter
		//
		// Project meta
		Meta types.ProjectMeta

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	ProjectRead struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectUpdate struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Project handle
		Handle string

		// Status POST parameter
		//
		// Project status
		Status string

		// Config POST parameter
		//
		// Project config
		Config types.ProjectConfig

		// Meta POST parameter
		//
		// Project meta
		Meta types.ProjectMeta

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	ProjectDelete struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectUndelete struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectListMembers struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectAddMember struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// UserID POST parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// RolePreset POST parameter
		//
		// Role preset
		RolePreset string
	}

	ProjectUpdateMember struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// RolePreset POST parameter
		//
		// Role preset
		RolePreset string
	}

	ProjectRemoveMember struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	ProjectGraph struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectGovernanceSave struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// StepKey PATH parameter
		//
		// Governance step key
		StepKey string

		// Values POST parameter
		//
		// Step form values
		Values map[string]interface{}
	}

	ProjectGovernanceTransition struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// StepKey PATH parameter
		//
		// Governance step key
		StepKey string

		// Action POST parameter
		//
		// Transition action: submit, approve, request-changes, reopen or recall
		Action string

		// Note POST parameter
		//
		// Review note (request-changes/reopen)
		Note string
	}

	ProjectCreateRevision struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectListRevisions struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectGetDeploymentPlan struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`
	}

	ProjectPublish struct {
		// ProjectID PATH parameter
		//
		// Project ID
		ProjectID uint64 `json:",string"`

		// Confirm POST parameter
		//
		// Confirm publish (required for dangerous changes)
		Confirm bool

		// Mappings POST parameter
		//
		// Per-module record migration mappings
		Mappings []types.ModuleMapping
	}
)

// NewProjectList request
func NewProjectList() *ProjectList {
	return &ProjectList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"handle":     r.Handle,
		"status":     r.Status,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectList) Fill(req *http.Request) (err error) {

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
		if val, ok := tmp["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
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
		if val, ok := tmp["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := tmp["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
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

// NewProjectCreate request
func NewProjectCreate() *ProjectCreate {
	return &ProjectCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle": r.Handle,
		"status": r.Status,
		"config": r.Config,
		"meta":   r.Meta,
		"labels": r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreate) GetConfig() types.ProjectConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreate) GetMeta() types.ProjectMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ProjectCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseProjectConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseProjectConfig(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseProjectMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseProjectMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
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

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseProjectConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseProjectConfig(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseProjectMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseProjectMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewProjectRead request
func NewProjectRead() *ProjectRead {
	return &ProjectRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectRead) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectRead) Fill(req *http.Request) (err error) {

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

// NewProjectUpdate request
func NewProjectUpdate() *ProjectUpdate {
	return &ProjectUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
		"handle":    r.Handle,
		"status":    r.Status,
		"config":    r.Config,
		"meta":      r.Meta,
		"updatedAt": r.UpdatedAt,
		"labels":    r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetConfig() types.ProjectConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetMeta() types.ProjectMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ProjectUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseProjectConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseProjectConfig(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseProjectMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseProjectMeta(val)
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

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
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

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseProjectConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseProjectConfig(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseProjectMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseProjectMeta(val)
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

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
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

// NewProjectDelete request
func NewProjectDelete() *ProjectDelete {
	return &ProjectDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectDelete) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectDelete) Fill(req *http.Request) (err error) {

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

// NewProjectUndelete request
func NewProjectUndelete() *ProjectUndelete {
	return &ProjectUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUndelete) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectUndelete) Fill(req *http.Request) (err error) {

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

// NewProjectListMembers request
func NewProjectListMembers() *ProjectListMembers {
	return &ProjectListMembers{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectListMembers) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectListMembers) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectListMembers) Fill(req *http.Request) (err error) {

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

// NewProjectAddMember request
func NewProjectAddMember() *ProjectAddMember {
	return &ProjectAddMember{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAddMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":  r.ProjectID,
		"userID":     r.UserID,
		"rolePreset": r.RolePreset,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAddMember) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAddMember) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectAddMember) GetRolePreset() string {
	return r.RolePreset
}

// Fill processes request and fills internal variables
func (r *ProjectAddMember) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["userID"]; ok && len(val) > 0 {
				r.UserID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["rolePreset"]; ok && len(val) > 0 {
				r.RolePreset, err = val[0], nil
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

		if val, ok := req.Form["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["rolePreset"]; ok && len(val) > 0 {
			r.RolePreset, err = val[0], nil
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

// NewProjectUpdateMember request
func NewProjectUpdateMember() *ProjectUpdateMember {
	return &ProjectUpdateMember{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdateMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":  r.ProjectID,
		"userID":     r.UserID,
		"rolePreset": r.RolePreset,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdateMember) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdateMember) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectUpdateMember) GetRolePreset() string {
	return r.RolePreset
}

// Fill processes request and fills internal variables
func (r *ProjectUpdateMember) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["rolePreset"]; ok && len(val) > 0 {
				r.RolePreset, err = val[0], nil
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

		if val, ok := req.Form["rolePreset"]; ok && len(val) > 0 {
			r.RolePreset, err = val[0], nil
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

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectRemoveMember request
func NewProjectRemoveMember() *ProjectRemoveMember {
	return &ProjectRemoveMember{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectRemoveMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
		"userID":    r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectRemoveMember) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectRemoveMember) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *ProjectRemoveMember) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "projectID")
		r.ProjectID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectGraph request
func NewProjectGraph() *ProjectGraph {
	return &ProjectGraph{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGraph) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGraph) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectGraph) Fill(req *http.Request) (err error) {

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

// NewProjectGovernanceSave request
func NewProjectGovernanceSave() *ProjectGovernanceSave {
	return &ProjectGovernanceSave{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceSave) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
		"stepKey":   r.StepKey,
		"values":    r.Values,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceSave) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceSave) GetStepKey() string {
	return r.StepKey
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceSave) GetValues() map[string]interface{} {
	return r.Values
}

// Fill processes request and fills internal variables
func (r *ProjectGovernanceSave) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["values[]"]; ok {
				r.Values, err = parseMapStringInterface(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["values"]; ok {
				r.Values, err = parseMapStringInterface(val)
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

		if val, ok := req.Form["values[]"]; ok {
			r.Values, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["values"]; ok {
			r.Values, err = parseMapStringInterface(val)
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

		val = chi.URLParam(req, "stepKey")
		r.StepKey, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectGovernanceTransition request
func NewProjectGovernanceTransition() *ProjectGovernanceTransition {
	return &ProjectGovernanceTransition{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceTransition) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
		"stepKey":   r.StepKey,
		"action":    r.Action,
		"note":      r.Note,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceTransition) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceTransition) GetStepKey() string {
	return r.StepKey
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceTransition) GetAction() string {
	return r.Action
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGovernanceTransition) GetNote() string {
	return r.Note
}

// Fill processes request and fills internal variables
func (r *ProjectGovernanceTransition) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["action"]; ok && len(val) > 0 {
				r.Action, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["note"]; ok && len(val) > 0 {
				r.Note, err = val[0], nil
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

		if val, ok := req.Form["action"]; ok && len(val) > 0 {
			r.Action, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["note"]; ok && len(val) > 0 {
			r.Note, err = val[0], nil
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

		val = chi.URLParam(req, "stepKey")
		r.StepKey, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectCreateRevision request
func NewProjectCreateRevision() *ProjectCreateRevision {
	return &ProjectCreateRevision{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreateRevision) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectCreateRevision) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectCreateRevision) Fill(req *http.Request) (err error) {

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

// NewProjectListRevisions request
func NewProjectListRevisions() *ProjectListRevisions {
	return &ProjectListRevisions{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectListRevisions) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectListRevisions) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectListRevisions) Fill(req *http.Request) (err error) {

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

// NewProjectGetDeploymentPlan request
func NewProjectGetDeploymentPlan() *ProjectGetDeploymentPlan {
	return &ProjectGetDeploymentPlan{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGetDeploymentPlan) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectGetDeploymentPlan) GetProjectID() uint64 {
	return r.ProjectID
}

// Fill processes request and fills internal variables
func (r *ProjectGetDeploymentPlan) Fill(req *http.Request) (err error) {

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

// NewProjectPublish request
func NewProjectPublish() *ProjectPublish {
	return &ProjectPublish{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPublish) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
		"confirm":   r.Confirm,
		"mappings":  r.Mappings,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPublish) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPublish) GetConfirm() bool {
	return r.Confirm
}

// Auditable returns all auditable/loggable parameters
func (r ProjectPublish) GetMappings() []types.ModuleMapping {
	return r.Mappings
}

// Fill processes request and fills internal variables
func (r *ProjectPublish) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["confirm"]; ok && len(val) > 0 {
				r.Confirm, err = payload.ParseBool(val[0]), nil
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

		if val, ok := req.Form["confirm"]; ok && len(val) > 0 {
			r.Confirm, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["mappings[]"]; ok && len(val) > 0  {
		//    r.Mappings, err = []types.ModuleMapping(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
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
