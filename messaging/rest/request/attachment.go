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

// Attachment original request parameters
type AttachmentOriginal struct {
	Download     bool
	Sign         string
	UserID       uint64 `json:",string"`
	Name         string
	AttachmentID uint64 `json:",string"`
}

func NewAttachmentOriginal() *AttachmentOriginal {
	return &AttachmentOriginal{}
}

func (r AttachmentOriginal) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["download"] = r.Download
	out["sign"] = r.Sign
	out["userID"] = r.UserID
	out["name"] = r.Name
	out["attachmentID"] = r.AttachmentID

	return out
}

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
		r.Download = parseBool(val)
	}
	if val, ok := get["sign"]; ok {
		r.Sign = val
	}
	if val, ok := get["userID"]; ok {
		r.UserID = parseUInt64(val)
	}
	r.Name = chi.URLParam(req, "name")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))

	return err
}

var _ RequestFiller = NewAttachmentOriginal()

// Attachment preview request parameters
type AttachmentPreview struct {
	Ext          string
	AttachmentID uint64 `json:",string"`
	Sign         string
	UserID       uint64 `json:",string"`
}

func NewAttachmentPreview() *AttachmentPreview {
	return &AttachmentPreview{}
}

func (r AttachmentPreview) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["ext"] = r.Ext
	out["attachmentID"] = r.AttachmentID
	out["sign"] = r.Sign
	out["userID"] = r.UserID

	return out
}

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

	r.Ext = chi.URLParam(req, "ext")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))
	if val, ok := get["sign"]; ok {
		r.Sign = val
	}
	if val, ok := get["userID"]; ok {
		r.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewAttachmentPreview()
