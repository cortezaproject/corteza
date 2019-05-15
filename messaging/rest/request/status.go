package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `status.go`, `status.util.go` or `status_test.go` to
	implement your API calls, helper functions and tests. The file `status.go`
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
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Status list request parameters
type StatusList struct {
}

func NewStatusList() *StatusList {
	return &StatusList{}
}

func (r StatusList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

func (sReq *StatusList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(sReq)

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

var _ RequestFiller = NewStatusList()

// Status set request parameters
type StatusSet struct {
	Icon    string
	Message string
	Expires string
}

func NewStatusSet() *StatusSet {
	return &StatusSet{}
}

func (r StatusSet) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["icon"] = r.Icon

	out["message"] = r.Message

	out["expires"] = r.Expires

	return out
}

func (sReq *StatusSet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(sReq)

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

	if val, ok := post["icon"]; ok {

		sReq.Icon = val
	}
	if val, ok := post["message"]; ok {

		sReq.Message = val
	}
	if val, ok := post["expires"]; ok {

		sReq.Expires = val
	}

	return err
}

var _ RequestFiller = NewStatusSet()

// Status delete request parameters
type StatusDelete struct {
}

func NewStatusDelete() *StatusDelete {
	return &StatusDelete{}
}

func (r StatusDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

func (sReq *StatusDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(sReq)

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

var _ RequestFiller = NewStatusDelete()
