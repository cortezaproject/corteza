package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `record.go`, `record.util.go` or `record_test.go` to
	implement your API calls, helper functions and tests. The file `record.go`
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

	"github.com/crusttech/crust/crm/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Record report request parameters
type RecordReport struct {
	Metrics    string
	Dimensions string
	Filter     string
	ModuleID   uint64 `json:",string"`
}

func NewRecordReport() *RecordReport {
	return &RecordReport{}
}

func (rReq *RecordReport) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

		rReq.Metrics = val
	}
	if val, ok := get["dimensions"]; ok {

		rReq.Dimensions = val
	}
	if val, ok := get["filter"]; ok {

		rReq.Filter = val
	}
	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordReport()

// Record list request parameters
type RecordList struct {
	Filter   string
	Page     int
	PerPage  int
	Sort     string
	ModuleID uint64 `json:",string"`
}

func NewRecordList() *RecordList {
	return &RecordList{}
}

func (rReq *RecordList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

		rReq.Filter = val
	}
	if val, ok := get["page"]; ok {

		rReq.Page = parseInt(val)
	}
	if val, ok := get["perPage"]; ok {

		rReq.PerPage = parseInt(val)
	}
	if val, ok := get["sort"]; ok {

		rReq.Sort = val
	}
	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordList()

// Record create request parameters
type RecordCreate struct {
	Values   types.RecordValueSet
	ModuleID uint64 `json:",string"`
}

func NewRecordCreate() *RecordCreate {
	return &RecordCreate{}
}

func (rReq *RecordCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordCreate()

// Record read request parameters
type RecordRead struct {
	RecordID uint64 `json:",string"`
	ModuleID uint64 `json:",string"`
}

func NewRecordRead() *RecordRead {
	return &RecordRead{}
}

func (rReq *RecordRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

	rReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))
	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordRead()

// Record update request parameters
type RecordUpdate struct {
	RecordID uint64 `json:",string"`
	ModuleID uint64 `json:",string"`
	Values   types.RecordValueSet
}

func NewRecordUpdate() *RecordUpdate {
	return &RecordUpdate{}
}

func (rReq *RecordUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

	rReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))
	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordUpdate()

// Record delete request parameters
type RecordDelete struct {
	RecordID uint64 `json:",string"`
	ModuleID uint64 `json:",string"`
}

func NewRecordDelete() *RecordDelete {
	return &RecordDelete{}
}

func (rReq *RecordDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

	rReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))
	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordDelete()

// Record upload request parameters
type RecordUpload struct {
	RecordID  uint64 `json:",string"`
	FieldName string
	ModuleID  uint64 `json:",string"`
	Upload    *multipart.FileHeader
}

func NewRecordUpload() *RecordUpload {
	return &RecordUpload{}
}

func (rReq *RecordUpload) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(rReq)

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

	rReq.RecordID = parseUInt64(chi.URLParam(r, "recordID"))
	rReq.FieldName = chi.URLParam(r, "fieldName")
	rReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	if _, rReq.Upload, err = r.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewRecordUpload()
