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

// AttachmentList request parameters
type AttachmentList struct {
	hasPageID bool
	rawPageID string
	PageID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasRecordID bool
	rawRecordID string
	RecordID    uint64 `json:",string"`

	hasFieldName bool
	rawFieldName string
	FieldName    string

	hasLimit bool
	rawLimit string
	Limit    uint

	hasOffset bool
	rawOffset string
	Offset    uint

	hasPage bool
	rawPage string
	Page    uint

	hasPerPage bool
	rawPerPage string
	PerPage    uint

	hasSign bool
	rawSign string
	Sign    string

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasKind bool
	rawKind string
	Kind    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewAttachmentList request
func NewAttachmentList() *AttachmentList {
	return &AttachmentList{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID
	out["moduleID"] = r.ModuleID
	out["recordID"] = r.RecordID
	out["fieldName"] = r.FieldName
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sign"] = r.Sign
	out["userID"] = r.UserID
	out["kind"] = r.Kind
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *AttachmentList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["pageID"]; ok {
		r.hasPageID = true
		r.rawPageID = val
		r.PageID = parseUInt64(val)
	}
	if val, ok := get["moduleID"]; ok {
		r.hasModuleID = true
		r.rawModuleID = val
		r.ModuleID = parseUInt64(val)
	}
	if val, ok := get["recordID"]; ok {
		r.hasRecordID = true
		r.rawRecordID = val
		r.RecordID = parseUInt64(val)
	}
	if val, ok := get["fieldName"]; ok {
		r.hasFieldName = true
		r.rawFieldName = val
		r.FieldName = val
	}
	if val, ok := get["limit"]; ok {
		r.hasLimit = true
		r.rawLimit = val
		r.Limit = parseUint(val)
	}
	if val, ok := get["offset"]; ok {
		r.hasOffset = true
		r.rawOffset = val
		r.Offset = parseUint(val)
	}
	if val, ok := get["page"]; ok {
		r.hasPage = true
		r.rawPage = val
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.hasPerPage = true
		r.rawPerPage = val
		r.PerPage = parseUint(val)
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
	r.hasKind = true
	r.rawKind = chi.URLParam(req, "kind")
	r.Kind = chi.URLParam(req, "kind")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewAttachmentList()

// AttachmentRead request parameters
type AttachmentRead struct {
	hasAttachmentID bool
	rawAttachmentID string
	AttachmentID    uint64 `json:",string"`

	hasKind bool
	rawKind string
	Kind    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasSign bool
	rawSign string
	Sign    string

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewAttachmentRead request
func NewAttachmentRead() *AttachmentRead {
	return &AttachmentRead{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["attachmentID"] = r.AttachmentID
	out["kind"] = r.Kind
	out["namespaceID"] = r.NamespaceID
	out["sign"] = r.Sign
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *AttachmentRead) Fill(req *http.Request) (err error) {
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

	r.hasAttachmentID = true
	r.rawAttachmentID = chi.URLParam(req, "attachmentID")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))
	r.hasKind = true
	r.rawKind = chi.URLParam(req, "kind")
	r.Kind = chi.URLParam(req, "kind")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
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

var _ RequestFiller = NewAttachmentRead()

// AttachmentDelete request parameters
type AttachmentDelete struct {
	hasAttachmentID bool
	rawAttachmentID string
	AttachmentID    uint64 `json:",string"`

	hasKind bool
	rawKind string
	Kind    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasSign bool
	rawSign string
	Sign    string

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewAttachmentDelete request
func NewAttachmentDelete() *AttachmentDelete {
	return &AttachmentDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r AttachmentDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["attachmentID"] = r.AttachmentID
	out["kind"] = r.Kind
	out["namespaceID"] = r.NamespaceID
	out["sign"] = r.Sign
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *AttachmentDelete) Fill(req *http.Request) (err error) {
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

	r.hasAttachmentID = true
	r.rawAttachmentID = chi.URLParam(req, "attachmentID")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))
	r.hasKind = true
	r.rawKind = chi.URLParam(req, "kind")
	r.Kind = chi.URLParam(req, "kind")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
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

var _ RequestFiller = NewAttachmentDelete()

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

	hasAttachmentID bool
	rawAttachmentID string
	AttachmentID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasKind bool
	rawKind string
	Kind    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
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
	out["attachmentID"] = r.AttachmentID
	out["name"] = r.Name
	out["kind"] = r.Kind
	out["namespaceID"] = r.NamespaceID

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
	r.hasAttachmentID = true
	r.rawAttachmentID = chi.URLParam(req, "attachmentID")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))
	r.hasName = true
	r.rawName = chi.URLParam(req, "name")
	r.Name = chi.URLParam(req, "name")
	r.hasKind = true
	r.rawKind = chi.URLParam(req, "kind")
	r.Kind = chi.URLParam(req, "kind")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewAttachmentOriginal()

// AttachmentPreview request parameters
type AttachmentPreview struct {
	hasAttachmentID bool
	rawAttachmentID string
	AttachmentID    uint64 `json:",string"`

	hasExt bool
	rawExt string
	Ext    string

	hasKind bool
	rawKind string
	Kind    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

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

	out["attachmentID"] = r.AttachmentID
	out["ext"] = r.Ext
	out["kind"] = r.Kind
	out["namespaceID"] = r.NamespaceID
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

	r.hasAttachmentID = true
	r.rawAttachmentID = chi.URLParam(req, "attachmentID")
	r.AttachmentID = parseUInt64(chi.URLParam(req, "attachmentID"))
	r.hasExt = true
	r.rawExt = chi.URLParam(req, "ext")
	r.Ext = chi.URLParam(req, "ext")
	r.hasKind = true
	r.rawKind = chi.URLParam(req, "kind")
	r.Kind = chi.URLParam(req, "kind")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
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

// HasPageID returns true if pageID was set
func (r *AttachmentList) HasPageID() bool {
	return r.hasPageID
}

// RawPageID returns raw value of pageID parameter
func (r *AttachmentList) RawPageID() string {
	return r.rawPageID
}

// GetPageID returns casted value of  pageID parameter
func (r *AttachmentList) GetPageID() uint64 {
	return r.PageID
}

// HasModuleID returns true if moduleID was set
func (r *AttachmentList) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *AttachmentList) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *AttachmentList) GetModuleID() uint64 {
	return r.ModuleID
}

// HasRecordID returns true if recordID was set
func (r *AttachmentList) HasRecordID() bool {
	return r.hasRecordID
}

// RawRecordID returns raw value of recordID parameter
func (r *AttachmentList) RawRecordID() string {
	return r.rawRecordID
}

// GetRecordID returns casted value of  recordID parameter
func (r *AttachmentList) GetRecordID() uint64 {
	return r.RecordID
}

// HasFieldName returns true if fieldName was set
func (r *AttachmentList) HasFieldName() bool {
	return r.hasFieldName
}

// RawFieldName returns raw value of fieldName parameter
func (r *AttachmentList) RawFieldName() string {
	return r.rawFieldName
}

// GetFieldName returns casted value of  fieldName parameter
func (r *AttachmentList) GetFieldName() string {
	return r.FieldName
}

// HasLimit returns true if limit was set
func (r *AttachmentList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *AttachmentList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *AttachmentList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *AttachmentList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *AttachmentList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *AttachmentList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *AttachmentList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *AttachmentList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *AttachmentList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *AttachmentList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *AttachmentList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *AttachmentList) GetPerPage() uint {
	return r.PerPage
}

// HasSign returns true if sign was set
func (r *AttachmentList) HasSign() bool {
	return r.hasSign
}

// RawSign returns raw value of sign parameter
func (r *AttachmentList) RawSign() string {
	return r.rawSign
}

// GetSign returns casted value of  sign parameter
func (r *AttachmentList) GetSign() string {
	return r.Sign
}

// HasUserID returns true if userID was set
func (r *AttachmentList) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *AttachmentList) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *AttachmentList) GetUserID() uint64 {
	return r.UserID
}

// HasKind returns true if kind was set
func (r *AttachmentList) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *AttachmentList) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *AttachmentList) GetKind() string {
	return r.Kind
}

// HasNamespaceID returns true if namespaceID was set
func (r *AttachmentList) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *AttachmentList) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *AttachmentList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasAttachmentID returns true if attachmentID was set
func (r *AttachmentRead) HasAttachmentID() bool {
	return r.hasAttachmentID
}

// RawAttachmentID returns raw value of attachmentID parameter
func (r *AttachmentRead) RawAttachmentID() string {
	return r.rawAttachmentID
}

// GetAttachmentID returns casted value of  attachmentID parameter
func (r *AttachmentRead) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// HasKind returns true if kind was set
func (r *AttachmentRead) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *AttachmentRead) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *AttachmentRead) GetKind() string {
	return r.Kind
}

// HasNamespaceID returns true if namespaceID was set
func (r *AttachmentRead) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *AttachmentRead) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *AttachmentRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasSign returns true if sign was set
func (r *AttachmentRead) HasSign() bool {
	return r.hasSign
}

// RawSign returns raw value of sign parameter
func (r *AttachmentRead) RawSign() string {
	return r.rawSign
}

// GetSign returns casted value of  sign parameter
func (r *AttachmentRead) GetSign() string {
	return r.Sign
}

// HasUserID returns true if userID was set
func (r *AttachmentRead) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *AttachmentRead) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *AttachmentRead) GetUserID() uint64 {
	return r.UserID
}

// HasAttachmentID returns true if attachmentID was set
func (r *AttachmentDelete) HasAttachmentID() bool {
	return r.hasAttachmentID
}

// RawAttachmentID returns raw value of attachmentID parameter
func (r *AttachmentDelete) RawAttachmentID() string {
	return r.rawAttachmentID
}

// GetAttachmentID returns casted value of  attachmentID parameter
func (r *AttachmentDelete) GetAttachmentID() uint64 {
	return r.AttachmentID
}

// HasKind returns true if kind was set
func (r *AttachmentDelete) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *AttachmentDelete) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *AttachmentDelete) GetKind() string {
	return r.Kind
}

// HasNamespaceID returns true if namespaceID was set
func (r *AttachmentDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *AttachmentDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *AttachmentDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasSign returns true if sign was set
func (r *AttachmentDelete) HasSign() bool {
	return r.hasSign
}

// RawSign returns raw value of sign parameter
func (r *AttachmentDelete) RawSign() string {
	return r.rawSign
}

// GetSign returns casted value of  sign parameter
func (r *AttachmentDelete) GetSign() string {
	return r.Sign
}

// HasUserID returns true if userID was set
func (r *AttachmentDelete) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *AttachmentDelete) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *AttachmentDelete) GetUserID() uint64 {
	return r.UserID
}

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

// HasKind returns true if kind was set
func (r *AttachmentOriginal) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *AttachmentOriginal) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *AttachmentOriginal) GetKind() string {
	return r.Kind
}

// HasNamespaceID returns true if namespaceID was set
func (r *AttachmentOriginal) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *AttachmentOriginal) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *AttachmentOriginal) GetNamespaceID() uint64 {
	return r.NamespaceID
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

// HasKind returns true if kind was set
func (r *AttachmentPreview) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *AttachmentPreview) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *AttachmentPreview) GetKind() string {
	return r.Kind
}

// HasNamespaceID returns true if namespaceID was set
func (r *AttachmentPreview) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *AttachmentPreview) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *AttachmentPreview) GetNamespaceID() uint64 {
	return r.NamespaceID
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
