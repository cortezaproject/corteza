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

// AuthInternalLogin request parameters
type AuthInternalLogin struct {
	hasEmail bool
	rawEmail string
	Email    string

	hasPassword bool
	rawPassword string
	Password    string
}

// NewAuthInternalLogin request
func NewAuthInternalLogin() *AuthInternalLogin {
	return &AuthInternalLogin{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalLogin) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email
	out["password"] = "*masked*sensitive*data*"

	return out
}

// Fill processes request and fills internal variables
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
		r.hasEmail = true
		r.rawEmail = val
		r.Email = val
	}
	if val, ok := post["password"]; ok {
		r.hasPassword = true
		r.rawPassword = val
		r.Password = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalLogin()

// AuthInternalSignup request parameters
type AuthInternalSignup struct {
	hasEmail bool
	rawEmail string
	Email    string

	hasUsername bool
	rawUsername string
	Username    string

	hasPassword bool
	rawPassword string
	Password    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasName bool
	rawName string
	Name    string
}

// NewAuthInternalSignup request
func NewAuthInternalSignup() *AuthInternalSignup {
	return &AuthInternalSignup{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email
	out["username"] = r.Username
	out["password"] = "*masked*sensitive*data*"

	out["handle"] = r.Handle
	out["name"] = r.Name

	return out
}

// Fill processes request and fills internal variables
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
		r.hasEmail = true
		r.rawEmail = val
		r.Email = val
	}
	if val, ok := post["username"]; ok {
		r.hasUsername = true
		r.rawUsername = val
		r.Username = val
	}
	if val, ok := post["password"]; ok {
		r.hasPassword = true
		r.rawPassword = val
		r.Password = val
	}
	if val, ok := post["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
	}
	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalSignup()

// AuthInternalRequestPasswordReset request parameters
type AuthInternalRequestPasswordReset struct {
	hasEmail bool
	rawEmail string
	Email    string
}

// NewAuthInternalRequestPasswordReset request
func NewAuthInternalRequestPasswordReset() *AuthInternalRequestPasswordReset {
	return &AuthInternalRequestPasswordReset{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalRequestPasswordReset) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email

	return out
}

// Fill processes request and fills internal variables
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
		r.hasEmail = true
		r.rawEmail = val
		r.Email = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalRequestPasswordReset()

// AuthInternalExchangePasswordResetToken request parameters
type AuthInternalExchangePasswordResetToken struct {
	hasToken bool
	rawToken string
	Token    string
}

// NewAuthInternalExchangePasswordResetToken request
func NewAuthInternalExchangePasswordResetToken() *AuthInternalExchangePasswordResetToken {
	return &AuthInternalExchangePasswordResetToken{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalExchangePasswordResetToken) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token

	return out
}

// Fill processes request and fills internal variables
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
		r.hasToken = true
		r.rawToken = val
		r.Token = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalExchangePasswordResetToken()

// AuthInternalResetPassword request parameters
type AuthInternalResetPassword struct {
	hasToken bool
	rawToken string
	Token    string

	hasPassword bool
	rawPassword string
	Password    string
}

// NewAuthInternalResetPassword request
func NewAuthInternalResetPassword() *AuthInternalResetPassword {
	return &AuthInternalResetPassword{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalResetPassword) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token
	out["password"] = "*masked*sensitive*data*"

	return out
}

// Fill processes request and fills internal variables
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
		r.hasToken = true
		r.rawToken = val
		r.Token = val
	}
	if val, ok := post["password"]; ok {
		r.hasPassword = true
		r.rawPassword = val
		r.Password = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalResetPassword()

// AuthInternalConfirmEmail request parameters
type AuthInternalConfirmEmail struct {
	hasToken bool
	rawToken string
	Token    string
}

// NewAuthInternalConfirmEmail request
func NewAuthInternalConfirmEmail() *AuthInternalConfirmEmail {
	return &AuthInternalConfirmEmail{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalConfirmEmail) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["token"] = r.Token

	return out
}

// Fill processes request and fills internal variables
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
		r.hasToken = true
		r.rawToken = val
		r.Token = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalConfirmEmail()

// AuthInternalChangePassword request parameters
type AuthInternalChangePassword struct {
	hasOldPassword bool
	rawOldPassword string
	OldPassword    string

	hasNewPassword bool
	rawNewPassword string
	NewPassword    string
}

// NewAuthInternalChangePassword request
func NewAuthInternalChangePassword() *AuthInternalChangePassword {
	return &AuthInternalChangePassword{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalChangePassword) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["oldPassword"] = "*masked*sensitive*data*"

	out["newPassword"] = "*masked*sensitive*data*"

	return out
}

// Fill processes request and fills internal variables
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
		r.hasOldPassword = true
		r.rawOldPassword = val
		r.OldPassword = val
	}
	if val, ok := post["newPassword"]; ok {
		r.hasNewPassword = true
		r.rawNewPassword = val
		r.NewPassword = val
	}

	return err
}

var _ RequestFiller = NewAuthInternalChangePassword()

// HasEmail returns true if email was set
func (r *AuthInternalLogin) HasEmail() bool {
	return r.hasEmail
}

// RawEmail returns raw value of email parameter
func (r *AuthInternalLogin) RawEmail() string {
	return r.rawEmail
}

// GetEmail returns casted value of  email parameter
func (r *AuthInternalLogin) GetEmail() string {
	return r.Email
}

// HasPassword returns true if password was set
func (r *AuthInternalLogin) HasPassword() bool {
	return r.hasPassword
}

// RawPassword returns raw value of password parameter
func (r *AuthInternalLogin) RawPassword() string {
	return r.rawPassword
}

// GetPassword returns casted value of  password parameter
func (r *AuthInternalLogin) GetPassword() string {
	return r.Password
}

// HasEmail returns true if email was set
func (r *AuthInternalSignup) HasEmail() bool {
	return r.hasEmail
}

// RawEmail returns raw value of email parameter
func (r *AuthInternalSignup) RawEmail() string {
	return r.rawEmail
}

// GetEmail returns casted value of  email parameter
func (r *AuthInternalSignup) GetEmail() string {
	return r.Email
}

// HasUsername returns true if username was set
func (r *AuthInternalSignup) HasUsername() bool {
	return r.hasUsername
}

// RawUsername returns raw value of username parameter
func (r *AuthInternalSignup) RawUsername() string {
	return r.rawUsername
}

// GetUsername returns casted value of  username parameter
func (r *AuthInternalSignup) GetUsername() string {
	return r.Username
}

// HasPassword returns true if password was set
func (r *AuthInternalSignup) HasPassword() bool {
	return r.hasPassword
}

// RawPassword returns raw value of password parameter
func (r *AuthInternalSignup) RawPassword() string {
	return r.rawPassword
}

// GetPassword returns casted value of  password parameter
func (r *AuthInternalSignup) GetPassword() string {
	return r.Password
}

// HasHandle returns true if handle was set
func (r *AuthInternalSignup) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *AuthInternalSignup) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *AuthInternalSignup) GetHandle() string {
	return r.Handle
}

// HasName returns true if name was set
func (r *AuthInternalSignup) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *AuthInternalSignup) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *AuthInternalSignup) GetName() string {
	return r.Name
}

// HasEmail returns true if email was set
func (r *AuthInternalRequestPasswordReset) HasEmail() bool {
	return r.hasEmail
}

// RawEmail returns raw value of email parameter
func (r *AuthInternalRequestPasswordReset) RawEmail() string {
	return r.rawEmail
}

// GetEmail returns casted value of  email parameter
func (r *AuthInternalRequestPasswordReset) GetEmail() string {
	return r.Email
}

// HasToken returns true if token was set
func (r *AuthInternalExchangePasswordResetToken) HasToken() bool {
	return r.hasToken
}

// RawToken returns raw value of token parameter
func (r *AuthInternalExchangePasswordResetToken) RawToken() string {
	return r.rawToken
}

// GetToken returns casted value of  token parameter
func (r *AuthInternalExchangePasswordResetToken) GetToken() string {
	return r.Token
}

// HasToken returns true if token was set
func (r *AuthInternalResetPassword) HasToken() bool {
	return r.hasToken
}

// RawToken returns raw value of token parameter
func (r *AuthInternalResetPassword) RawToken() string {
	return r.rawToken
}

// GetToken returns casted value of  token parameter
func (r *AuthInternalResetPassword) GetToken() string {
	return r.Token
}

// HasPassword returns true if password was set
func (r *AuthInternalResetPassword) HasPassword() bool {
	return r.hasPassword
}

// RawPassword returns raw value of password parameter
func (r *AuthInternalResetPassword) RawPassword() string {
	return r.rawPassword
}

// GetPassword returns casted value of  password parameter
func (r *AuthInternalResetPassword) GetPassword() string {
	return r.Password
}

// HasToken returns true if token was set
func (r *AuthInternalConfirmEmail) HasToken() bool {
	return r.hasToken
}

// RawToken returns raw value of token parameter
func (r *AuthInternalConfirmEmail) RawToken() string {
	return r.rawToken
}

// GetToken returns casted value of  token parameter
func (r *AuthInternalConfirmEmail) GetToken() string {
	return r.Token
}

// HasOldPassword returns true if oldPassword was set
func (r *AuthInternalChangePassword) HasOldPassword() bool {
	return r.hasOldPassword
}

// RawOldPassword returns raw value of oldPassword parameter
func (r *AuthInternalChangePassword) RawOldPassword() string {
	return r.rawOldPassword
}

// GetOldPassword returns casted value of  oldPassword parameter
func (r *AuthInternalChangePassword) GetOldPassword() string {
	return r.OldPassword
}

// HasNewPassword returns true if newPassword was set
func (r *AuthInternalChangePassword) HasNewPassword() bool {
	return r.hasNewPassword
}

// RawNewPassword returns raw value of newPassword parameter
func (r *AuthInternalChangePassword) RawNewPassword() string {
	return r.rawNewPassword
}

// GetNewPassword returns casted value of  newPassword parameter
func (r *AuthInternalChangePassword) GetNewPassword() string {
	return r.NewPassword
}
