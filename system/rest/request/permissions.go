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
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
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
	PermissionsList struct {
	}

	PermissionsEffective struct {
		// Resource GET parameter
		//
		// Show only rules for a specific resource
		Resource string
	}

	PermissionsEvaluate struct {
		// Resource GET parameter
		//
		// Show only rules for a specific resource
		Resource []string

		// UserID GET parameter
		//
		//
		UserID uint64 `json:",string"`

		// RoleID GET parameter
		//
		//
		RoleID []uint64
	}

	PermissionsRead struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	PermissionsDelete struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`
	}

	PermissionsUpdate struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`

		// Rules POST parameter
		//
		// List of permission rules to set
		Rules rbac.RuleSet
	}

	PermissionsClone struct {
		// RoleID PATH parameter
		//
		// Role ID
		RoleID uint64 `json:",string"`

		// CloneToRoleID GET parameter
		//
		// Clone set of rules to roleID
		CloneToRoleID []string
	}
)

// NewPermissionsList request
func NewPermissionsList() *PermissionsList {
	return &PermissionsList{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsList) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *PermissionsList) Fill(req *http.Request) (err error) {

	return err
}

// NewPermissionsEffective request
func NewPermissionsEffective() *PermissionsEffective {
	return &PermissionsEffective{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEffective) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource": r.Resource,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEffective) GetResource() string {
	return r.Resource
}

// Fill processes request and fills internal variables
func (r *PermissionsEffective) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewPermissionsEvaluate request
func NewPermissionsEvaluate() *PermissionsEvaluate {
	return &PermissionsEvaluate{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEvaluate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource": r.Resource,
		"userID":   r.UserID,
		"roleID":   r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEvaluate) GetResource() []string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEvaluate) GetUserID() uint64 {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEvaluate) GetRoleID() []uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *PermissionsEvaluate) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["resource[]"]; ok {
			r.Resource, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["resource"]; ok {
			r.Resource, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["roleID[]"]; ok {
			r.RoleID, err = payload.ParseUint64s(val), nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["roleID"]; ok {
			r.RoleID, err = payload.ParseUint64s(val), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewPermissionsRead request
func NewPermissionsRead() *PermissionsRead {
	return &PermissionsRead{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsRead) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *PermissionsRead) Fill(req *http.Request) (err error) {

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

// NewPermissionsDelete request
func NewPermissionsDelete() *PermissionsDelete {
	return &PermissionsDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsDelete) GetRoleID() uint64 {
	return r.RoleID
}

// Fill processes request and fills internal variables
func (r *PermissionsDelete) Fill(req *http.Request) (err error) {

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

// NewPermissionsUpdate request
func NewPermissionsUpdate() *PermissionsUpdate {
	return &PermissionsUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID": r.RoleID,
		"rules":  r.Rules,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsUpdate) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsUpdate) GetRules() rbac.RuleSet {
	return r.Rules
}

// Fill processes request and fills internal variables
func (r *PermissionsUpdate) Fill(req *http.Request) (err error) {

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

		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["rules[]"]; ok && len(val) > 0  {
		//    r.Rules, err = rbac.RuleSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
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

// NewPermissionsClone request
func NewPermissionsClone() *PermissionsClone {
	return &PermissionsClone{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsClone) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"roleID":        r.RoleID,
		"cloneToRoleID": r.CloneToRoleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsClone) GetRoleID() uint64 {
	return r.RoleID
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsClone) GetCloneToRoleID() []string {
	return r.CloneToRoleID
}

// Fill processes request and fills internal variables
func (r *PermissionsClone) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["cloneToRoleID[]"]; ok {
			r.CloneToRoleID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["cloneToRoleID"]; ok {
			r.CloneToRoleID, err = val, nil
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
