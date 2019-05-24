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

	"github.com/cortezaproject/corteza-server/internal/permissions"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Permissions list request parameters
type PermissionsList struct {
}

func NewPermissionsList() *PermissionsList {
	return &PermissionsList{}
}

func (r PermissionsList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

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

// Permissions effective request parameters
type PermissionsEffective struct {
	Resource string
}

func NewPermissionsEffective() *PermissionsEffective {
	return &PermissionsEffective{}
}

func (r PermissionsEffective) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource

	return out
}

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
		r.Resource = val
	}

	return err
}

var _ RequestFiller = NewPermissionsEffective()

// Permissions read request parameters
type PermissionsRead struct {
	RoleID uint64 `json:",string"`
}

func NewPermissionsRead() *PermissionsRead {
	return &PermissionsRead{}
}

func (r PermissionsRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionsRead()

// Permissions delete request parameters
type PermissionsDelete struct {
	RoleID uint64 `json:",string"`
}

func NewPermissionsDelete() *PermissionsDelete {
	return &PermissionsDelete{}
}

func (r PermissionsDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionsDelete()

// Permissions update request parameters
type PermissionsUpdate struct {
	RoleID uint64 `json:",string"`
	Rules  permissions.RuleSet
}

func NewPermissionsUpdate() *PermissionsUpdate {
	return &PermissionsUpdate{}
}

func (r PermissionsUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["rules"] = r.Rules

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionsUpdate()
