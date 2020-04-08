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

// StatusList request parameters
type StatusList struct {
}

// NewStatusList request
func NewStatusList() *StatusList {
	return &StatusList{}
}

// Auditable returns all auditable/loggable parameters
func (r StatusList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *StatusList) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewStatusList()

// StatusSet request parameters
type StatusSet struct {
	hasIcon bool
	rawIcon string
	Icon    string

	hasMessage bool
	rawMessage string
	Message    string

	hasExpires bool
	rawExpires string
	Expires    string
}

// NewStatusSet request
func NewStatusSet() *StatusSet {
	return &StatusSet{}
}

// Auditable returns all auditable/loggable parameters
func (r StatusSet) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["icon"] = r.Icon
	out["message"] = r.Message
	out["expires"] = r.Expires

	return out
}

// Fill processes request and fills internal variables
func (r *StatusSet) Fill(req *http.Request) (err error) {
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

	if val, ok := post["icon"]; ok {
		r.hasIcon = true
		r.rawIcon = val
		r.Icon = val
	}
	if val, ok := post["message"]; ok {
		r.hasMessage = true
		r.rawMessage = val
		r.Message = val
	}
	if val, ok := post["expires"]; ok {
		r.hasExpires = true
		r.rawExpires = val
		r.Expires = val
	}

	return err
}

var _ RequestFiller = NewStatusSet()

// StatusDelete request parameters
type StatusDelete struct {
}

// NewStatusDelete request
func NewStatusDelete() *StatusDelete {
	return &StatusDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r StatusDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *StatusDelete) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewStatusDelete()

// HasIcon returns true if icon was set
func (r *StatusSet) HasIcon() bool {
	return r.hasIcon
}

// RawIcon returns raw value of icon parameter
func (r *StatusSet) RawIcon() string {
	return r.rawIcon
}

// GetIcon returns casted value of  icon parameter
func (r *StatusSet) GetIcon() string {
	return r.Icon
}

// HasMessage returns true if message was set
func (r *StatusSet) HasMessage() bool {
	return r.hasMessage
}

// RawMessage returns raw value of message parameter
func (r *StatusSet) RawMessage() string {
	return r.rawMessage
}

// GetMessage returns casted value of  message parameter
func (r *StatusSet) GetMessage() string {
	return r.Message
}

// HasExpires returns true if expires was set
func (r *StatusSet) HasExpires() bool {
	return r.hasExpires
}

// RawExpires returns raw value of expires parameter
func (r *StatusSet) RawExpires() string {
	return r.rawExpires
}

// GetExpires returns casted value of  expires parameter
func (r *StatusSet) GetExpires() string {
	return r.Expires
}
