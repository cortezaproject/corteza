package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `permissions.go`, `permissions.util.go` or `permissions_test.go` to
	implement your API calls, helper functions and tests. The file `permissions.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// PermissionsList request parameters
type PermissionsList struct {
}

// NewPermissionsList request
func NewPermissionsList() *PermissionsList {
	return &PermissionsList{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *PermissionsList) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	return err
}

var _ RequestFiller = NewPermissionsList()

// PermissionsEffective request parameters
type PermissionsEffective struct {
	hasResource bool
	rawResource string
	Resource    string
}

// NewPermissionsEffective request
func NewPermissionsEffective() *PermissionsEffective {
	return &PermissionsEffective{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsEffective) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource

	return out
}

// Fill processes request and fills internal variables
func (r *PermissionsEffective) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["resource"]; ok {
		r.hasResource = true
		r.rawResource = val
		r.Resource = val
	}

	return err
}

var _ RequestFiller = NewPermissionsEffective()

// PermissionsRead request parameters
type PermissionsRead struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewPermissionsRead request
func NewPermissionsRead() *PermissionsRead {
	return &PermissionsRead{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *PermissionsRead) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionsRead()

// PermissionsDelete request parameters
type PermissionsDelete struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewPermissionsDelete request
func NewPermissionsDelete() *PermissionsDelete {
	return &PermissionsDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *PermissionsDelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionsDelete()

// PermissionsUpdate request parameters
type PermissionsUpdate struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasRules bool
	rawRules string
	Rules    permissions.RuleSet
}

// NewPermissionsUpdate request
func NewPermissionsUpdate() *PermissionsUpdate {
	return &PermissionsUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r PermissionsUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["rules"] = r.Rules

	return out
}

// Fill processes request and fills internal variables
func (r *PermissionsUpdate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionsUpdate()

// HasResource returns true if resource was set
func (r *PermissionsEffective) HasResource() bool {
	return r.hasResource
}

// RawResource returns raw value of resource parameter
func (r *PermissionsEffective) RawResource() string {
	return r.rawResource
}

// GetResource returns casted value of  resource parameter
func (r *PermissionsEffective) GetResource() string {
	return r.Resource
}

// HasRoleID returns true if roleID was set
func (r *PermissionsRead) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *PermissionsRead) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *PermissionsRead) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *PermissionsDelete) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *PermissionsDelete) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *PermissionsDelete) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *PermissionsUpdate) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *PermissionsUpdate) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *PermissionsUpdate) GetRoleID() uint64 {
	return r.RoleID
}

// HasRules returns true if rules was set
func (r *PermissionsUpdate) HasRules() bool {
	return r.hasRules
}

// RawRules returns raw value of rules parameter
func (r *PermissionsUpdate) RawRules() string {
	return r.rawRules
}

// GetRules returns casted value of  rules parameter
func (r *PermissionsUpdate) GetRules() permissions.RuleSet {
	return r.Rules
}
