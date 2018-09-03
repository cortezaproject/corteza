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
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}

// Auth login request parameters
type AuthLogin struct {
	Username string
	Password string
}

func NewAuthLogin() *AuthLogin {
	return &AuthLogin{}
}

func (a *AuthLogin) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(a)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}

		return err
	}

	r.ParseForm()
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

	if val, ok := post["username"]; ok {

		a.Username = val
	}
	if val, ok := post["password"]; ok {

		a.Password = val
	}

	return err
}

var _ RequestFiller = NewAuthLogin()

// Auth create request parameters
type AuthCreate struct {
	Name     string
	Email    string
	Username string
	Password string
}

func NewAuthCreate() *AuthCreate {
	return &AuthCreate{}
}

func (a *AuthCreate) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(a)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			err = errors.Wrap(err, "error parsing http request body")
		}

		return err
	}

	r.ParseForm()
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

		a.Name = val
	}
	if val, ok := post["email"]; ok {

		a.Email = val
	}
	if val, ok := post["username"]; ok {

		a.Username = val
	}
	if val, ok := post["password"]; ok {

		a.Password = val
	}

	return err
}

var _ RequestFiller = NewAuthCreate()
