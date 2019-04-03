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
	SelfID uint64 `json:",string"`
}

func NewPageList() *PageList {
	return &PageList{}
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
}

func NewPageCreate() *PageCreate {
	return &PageCreate{}
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

	return err
}

var _ RequestFiller = NewPageCreate()

// Page read request parameters
type PageRead struct {
	PageID uint64 `json:",string"`
}

func NewPageRead() *PageRead {
	return &PageRead{}
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

	return err
}

var _ RequestFiller = NewPageRead()

// Page tree request parameters
type PageTree struct {
}

func NewPageTree() *PageTree {
	return &PageTree{}
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

	return err
}

var _ RequestFiller = NewPageTree()

// Page update request parameters
type PageUpdate struct {
	PageID      uint64 `json:",string"`
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
	SelfID  uint64 `json:",string"`
	PageIDs []string
}

func NewPageReorder() *PageReorder {
	return &PageReorder{}
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

	return err
}

var _ RequestFiller = NewPageReorder()

// Page delete request parameters
type PageDelete struct {
	PageID uint64 `json:",string"`
}

func NewPageDelete() *PageDelete {
	return &PageDelete{}
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

	return err
}

var _ RequestFiller = NewPageDelete()

// Page upload request parameters
type PageUpload struct {
	PageID uint64 `json:",string"`
	Upload *multipart.FileHeader
}

func NewPageUpload() *PageUpload {
	return &PageUpload{}
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
	if _, pReq.Upload, err = r.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewPageUpload()
