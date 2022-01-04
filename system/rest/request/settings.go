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
	"github.com/cortezaproject/corteza-server/system/types"
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
	SettingsList struct {
		// Prefix GET parameter
		//
		// Key prefix
		Prefix string
	}

	SettingsUpdate struct {
		// Values POST parameter
		//
		// Array of new settings: `[{ name: ..., value: ... }]`. Omit value to remove setting
		Values types.SettingValueSet
	}

	SettingsGet struct {
		// Key PATH parameter
		//
		// Setting key
		Key string

		// OwnerID GET parameter
		//
		// Owner ID
		OwnerID uint64 `json:",string"`
	}

	SettingsSet struct {
		// Key PATH parameter
		//
		// Key
		Key string

		// Upload POST parameter
		//
		// File to upload
		Upload *multipart.FileHeader

		// OwnerID POST parameter
		//
		// Owner ID
		OwnerID uint64 `json:",string"`
	}

	SettingsCurrent struct {
	}
)

// NewSettingsList request
func NewSettingsList() *SettingsList {
	return &SettingsList{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"prefix": r.Prefix,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsList) GetPrefix() string {
	return r.Prefix
}

// Fill processes request and fills internal variables
func (r *SettingsList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["prefix"]; ok && len(val) > 0 {
			r.Prefix, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewSettingsUpdate request
func NewSettingsUpdate() *SettingsUpdate {
	return &SettingsUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"values": r.Values,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsUpdate) GetValues() types.SettingValueSet {
	return r.Values
}

// Fill processes request and fills internal variables
func (r *SettingsUpdate) Fill(req *http.Request) (err error) {

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
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["values[]"]; ok && len(val) > 0  {
		//    r.Values, err = types.SettingValueSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	return err
}

// NewSettingsGet request
func NewSettingsGet() *SettingsGet {
	return &SettingsGet{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsGet) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"key":     r.Key,
		"ownerID": r.OwnerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsGet) GetKey() string {
	return r.Key
}

// Auditable returns all auditable/loggable parameters
func (r SettingsGet) GetOwnerID() uint64 {
	return r.OwnerID
}

// Fill processes request and fills internal variables
func (r *SettingsGet) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "key")
		r.Key, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSettingsSet request
func NewSettingsSet() *SettingsSet {
	return &SettingsSet{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsSet) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"key":     r.Key,
		"upload":  r.Upload,
		"ownerID": r.OwnerID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsSet) GetKey() string {
	return r.Key
}

// Auditable returns all auditable/loggable parameters
func (r SettingsSet) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Auditable returns all auditable/loggable parameters
func (r SettingsSet) GetOwnerID() uint64 {
	return r.OwnerID
}

// Fill processes request and fills internal variables
func (r *SettingsSet) Fill(req *http.Request) (err error) {

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
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			// Ignoring upload as its handled in the POST params section

			if val, ok := req.MultipartForm.Value["ownerID"]; ok && len(val) > 0 {
				r.OwnerID, err = payload.ParseUint64(val[0]), nil
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

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

		if val, ok := req.Form["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "key")
		r.Key, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSettingsCurrent request
func NewSettingsCurrent() *SettingsCurrent {
	return &SettingsCurrent{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsCurrent) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *SettingsCurrent) Fill(req *http.Request) (err error) {

	return err
}
