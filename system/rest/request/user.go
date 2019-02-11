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
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
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

func (u *UserList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	if val, ok := get["query"]; ok {

		u.Query = val
	}
	if val, ok := get["username"]; ok {

		u.Username = val
	}
	if val, ok := get["email"]; ok {

		u.Email = val
	}

	return err
}

var _ RequestFiller = NewUserList()

// User create request parameters
type UserCreate struct {
	Email  string
	Name   string
	Handle string
	Kind   string
}

func NewUserCreate() *UserCreate {
	return &UserCreate{}
}

func (u *UserCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	if val, ok := post["email"]; ok {

		u.Email = val
	}
	if val, ok := post["name"]; ok {

		u.Name = val
	}
	if val, ok := post["handle"]; ok {

		u.Handle = val
	}
	if val, ok := post["kind"]; ok {

		u.Kind = val
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
	Kind   string
}

func NewUserUpdate() *UserUpdate {
	return &UserUpdate{}
}

func (u *UserUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	u.UserID = parseUInt64(chi.URLParam(r, "userID"))
	if val, ok := post["email"]; ok {

		u.Email = val
	}
	if val, ok := post["name"]; ok {

		u.Name = val
	}
	if val, ok := post["handle"]; ok {

		u.Handle = val
	}
	if val, ok := post["kind"]; ok {

		u.Kind = val
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

func (u *UserRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	u.UserID = parseUInt64(chi.URLParam(r, "userID"))

	return err
}

var _ RequestFiller = NewUserRead()

// User remove request parameters
type UserRemove struct {
	UserID uint64 `json:",string"`
}

func NewUserRemove() *UserRemove {
	return &UserRemove{}
}

func (u *UserRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	u.UserID = parseUInt64(chi.URLParam(r, "userID"))

	return err
}

var _ RequestFiller = NewUserRemove()

// User suspend request parameters
type UserSuspend struct {
	UserID uint64 `json:",string"`
}

func NewUserSuspend() *UserSuspend {
	return &UserSuspend{}
}

func (u *UserSuspend) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	u.UserID = parseUInt64(chi.URLParam(r, "userID"))

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

func (u *UserUnsuspend) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(u)

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

	u.UserID = parseUInt64(chi.URLParam(r, "userID"))

	return err
}

var _ RequestFiller = NewUserUnsuspend()
