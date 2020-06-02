package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `page.go`, `page.util.go` or `page_test.go` to
	implement your API calls, helper functions and tests. The file `page.go`
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

	sqlxTypes "github.com/jmoiron/sqlx/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// PageList request parameters
type PageList struct {
	hasSelfID bool
	rawSelfID string
	SelfID    uint64 `json:",string"`

	hasQuery bool
	rawQuery string
	Query    string

	hasHandle bool
	rawHandle string
	Handle    string

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

	hasSort bool
	rawSort string
	Sort    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewPageList request
func NewPageList() *PageList {
	return &PageList{}
}

// Auditable returns all auditable/loggable parameters
func (r PageList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["selfID"] = r.SelfID
	out["query"] = r.Query
	out["handle"] = r.Handle
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *PageList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["selfID"]; ok {
		r.hasSelfID = true
		r.rawSelfID = val
		r.SelfID = parseUInt64(val)
	}
	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}
	if val, ok := get["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
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
	if val, ok := get["sort"]; ok {
		r.hasSort = true
		r.rawSort = val
		r.Sort = val
	}
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageList()

// PageCreate request parameters
type PageCreate struct {
	hasSelfID bool
	rawSelfID string
	SelfID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasTitle bool
	rawTitle string
	Title    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasDescription bool
	rawDescription string
	Description    string

	hasWeight bool
	rawWeight string
	Weight    int

	hasVisible bool
	rawVisible string
	Visible    bool

	hasBlocks bool
	rawBlocks string
	Blocks    sqlxTypes.JSONText

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewPageCreate request
func NewPageCreate() *PageCreate {
	return &PageCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["selfID"] = r.SelfID
	out["moduleID"] = r.ModuleID
	out["title"] = r.Title
	out["handle"] = r.Handle
	out["description"] = r.Description
	out["weight"] = r.Weight
	out["visible"] = r.Visible
	out["blocks"] = r.Blocks
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *PageCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["selfID"]; ok {
		r.hasSelfID = true
		r.rawSelfID = val
		r.SelfID = parseUInt64(val)
	}
	if val, ok := post["moduleID"]; ok {
		r.hasModuleID = true
		r.rawModuleID = val
		r.ModuleID = parseUInt64(val)
	}
	if val, ok := post["title"]; ok {
		r.hasTitle = true
		r.rawTitle = val
		r.Title = val
	}
	if val, ok := post["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
	}
	if val, ok := post["description"]; ok {
		r.hasDescription = true
		r.rawDescription = val
		r.Description = val
	}
	if val, ok := post["weight"]; ok {
		r.hasWeight = true
		r.rawWeight = val
		r.Weight = parseInt(val)
	}
	if val, ok := post["visible"]; ok {
		r.hasVisible = true
		r.rawVisible = val
		r.Visible = parseBool(val)
	}
	if val, ok := post["blocks"]; ok {
		r.hasBlocks = true
		r.rawBlocks = val

		if r.Blocks, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageCreate()

// PageRead request parameters
type PageRead struct {
	hasPageID bool
	rawPageID string
	PageID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewPageRead request
func NewPageRead() *PageRead {
	return &PageRead{}
}

// Auditable returns all auditable/loggable parameters
func (r PageRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *PageRead) Fill(req *http.Request) (err error) {
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

	r.hasPageID = true
	r.rawPageID = chi.URLParam(req, "pageID")
	r.PageID = parseUInt64(chi.URLParam(req, "pageID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageRead()

// PageTree request parameters
type PageTree struct {
	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewPageTree request
func NewPageTree() *PageTree {
	return &PageTree{}
}

// Auditable returns all auditable/loggable parameters
func (r PageTree) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *PageTree) Fill(req *http.Request) (err error) {
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

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageTree()

// PageUpdate request parameters
type PageUpdate struct {
	hasPageID bool
	rawPageID string
	PageID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasSelfID bool
	rawSelfID string
	SelfID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasTitle bool
	rawTitle string
	Title    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasDescription bool
	rawDescription string
	Description    string

	hasWeight bool
	rawWeight string
	Weight    int

	hasVisible bool
	rawVisible string
	Visible    bool

	hasBlocks bool
	rawBlocks string
	Blocks    sqlxTypes.JSONText
}

// NewPageUpdate request
func NewPageUpdate() *PageUpdate {
	return &PageUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID
	out["namespaceID"] = r.NamespaceID
	out["selfID"] = r.SelfID
	out["moduleID"] = r.ModuleID
	out["title"] = r.Title
	out["handle"] = r.Handle
	out["description"] = r.Description
	out["weight"] = r.Weight
	out["visible"] = r.Visible
	out["blocks"] = r.Blocks

	return out
}

// Fill processes request and fills internal variables
func (r *PageUpdate) Fill(req *http.Request) (err error) {
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

	r.hasPageID = true
	r.rawPageID = chi.URLParam(req, "pageID")
	r.PageID = parseUInt64(chi.URLParam(req, "pageID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["selfID"]; ok {
		r.hasSelfID = true
		r.rawSelfID = val
		r.SelfID = parseUInt64(val)
	}
	if val, ok := post["moduleID"]; ok {
		r.hasModuleID = true
		r.rawModuleID = val
		r.ModuleID = parseUInt64(val)
	}
	if val, ok := post["title"]; ok {
		r.hasTitle = true
		r.rawTitle = val
		r.Title = val
	}
	if val, ok := post["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
	}
	if val, ok := post["description"]; ok {
		r.hasDescription = true
		r.rawDescription = val
		r.Description = val
	}
	if val, ok := post["weight"]; ok {
		r.hasWeight = true
		r.rawWeight = val
		r.Weight = parseInt(val)
	}
	if val, ok := post["visible"]; ok {
		r.hasVisible = true
		r.rawVisible = val
		r.Visible = parseBool(val)
	}
	if val, ok := post["blocks"]; ok {
		r.hasBlocks = true
		r.rawBlocks = val

		if r.Blocks, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewPageUpdate()

// PageReorder request parameters
type PageReorder struct {
	hasSelfID bool
	rawSelfID string
	SelfID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasPageIDs bool
	rawPageIDs []string
	PageIDs    []string
}

// NewPageReorder request
func NewPageReorder() *PageReorder {
	return &PageReorder{}
}

// Auditable returns all auditable/loggable parameters
func (r PageReorder) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["selfID"] = r.SelfID
	out["namespaceID"] = r.NamespaceID
	out["pageIDs"] = r.PageIDs

	return out
}

// Fill processes request and fills internal variables
func (r *PageReorder) Fill(req *http.Request) (err error) {
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

	r.hasSelfID = true
	r.rawSelfID = chi.URLParam(req, "selfID")
	r.SelfID = parseUInt64(chi.URLParam(req, "selfID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	if val, ok := req.Form["pageIDs"]; ok {
		r.hasPageIDs = true
		r.rawPageIDs = val
		r.PageIDs = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewPageReorder()

// PageDelete request parameters
type PageDelete struct {
	hasPageID bool
	rawPageID string
	PageID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewPageDelete request
func NewPageDelete() *PageDelete {
	return &PageDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r PageDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *PageDelete) Fill(req *http.Request) (err error) {
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

	r.hasPageID = true
	r.rawPageID = chi.URLParam(req, "pageID")
	r.PageID = parseUInt64(chi.URLParam(req, "pageID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageDelete()

// PageUpload request parameters
type PageUpload struct {
	hasPageID bool
	rawPageID string
	PageID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasUpload bool
	rawUpload string
	Upload    *multipart.FileHeader
}

// NewPageUpload request
func NewPageUpload() *PageUpload {
	return &PageUpload{}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpload) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID
	out["namespaceID"] = r.NamespaceID
	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	return out
}

// Fill processes request and fills internal variables
func (r *PageUpload) Fill(req *http.Request) (err error) {
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

	r.hasPageID = true
	r.rawPageID = chi.URLParam(req, "pageID")
	r.PageID = parseUInt64(chi.URLParam(req, "pageID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	return err
}

var _ RequestFiller = NewPageUpload()

// PageTriggerScript request parameters
type PageTriggerScript struct {
	hasPageID bool
	rawPageID string
	PageID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string
}

// NewPageTriggerScript request
func NewPageTriggerScript() *PageTriggerScript {
	return &PageTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r PageTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID
	out["namespaceID"] = r.NamespaceID
	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
func (r *PageTriggerScript) Fill(req *http.Request) (err error) {
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

	r.hasPageID = true
	r.rawPageID = chi.URLParam(req, "pageID")
	r.PageID = parseUInt64(chi.URLParam(req, "pageID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewPageTriggerScript()

// HasSelfID returns true if selfID was set
func (r *PageList) HasSelfID() bool {
	return r.hasSelfID
}

// RawSelfID returns raw value of selfID parameter
func (r *PageList) RawSelfID() string {
	return r.rawSelfID
}

// GetSelfID returns casted value of  selfID parameter
func (r *PageList) GetSelfID() uint64 {
	return r.SelfID
}

// HasQuery returns true if query was set
func (r *PageList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *PageList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *PageList) GetQuery() string {
	return r.Query
}

// HasHandle returns true if handle was set
func (r *PageList) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *PageList) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *PageList) GetHandle() string {
	return r.Handle
}

// HasLimit returns true if limit was set
func (r *PageList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *PageList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *PageList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *PageList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *PageList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *PageList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *PageList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *PageList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *PageList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *PageList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *PageList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *PageList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *PageList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *PageList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *PageList) GetSort() string {
	return r.Sort
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageList) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageList) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasSelfID returns true if selfID was set
func (r *PageCreate) HasSelfID() bool {
	return r.hasSelfID
}

// RawSelfID returns raw value of selfID parameter
func (r *PageCreate) RawSelfID() string {
	return r.rawSelfID
}

// GetSelfID returns casted value of  selfID parameter
func (r *PageCreate) GetSelfID() uint64 {
	return r.SelfID
}

// HasModuleID returns true if moduleID was set
func (r *PageCreate) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *PageCreate) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *PageCreate) GetModuleID() uint64 {
	return r.ModuleID
}

// HasTitle returns true if title was set
func (r *PageCreate) HasTitle() bool {
	return r.hasTitle
}

// RawTitle returns raw value of title parameter
func (r *PageCreate) RawTitle() string {
	return r.rawTitle
}

// GetTitle returns casted value of  title parameter
func (r *PageCreate) GetTitle() string {
	return r.Title
}

// HasHandle returns true if handle was set
func (r *PageCreate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *PageCreate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *PageCreate) GetHandle() string {
	return r.Handle
}

// HasDescription returns true if description was set
func (r *PageCreate) HasDescription() bool {
	return r.hasDescription
}

// RawDescription returns raw value of description parameter
func (r *PageCreate) RawDescription() string {
	return r.rawDescription
}

// GetDescription returns casted value of  description parameter
func (r *PageCreate) GetDescription() string {
	return r.Description
}

// HasWeight returns true if weight was set
func (r *PageCreate) HasWeight() bool {
	return r.hasWeight
}

// RawWeight returns raw value of weight parameter
func (r *PageCreate) RawWeight() string {
	return r.rawWeight
}

// GetWeight returns casted value of  weight parameter
func (r *PageCreate) GetWeight() int {
	return r.Weight
}

// HasVisible returns true if visible was set
func (r *PageCreate) HasVisible() bool {
	return r.hasVisible
}

// RawVisible returns raw value of visible parameter
func (r *PageCreate) RawVisible() string {
	return r.rawVisible
}

// GetVisible returns casted value of  visible parameter
func (r *PageCreate) GetVisible() bool {
	return r.Visible
}

// HasBlocks returns true if blocks was set
func (r *PageCreate) HasBlocks() bool {
	return r.hasBlocks
}

// RawBlocks returns raw value of blocks parameter
func (r *PageCreate) RawBlocks() string {
	return r.rawBlocks
}

// GetBlocks returns casted value of  blocks parameter
func (r *PageCreate) GetBlocks() sqlxTypes.JSONText {
	return r.Blocks
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageCreate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageCreate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasPageID returns true if pageID was set
func (r *PageRead) HasPageID() bool {
	return r.hasPageID
}

// RawPageID returns raw value of pageID parameter
func (r *PageRead) RawPageID() string {
	return r.rawPageID
}

// GetPageID returns casted value of  pageID parameter
func (r *PageRead) GetPageID() uint64 {
	return r.PageID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageRead) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageRead) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageTree) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageTree) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageTree) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasPageID returns true if pageID was set
func (r *PageUpdate) HasPageID() bool {
	return r.hasPageID
}

// RawPageID returns raw value of pageID parameter
func (r *PageUpdate) RawPageID() string {
	return r.rawPageID
}

// GetPageID returns casted value of  pageID parameter
func (r *PageUpdate) GetPageID() uint64 {
	return r.PageID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageUpdate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageUpdate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasSelfID returns true if selfID was set
func (r *PageUpdate) HasSelfID() bool {
	return r.hasSelfID
}

// RawSelfID returns raw value of selfID parameter
func (r *PageUpdate) RawSelfID() string {
	return r.rawSelfID
}

// GetSelfID returns casted value of  selfID parameter
func (r *PageUpdate) GetSelfID() uint64 {
	return r.SelfID
}

// HasModuleID returns true if moduleID was set
func (r *PageUpdate) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *PageUpdate) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *PageUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// HasTitle returns true if title was set
func (r *PageUpdate) HasTitle() bool {
	return r.hasTitle
}

// RawTitle returns raw value of title parameter
func (r *PageUpdate) RawTitle() string {
	return r.rawTitle
}

// GetTitle returns casted value of  title parameter
func (r *PageUpdate) GetTitle() string {
	return r.Title
}

// HasHandle returns true if handle was set
func (r *PageUpdate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *PageUpdate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *PageUpdate) GetHandle() string {
	return r.Handle
}

// HasDescription returns true if description was set
func (r *PageUpdate) HasDescription() bool {
	return r.hasDescription
}

// RawDescription returns raw value of description parameter
func (r *PageUpdate) RawDescription() string {
	return r.rawDescription
}

// GetDescription returns casted value of  description parameter
func (r *PageUpdate) GetDescription() string {
	return r.Description
}

// HasWeight returns true if weight was set
func (r *PageUpdate) HasWeight() bool {
	return r.hasWeight
}

// RawWeight returns raw value of weight parameter
func (r *PageUpdate) RawWeight() string {
	return r.rawWeight
}

// GetWeight returns casted value of  weight parameter
func (r *PageUpdate) GetWeight() int {
	return r.Weight
}

// HasVisible returns true if visible was set
func (r *PageUpdate) HasVisible() bool {
	return r.hasVisible
}

// RawVisible returns raw value of visible parameter
func (r *PageUpdate) RawVisible() string {
	return r.rawVisible
}

// GetVisible returns casted value of  visible parameter
func (r *PageUpdate) GetVisible() bool {
	return r.Visible
}

// HasBlocks returns true if blocks was set
func (r *PageUpdate) HasBlocks() bool {
	return r.hasBlocks
}

// RawBlocks returns raw value of blocks parameter
func (r *PageUpdate) RawBlocks() string {
	return r.rawBlocks
}

// GetBlocks returns casted value of  blocks parameter
func (r *PageUpdate) GetBlocks() sqlxTypes.JSONText {
	return r.Blocks
}

// HasSelfID returns true if selfID was set
func (r *PageReorder) HasSelfID() bool {
	return r.hasSelfID
}

// RawSelfID returns raw value of selfID parameter
func (r *PageReorder) RawSelfID() string {
	return r.rawSelfID
}

// GetSelfID returns casted value of  selfID parameter
func (r *PageReorder) GetSelfID() uint64 {
	return r.SelfID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageReorder) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageReorder) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageReorder) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasPageIDs returns true if pageIDs was set
func (r *PageReorder) HasPageIDs() bool {
	return r.hasPageIDs
}

// RawPageIDs returns raw value of pageIDs parameter
func (r *PageReorder) RawPageIDs() []string {
	return r.rawPageIDs
}

// GetPageIDs returns casted value of  pageIDs parameter
func (r *PageReorder) GetPageIDs() []string {
	return r.PageIDs
}

// HasPageID returns true if pageID was set
func (r *PageDelete) HasPageID() bool {
	return r.hasPageID
}

// RawPageID returns raw value of pageID parameter
func (r *PageDelete) RawPageID() string {
	return r.rawPageID
}

// GetPageID returns casted value of  pageID parameter
func (r *PageDelete) GetPageID() uint64 {
	return r.PageID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasPageID returns true if pageID was set
func (r *PageUpload) HasPageID() bool {
	return r.hasPageID
}

// RawPageID returns raw value of pageID parameter
func (r *PageUpload) RawPageID() string {
	return r.rawPageID
}

// GetPageID returns casted value of  pageID parameter
func (r *PageUpload) GetPageID() uint64 {
	return r.PageID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageUpload) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageUpload) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageUpload) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasUpload returns true if upload was set
func (r *PageUpload) HasUpload() bool {
	return r.hasUpload
}

// RawUpload returns raw value of upload parameter
func (r *PageUpload) RawUpload() string {
	return r.rawUpload
}

// GetUpload returns casted value of  upload parameter
func (r *PageUpload) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// HasPageID returns true if pageID was set
func (r *PageTriggerScript) HasPageID() bool {
	return r.hasPageID
}

// RawPageID returns raw value of pageID parameter
func (r *PageTriggerScript) RawPageID() string {
	return r.rawPageID
}

// GetPageID returns casted value of  pageID parameter
func (r *PageTriggerScript) GetPageID() uint64 {
	return r.PageID
}

// HasNamespaceID returns true if namespaceID was set
func (r *PageTriggerScript) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *PageTriggerScript) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *PageTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasScript returns true if script was set
func (r *PageTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *PageTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *PageTriggerScript) GetScript() string {
	return r.Script
}
