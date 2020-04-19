package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `settings.go`, `settings.util.go` or `settings_test.go` to
	implement your API calls, helper functions and tests. The file `settings.go`
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

	"github.com/cortezaproject/corteza-server/pkg/settings"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// SettingsList request parameters
type SettingsList struct {
	hasPrefix bool
	rawPrefix string
	Prefix    string
}

// NewSettingsList request
func NewSettingsList() *SettingsList {
	return &SettingsList{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["prefix"] = r.Prefix

	return out
}

// Fill processes request and fills internal variables
func (r *SettingsList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["prefix"]; ok {
		r.hasPrefix = true
		r.rawPrefix = val
		r.Prefix = val
	}

	return err
}

var _ RequestFiller = NewSettingsList()

// SettingsUpdate request parameters
type SettingsUpdate struct {
	hasValues bool
	rawValues string
	Values    settings.ValueSet
}

// NewSettingsUpdate request
func NewSettingsUpdate() *SettingsUpdate {
	return &SettingsUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["values"] = r.Values

	return out
}

// Fill processes request and fills internal variables
func (r *SettingsUpdate) Fill(req *http.Request) (err error) {
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

	return err
}

var _ RequestFiller = NewSettingsUpdate()

// SettingsGet request parameters
type SettingsGet struct {
	hasOwnerID bool
	rawOwnerID string
	OwnerID    uint64 `json:",string"`

	hasKey bool
	rawKey string
	Key    string
}

// NewSettingsGet request
func NewSettingsGet() *SettingsGet {
	return &SettingsGet{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsGet) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["ownerID"] = r.OwnerID
	out["key"] = r.Key

	return out
}

// Fill processes request and fills internal variables
func (r *SettingsGet) Fill(req *http.Request) (err error) {
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

	if val, ok := get["ownerID"]; ok {
		r.hasOwnerID = true
		r.rawOwnerID = val
		r.OwnerID = parseUInt64(val)
	}
	r.hasKey = true
	r.rawKey = chi.URLParam(req, "key")
	r.Key = chi.URLParam(req, "key")

	return err
}

var _ RequestFiller = NewSettingsGet()

// SettingsSet request parameters
type SettingsSet struct {
	hasKey bool
	rawKey string
	Key    string

	hasUpload bool
	rawUpload string
	Upload    *multipart.FileHeader

	hasOwnerID bool
	rawOwnerID string
	OwnerID    uint64 `json:",string"`
}

// NewSettingsSet request
func NewSettingsSet() *SettingsSet {
	return &SettingsSet{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsSet) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["key"] = r.Key
	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	out["ownerID"] = r.OwnerID

	return out
}

// Fill processes request and fills internal variables
func (r *SettingsSet) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseMultipartForm(32 << 20); err != nil {
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

	r.hasKey = true
	r.rawKey = chi.URLParam(req, "key")
	r.Key = chi.URLParam(req, "key")
	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	if val, ok := post["ownerID"]; ok {
		r.hasOwnerID = true
		r.rawOwnerID = val
		r.OwnerID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewSettingsSet()

// SettingsCurrent request parameters
type SettingsCurrent struct {
}

// NewSettingsCurrent request
func NewSettingsCurrent() *SettingsCurrent {
	return &SettingsCurrent{}
}

// Auditable returns all auditable/loggable parameters
func (r SettingsCurrent) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	return out
}

// Fill processes request and fills internal variables
func (r *SettingsCurrent) Fill(req *http.Request) (err error) {
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

	return err
}

var _ RequestFiller = NewSettingsCurrent()

// HasPrefix returns true if prefix was set
func (r *SettingsList) HasPrefix() bool {
	return r.hasPrefix
}

// RawPrefix returns raw value of prefix parameter
func (r *SettingsList) RawPrefix() string {
	return r.rawPrefix
}

// GetPrefix returns casted value of  prefix parameter
func (r *SettingsList) GetPrefix() string {
	return r.Prefix
}

// HasValues returns true if values was set
func (r *SettingsUpdate) HasValues() bool {
	return r.hasValues
}

// RawValues returns raw value of values parameter
func (r *SettingsUpdate) RawValues() string {
	return r.rawValues
}

// GetValues returns casted value of  values parameter
func (r *SettingsUpdate) GetValues() settings.ValueSet {
	return r.Values
}

// HasOwnerID returns true if ownerID was set
func (r *SettingsGet) HasOwnerID() bool {
	return r.hasOwnerID
}

// RawOwnerID returns raw value of ownerID parameter
func (r *SettingsGet) RawOwnerID() string {
	return r.rawOwnerID
}

// GetOwnerID returns casted value of  ownerID parameter
func (r *SettingsGet) GetOwnerID() uint64 {
	return r.OwnerID
}

// HasKey returns true if key was set
func (r *SettingsGet) HasKey() bool {
	return r.hasKey
}

// RawKey returns raw value of key parameter
func (r *SettingsGet) RawKey() string {
	return r.rawKey
}

// GetKey returns casted value of  key parameter
func (r *SettingsGet) GetKey() string {
	return r.Key
}

// HasKey returns true if key was set
func (r *SettingsSet) HasKey() bool {
	return r.hasKey
}

// RawKey returns raw value of key parameter
func (r *SettingsSet) RawKey() string {
	return r.rawKey
}

// GetKey returns casted value of  key parameter
func (r *SettingsSet) GetKey() string {
	return r.Key
}

// HasUpload returns true if upload was set
func (r *SettingsSet) HasUpload() bool {
	return r.hasUpload
}

// RawUpload returns raw value of upload parameter
func (r *SettingsSet) RawUpload() string {
	return r.rawUpload
}

// GetUpload returns casted value of  upload parameter
func (r *SettingsSet) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// HasOwnerID returns true if ownerID was set
func (r *SettingsSet) HasOwnerID() bool {
	return r.hasOwnerID
}

// RawOwnerID returns raw value of ownerID parameter
func (r *SettingsSet) RawOwnerID() string {
	return r.rawOwnerID
}

// GetOwnerID returns casted value of  ownerID parameter
func (r *SettingsSet) GetOwnerID() uint64 {
	return r.OwnerID
}
