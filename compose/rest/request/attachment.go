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

// Attachment list request parameters
type AttachmentList struct {
	PageID      uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	RecordID    uint64 `json:",string"`
	FieldName   string
	Page        uint
	PerPage     uint
	Sign        string
	UserID      uint64 `json:",string"`
	Kind        string
	NamespaceID uint64 `json:",string"`
}

func NewAttachmentList() *AttachmentList {
	return &AttachmentList{}
}

func (r AttachmentList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID

	out["moduleID"] = r.ModuleID

	out["recordID"] = r.RecordID

	out["fieldName"] = r.FieldName

	out["page"] = r.Page

	out["perPage"] = r.PerPage

	out["sign"] = r.Sign

	out["userID"] = r.UserID

	out["kind"] = r.Kind

	out["namespaceID"] = r.NamespaceID

	return out
}

func (aReq *AttachmentList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(aReq)

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

	if val, ok := get["pageID"]; ok {

		aReq.PageID = parseUInt64(val)
	}
	if val, ok := get["moduleID"]; ok {

		aReq.ModuleID = parseUInt64(val)
	}
	if val, ok := get["recordID"]; ok {

		aReq.RecordID = parseUInt64(val)
	}
	if val, ok := get["fieldName"]; ok {

		aReq.FieldName = val
	}
	if val, ok := get["page"]; ok {

		aReq.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {

		aReq.PerPage = parseUint(val)
	}
	if val, ok := get["sign"]; ok {

		aReq.Sign = val
	}
	if val, ok := get["userID"]; ok {

		aReq.UserID = parseUInt64(val)
	}
	aReq.Kind = chi.URLParam(r, "kind")
	aReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewAttachmentList()

// Attachment read request parameters
type AttachmentRead struct {
	AttachmentID uint64 `json:",string"`
	Kind         string
	NamespaceID  uint64 `json:",string"`
	Sign         string
	UserID       uint64 `json:",string"`
}

func NewAttachmentRead() *AttachmentRead {
	return &AttachmentRead{}
}

func (r AttachmentRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["attachmentID"] = r.AttachmentID

	out["kind"] = r.Kind

	out["namespaceID"] = r.NamespaceID

	out["sign"] = r.Sign

	out["userID"] = r.UserID

	return out
}

func (aReq *AttachmentRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(aReq)

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

	aReq.AttachmentID = parseUInt64(chi.URLParam(r, "attachmentID"))
	aReq.Kind = chi.URLParam(r, "kind")
	aReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))
	if val, ok := get["sign"]; ok {

		aReq.Sign = val
	}
	if val, ok := get["userID"]; ok {

		aReq.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewAttachmentRead()

// Attachment delete request parameters
type AttachmentDelete struct {
	AttachmentID uint64 `json:",string"`
	Kind         string
	NamespaceID  uint64 `json:",string"`
	Sign         string
	UserID       uint64 `json:",string"`
}

func NewAttachmentDelete() *AttachmentDelete {
	return &AttachmentDelete{}
}

func (r AttachmentDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["attachmentID"] = r.AttachmentID

	out["kind"] = r.Kind

	out["namespaceID"] = r.NamespaceID

	out["sign"] = r.Sign

	out["userID"] = r.UserID

	return out
}

func (aReq *AttachmentDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(aReq)

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

	aReq.AttachmentID = parseUInt64(chi.URLParam(r, "attachmentID"))
	aReq.Kind = chi.URLParam(r, "kind")
	aReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))
	if val, ok := get["sign"]; ok {

		aReq.Sign = val
	}
	if val, ok := get["userID"]; ok {

		aReq.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewAttachmentDelete()

// Attachment original request parameters
type AttachmentOriginal struct {
	Download     bool
	Sign         string
	UserID       uint64 `json:",string"`
	AttachmentID uint64 `json:",string"`
	Name         string
	Kind         string
	NamespaceID  uint64 `json:",string"`
}

func NewAttachmentOriginal() *AttachmentOriginal {
	return &AttachmentOriginal{}
}

func (r AttachmentOriginal) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["download"] = r.Download

	out["sign"] = r.Sign

	out["userID"] = r.UserID

	out["attachmentID"] = r.AttachmentID

	out["name"] = r.Name

	out["kind"] = r.Kind

	out["namespaceID"] = r.NamespaceID

	return out
}

func (aReq *AttachmentOriginal) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(aReq)

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

	if val, ok := get["download"]; ok {

		aReq.Download = parseBool(val)
	}
	if val, ok := get["sign"]; ok {

		aReq.Sign = val
	}
	if val, ok := get["userID"]; ok {

		aReq.UserID = parseUInt64(val)
	}
	aReq.AttachmentID = parseUInt64(chi.URLParam(r, "attachmentID"))
	aReq.Name = chi.URLParam(r, "name")
	aReq.Kind = chi.URLParam(r, "kind")
	aReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewAttachmentOriginal()

// Attachment preview request parameters
type AttachmentPreview struct {
	AttachmentID uint64 `json:",string"`
	Ext          string
	Kind         string
	NamespaceID  uint64 `json:",string"`
	Sign         string
	UserID       uint64 `json:",string"`
}

func NewAttachmentPreview() *AttachmentPreview {
	return &AttachmentPreview{}
}

func (r AttachmentPreview) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["attachmentID"] = r.AttachmentID

	out["ext"] = r.Ext

	out["kind"] = r.Kind

	out["namespaceID"] = r.NamespaceID

	out["sign"] = r.Sign

	out["userID"] = r.UserID

	return out
}

func (aReq *AttachmentPreview) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(aReq)

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

	aReq.AttachmentID = parseUInt64(chi.URLParam(r, "attachmentID"))
	aReq.Ext = chi.URLParam(r, "ext")
	aReq.Kind = chi.URLParam(r, "kind")
	aReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))
	if val, ok := get["sign"]; ok {

		aReq.Sign = val
	}
	if val, ok := get["userID"]; ok {

		aReq.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewAttachmentPreview()
