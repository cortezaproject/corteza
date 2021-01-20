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
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

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
