package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `auth.go`, `auth.util.go` or `auth_test.go` to
	implement your API calls, helper functions and tests. The file `auth.go`
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

// AuthSettings request parameters
type AuthSettings struct {
}

// NewAuthSettings request
func NewAuthSettings() *AuthSettings {
	return &AuthSettings{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthSettings) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *AuthSettings) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewAuthSettings()

// AuthCheck request parameters
type AuthCheck struct {
}

// NewAuthCheck request
func NewAuthCheck() *AuthCheck {
	return &AuthCheck{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthCheck) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *AuthCheck) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewAuthCheck()

// AuthExchangeAuthToken request parameters
type AuthExchangeAuthToken struct {
	hasToken bool
	rawToken string
	Token    string
}

// NewAuthExchangeAuthToken request
func NewAuthExchangeAuthToken() *AuthExchangeAuthToken {
	return &AuthExchangeAuthToken{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthExchangeAuthToken) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token

	return out
}

// Fill processes request and fills internal variables
func (r *AuthExchangeAuthToken) Fill(req *http.Request) (err error) {
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

	if val, ok := post["token"]; ok {
		r.hasToken = true
		r.rawToken = val
		r.Token = val
	}

	return err
}

var _ RequestFiller = NewAuthExchangeAuthToken()

// AuthLogout request parameters
type AuthLogout struct {
}

// NewAuthLogout request
func NewAuthLogout() *AuthLogout {
	return &AuthLogout{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthLogout) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *AuthLogout) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewAuthLogout()

// HasToken returns true if token was set
func (r *AuthExchangeAuthToken) HasToken() bool {
	return r.hasToken
}

// RawToken returns raw value of token parameter
func (r *AuthExchangeAuthToken) RawToken() string {
	return r.rawToken
}

// GetToken returns casted value of  token parameter
func (r *AuthExchangeAuthToken) GetToken() string {
	return r.Token
}
