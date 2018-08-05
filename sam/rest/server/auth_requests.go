package server

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
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Auth login request parameters
type AuthLoginRequest struct {
	Username string
	Password string
}

func (AuthLoginRequest) new() *AuthLoginRequest {
	return &AuthLoginRequest{}
}

func (a *AuthLoginRequest) Fill(r *http.Request) error {
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

	a.Username = post["username"]

	a.Password = post["password"]
	return nil
}

var _ RequestFiller = AuthLoginRequest{}.new()

// Auth create request parameters
type AuthCreateRequest struct {
	Name     string
	Email    string
	Username string
	Password string
}

func (AuthCreateRequest) new() *AuthCreateRequest {
	return &AuthCreateRequest{}
}

func (a *AuthCreateRequest) Fill(r *http.Request) error {
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

	a.Name = post["name"]

	a.Email = post["email"]

	a.Username = post["username"]

	a.Password = post["password"]
	return nil
}

var _ RequestFiller = AuthCreateRequest{}.new()
