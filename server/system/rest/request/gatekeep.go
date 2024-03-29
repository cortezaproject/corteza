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
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-chi/chi/v5"
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
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	GatekeepLock struct {
		// Resource POST parameter
		//
		// Resource to lock
		Resource string

		// ExpiresIn POST parameter
		//
		// Standard go duration string (e.g. 1h30m)
		ExpiresIn string

		// UserID POST parameter
		//
		// User ID that locked the resource
		UserID uint64 `json:",string"`
	}

	GatekeepUnlock struct {
		// Resource POST parameter
		//
		// Resource to unlock
		Resource string

		// UserID POST parameter
		//
		// User ID that locked the resource
		UserID uint64 `json:",string"`
	}

	GatekeepCheck struct {
		// LockID PATH parameter
		//
		// Lock to check
		LockID uint64 `json:",string"`

		// Resource POST parameter
		//
		// Resource to lock
		Resource string

		// UserID POST parameter
		//
		// User ID that locked the resource
		UserID uint64 `json:",string"`
	}
)

// NewGatekeepLock request
func NewGatekeepLock() *GatekeepLock {
	return &GatekeepLock{}
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepLock) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource":  r.Resource,
		"expiresIn": r.ExpiresIn,
		"userID":    r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepLock) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepLock) GetExpiresIn() string {
	return r.ExpiresIn
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepLock) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *GatekeepLock) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["resource"]; ok && len(val) > 0 {
				r.Resource, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["expiresIn"]; ok && len(val) > 0 {
				r.ExpiresIn, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["userID"]; ok && len(val) > 0 {
				r.UserID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["expiresIn"]; ok && len(val) > 0 {
			r.ExpiresIn, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewGatekeepUnlock request
func NewGatekeepUnlock() *GatekeepUnlock {
	return &GatekeepUnlock{}
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepUnlock) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource": r.Resource,
		"userID":   r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepUnlock) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepUnlock) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *GatekeepUnlock) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["resource"]; ok && len(val) > 0 {
				r.Resource, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["userID"]; ok && len(val) > 0 {
				r.UserID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewGatekeepCheck request
func NewGatekeepCheck() *GatekeepCheck {
	return &GatekeepCheck{}
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepCheck) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"lockID":   r.LockID,
		"resource": r.Resource,
		"userID":   r.UserID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepCheck) GetLockID() uint64 {
	return r.LockID
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepCheck) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r GatekeepCheck) GetUserID() uint64 {
	return r.UserID
}

// Fill processes request and fills internal variables
func (r *GatekeepCheck) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["resource"]; ok && len(val) > 0 {
				r.Resource, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["userID"]; ok && len(val) > 0 {
				r.UserID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["userID"]; ok && len(val) > 0 {
			r.UserID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "lockID")
		r.LockID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
