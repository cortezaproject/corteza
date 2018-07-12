package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `user.go`, `user.util.go` or `user_test.go` to
	implement your API calls, helper functions and tests. The file `user.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// User login request parameters
type userLoginRequest struct {
	username string
	password string
}

func (userLoginRequest) new() *userLoginRequest {
	return &userLoginRequest{}
}

func (u *userLoginRequest) Fill(r *http.Request) error {
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

	u.username = post["username"]

	u.password = post["password"]
	return nil
}

var _ RequestFiller = userLoginRequest{}.new()

// User search request parameters
type userSearchRequest struct {
	query string
}

func (userSearchRequest) new() *userSearchRequest {
	return &userSearchRequest{}
}

func (u *userSearchRequest) Fill(r *http.Request) error {
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

	u.query = get["query"]
	return nil
}

var _ RequestFiller = userSearchRequest{}.new()
