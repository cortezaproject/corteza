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
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// User list request parameters
type UserList struct {
	Query    string
	Username string
	Email    string
}

func NewUserList() *UserList {
	return &UserList{}
}

func (r UserList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["username"] = r.Username
	out["email"] = r.Email

	return out
}

func (r *UserList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["query"]; ok {
		r.Query = val
	}
	if val, ok := get["username"]; ok {
		r.Username = val
	}
	if val, ok := get["email"]; ok {
		r.Email = val
	}

	return err
}

var _ RequestFiller = NewUserList()

// User create request parameters
type UserCreate struct {
	Email  string
	Name   string
	Handle string
	Kind   types.UserKind
}

func NewUserCreate() *UserCreate {
	return &UserCreate{}
}

func (r UserCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["kind"] = r.Kind

	return out
}

func (r *UserCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["email"]; ok {
		r.Email = val
	}
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["handle"]; ok {
		r.Handle = val
	}
	if val, ok := post["kind"]; ok {
		r.Kind = types.UserKind(val)
	}

	return err
}

var _ RequestFiller = NewUserCreate()

// User update request parameters
type UserUpdate struct {
	UserID uint64 `json:",string"`
	Email  string
	Name   string
	Handle string
	Kind   types.UserKind
}

func NewUserUpdate() *UserUpdate {
	return &UserUpdate{}
}

func (r UserUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID
	out["email"] = r.Email
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["kind"] = r.Kind

	return out
}

func (r *UserUpdate) Fill(req *http.Request) (err error) {
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

	r.UserID = parseUInt64(chi.URLParam(req, "userID"))
	if val, ok := post["email"]; ok {
		r.Email = val
	}
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["handle"]; ok {
		r.Handle = val
	}
	if val, ok := post["kind"]; ok {
		r.Kind = types.UserKind(val)
	}

	return err
}

var _ RequestFiller = NewUserUpdate()

// User read request parameters
type UserRead struct {
	UserID uint64 `json:",string"`
}

func NewUserRead() *UserRead {
	return &UserRead{}
}

func (r UserRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

func (r *UserRead) Fill(req *http.Request) (err error) {
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

	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserRead()

// User delete request parameters
type UserDelete struct {
	UserID uint64 `json:",string"`
}

func NewUserDelete() *UserDelete {
	return &UserDelete{}
}

func (r UserDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

func (r *UserDelete) Fill(req *http.Request) (err error) {
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

	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserDelete()

// User suspend request parameters
type UserSuspend struct {
	UserID uint64 `json:",string"`
}

func NewUserSuspend() *UserSuspend {
	return &UserSuspend{}
}

func (r UserSuspend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

func (r *UserSuspend) Fill(req *http.Request) (err error) {
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

	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserSuspend()

// User unsuspend request parameters
type UserUnsuspend struct {
	UserID uint64 `json:",string"`
}

func NewUserUnsuspend() *UserUnsuspend {
	return &UserUnsuspend{}
}

func (r UserUnsuspend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

func (r *UserUnsuspend) Fill(req *http.Request) (err error) {
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

	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserUnsuspend()
