package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `application.go`, `application.util.go` or `application_test.go` to
	implement your API calls, helper functions and tests. The file `application.go`
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

	sqlxTypes "github.com/jmoiron/sqlx/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Application list request parameters
type ApplicationList struct {
}

func NewApplicationList() *ApplicationList {
	return &ApplicationList{}
}

func (r ApplicationList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

func (r *ApplicationList) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewApplicationList()

// Application create request parameters
type ApplicationCreate struct {
	Name    string
	Enabled bool
	Unify   sqlxTypes.JSONText
	Config  sqlxTypes.JSONText
}

func NewApplicationCreate() *ApplicationCreate {
	return &ApplicationCreate{}
}

func (r ApplicationCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["enabled"] = r.Enabled
	out["unify"] = r.Unify
	out["config"] = r.Config

	return out
}

func (r *ApplicationCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	if val, ok := post["unify"]; ok {

		if r.Unify, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["config"]; ok {

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewApplicationCreate()

// Application update request parameters
type ApplicationUpdate struct {
	ApplicationID uint64 `json:",string"`
	Name          string
	Enabled       bool
	Unify         sqlxTypes.JSONText
	Config        sqlxTypes.JSONText
}

func NewApplicationUpdate() *ApplicationUpdate {
	return &ApplicationUpdate{}
}

func (r ApplicationUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID
	out["name"] = r.Name
	out["enabled"] = r.Enabled
	out["unify"] = r.Unify
	out["config"] = r.Config

	return out
}

func (r *ApplicationUpdate) Fill(req *http.Request) (err error) {
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

	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["enabled"]; ok {
		r.Enabled = parseBool(val)
	}
	if val, ok := post["unify"]; ok {

		if r.Unify, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["config"]; ok {

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewApplicationUpdate()

// Application read request parameters
type ApplicationRead struct {
	ApplicationID uint64 `json:",string"`
}

func NewApplicationRead() *ApplicationRead {
	return &ApplicationRead{}
}

func (r ApplicationRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID

	return out
}

func (r *ApplicationRead) Fill(req *http.Request) (err error) {
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

	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))

	return err
}

var _ RequestFiller = NewApplicationRead()

// Application delete request parameters
type ApplicationDelete struct {
	ApplicationID uint64 `json:",string"`
}

func NewApplicationDelete() *ApplicationDelete {
	return &ApplicationDelete{}
}

func (r ApplicationDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID

	return out
}

func (r *ApplicationDelete) Fill(req *http.Request) (err error) {
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

	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))

	return err
}

var _ RequestFiller = NewApplicationDelete()
