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

func (r AuthInternalLogin) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email
	out["password"] = "*masked*sensitive*data*"

	return out
}

func (r *AuthInternalLogin) Fill(req *http.Request) (err error) {
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
	if val, ok := post["password"]; ok {
		r.Password = val
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

func (r AuthInternalSignup) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email
	out["username"] = r.Username
	out["password"] = "*masked*sensitive*data*"

	out["handle"] = r.Handle
	out["name"] = r.Name

	return out
}

func (r *AuthInternalSignup) Fill(req *http.Request) (err error) {
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
	if val, ok := post["username"]; ok {
		r.Username = val
	}
	if val, ok := post["password"]; ok {
		r.Password = val
	}
	if val, ok := post["handle"]; ok {
		r.Handle = val
	}
	if val, ok := post["name"]; ok {
		r.Name = val
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

func (r AuthInternalRequestPasswordReset) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email

	return out
}

func (r *AuthInternalRequestPasswordReset) Fill(req *http.Request) (err error) {
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

	return err
}

var _ RequestFiller = NewAuthInternalRequestPasswordReset()

// AuthInternal exchangePasswordResetToken request parameters
type AuthInternalExchangePasswordResetToken struct {
	Token string
}

func NewAuthInternalExchangePasswordResetToken() *AuthInternalExchangePasswordResetToken {
	return &AuthInternalExchangePasswordResetToken{}
}

func (r AuthInternalExchangePasswordResetToken) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token

	return out
}

func (r *AuthInternalExchangePasswordResetToken) Fill(req *http.Request) (err error) {
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
		r.Token = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalExchangePasswordResetToken()

// AuthInternal resetPassword request parameters
type AuthInternalResetPassword struct {
	Token    string
	Password string
}

func NewAuthInternalResetPassword() *AuthInternalResetPassword {
	return &AuthInternalResetPassword{}
}

func (r AuthInternalResetPassword) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token
	out["password"] = "*masked*sensitive*data*"

	return out
}

func (r *AuthInternalResetPassword) Fill(req *http.Request) (err error) {
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
		r.Token = val
	}
	if val, ok := post["password"]; ok {
		r.Password = val
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

func (r AuthInternalConfirmEmail) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token

	return out
}

func (r *AuthInternalConfirmEmail) Fill(req *http.Request) (err error) {
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
		r.Token = val
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

func (r AuthInternalChangePassword) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["oldPassword"] = "*masked*sensitive*data*"

	out["newPassword"] = "*masked*sensitive*data*"

	return out
}

func (r *AuthInternalChangePassword) Fill(req *http.Request) (err error) {
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

	if val, ok := post["oldPassword"]; ok {
		r.OldPassword = val
	}
	if val, ok := post["newPassword"]; ok {
		r.NewPassword = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalChangePassword()
