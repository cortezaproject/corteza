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

// Page list request parameters
type PageList struct {
	SelfID      uint64 `json:",string"`
	Query       string
	Page        uint
	PerPage     uint
	NamespaceID uint64 `json:",string"`
}

func NewPageList() *PageList {
	return &PageList{}
}

func (r PageList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["selfID"] = r.SelfID

	out["query"] = r.Query

	out["page"] = r.Page

	out["perPage"] = r.PerPage

	out["namespaceID"] = r.NamespaceID

	return out
}

func (pReq *PageList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	if val, ok := get["selfID"]; ok {

		pReq.SelfID = parseUInt64(val)
	}
	if val, ok := get["query"]; ok {

		pReq.Query = val
	}
	if val, ok := get["page"]; ok {

		pReq.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {

		pReq.PerPage = parseUint(val)
	}
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageList()

// Page create request parameters
type PageCreate struct {
	SelfID      uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	Title       string
	Description string
	Visible     bool
	Blocks      sqlxTypes.JSONText
	NamespaceID uint64 `json:",string"`
}

func NewPageCreate() *PageCreate {
	return &PageCreate{}
}

func (r PageCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["selfID"] = r.SelfID

	out["moduleID"] = r.ModuleID

	out["title"] = r.Title

	out["description"] = r.Description

	out["visible"] = r.Visible

	out["blocks"] = r.Blocks

	out["namespaceID"] = r.NamespaceID

	return out
}

func (pReq *PageCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	if val, ok := post["selfID"]; ok {

		pReq.SelfID = parseUInt64(val)
	}
	if val, ok := post["moduleID"]; ok {

		pReq.ModuleID = parseUInt64(val)
	}
	if val, ok := post["title"]; ok {

		pReq.Title = val
	}
	if val, ok := post["description"]; ok {

		pReq.Description = val
	}
	if val, ok := post["visible"]; ok {

		pReq.Visible = parseBool(val)
	}
	if val, ok := post["blocks"]; ok {

		if pReq.Blocks, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageCreate()

// Page read request parameters
type PageRead struct {
	PageID      uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewPageRead() *PageRead {
	return &PageRead{}
}

func (r PageRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID

	out["namespaceID"] = r.NamespaceID

	return out
}

func (pReq *PageRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	pReq.PageID = parseUInt64(chi.URLParam(r, "pageID"))
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageRead()

// Page tree request parameters
type PageTree struct {
	NamespaceID uint64 `json:",string"`
}

func NewPageTree() *PageTree {
	return &PageTree{}
}

func (r PageTree) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID

	return out
}

func (pReq *PageTree) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageTree()

// Page update request parameters
type PageUpdate struct {
	PageID      uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	SelfID      uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	Title       string
	Description string
	Visible     bool
	Blocks      sqlxTypes.JSONText
}

func NewPageUpdate() *PageUpdate {
	return &PageUpdate{}
}

func (r PageUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID

	out["namespaceID"] = r.NamespaceID

	out["selfID"] = r.SelfID

	out["moduleID"] = r.ModuleID

	out["title"] = r.Title

	out["description"] = r.Description

	out["visible"] = r.Visible

	out["blocks"] = r.Blocks

	return out
}

func (pReq *PageUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	pReq.PageID = parseUInt64(chi.URLParam(r, "pageID"))
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))
	if val, ok := post["selfID"]; ok {

		pReq.SelfID = parseUInt64(val)
	}
	if val, ok := post["moduleID"]; ok {

		pReq.ModuleID = parseUInt64(val)
	}
	if val, ok := post["title"]; ok {

		pReq.Title = val
	}
	if val, ok := post["description"]; ok {

		pReq.Description = val
	}
	if val, ok := post["visible"]; ok {

		pReq.Visible = parseBool(val)
	}
	if val, ok := post["blocks"]; ok {

		if pReq.Blocks, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewPageUpdate()

// Page reorder request parameters
type PageReorder struct {
	SelfID      uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	PageIDs     []string
}

func NewPageReorder() *PageReorder {
	return &PageReorder{}
}

func (r PageReorder) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["selfID"] = r.SelfID

	out["namespaceID"] = r.NamespaceID

	out["pageIDs"] = r.PageIDs

	return out
}

func (pReq *PageReorder) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	pReq.SelfID = parseUInt64(chi.URLParam(r, "selfID"))
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageReorder()

// Page delete request parameters
type PageDelete struct {
	PageID      uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewPageDelete() *PageDelete {
	return &PageDelete{}
}

func (r PageDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID

	out["namespaceID"] = r.NamespaceID

	return out
}

func (pReq *PageDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

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

	pReq.PageID = parseUInt64(chi.URLParam(r, "pageID"))
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewPageDelete()

// Page upload request parameters
type PageUpload struct {
	PageID      uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	Upload      *multipart.FileHeader
}

func NewPageUpload() *PageUpload {
	return &PageUpload{}
}

func (r PageUpload) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["pageID"] = r.PageID

	out["namespaceID"] = r.NamespaceID

	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	return out
}

func (pReq *PageUpload) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(pReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseMultipartForm(32 << 20); err != nil {
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

	pReq.PageID = parseUInt64(chi.URLParam(r, "pageID"))
	pReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))
	if _, pReq.Upload, err = r.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewPageUpload()
