package request

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

// User search request parameters
type UserSearch struct {
	Query string
}

func NewUserSearch() *UserSearch {
	return &UserSearch{}
}

func (u *UserSearch) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	if val, ok := get["query"]; ok {

		u.Query = val
	}

	return err
}

var _ RequestFiller = NewUserSearch()

// User message request parameters
type UserMessage struct {
	UserID  uint64
	Message string
}

func NewUserMessage() *UserMessage {
	return &UserMessage{}
}

func (u *UserMessage) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	u.UserID = parseUInt64(chi.URLParam(r, "userID"))
	if val, ok := post["message"]; ok {

		u.Message = val
	}

	return err
}

var _ RequestFiller = NewUserMessage()
