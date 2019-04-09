package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `auth_internal.go`, `auth_internal.util.go` or `auth_internal_test.go` to
	implement your API calls, helper functions and tests. The file `auth_internal.go`
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

// AuthInternal login request parameters
type AuthInternalLogin struct {
	Email    string
	Password string
}

func NewAuthInternalLogin() *AuthInternalLogin {
	return &AuthInternalLogin{}
}

func (auReq *AuthInternalLogin) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(auReq)

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

		auReq.Email = val
	}
	if val, ok := post["password"]; ok {

		auReq.Password = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalLogin()

// AuthInternal signup request parameters
type AuthInternalSignup struct {
	Email    string
	Username string
	Password string
	Handle   string
	Name     string
}

func NewAuthInternalSignup() *AuthInternalSignup {
	return &AuthInternalSignup{}
}

func (auReq *AuthInternalSignup) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(auReq)

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

		auReq.Email = val
	}
	if val, ok := post["username"]; ok {

		auReq.Username = val
	}
	if val, ok := post["password"]; ok {

		auReq.Password = val
	}
	if val, ok := post["handle"]; ok {

		auReq.Handle = val
	}
	if val, ok := post["name"]; ok {

		auReq.Name = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalSignup()

// AuthInternal requestPasswordReset request parameters
type AuthInternalRequestPasswordReset struct {
	Email string
}

func NewAuthInternalRequestPasswordReset() *AuthInternalRequestPasswordReset {
	return &AuthInternalRequestPasswordReset{}
}

func (auReq *AuthInternalRequestPasswordReset) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(auReq)

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

		auReq.Email = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalRequestPasswordReset()

// AuthInternal resetPassword request parameters
type AuthInternalResetPassword struct {
	Token    string
	Password string
}

func NewAuthInternalResetPassword() *AuthInternalResetPassword {
	return &AuthInternalResetPassword{}
}

func (auReq *AuthInternalResetPassword) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(auReq)

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

	if val, ok := post["token"]; ok {

		auReq.Token = val
	}
	if val, ok := post["password"]; ok {

		auReq.Password = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalResetPassword()

// AuthInternal confirmEmail request parameters
type AuthInternalConfirmEmail struct {
	Token string
}

func NewAuthInternalConfirmEmail() *AuthInternalConfirmEmail {
	return &AuthInternalConfirmEmail{}
}

func (auReq *AuthInternalConfirmEmail) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(auReq)

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

	if val, ok := post["token"]; ok {

		auReq.Token = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalConfirmEmail()

// AuthInternal changePassword request parameters
type AuthInternalChangePassword struct {
	OldPassword string
	NewPassword string
}

func NewAuthInternalChangePassword() *AuthInternalChangePassword {
	return &AuthInternalChangePassword{}
}

func (auReq *AuthInternalChangePassword) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(auReq)

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

	if val, ok := post["oldPassword"]; ok {

		auReq.OldPassword = val
	}
	if val, ok := post["newPassword"]; ok {

		auReq.NewPassword = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalChangePassword()
