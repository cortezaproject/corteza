package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `attachment.go`, `attachment.util.go` or `attachment_test.go` to
	implement your API calls, helper functions and tests. The file `attachment.go`
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

// AttachmentOriginal request parameters
type AttachmentOriginal struct {
	hasDownload bool
	rawDownload string
	Download    bool

	hasSign bool
	rawSign string
	Sign    string

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasAttachmentID bool
	rawAttachmentID string
	AttachmentID    uint64 `json:",string"`
}

// NewAttachmentOriginal request
func NewAttachmentOriginal() *AttachmentOriginal {
	return &AttachmentOriginal{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentOriginal) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["download"] = r.Download
	out["sign"] = r.Sign
	out["userID"] = r.UserID
	out["name"] = r.Name
	out["attachmentID"] = r.AttachmentID

	return out
}

// Fill processes request and fills internal variables
func (r *AttachmentOriginal) Fill(req *http.Request) (err error) {
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

	if val, ok := get["download"]; ok {
		r.hasDownload = true
		r.rawDownload = val
		r.Download = parseBool(val)
	}
	if val, ok := get["sign"]; ok {
		r.hasSign = true
		r.rawSign = val
		r.Sign = val
	}
	if val, ok := get["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseUInt64(val)
	}
	r.hasName = true
	r.rawName = chi.URLParam(req, "name")
	r.Name = chi.URLParam(req, "name")
	r.hasAttachmentID = true
	r.rawAttachmentID = chi.URLParam(req, "attachmentID")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))

	return err
}

var _ RequestFiller = NewAttachmentOriginal()

// AttachmentPreview request parameters
type AttachmentPreview struct {
	hasExt bool
	rawExt string
	Ext    string

	hasAttachmentID bool
	rawAttachmentID string
	AttachmentID    uint64 `json:",string"`

	hasSign bool
	rawSign string
	Sign    string

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewAttachmentPreview request
func NewAttachmentPreview() *AttachmentPreview {
	return &AttachmentPreview{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentPreview) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["ext"] = r.Ext
	out["attachmentID"] = r.AttachmentID
	out["sign"] = r.Sign
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *AttachmentPreview) Fill(req *http.Request) (err error) {
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

	r.hasExt = true
	r.rawExt = chi.URLParam(req, "ext")
	r.Ext = chi.URLParam(req, "ext")
	r.hasAttachmentID = true
	r.rawAttachmentID = chi.URLParam(req, "attachmentID")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))
	if val, ok := get["sign"]; ok {
		r.hasSign = true
		r.rawSign = val
		r.Sign = val
	}
	if val, ok := get["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewAttachmentPreview()

// HasDownload returns true if download was set
func (r *AttachmentOriginal) HasDownload() bool {
	return r.hasDownload
}

// RawDownload returns raw value of download parameter
func (r *AttachmentOriginal) RawDownload() string {
	return r.rawDownload
}

// GetDownload returns casted value of  download parameter
func (r *AttachmentOriginal) GetDownload() bool {
	return r.Download
}

// HasSign returns true if sign was set
func (r *AttachmentOriginal) HasSign() bool {
	return r.hasSign
}

// RawSign returns raw value of sign parameter
func (r *AttachmentOriginal) RawSign() string {
	return r.rawSign
}

// GetSign returns casted value of  sign parameter
func (r *AttachmentOriginal) GetSign() string {
	return r.Sign
}

// HasUserID returns true if userID was set
func (r *AttachmentOriginal) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *AttachmentOriginal) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *AttachmentOriginal) GetUserID() uint64 {
	return r.UserID
}

// HasName returns true if name was set
func (r *AttachmentOriginal) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *AttachmentOriginal) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *AttachmentOriginal) GetName() string {
	return r.Name
}

// HasAttachmentID returns true if attachmentID was set
func (r *AttachmentOriginal) HasAttachmentID() bool {
	return r.hasAttachmentID
}

// RawAttachmentID returns raw value of attachmentID parameter
func (r *AttachmentOriginal) RawAttachmentID() string {
	return r.rawAttachmentID
}

// GetAttachmentID returns casted value of  attachmentID parameter
func (r *AttachmentOriginal) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// HasExt returns true if ext was set
func (r *AttachmentPreview) HasExt() bool {
	return r.hasExt
}

// RawExt returns raw value of ext parameter
func (r *AttachmentPreview) RawExt() string {
	return r.rawExt
}

// GetExt returns casted value of  ext parameter
func (r *AttachmentPreview) GetExt() string {
	return r.Ext
}

// HasAttachmentID returns true if attachmentID was set
func (r *AttachmentPreview) HasAttachmentID() bool {
	return r.hasAttachmentID
}

// RawAttachmentID returns raw value of attachmentID parameter
func (r *AttachmentPreview) RawAttachmentID() string {
	return r.rawAttachmentID
}

// GetAttachmentID returns casted value of  attachmentID parameter
func (r *AttachmentPreview) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// HasSign returns true if sign was set
func (r *AttachmentPreview) HasSign() bool {
	return r.hasSign
}

// RawSign returns raw value of sign parameter
func (r *AttachmentPreview) RawSign() string {
	return r.rawSign
}

// GetSign returns casted value of  sign parameter
func (r *AttachmentPreview) GetSign() string {
	return r.Sign
}

// HasUserID returns true if userID was set
func (r *AttachmentPreview) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *AttachmentPreview) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *AttachmentPreview) GetUserID() uint64 {
	return r.UserID
}
