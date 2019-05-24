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

	"github.com/cortezaproject/corteza-server/compose/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Record report request parameters
type RecordReport struct {
	Metrics     string
	Dimensions  string
	Filter      string
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
}

func NewRecordReport() *RecordReport {
	return &RecordReport{}
}

func (r RecordReport) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["metrics"] = r.Metrics
	out["dimensions"] = r.Dimensions
	out["filter"] = r.Filter
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

func (r *RecordReport) Fill(req *http.Request) (err error) {
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

	if val, ok := get["metrics"]; ok {
		r.Metrics = val
	}
	if val, ok := get["dimensions"]; ok {
		r.Dimensions = val
	}
	if val, ok := get["filter"]; ok {
		r.Filter = val
	}
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordReport()

// Record list request parameters
type RecordList struct {
	Filter      string
	Page        uint
	PerPage     uint
	Sort        string
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
}

func NewRecordList() *RecordList {
	return &RecordList{}
}

func (r RecordList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["filter"] = r.Filter
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

func (r *RecordList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["filter"]; ok {
		r.Filter = val
	}
	if val, ok := get["page"]; ok {
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.PerPage = parseUint(val)
	}
	if val, ok := get["sort"]; ok {
		r.Sort = val
	}
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordList()

// Record create request parameters
type RecordCreate struct {
	Values      types.RecordValueSet
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
}

func NewRecordCreate() *RecordCreate {
	return &RecordCreate{}
}

func (r RecordCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["values"] = r.Values
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

func (r *RecordCreate) Fill(req *http.Request) (err error) {
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

	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordCreate()

// Record read request parameters
type RecordRead struct {
	RecordID    uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
}

func NewRecordRead() *RecordRead {
	return &RecordRead{}
}

func (r RecordRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

func (r *RecordRead) Fill(req *http.Request) (err error) {
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

	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordRead()

// Record update request parameters
type RecordUpdate struct {
	RecordID    uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	Values      types.RecordValueSet
}

func NewRecordUpdate() *RecordUpdate {
	return &RecordUpdate{}
}

func (r RecordUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["values"] = r.Values

	return out
}

func (r *RecordUpdate) Fill(req *http.Request) (err error) {
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

	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordUpdate()

// Record delete request parameters
type RecordDelete struct {
	RecordID    uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
}

func NewRecordDelete() *RecordDelete {
	return &RecordDelete{}
}

func (r RecordDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

func (r *RecordDelete) Fill(req *http.Request) (err error) {
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

	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordDelete()

// Record upload request parameters
type RecordUpload struct {
	RecordID    uint64 `json:",string"`
	FieldName   string
	NamespaceID uint64 `json:",string"`
	ModuleID    uint64 `json:",string"`
	Upload      *multipart.FileHeader
}

func NewRecordUpload() *RecordUpload {
	return &RecordUpload{}
}

func (r RecordUpload) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["fieldName"] = r.FieldName
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	return out
}

func (r *RecordUpload) Fill(req *http.Request) (err error) {
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

	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.FieldName = chi.URLParam(req, "fieldName")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error procesing uploaded file")
	}

	return err
}

var _ RequestFiller = NewRecordUpload()
