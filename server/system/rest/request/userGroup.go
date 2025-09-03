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
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/types"
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
	UserGroupList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// MemberID GET parameter
		//
		// Search roles for member
		MemberID uint64 `json:",string"`

		// UserGroupID GET parameter
		//
		// Search roles by ID
		UserGroupID []id.ID

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted user groups
		Deleted uint

		// Archived GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) archived user groups
		Archived uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

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

	UserGroupCreate struct {
		// Name POST parameter
		//
		// Name of user group
		Name string

		// Handle POST parameter
		//
		// Handle for user group
		Handle string

		// Members POST parameter
		//
		// user group member IDs
		Members []string

		// SelfID POST parameter
		//
		// Parent user group ID
		SelfID uint64 `json:",string"`

		// Meta POST parameter
		//
		// Meta
		Meta *types.UserGroupMeta

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	UserGroupUpdate struct {
		// UserGroupID PATH parameter
		//
		// User Group ID
		UserGroupID id.ID

		// Name POST parameter
		//
		// Name of user group
		Name string

		// Handle POST parameter
		//
		// Handle for user group
		Handle string

		// Members POST parameter
		//
		// user group member IDs
		Members []string

		// SelfID POST parameter
		//
		// Parent user group ID
		SelfID uint64 `json:",string"`

		// Meta POST parameter
		//
		// Meta
		Meta *types.UserGroupMeta

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	UserGroupRead struct {
		// UserGroupID PATH parameter
		//
		// User group ID
		UserGroupID id.ID
	}

	UserGroupDelete struct {
		// UserGroupID PATH parameter
		//
		// User Group ID
		UserGroupID id.ID
	}

	UserGroupUndelete struct {
		// UserGroupID PATH parameter
		//
		// User group ID
		UserGroupID id.ID
	}

	UserGroupMemberList struct {
		// UserGroupID PATH parameter
		//
		// Source User Group ID
		UserGroupID id.ID
	}

	UserGroupMemberAdd struct {
		// UserGroupID PATH parameter
		//
		// Source User Group ID
		UserGroupID id.ID

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}
)

// NewUserGroupList request
func NewUserGroupList() *UserGroupList {
	return &UserGroupList{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":       r.Query,
		"memberID":    r.MemberID,
		"userGroupID": r.UserGroupID,
		"deleted":     r.Deleted,
		"archived":    r.Archived,
		"labels":      r.Labels,
		"limit":       r.Limit,
		"incTotal":    r.IncTotal,
		"pageCursor":  r.PageCursor,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetMemberID() uint64 {
	return r.MemberID
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetUserGroupID() []id.ID {
	return r.UserGroupID
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetArchived() uint {
	return r.Archived
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *UserGroupList) Fill(req *http.Request) (err error) {

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
		if val, ok := tmp["userGroupID[]"]; ok {
			r.UserGroupID, err = id.ParseIDs(val)
			if err != nil {
				return err
			}
		} else if val, ok := tmp["userGroupID"]; ok {
			r.UserGroupID, err = id.ParseIDs(val)
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

// NewUserGroupCreate request
func NewUserGroupCreate() *UserGroupCreate {
	return &UserGroupCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":    r.Name,
		"handle":  r.Handle,
		"members": r.Members,
		"selfID":  r.SelfID,
		"meta":    r.Meta,
		"labels":  r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) GetMembers() []string {
	return r.Members
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) GetSelfID() uint64 {
	return r.SelfID
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) GetMeta() *types.UserGroupMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *UserGroupCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["selfID"]; ok && len(val) > 0 {
				r.SelfID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseUserGroupMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseUserGroupMeta(val)
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

		if val, ok := req.Form["selfID"]; ok && len(val) > 0 {
			r.SelfID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseUserGroupMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseUserGroupMeta(val)
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

// NewUserGroupUpdate request
func NewUserGroupUpdate() *UserGroupUpdate {
	return &UserGroupUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userGroupID": r.UserGroupID,
		"name":        r.Name,
		"handle":      r.Handle,
		"members":     r.Members,
		"selfID":      r.SelfID,
		"meta":        r.Meta,
		"labels":      r.Labels,
		"updatedAt":   r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetUserGroupID() id.ID {
	return r.UserGroupID
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetMembers() []string {
	return r.Members
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetSelfID() uint64 {
	return r.SelfID
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetMeta() *types.UserGroupMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *UserGroupUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["selfID"]; ok && len(val) > 0 {
				r.SelfID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseUserGroupMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseUserGroupMeta(val)
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

		if val, ok := req.Form["selfID"]; ok && len(val) > 0 {
			r.SelfID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseUserGroupMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseUserGroupMeta(val)
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

		val = chi.URLParam(req, "userGroupID")
		r.UserGroupID, err = id.ParseID(val)
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserGroupRead request
func NewUserGroupRead() *UserGroupRead {
	return &UserGroupRead{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userGroupID": r.UserGroupID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupRead) GetUserGroupID() id.ID {
	return r.UserGroupID
}

// Fill processes request and fills internal variables
func (r *UserGroupRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userGroupID")
		r.UserGroupID, err = id.ParseID(val)
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserGroupDelete request
func NewUserGroupDelete() *UserGroupDelete {
	return &UserGroupDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userGroupID": r.UserGroupID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupDelete) GetUserGroupID() id.ID {
	return r.UserGroupID
}

// Fill processes request and fills internal variables
func (r *UserGroupDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userGroupID")
		r.UserGroupID, err = id.ParseID(val)
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserGroupUndelete request
func NewUserGroupUndelete() *UserGroupUndelete {
	return &UserGroupUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userGroupID": r.UserGroupID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupUndelete) GetUserGroupID() id.ID {
	return r.UserGroupID
}

// Fill processes request and fills internal variables
func (r *UserGroupUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userGroupID")
		r.UserGroupID, err = id.ParseID(val)
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserGroupMemberList request
func NewUserGroupMemberList() *UserGroupMemberList {
	return &UserGroupMemberList{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupMemberList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userGroupID": r.UserGroupID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupMemberList) GetUserGroupID() id.ID {
	return r.UserGroupID
}

// Fill processes request and fills internal variables
func (r *UserGroupMemberList) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userGroupID")
		r.UserGroupID, err = id.ParseID(val)
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserGroupMemberAdd request
func NewUserGroupMemberAdd() *UserGroupMemberAdd {
	return &UserGroupMemberAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupMemberAdd) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userGroupID": r.UserGroupID,
		"userID":      r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupMemberAdd) GetUserGroupID() id.ID {
	return r.UserGroupID
}

// Auditable returns all auditable/loggable parameters
func (r UserGroupMemberAdd) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserGroupMemberAdd) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userGroupID")
		r.UserGroupID, err = id.ParseID(val)
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
