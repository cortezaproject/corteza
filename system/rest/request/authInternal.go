package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
)

type (
	// Internal API interface
	AuthInternalLogin struct {
		// Email POST parameter
		//
		// Email
		Email string

		// Password POST parameter
		//
		// Password
		Password string
	}

	AuthInternalSignup struct {
		// Email POST parameter
		//
		// Email
		Email string

		// Username POST parameter
		//
		// Username
		Username string

		// Password POST parameter
		//
		// Password
		Password string

		// Handle POST parameter
		//
		// User handle
		Handle string

		// Name POST parameter
		//
		// Display name
		Name string
	}

	AuthInternalRequestPasswordReset struct {
		// Email POST parameter
		//
		// Email
		Email string
	}

	AuthInternalExchangePasswordResetToken struct {
		// Token POST parameter
		//
		// Token
		Token string
	}

	AuthInternalResetPassword struct {
		// Token POST parameter
		//
		// Token
		Token string

		// Password POST parameter
		//
		// Password
		Password string
	}

	AuthInternalConfirmEmail struct {
		// Token POST parameter
		//
		// Token
		Token string
	}

	AuthInternalChangePassword struct {
		// OldPassword POST parameter
		//
		// Old password
		OldPassword string

		// NewPassword POST parameter
		//
		// New password
		NewPassword string
	}
)

// NewAuthInternalLogin request
func NewAuthInternalLogin() *AuthInternalLogin {
	return &AuthInternalLogin{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalLogin) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"email": r.Email,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalLogin) GetEmail() string {
	return r.Email
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalLogin) GetPassword() string {
	return r.Password
}

// Fill processes request and fills internal variables
func (r *AuthInternalLogin) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["email"]; ok && len(val) > 0 {
			r.Email, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["password"]; ok && len(val) > 0 {
			r.Password, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthInternalSignup request
func NewAuthInternalSignup() *AuthInternalSignup {
	return &AuthInternalSignup{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"email":    r.Email,
		"username": r.Username,
		"handle":   r.Handle,
		"name":     r.Name,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) GetEmail() string {
	return r.Email
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) GetUsername() string {
	return r.Username
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) GetPassword() string {
	return r.Password
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalSignup) GetName() string {
	return r.Name
}

// Fill processes request and fills internal variables
func (r *AuthInternalSignup) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["email"]; ok && len(val) > 0 {
			r.Email, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["username"]; ok && len(val) > 0 {
			r.Username, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["password"]; ok && len(val) > 0 {
			r.Password, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthInternalRequestPasswordReset request
func NewAuthInternalRequestPasswordReset() *AuthInternalRequestPasswordReset {
	return &AuthInternalRequestPasswordReset{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalRequestPasswordReset) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"email": r.Email,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalRequestPasswordReset) GetEmail() string {
	return r.Email
}

// Fill processes request and fills internal variables
func (r *AuthInternalRequestPasswordReset) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["email"]; ok && len(val) > 0 {
			r.Email, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthInternalExchangePasswordResetToken request
func NewAuthInternalExchangePasswordResetToken() *AuthInternalExchangePasswordResetToken {
	return &AuthInternalExchangePasswordResetToken{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalExchangePasswordResetToken) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"token": r.Token,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalExchangePasswordResetToken) GetToken() string {
	return r.Token
}

// Fill processes request and fills internal variables
func (r *AuthInternalExchangePasswordResetToken) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["token"]; ok && len(val) > 0 {
			r.Token, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthInternalResetPassword request
func NewAuthInternalResetPassword() *AuthInternalResetPassword {
	return &AuthInternalResetPassword{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalResetPassword) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"token": r.Token,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalResetPassword) GetToken() string {
	return r.Token
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalResetPassword) GetPassword() string {
	return r.Password
}

// Fill processes request and fills internal variables
func (r *AuthInternalResetPassword) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["token"]; ok && len(val) > 0 {
			r.Token, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["password"]; ok && len(val) > 0 {
			r.Password, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthInternalConfirmEmail request
func NewAuthInternalConfirmEmail() *AuthInternalConfirmEmail {
	return &AuthInternalConfirmEmail{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalConfirmEmail) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"token": r.Token,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalConfirmEmail) GetToken() string {
	return r.Token
}

// Fill processes request and fills internal variables
func (r *AuthInternalConfirmEmail) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["token"]; ok && len(val) > 0 {
			r.Token, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthInternalChangePassword request
func NewAuthInternalChangePassword() *AuthInternalChangePassword {
	return &AuthInternalChangePassword{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalChangePassword) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalChangePassword) GetOldPassword() string {
	return r.OldPassword
}

// Auditable returns all auditable/loggable parameters
func (r AuthInternalChangePassword) GetNewPassword() string {
	return r.NewPassword
}

// Fill processes request and fills internal variables
func (r *AuthInternalChangePassword) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["oldPassword"]; ok && len(val) > 0 {
			r.OldPassword, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["newPassword"]; ok && len(val) > 0 {
			r.NewPassword, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
