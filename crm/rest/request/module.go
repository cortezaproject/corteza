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

// Module chart request parameters
type ModuleChart struct {
	Name        string
	Description string
	XAxis       string
	XMin        string
	XMax        string
	YAxis       string
	GroupBy     string
	Sum         string
	Count       string
	Kind        string
	ModuleID    uint64 `json:",string"`
}

func NewModuleChart() *ModuleChart {
	return &ModuleChart{}
}

func (m *ModuleChart) Fill(r *http.Request) (err error) {
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

	if val, ok := get["name"]; ok {

		m.Name = val
	}
	if val, ok := get["description"]; ok {

		m.Description = val
	}
	if val, ok := get["xAxis"]; ok {

		m.XAxis = val
	}
	if val, ok := get["xMin"]; ok {

		m.XMin = val
	}
	if val, ok := get["xMax"]; ok {

		m.XMax = val
	}
	if val, ok := get["yAxis"]; ok {

		m.YAxis = val
	}
	if val, ok := get["groupBy"]; ok {

		m.GroupBy = val
	}
	if val, ok := get["sum"]; ok {

		m.Sum = val
	}
	if val, ok := get["count"]; ok {

		m.Count = val
	}
	if val, ok := get["kind"]; ok {

		m.Kind = val
	}
	m.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleChart()

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

// Module content/list request parameters
type ModuleContentList struct {
	Query    string
	Page     int
	PerPage  int
	Sort     string
	ModuleID uint64 `json:",string"`
}

func NewModuleContentList() *ModuleContentList {
	return &ModuleContentList{}
}

func (m *ModuleContentList) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewModuleContentList()

// Module content/create request parameters
type ModuleContentCreate struct {
	ModuleID uint64 `json:",string"`
	Fields   sqlxTypes.JSONText
}

func NewModuleContentCreate() *ModuleContentCreate {
	return &ModuleContentCreate{}
}

func (m *ModuleContentCreate) Fill(r *http.Request) (err error) {
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

var _ RequestFiller = NewModuleContentCreate()

// Module content/read request parameters
type ModuleContentRead struct {
	ModuleID  uint64 `json:",string"`
	ContentID uint64 `json:",string"`
}

func NewModuleContentRead() *ModuleContentRead {
	return &ModuleContentRead{}
}

func (m *ModuleContentRead) Fill(r *http.Request) (err error) {
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
	m.ContentID = parseUInt64(chi.URLParam(r, "contentID"))

	return err
}

var _ RequestFiller = NewModuleContentRead()

// Module content/edit request parameters
type ModuleContentEdit struct {
	ModuleID  uint64 `json:",string"`
	ContentID uint64 `json:",string"`
	Fields    sqlxTypes.JSONText
}

func NewModuleContentEdit() *ModuleContentEdit {
	return &ModuleContentEdit{}
}

func (m *ModuleContentEdit) Fill(r *http.Request) (err error) {
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
	m.ContentID = parseUInt64(chi.URLParam(r, "contentID"))
	if val, ok := post["fields"]; ok {

		if m.Fields, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleContentEdit()

// Module content/delete request parameters
type ModuleContentDelete struct {
	ModuleID  uint64 `json:",string"`
	ContentID uint64 `json:",string"`
}

func NewModuleContentDelete() *ModuleContentDelete {
	return &ModuleContentDelete{}
}

func (m *ModuleContentDelete) Fill(r *http.Request) (err error) {
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
	m.ContentID = parseUInt64(chi.URLParam(r, "contentID"))

	return err
}

var _ RequestFiller = NewModuleContentDelete()
