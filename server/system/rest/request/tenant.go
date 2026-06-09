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
	TenantList struct {
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
		// Exclude (0, default), include (1) or return only (2) deleted tenants
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

	TenantCreate struct {
		// Handle POST parameter
		//
		// Tenant handle
		Handle string

		// Status POST parameter
		//
		// Tenant status
		Status string

		// Config POST parameter
		//
		// Tenant config
		Config types.TenantConfig

		// Meta POST parameter
		//
		// Tenant meta
		Meta types.TenantMeta

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	TenantRead struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantUpdate struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Tenant handle
		Handle string

		// Status POST parameter
		//
		// Tenant status
		Status string

		// Config POST parameter
		//
		// Tenant config
		Config types.TenantConfig

		// Meta POST parameter
		//
		// Tenant meta
		Meta types.TenantMeta

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	TenantDelete struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantUndelete struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantSuspend struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantActivate struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantArchive struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantListMembers struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`
	}

	TenantAddMember struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`

		// UserID POST parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// Role POST parameter
		//
		// Member role
		Role string
	}

	TenantUpdateMember struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// Role POST parameter
		//
		// Member role
		Role string

		// Status POST parameter
		//
		// Member status
		Status string
	}

	TenantRemoveMember struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	TenantSuspendMember struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	TenantActivateMember struct {
		// TenantID PATH parameter
		//
		// Tenant ID
		TenantID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}
)

// NewTenantList request
func NewTenantList() *TenantList {
	return &TenantList{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) Auditable() map[string]interface{} {
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
func (r TenantList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r TenantList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *TenantList) Fill(req *http.Request) (err error) {

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

// NewTenantCreate request
func NewTenantCreate() *TenantCreate {
	return &TenantCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle": r.Handle,
		"status": r.Status,
		"config": r.Config,
		"meta":   r.Meta,
		"labels": r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TenantCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r TenantCreate) GetConfig() types.TenantConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r TenantCreate) GetMeta() types.TenantMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r TenantCreate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *TenantCreate) Fill(req *http.Request) (err error) {

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
				r.Config, err = types.ParseTenantConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseTenantConfig(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseTenantMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseTenantMeta(val)
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
			r.Config, err = types.ParseTenantConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseTenantConfig(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTenantMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTenantMeta(val)
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

// NewTenantRead request
func NewTenantRead() *TenantRead {
	return &TenantRead{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantRead) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantUpdate request
func NewTenantUpdate() *TenantUpdate {
	return &TenantUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID":  r.TenantID,
		"handle":    r.Handle,
		"status":    r.Status,
		"config":    r.Config,
		"meta":      r.Meta,
		"updatedAt": r.UpdatedAt,
		"labels":    r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetTenantID() uint64 {
	return r.TenantID
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetConfig() types.TenantConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetMeta() types.TenantMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *TenantUpdate) Fill(req *http.Request) (err error) {

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
				r.Config, err = types.ParseTenantConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseTenantConfig(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseTenantMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseTenantMeta(val)
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
			r.Config, err = types.ParseTenantConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseTenantConfig(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTenantMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTenantMeta(val)
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

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantDelete request
func NewTenantDelete() *TenantDelete {
	return &TenantDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantDelete) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantUndelete request
func NewTenantUndelete() *TenantUndelete {
	return &TenantUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantUndelete) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantSuspend request
func NewTenantSuspend() *TenantSuspend {
	return &TenantSuspend{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantSuspend) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantSuspend) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantSuspend) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantActivate request
func NewTenantActivate() *TenantActivate {
	return &TenantActivate{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantActivate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantActivate) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantActivate) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantArchive request
func NewTenantArchive() *TenantArchive {
	return &TenantArchive{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantArchive) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantArchive) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantArchive) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantListMembers request
func NewTenantListMembers() *TenantListMembers {
	return &TenantListMembers{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantListMembers) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantListMembers) GetTenantID() uint64 {
	return r.TenantID
}

// Fill processes request and fills internal variables
func (r *TenantListMembers) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantAddMember request
func NewTenantAddMember() *TenantAddMember {
	return &TenantAddMember{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantAddMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
		"userID":   r.UserID,
		"role":     r.Role,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantAddMember) GetTenantID() uint64 {
	return r.TenantID
}

// Auditable returns all auditable/loggable parameters
func (r TenantAddMember) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r TenantAddMember) GetRole() string {
	return r.Role
}

// Fill processes request and fills internal variables
func (r *TenantAddMember) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["role"]; ok && len(val) > 0 {
				r.Role, err = val[0], nil
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

		if val, ok := req.Form["role"]; ok && len(val) > 0 {
			r.Role, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTenantUpdateMember request
func NewTenantUpdateMember() *TenantUpdateMember {
	return &TenantUpdateMember{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdateMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
		"userID":   r.UserID,
		"role":     r.Role,
		"status":   r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdateMember) GetTenantID() uint64 {
	return r.TenantID
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdateMember) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdateMember) GetRole() string {
	return r.Role
}

// Auditable returns all auditable/loggable parameters
func (r TenantUpdateMember) GetStatus() string {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *TenantUpdateMember) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["role"]; ok && len(val) > 0 {
				r.Role, err = val[0], nil
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
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["role"]; ok && len(val) > 0 {
			r.Role, err = val[0], nil
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
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
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

// NewTenantRemoveMember request
func NewTenantRemoveMember() *TenantRemoveMember {
	return &TenantRemoveMember{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantRemoveMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
		"userID":   r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantRemoveMember) GetTenantID() uint64 {
	return r.TenantID
}

// Auditable returns all auditable/loggable parameters
func (r TenantRemoveMember) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *TenantRemoveMember) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
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

// NewTenantSuspendMember request
func NewTenantSuspendMember() *TenantSuspendMember {
	return &TenantSuspendMember{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantSuspendMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
		"userID":   r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantSuspendMember) GetTenantID() uint64 {
	return r.TenantID
}

// Auditable returns all auditable/loggable parameters
func (r TenantSuspendMember) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *TenantSuspendMember) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
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

// NewTenantActivateMember request
func NewTenantActivateMember() *TenantActivateMember {
	return &TenantActivateMember{}
}

// Auditable returns all auditable/loggable parameters
func (r TenantActivateMember) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"tenantID": r.TenantID,
		"userID":   r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TenantActivateMember) GetTenantID() uint64 {
	return r.TenantID
}

// Auditable returns all auditable/loggable parameters
func (r TenantActivateMember) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *TenantActivateMember) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "tenantID")
		r.TenantID, err = payload.ParseUint64(val), nil
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
