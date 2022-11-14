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
	UserList struct {
		// UserID GET parameter
		//
		// Filter by user ID
		UserID []string

		// RoleID GET parameter
		//
		// Filter by role membership
		RoleID []string

		// Query GET parameter
		//
		// Search query to match against users
		Query string

		// Username GET parameter
		//
		// Search username to match against users
		Username string

		// Email GET parameter
		//
		// Search email to match against users
		Email string

		// Handle GET parameter
		//
		// Search handle to match against users
		Handle string

		// Kind GET parameter
		//
		// Kind (normal, bot)
		Kind types.UserKind

		// IncDeleted GET parameter
		//
		// [Deprecated] Include deleted users (requires 'access' permission)
		IncDeleted bool

		// IncSuspended GET parameter
		//
		// [Deprecated] Include suspended users
		IncSuspended bool

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted users
		Deleted uint

		// Suspended GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) suspended users
		Suspended uint

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

	UserCreate struct {
		// Email POST parameter
		//
		// Email
		Email string

		// Name POST parameter
		//
		// Name
		Name string

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Kind POST parameter
		//
		// Kind (normal, bot)
		Kind types.UserKind

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	UserUpdate struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// Email POST parameter
		//
		// Email
		Email string

		// Name POST parameter
		//
		// Name
		Name string

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Kind POST parameter
		//
		// Kind (normal, bot)
		Kind types.UserKind

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	UserPartialUpdate struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserRead struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserDelete struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserSuspend struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserUnsuspend struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserUndelete struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserSetPassword struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`

		// Password POST parameter
		//
		// New password or empty to unset
		Password string
	}

	UserMembershipList struct {
		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserMembershipAdd struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserMembershipRemove struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`

		// UserID PATH parameter
		//
		// User ID
		UserID uint64 `json:",string"`
	}

	UserTriggerScript struct {
		// UserID PATH parameter
		//
		// ID
		UserID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string

		// Args POST parameter
		//
		// Arguments to pass to the script
		Args map[string]interface{}
	}

	UserSessionsRemove struct {
		// UserID PATH parameter
		//
		// ID
		UserID uint64 `json:",string"`
	}

	UserListCredentials struct {
		// UserID PATH parameter
		//
		// ID
		UserID uint64 `json:",string"`
	}

	UserDeleteCredentials struct {
		// UserID PATH parameter
		//
		// ID
		UserID uint64 `json:",string"`

		// CredentialsID PATH parameter
		//
		// Credentials ID
		CredentialsID uint64 `json:",string"`
	}

	UserExport struct {
		// Filename PATH parameter
		//
		// Output filename
		Filename string

		// InclRoleMembership GET parameter
		//
		// Include role membership
		InclRoleMembership bool

		// InclRoles GET parameter
		//
		// Include roles
		InclRoles bool
	}

	UserImport struct {
		// Upload POST parameter
		//
		// File import
		Upload *multipart.FileHeader
	}
)

// NewUserList request
func NewUserList() *UserList {
	return &UserList{}
}

// Auditable returns all auditable/loggable parameters
func (r UserList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID":       r.UserID,
		"roleID":       r.RoleID,
		"query":        r.Query,
		"username":     r.Username,
		"email":        r.Email,
		"handle":       r.Handle,
		"kind":         r.Kind,
		"incDeleted":   r.IncDeleted,
		"incSuspended": r.IncSuspended,
		"deleted":      r.Deleted,
		"suspended":    r.Suspended,
		"labels":       r.Labels,
		"limit":        r.Limit,
		"incTotal":     r.IncTotal,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetUserID() []string {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetRoleID() []string {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetUsername() string {
	return r.Username
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetEmail() string {
	return r.Email
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetKind() types.UserKind {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetIncDeleted() bool {
	return r.IncDeleted
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetIncSuspended() bool {
	return r.IncSuspended
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetSuspended() uint {
	return r.Suspended
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r UserList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *UserList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["userID[]"]; ok {
			r.UserID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["userID"]; ok {
			r.UserID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["roleID[]"]; ok {
			r.RoleID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["roleID"]; ok {
			r.RoleID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["username"]; ok && len(val) > 0 {
			r.Username, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["email"]; ok && len(val) > 0 {
			r.Email, err = val[0], nil
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
		if val, ok := tmp["kind"]; ok && len(val) > 0 {
			r.Kind, err = types.UserKind(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incDeleted"]; ok && len(val) > 0 {
			r.IncDeleted, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incSuspended"]; ok && len(val) > 0 {
			r.IncSuspended, err = payload.ParseBool(val[0]), nil
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
		if val, ok := tmp["suspended"]; ok && len(val) > 0 {
			r.Suspended, err = payload.ParseUint(val[0]), nil
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

// NewUserCreate request
func NewUserCreate() *UserCreate {
	return &UserCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"email":  r.Email,
		"name":   r.Name,
		"handle": r.Handle,
		"kind":   r.Kind,
		"labels": r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) GetEmail() string {
	return r.Email
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) GetKind() types.UserKind {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *UserCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["email"]; ok && len(val) > 0 {
				r.Email, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["kind"]; ok && len(val) > 0 {
				r.Kind, err = types.UserKind(val[0]), nil
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

		if val, ok := req.Form["email"]; ok && len(val) > 0 {
			r.Email, err = val[0], nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = types.UserKind(val[0]), nil
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

// NewUserUpdate request
func NewUserUpdate() *UserUpdate {
	return &UserUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
		"email":  r.Email,
		"name":   r.Name,
		"handle": r.Handle,
		"kind":   r.Kind,
		"labels": r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) GetEmail() string {
	return r.Email
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) GetKind() types.UserKind {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *UserUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["email"]; ok && len(val) > 0 {
				r.Email, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["kind"]; ok && len(val) > 0 {
				r.Kind, err = types.UserKind(val[0]), nil
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

		if val, ok := req.Form["email"]; ok && len(val) > 0 {
			r.Email, err = val[0], nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = types.UserKind(val[0]), nil
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

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserPartialUpdate request
func NewUserPartialUpdate() *UserPartialUpdate {
	return &UserPartialUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserPartialUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserPartialUpdate) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserPartialUpdate) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserRead request
func NewUserRead() *UserRead {
	return &UserRead{}
}

// Auditable returns all auditable/loggable parameters
func (r UserRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserRead) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserDelete request
func NewUserDelete() *UserDelete {
	return &UserDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r UserDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserDelete) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserSuspend request
func NewUserSuspend() *UserSuspend {
	return &UserSuspend{}
}

// Auditable returns all auditable/loggable parameters
func (r UserSuspend) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserSuspend) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserSuspend) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserUnsuspend request
func NewUserUnsuspend() *UserUnsuspend {
	return &UserUnsuspend{}
}

// Auditable returns all auditable/loggable parameters
func (r UserUnsuspend) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserUnsuspend) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserUnsuspend) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserUndelete request
func NewUserUndelete() *UserUndelete {
	return &UserUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r UserUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserUndelete) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserSetPassword request
func NewUserSetPassword() *UserSetPassword {
	return &UserSetPassword{}
}

// Auditable returns all auditable/loggable parameters
func (r UserSetPassword) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserSetPassword) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r UserSetPassword) GetPassword() string {
	return r.Password
}

// Fill processes request and fills internal variables
func (r *UserSetPassword) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["password"]; ok && len(val) > 0 {
				r.Password, err = val[0], nil
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

		if val, ok := req.Form["password"]; ok && len(val) > 0 {
			r.Password, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserMembershipList request
func NewUserMembershipList() *UserMembershipList {
	return &UserMembershipList{}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipList) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserMembershipList) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserMembershipAdd request
func NewUserMembershipAdd() *UserMembershipAdd {
	return &UserMembershipAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipAdd) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipAdd) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipAdd) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserMembershipAdd) Fill(req *http.Request) (err error) {

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

// NewUserMembershipRemove request
func NewUserMembershipRemove() *UserMembershipRemove {
	return &UserMembershipRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipRemove) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipRemove) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserMembershipRemove) Fill(req *http.Request) (err error) {

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

// NewUserTriggerScript request
func NewUserTriggerScript() *UserTriggerScript {
	return &UserTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r UserTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
		"script": r.Script,
		"args":   r.Args,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserTriggerScript) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r UserTriggerScript) GetScript() string {
	return r.Script
}

// Auditable returns all auditable/loggable parameters
func (r UserTriggerScript) GetArgs() map[string]interface{} {
	return r.Args
}

// Fill processes request and fills internal variables
func (r *UserTriggerScript) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["script"]; ok && len(val) > 0 {
				r.Script, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["args[]"]; ok {
				r.Args, err = parseMapStringInterface(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["args"]; ok {
				r.Args, err = parseMapStringInterface(val)
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

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserSessionsRemove request
func NewUserSessionsRemove() *UserSessionsRemove {
	return &UserSessionsRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r UserSessionsRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserSessionsRemove) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserSessionsRemove) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserListCredentials request
func NewUserListCredentials() *UserListCredentials {
	return &UserListCredentials{}
}

// Auditable returns all auditable/loggable parameters
func (r UserListCredentials) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserListCredentials) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *UserListCredentials) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserDeleteCredentials request
func NewUserDeleteCredentials() *UserDeleteCredentials {
	return &UserDeleteCredentials{}
}

// Auditable returns all auditable/loggable parameters
func (r UserDeleteCredentials) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID":        r.UserID,
		"credentialsID": r.CredentialsID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserDeleteCredentials) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r UserDeleteCredentials) GetCredentialsID() uint64 {
	return r.CredentialsID
}

// Fill processes request and fills internal variables
func (r *UserDeleteCredentials) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "userID")
		r.UserID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "credentialsID")
		r.CredentialsID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserExport request
func NewUserExport() *UserExport {
	return &UserExport{}
}

// Auditable returns all auditable/loggable parameters
func (r UserExport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"filename":           r.Filename,
		"inclRoleMembership": r.InclRoleMembership,
		"inclRoles":          r.InclRoles,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserExport) GetFilename() string {
	return r.Filename
}

// Auditable returns all auditable/loggable parameters
func (r UserExport) GetInclRoleMembership() bool {
	return r.InclRoleMembership
}

// Auditable returns all auditable/loggable parameters
func (r UserExport) GetInclRoles() bool {
	return r.InclRoles
}

// Fill processes request and fills internal variables
func (r *UserExport) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["inclRoleMembership"]; ok && len(val) > 0 {
			r.InclRoleMembership, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["inclRoles"]; ok && len(val) > 0 {
			r.InclRoles, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "filename")
		r.Filename, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewUserImport request
func NewUserImport() *UserImport {
	return &UserImport{}
}

// Auditable returns all auditable/loggable parameters
func (r UserImport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"upload": r.Upload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r UserImport) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *UserImport) Fill(req *http.Request) (err error) {

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

			// Ignoring upload as its handled in the POST params section
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

	}

	return err
}
