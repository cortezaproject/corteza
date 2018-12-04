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
	"encoding/json"
	"github.com/crusttech/crust/internal/rbac"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}
var _ = multipart.FileHeader{}
var _ = rbac.Operation{}

// Page list request parameters
type PageList struct {
	SelfID uint64 `json:",string"`
}

func NewPageList() *PageList {
	return &PageList{}
}

func (p *PageList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

		p.SelfID = parseUInt64(val)
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
	Blocks      types.JSONText
}

func NewPageCreate() *PageCreate {
	return &PageCreate{}
}

func (p *PageCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

		p.SelfID = parseUInt64(val)
	}
	if val, ok := post["moduleID"]; ok {

		p.ModuleID = parseUInt64(val)
	}
	if val, ok := post["title"]; ok {

		p.Title = val
	}
	if val, ok := post["description"]; ok {

		p.Description = val
	}
	if val, ok := post["visible"]; ok {

		p.Visible = parseBool(val)
	}
	if val, ok := post["blocks"]; ok {

		if p.Blocks, err = parseJSONText(val); err != nil {
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

func (p *PageRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

	p.PageID = parseUInt64(chi.URLParam(r, "pageID"))

	return err
}

var _ RequestFiller = NewPageRead()

// Page tree request parameters
type PageTree struct {
}

func NewPageTree() *PageTree {
	return &PageTree{}
}

func (p *PageTree) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

// Page edit request parameters
type PageEdit struct {
	PageID      uint64 `json:",string"`
	SelfID      uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	Title       string
	Description string
	Visible     bool
	Blocks      types.JSONText
}

func NewPageEdit() *PageEdit {
	return &PageEdit{}
}

func (p *PageEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

	p.PageID = parseUInt64(chi.URLParam(r, "pageID"))
	if val, ok := post["selfID"]; ok {

		p.SelfID = parseUInt64(val)
	}
	if val, ok := post["moduleID"]; ok {

		p.ModuleID = parseUInt64(val)
	}
	if val, ok := post["title"]; ok {

		p.Title = val
	}
	if val, ok := post["description"]; ok {

		p.Description = val
	}
	if val, ok := post["visible"]; ok {

		p.Visible = parseBool(val)
	}
	if val, ok := post["blocks"]; ok {

		if p.Blocks, err = parseJSONText(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewPageEdit()

// Page reorder request parameters
type PageReorder struct {
	SelfID  uint64 `json:",string"`
	PageIDs []string
}

func NewPageReorder() *PageReorder {
	return &PageReorder{}
}

func (p *PageReorder) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

	p.SelfID = parseUInt64(chi.URLParam(r, "selfID"))

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

func (p *PageDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(p)

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

	p.PageID = parseUInt64(chi.URLParam(r, "pageID"))

	return err
}

var _ RequestFiller = NewPageDelete()
