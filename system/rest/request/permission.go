package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `permission.go`, `permission.util.go` or `permission_test.go` to
	implement your API calls, helper functions and tests. The file `permission.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rules"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Permission read request parameters
type PermissionRead struct {
	RoleID uint64 `json:",string"`
}

func NewPermissionRead() *PermissionRead {
	return &PermissionRead{}
}

func (pe *PermissionRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pe)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	pe.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionRead()

// Permission delete request parameters
type PermissionDelete struct {
	RoleID uint64 `json:",string"`
}

func NewPermissionDelete() *PermissionDelete {
	return &PermissionDelete{}
}

func (pe *PermissionDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pe)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	pe.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionDelete()

// Permission update request parameters
type PermissionUpdate struct {
	RoleID      uint64 `json:",string"`
	Permissions []rules.Rule
}

func NewPermissionUpdate() *PermissionUpdate {
	return &PermissionUpdate{}
}

func (pe *PermissionUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pe)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	pe.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

	return err
}

var _ RequestFiller = NewPermissionUpdate()
