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
	AuthSettings struct {
	}

	AuthCheck struct {
	}

	AuthImpersonate struct {
		// UserID POST parameter
		//
		// ID of the impersonated user
		UserID uint64 `json:",string"`
	}

	AuthExchangeAuthToken struct {
		// Token POST parameter
		//
		// Token to be exchanged for JWT
		Token string
	}

	AuthLogout struct {
	}
)

// NewAuthSettings request
func NewAuthSettings() *AuthSettings {
	return &AuthSettings{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthSettings) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *AuthSettings) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	return err
}

// NewAuthCheck request
func NewAuthCheck() *AuthCheck {
	return &AuthCheck{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthCheck) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *AuthCheck) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	return err
}

// NewAuthImpersonate request
func NewAuthImpersonate() *AuthImpersonate {
	return &AuthImpersonate{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthImpersonate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"userID": r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthImpersonate) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *AuthImpersonate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewAuthExchangeAuthToken request
func NewAuthExchangeAuthToken() *AuthExchangeAuthToken {
	return &AuthExchangeAuthToken{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthExchangeAuthToken) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"token": r.Token,
	}
}

// Auditable returns all auditable/loggable parameters
func (r AuthExchangeAuthToken) GetToken() string {
	return r.Token
}

// Fill processes request and fills internal variables
func (r *AuthExchangeAuthToken) Fill(req *http.Request) (err error) {
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

// NewAuthLogout request
func NewAuthLogout() *AuthLogout {
	return &AuthLogout{}
}

// Auditable returns all auditable/loggable parameters
func (r AuthLogout) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *AuthLogout) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	return err
}
