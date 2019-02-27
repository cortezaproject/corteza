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

func (mReq *ModuleList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

		mReq.Query = val
	}

	return err
}

var _ RequestFiller = NewModuleList()

// Module create request parameters
type ModuleCreate struct {
	Name   string
	Fields types.ModuleFieldSet
	Meta   sqlxTypes.JSONText
}

func NewModuleCreate() *ModuleCreate {
	return &ModuleCreate{}
}

func (mReq *ModuleCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

		mReq.Name = val
	}
	if val, ok := post["meta"]; ok {

		if mReq.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
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

func (mReq *ModuleRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleRead()

// Module attachments request parameters
type ModuleAttachments struct {
	Filter  string
	Page    int
	PerPage int
	Sort    string
}

func NewModuleAttachments() *ModuleAttachments {
	return &ModuleAttachments{}
}

func (mReq *ModuleAttachments) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	if val, ok := get["filter"]; ok {

		mReq.Filter = val
	}
	if val, ok := get["page"]; ok {

		mReq.Page = parseInt(val)
	}
	if val, ok := get["perPage"]; ok {

		mReq.PerPage = parseInt(val)
	}
	if val, ok := get["sort"]; ok {

		mReq.Sort = val
	}

	return err
}

var _ RequestFiller = NewModuleAttachments()

// Module update request parameters
type ModuleUpdate struct {
	ModuleID uint64 `json:",string"`
	Name     string
	Fields   types.ModuleFieldSet
	Meta     sqlxTypes.JSONText
}

func NewModuleUpdate() *ModuleUpdate {
	return &ModuleUpdate{}
}

func (mReq *ModuleUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	if val, ok := post["name"]; ok {

		mReq.Name = val
	}
	if val, ok := post["meta"]; ok {

		if mReq.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleUpdate()

// Module delete request parameters
type ModuleDelete struct {
	ModuleID uint64 `json:",string"`
}

func NewModuleDelete() *ModuleDelete {
	return &ModuleDelete{}
}

func (mReq *ModuleDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

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

func (mReq *ModuleRecordReport) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

		mReq.Metrics = val
	}
	if val, ok := get["dimensions"]; ok {

		mReq.Dimensions = val
	}
	if val, ok := get["filter"]; ok {

		mReq.Filter = val
	}
	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleRecordReport()

// Module record/list request parameters
type ModuleRecordList struct {
	Filter   string
	Page     int
	PerPage  int
	Sort     string
	ModuleID uint64 `json:",string"`
}

func NewModuleRecordList() *ModuleRecordList {
	return &ModuleRecordList{}
}

func (mReq *ModuleRecordList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	if val, ok := get["filter"]; ok {

		mReq.Filter = val
	}
	if val, ok := get["page"]; ok {

		mReq.Page = parseInt(val)
	}
	if val, ok := get["perPage"]; ok {

		mReq.PerPage = parseInt(val)
	}
	if val, ok := get["sort"]; ok {

		mReq.Sort = val
	}
	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewModuleRecordList()

// Module record/create request parameters
type ModuleRecordCreate struct {
	ModuleID uint64 `json:",string"`
	Values   types.RecordValueSet
}

func NewModuleRecordCreate() *ModuleRecordCreate {
	return &ModuleRecordCreate{}
}

func (mReq *ModuleRecordCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

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

func (mReq *ModuleRecordRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	mReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))

	return err
}

var _ RequestFiller = NewModuleRecordRead()

// Module record/update request parameters
type ModuleRecordUpdate struct {
	ModuleID uint64 `json:",string"`
	RecordID uint64 `json:",string"`
	Values   types.RecordValueSet
}

func NewModuleRecordUpdate() *ModuleRecordUpdate {
	return &ModuleRecordUpdate{}
}

func (mReq *ModuleRecordUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	mReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))

	return err
}

var _ RequestFiller = NewModuleRecordUpdate()

// Module record/delete request parameters
type ModuleRecordDelete struct {
	ModuleID uint64 `json:",string"`
	RecordID uint64 `json:",string"`
}

func NewModuleRecordDelete() *ModuleRecordDelete {
	return &ModuleRecordDelete{}
}

func (mReq *ModuleRecordDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	mReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))

	return err
}

var _ RequestFiller = NewModuleRecordDelete()
