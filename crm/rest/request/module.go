package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Module list request parameters
type ModuleList struct {
	Query string
}

func NewModuleList() *ModuleList {
	return &ModuleList{}
}

func (m *ModuleList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	if val, ok := get["query"]; ok {

		m.Query = val
	}

	return err
}

var _ RequestFiller = NewModuleList()

// Module create request parameters
type ModuleCreate struct {
	Name   string
	Fields types.ModuleFieldSet
}

func NewModuleCreate() *ModuleCreate {
	return &ModuleCreate{}
}

func (m *ModuleCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	if val, ok := post["name"]; ok {

		m.Name = val
	}

	return err
}

var _ RequestFiller = NewModuleCreate()

// Module read request parameters
type ModuleRead struct {
	ModuleID uint64 `json:",string"`
}

func NewModuleRead() *ModuleRead {
	return &ModuleRead{}
}

func (m *ModuleRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleRead()

// Module edit request parameters
type ModuleEdit struct {
	ModuleID uint64 `json:",string"`
	Name     string
	Fields   types.ModuleFieldSet
}

func NewModuleEdit() *ModuleEdit {
	return &ModuleEdit{}
}

func (m *ModuleEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	if val, ok := post["name"]; ok {

		m.Name = val
	}

	return err
}

var _ RequestFiller = NewModuleEdit()

// Module delete request parameters
type ModuleDelete struct {
	ModuleID uint64 `json:",string"`
}

func NewModuleDelete() *ModuleDelete {
	return &ModuleDelete{}
}

func (m *ModuleDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleDelete()

// Module record/report request parameters
type ModuleRecordReport struct {
	Metrics    string
	Dimensions string
	Filter     string
	ModuleID   uint64 `json:",string"`
}

func NewModuleRecordReport() *ModuleRecordReport {
	return &ModuleRecordReport{}
}

func (m *ModuleRecordReport) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	if val, ok := get["metrics"]; ok {

		m.Metrics = val
	}
	if val, ok := get["dimensions"]; ok {

		m.Dimensions = val
	}
	if val, ok := get["filter"]; ok {

		m.Filter = val
	}
	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleRecordReport()

// Module record/list request parameters
type ModuleRecordList struct {
	Query    string
	Page     int
	PerPage  int
	Sort     string
	ModuleID uint64 `json:",string"`
}

func NewModuleRecordList() *ModuleRecordList {
	return &ModuleRecordList{}
}

func (m *ModuleRecordList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	if val, ok := get["query"]; ok {

		m.Query = val
	}
	if val, ok := get["page"]; ok {

		m.Page = parseInt(val)
	}
	if val, ok := get["perPage"]; ok {

		m.PerPage = parseInt(val)
	}
	if val, ok := get["sort"]; ok {

		m.Sort = val
	}
	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleRecordList()

// Module record/create request parameters
type ModuleRecordCreate struct {
	ModuleID uint64 `json:",string"`
	Fields   sqlxTypes.JSONText
}

func NewModuleRecordCreate() *ModuleRecordCreate {
	return &ModuleRecordCreate{}
}

func (m *ModuleRecordCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleRecordCreate()

// Module record/read request parameters
type ModuleRecordRead struct {
	ModuleID uint64 `json:",string"`
	RecordID uint64 `json:",string"`
}

func NewModuleRecordRead() *ModuleRecordRead {
	return &ModuleRecordRead{}
}

func (m *ModuleRecordRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	m.RecordID = parseUInt64(chi.URLParam(r, "recordID"))

	return err
}

var _ RequestFiller = NewModuleRecordRead()

// Module record/edit request parameters
type ModuleRecordEdit struct {
	ModuleID uint64 `json:",string"`
	RecordID uint64 `json:",string"`
	Fields   sqlxTypes.JSONText
}

func NewModuleRecordEdit() *ModuleRecordEdit {
	return &ModuleRecordEdit{}
}

func (m *ModuleRecordEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	m.RecordID = parseUInt64(chi.URLParam(r, "recordID"))
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleRecordEdit()

// Module record/delete request parameters
type ModuleRecordDelete struct {
	ModuleID uint64 `json:",string"`
	RecordID uint64 `json:",string"`
}

func NewModuleRecordDelete() *ModuleRecordDelete {
	return &ModuleRecordDelete{}
}

func (m *ModuleRecordDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(m)

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

	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	m.RecordID = parseUInt64(chi.URLParam(r, "recordID"))

	return err
}

var _ RequestFiller = NewModuleRecordDelete()
