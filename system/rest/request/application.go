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

func (apReq *ApplicationList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(apReq)

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

func (apReq *ApplicationCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(apReq)

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

	if val, ok := post["name"]; ok {

		apReq.Name = val
	}
	if val, ok := post["enabled"]; ok {

		apReq.Enabled = parseBool(val)
	}
	if val, ok := post["unify"]; ok {

		if apReq.Unify, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["config"]; ok {

		if apReq.Config, err = parseJSONTextWithErr(val); err != nil {
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

func (apReq *ApplicationUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(apReq)

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

	apReq.ApplicationID = parseUInt64(chi.URLParam(r, "applicationID"))
	if val, ok := post["name"]; ok {

		apReq.Name = val
	}
	if val, ok := post["enabled"]; ok {

		apReq.Enabled = parseBool(val)
	}
	if val, ok := post["unify"]; ok {

		if apReq.Unify, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["config"]; ok {

		if apReq.Config, err = parseJSONTextWithErr(val); err != nil {
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

func (apReq *ApplicationRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(apReq)

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

	apReq.ApplicationID = parseUInt64(chi.URLParam(r, "applicationID"))

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

func (apReq *ApplicationDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(apReq)

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

	apReq.ApplicationID = parseUInt64(chi.URLParam(r, "applicationID"))

	return err
}

var _ RequestFiller = NewApplicationDelete()
