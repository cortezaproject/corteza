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
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
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
	RoleList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// MemberID GET parameter
		//
		// Search roles for member
		MemberID uint64 `json:",string"`

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted roles
		Deleted uint

		// Archived GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) archived roles
		Archived uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	RoleCreate struct {
		// Name POST parameter
		//
		// Name of role
		Name string

		// Handle POST parameter
		//
		// Handle for role
		Handle string

		// Members POST parameter
		//
		// role member IDs
		Members []string

		// Meta POST parameter
		//
		// Meta
		Meta *types.RoleMeta

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	RoleUpdate struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`

		// Name POST parameter
		//
		// Name of role
		Name string

		// Handle POST parameter
		//
		// Handle for role
		Handle string

		// Members POST parameter
		//
		// role member IDs
		Members []string

		// Meta POST parameter
		//
		// Meta
		Meta *types.RoleMeta

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	RoleRead struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	RoleDelete struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	RoleArchive struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	RoleUnarchive struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	RoleUndelete struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	RoleMove struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`

		// OrganisationID POST parameter
		//
		// Role ID
		OrganisationID uint64 `json:",string"`
	}

	RoleMerge struct {
		// RoleID PATH parameter
		//
		// Source Role ID
		RoleID uint64 `json:",string"`

		// Destination POST parameter
		//
		// Destination Role ID
		Destination uint64 `json:",string"`
	}

	RoleMemberList struct {
		// RoleID PATH parameter
		//
		// Source Role ID
		RoleID uint64 `json:",string"`
	}

	RoleMemberAdd struct {
		// RoleID PATH parameter
		//
		// Source Role ID
		RoleID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	RoleMemberRemove struct {
		// RoleID PATH parameter
		//
		// Source Role ID
		RoleID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	RoleTriggerScript struct {
		// RoleID PATH parameter
		//
		// ID
		RoleID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string

		// Args POST parameter
		//
		// Arguments to pass to the script
		Args map[string]interface{}
	}
)

// NewRoleList request
func NewRoleList() *RoleList {
	return &RoleList{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"memberID":   r.MemberID,
		"deleted":    r.Deleted,
		"archived":   r.Archived,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetMemberID() uint64 {
	return r.MemberID
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetArchived() uint {
	return r.Archived
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *RoleList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["memberID"]; ok && len(val) > 0 {
			r.MemberID, err = payload.ParseUint64(val[0]), nil
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
		if val, ok := tmp["archived"]; ok && len(val) > 0 {
			r.Archived, err = payload.ParseUint(val[0]), nil
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

// NewRoleCreate request
func NewRoleCreate() *RoleCreate {
	return &RoleCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":    r.Name,
		"handle":  r.Handle,
		"members": r.Members,
		"meta":    r.Meta,
		"labels":  r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) GetMembers() []string {
	return r.Members
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) GetMeta() *types.RoleMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *RoleCreate) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseRoleMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseRoleMeta(val)
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["members[]"]; ok && len(val) > 0  {
		//    r.Members, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseRoleMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseRoleMeta(val)
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

// NewRoleUpdate request
func NewRoleUpdate() *RoleUpdate {
	return &RoleUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID":  r.RoleID,
		"name":    r.Name,
		"handle":  r.Handle,
		"members": r.Members,
		"meta":    r.Meta,
		"labels":  r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) GetMembers() []string {
	return r.Members
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) GetMeta() *types.RoleMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *RoleUpdate) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseRoleMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseRoleMeta(val)
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["members[]"]; ok && len(val) > 0  {
		//    r.Members, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseRoleMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseRoleMeta(val)
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

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleRead request
func NewRoleRead() *RoleRead {
	return &RoleRead{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleRead) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *RoleRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleDelete request
func NewRoleDelete() *RoleDelete {
	return &RoleDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleDelete) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *RoleDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleArchive request
func NewRoleArchive() *RoleArchive {
	return &RoleArchive{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleArchive) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleArchive) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *RoleArchive) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleUnarchive request
func NewRoleUnarchive() *RoleUnarchive {
	return &RoleUnarchive{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUnarchive) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUnarchive) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *RoleUnarchive) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleUndelete request
func NewRoleUndelete() *RoleUndelete {
	return &RoleUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUndelete) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *RoleUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleMove request
func NewRoleMove() *RoleMove {
	return &RoleMove{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID":         r.RoleID,
		"organisationID": r.OrganisationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMove) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r RoleMove) GetOrganisationID() uint64 {
	return r.OrganisationID
}

// Fill processes request and fills internal variables
func (r *RoleMove) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["organisationID"]; ok && len(val) > 0 {
				r.OrganisationID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["organisationID"]; ok && len(val) > 0 {
			r.OrganisationID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleMerge request
func NewRoleMerge() *RoleMerge {
	return &RoleMerge{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMerge) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID":      r.RoleID,
		"destination": r.Destination,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMerge) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r RoleMerge) GetDestination() uint64 {
	return r.Destination
}

// Fill processes request and fills internal variables
func (r *RoleMerge) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["destination"]; ok && len(val) > 0 {
				r.Destination, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["destination"]; ok && len(val) > 0 {
			r.Destination, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleMemberList request
func NewRoleMemberList() *RoleMemberList {
	return &RoleMemberList{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberList) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *RoleMemberList) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRoleMemberAdd request
func NewRoleMemberAdd() *RoleMemberAdd {
	return &RoleMemberAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberAdd) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberAdd) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberAdd) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *RoleMemberAdd) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
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

// NewRoleMemberRemove request
func NewRoleMemberRemove() *RoleMemberRemove {
	return &RoleMemberRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberRemove) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberRemove) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *RoleMemberRemove) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
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

// NewRoleTriggerScript request
func NewRoleTriggerScript() *RoleTriggerScript {
	return &RoleTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
		"script": r.Script,
		"args":   r.Args,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RoleTriggerScript) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r RoleTriggerScript) GetScript() string {
	return r.Script
}

// Auditable returns all auditable/loggable parameters
func (r RoleTriggerScript) GetArgs() map[string]interface{} {
	return r.Args
}

// Fill processes request and fills internal variables
func (r *RoleTriggerScript) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["script"]; ok && len(val) > 0 {
				r.Script, err = val[0], nil
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

		if val, ok := req.Form["script"]; ok && len(val) > 0 {
			r.Script, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["args[]"]; ok {
			r.Args, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["args"]; ok {
			r.Args, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "roleID")
		r.RoleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
