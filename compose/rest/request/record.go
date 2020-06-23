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

// RecordReport request parameters
type RecordReport struct {
	hasMetrics bool
	rawMetrics string
	Metrics    string

	hasDimensions bool
	rawDimensions string
	Dimensions    string

	hasFilter bool
	rawFilter string
	Filter    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordReport request
func NewRecordReport() *RecordReport {
	return &RecordReport{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["metrics"] = r.Metrics
	out["dimensions"] = r.Dimensions
	out["filter"] = r.Filter
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
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
		r.hasMetrics = true
		r.rawMetrics = val
		r.Metrics = val
	}
	if val, ok := get["dimensions"]; ok {
		r.hasDimensions = true
		r.rawDimensions = val
		r.Dimensions = val
	}
	if val, ok := get["filter"]; ok {
		r.hasFilter = true
		r.rawFilter = val
		r.Filter = val
	}
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordReport()

// RecordList request parameters
type RecordList struct {
	hasQuery bool
	rawQuery string
	Query    string

	hasFilter bool
	rawFilter string
	Filter    string

	hasDeleted bool
	rawDeleted string
	Deleted    uint

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

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordList request
func NewRecordList() *RecordList {
	return &RecordList{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["filter"] = r.Filter
	out["deleted"] = r.Deleted
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
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

	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}
	if val, ok := get["filter"]; ok {
		r.hasFilter = true
		r.rawFilter = val
		r.Filter = val
	}
	if val, ok := get["deleted"]; ok {
		r.hasDeleted = true
		r.rawDeleted = val
		r.Deleted = parseUint(val)
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
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordList()

// RecordImportInit request parameters
type RecordImportInit struct {
	hasUpload bool
	rawUpload string
	Upload    *multipart.FileHeader

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordImportInit request
func NewRecordImportInit() *RecordImportInit {
	return &RecordImportInit{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportInit) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
func (r *RecordImportInit) Fill(req *http.Request) (err error) {
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

	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordImportInit()

// RecordImportRun request parameters
type RecordImportRun struct {
	hasSessionID bool
	rawSessionID string
	SessionID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasFields bool
	rawFields string
	Fields    json.RawMessage

	hasOnError bool
	rawOnError string
	OnError    string
}

// NewRecordImportRun request
func NewRecordImportRun() *RecordImportRun {
	return &RecordImportRun{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["sessionID"] = r.SessionID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["fields"] = r.Fields
	out["onError"] = r.OnError

	return out
}

// Fill processes request and fills internal variables
func (r *RecordImportRun) Fill(req *http.Request) (err error) {
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

	r.hasSessionID = true
	r.rawSessionID = chi.URLParam(req, "sessionID")
	r.SessionID = parseUInt64(chi.URLParam(req, "sessionID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	if val, ok := post["fields"]; ok {
		r.hasFields = true
		r.rawFields = val
		r.Fields = json.RawMessage(val)
	}
	if val, ok := post["onError"]; ok {
		r.hasOnError = true
		r.rawOnError = val
		r.OnError = val
	}

	return err
}

var _ RequestFiller = NewRecordImportRun()

// RecordImportProgress request parameters
type RecordImportProgress struct {
	hasSessionID bool
	rawSessionID string
	SessionID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordImportProgress request
func NewRecordImportProgress() *RecordImportProgress {
	return &RecordImportProgress{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportProgress) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["sessionID"] = r.SessionID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
func (r *RecordImportProgress) Fill(req *http.Request) (err error) {
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

	r.hasSessionID = true
	r.rawSessionID = chi.URLParam(req, "sessionID")
	r.SessionID = parseUInt64(chi.URLParam(req, "sessionID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordImportProgress()

// RecordExport request parameters
type RecordExport struct {
	hasFilter bool
	rawFilter string
	Filter    string

	hasFields bool
	rawFields []string
	Fields    []string

	hasTimezone bool
	rawTimezone string
	Timezone    string

	hasFilename bool
	rawFilename string
	Filename    string

	hasExt bool
	rawExt string
	Ext    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordExport request
func NewRecordExport() *RecordExport {
	return &RecordExport{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["filter"] = r.Filter
	out["fields"] = r.Fields
	out["timezone"] = r.Timezone
	out["filename"] = r.Filename
	out["ext"] = r.Ext
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
func (r *RecordExport) Fill(req *http.Request) (err error) {
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
		r.hasFilter = true
		r.rawFilter = val
		r.Filter = val
	}

	if val, ok := urlQuery["fields[]"]; ok {
		r.hasFields = true
		r.rawFields = val
		r.Fields = parseStrings(val)
	} else if val, ok = urlQuery["fields"]; ok {
		r.hasFields = true
		r.rawFields = val
		r.Fields = parseStrings(val)
	}

	if val, ok := get["timezone"]; ok {
		r.hasTimezone = true
		r.rawTimezone = val
		r.Timezone = val
	}
	r.hasFilename = true
	r.rawFilename = chi.URLParam(req, "filename")
	r.Filename = chi.URLParam(req, "filename")
	r.hasExt = true
	r.rawExt = chi.URLParam(req, "ext")
	r.Ext = chi.URLParam(req, "ext")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordExport()

// RecordExec request parameters
type RecordExec struct {
	hasProcedure bool
	rawProcedure string
	Procedure    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasArgs bool
	rawArgs []string
	Args    []ProcedureArg
}

// NewRecordExec request
func NewRecordExec() *RecordExec {
	return &RecordExec{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordExec) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["procedure"] = r.Procedure
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["args"] = r.Args

	return out
}

// Fill processes request and fills internal variables
func (r *RecordExec) Fill(req *http.Request) (err error) {
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

	r.hasProcedure = true
	r.rawProcedure = chi.URLParam(req, "procedure")
	r.Procedure = chi.URLParam(req, "procedure")
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordExec()

// RecordCreate request parameters
type RecordCreate struct {
	hasValues bool
	rawValues string
	Values    types.RecordValueSet

	hasRecords bool
	rawRecords string
	Records    types.RecordBulkSet

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordCreate request
func NewRecordCreate() *RecordCreate {
	return &RecordCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["values"] = r.Values
	out["records"] = r.Records
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordCreate()

// RecordRead request parameters
type RecordRead struct {
	hasRecordID bool
	rawRecordID string
	RecordID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordRead request
func NewRecordRead() *RecordRead {
	return &RecordRead{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasRecordID = true
	r.rawRecordID = chi.URLParam(req, "recordID")
	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordRead()

// RecordUpdate request parameters
type RecordUpdate struct {
	hasRecordID bool
	rawRecordID string
	RecordID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasValues bool
	rawValues string
	Values    types.RecordValueSet

	hasRecords bool
	rawRecords string
	Records    types.RecordBulkSet
}

// NewRecordUpdate request
func NewRecordUpdate() *RecordUpdate {
	return &RecordUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["values"] = r.Values
	out["records"] = r.Records

	return out
}

// Fill processes request and fills internal variables
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

	r.hasRecordID = true
	r.rawRecordID = chi.URLParam(req, "recordID")
	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordUpdate()

// RecordBulkDelete request parameters
type RecordBulkDelete struct {
	hasRecordIDs bool
	rawRecordIDs []string
	RecordIDs    []string

	hasTruncate bool
	rawTruncate string
	Truncate    bool

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordBulkDelete request
func NewRecordBulkDelete() *RecordBulkDelete {
	return &RecordBulkDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordBulkDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordIDs"] = r.RecordIDs
	out["truncate"] = r.Truncate
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
func (r *RecordBulkDelete) Fill(req *http.Request) (err error) {
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

	if val, ok := req.Form["recordIDs"]; ok {
		r.hasRecordIDs = true
		r.rawRecordIDs = val
		r.RecordIDs = parseStrings(val)
	}

	if val, ok := post["truncate"]; ok {
		r.hasTruncate = true
		r.rawTruncate = val
		r.Truncate = parseBool(val)
	}
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordBulkDelete()

// RecordDelete request parameters
type RecordDelete struct {
	hasRecordID bool
	rawRecordID string
	RecordID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordDelete request
func NewRecordDelete() *RecordDelete {
	return &RecordDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasRecordID = true
	r.rawRecordID = chi.URLParam(req, "recordID")
	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordDelete()

// RecordUpload request parameters
type RecordUpload struct {
	hasRecordID bool
	rawRecordID string
	RecordID    uint64 `json:",string"`

	hasFieldName bool
	rawFieldName string
	FieldName    string

	hasUpload bool
	rawUpload string
	Upload    *multipart.FileHeader

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordUpload request
func NewRecordUpload() *RecordUpload {
	return &RecordUpload{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["fieldName"] = r.FieldName
	out["upload.size"] = r.Upload.Size
	out["upload.filename"] = r.Upload.Filename

	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
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

	if val, ok := post["recordID"]; ok {
		r.hasRecordID = true
		r.rawRecordID = val
		r.RecordID = parseUInt64(val)
	}
	if val, ok := post["fieldName"]; ok {
		r.hasFieldName = true
		r.rawFieldName = val
		r.FieldName = val
	}
	if _, r.Upload, err = req.FormFile("upload"); err != nil {
		return errors.Wrap(err, "error processing uploaded file")
	}

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordUpload()

// RecordTriggerScript request parameters
type RecordTriggerScript struct {
	hasRecordID bool
	rawRecordID string
	RecordID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string

	hasValues bool
	rawValues string
	Values    types.RecordValueSet
}

// NewRecordTriggerScript request
func NewRecordTriggerScript() *RecordTriggerScript {
	return &RecordTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["recordID"] = r.RecordID
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID
	out["script"] = r.Script
	out["values"] = r.Values

	return out
}

// Fill processes request and fills internal variables
func (r *RecordTriggerScript) Fill(req *http.Request) (err error) {
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

	r.hasRecordID = true
	r.rawRecordID = chi.URLParam(req, "recordID")
	r.RecordID = parseUInt64(chi.URLParam(req, "recordID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewRecordTriggerScript()

// RecordTriggerScriptOnList request parameters
type RecordTriggerScriptOnList struct {
	hasScript bool
	rawScript string
	Script    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`
}

// NewRecordTriggerScriptOnList request
func NewRecordTriggerScriptOnList() *RecordTriggerScriptOnList {
	return &RecordTriggerScriptOnList{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScriptOnList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["script"] = r.Script
	out["namespaceID"] = r.NamespaceID
	out["moduleID"] = r.ModuleID

	return out
}

// Fill processes request and fills internal variables
func (r *RecordTriggerScriptOnList) Fill(req *http.Request) (err error) {
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

	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))

	return err
}

var _ RequestFiller = NewRecordTriggerScriptOnList()

// HasMetrics returns true if metrics was set
func (r *RecordReport) HasMetrics() bool {
	return r.hasMetrics
}

// RawMetrics returns raw value of metrics parameter
func (r *RecordReport) RawMetrics() string {
	return r.rawMetrics
}

// GetMetrics returns casted value of  metrics parameter
func (r *RecordReport) GetMetrics() string {
	return r.Metrics
}

// HasDimensions returns true if dimensions was set
func (r *RecordReport) HasDimensions() bool {
	return r.hasDimensions
}

// RawDimensions returns raw value of dimensions parameter
func (r *RecordReport) RawDimensions() string {
	return r.rawDimensions
}

// GetDimensions returns casted value of  dimensions parameter
func (r *RecordReport) GetDimensions() string {
	return r.Dimensions
}

// HasFilter returns true if filter was set
func (r *RecordReport) HasFilter() bool {
	return r.hasFilter
}

// RawFilter returns raw value of filter parameter
func (r *RecordReport) RawFilter() string {
	return r.rawFilter
}

// GetFilter returns casted value of  filter parameter
func (r *RecordReport) GetFilter() string {
	return r.Filter
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordReport) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordReport) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordReport) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordReport) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordReport) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordReport) GetModuleID() uint64 {
	return r.ModuleID
}

// HasQuery returns true if query was set
func (r *RecordList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *RecordList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *RecordList) GetQuery() string {
	return r.Query
}

// HasFilter returns true if filter was set
func (r *RecordList) HasFilter() bool {
	return r.hasFilter
}

// RawFilter returns raw value of filter parameter
func (r *RecordList) RawFilter() string {
	return r.rawFilter
}

// GetFilter returns casted value of  filter parameter
func (r *RecordList) GetFilter() string {
	return r.Filter
}

// HasDeleted returns true if deleted was set
func (r *RecordList) HasDeleted() bool {
	return r.hasDeleted
}

// RawDeleted returns raw value of deleted parameter
func (r *RecordList) RawDeleted() string {
	return r.rawDeleted
}

// GetDeleted returns casted value of  deleted parameter
func (r *RecordList) GetDeleted() uint {
	return r.Deleted
}

// HasLimit returns true if limit was set
func (r *RecordList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *RecordList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *RecordList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *RecordList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *RecordList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *RecordList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *RecordList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *RecordList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *RecordList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *RecordList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *RecordList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *RecordList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *RecordList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *RecordList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *RecordList) GetSort() string {
	return r.Sort
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordList) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordList) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordList) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordList) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordList) GetModuleID() uint64 {
	return r.ModuleID
}

// HasUpload returns true if upload was set
func (r *RecordImportInit) HasUpload() bool {
	return r.hasUpload
}

// RawUpload returns raw value of upload parameter
func (r *RecordImportInit) RawUpload() string {
	return r.rawUpload
}

// GetUpload returns casted value of  upload parameter
func (r *RecordImportInit) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordImportInit) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordImportInit) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordImportInit) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordImportInit) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordImportInit) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordImportInit) GetModuleID() uint64 {
	return r.ModuleID
}

// HasSessionID returns true if sessionID was set
func (r *RecordImportRun) HasSessionID() bool {
	return r.hasSessionID
}

// RawSessionID returns raw value of sessionID parameter
func (r *RecordImportRun) RawSessionID() string {
	return r.rawSessionID
}

// GetSessionID returns casted value of  sessionID parameter
func (r *RecordImportRun) GetSessionID() uint64 {
	return r.SessionID
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordImportRun) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordImportRun) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordImportRun) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordImportRun) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordImportRun) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordImportRun) GetModuleID() uint64 {
	return r.ModuleID
}

// HasFields returns true if fields was set
func (r *RecordImportRun) HasFields() bool {
	return r.hasFields
}

// RawFields returns raw value of fields parameter
func (r *RecordImportRun) RawFields() string {
	return r.rawFields
}

// GetFields returns casted value of  fields parameter
func (r *RecordImportRun) GetFields() json.RawMessage {
	return r.Fields
}

// HasOnError returns true if onError was set
func (r *RecordImportRun) HasOnError() bool {
	return r.hasOnError
}

// RawOnError returns raw value of onError parameter
func (r *RecordImportRun) RawOnError() string {
	return r.rawOnError
}

// GetOnError returns casted value of  onError parameter
func (r *RecordImportRun) GetOnError() string {
	return r.OnError
}

// HasSessionID returns true if sessionID was set
func (r *RecordImportProgress) HasSessionID() bool {
	return r.hasSessionID
}

// RawSessionID returns raw value of sessionID parameter
func (r *RecordImportProgress) RawSessionID() string {
	return r.rawSessionID
}

// GetSessionID returns casted value of  sessionID parameter
func (r *RecordImportProgress) GetSessionID() uint64 {
	return r.SessionID
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordImportProgress) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordImportProgress) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordImportProgress) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordImportProgress) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordImportProgress) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordImportProgress) GetModuleID() uint64 {
	return r.ModuleID
}

// HasFilter returns true if filter was set
func (r *RecordExport) HasFilter() bool {
	return r.hasFilter
}

// RawFilter returns raw value of filter parameter
func (r *RecordExport) RawFilter() string {
	return r.rawFilter
}

// GetFilter returns casted value of  filter parameter
func (r *RecordExport) GetFilter() string {
	return r.Filter
}

// HasFields returns true if fields was set
func (r *RecordExport) HasFields() bool {
	return r.hasFields
}

// RawFields returns raw value of fields parameter
func (r *RecordExport) RawFields() []string {
	return r.rawFields
}

// GetFields returns casted value of  fields parameter
func (r *RecordExport) GetFields() []string {
	return r.Fields
}

// HasTimezone returns true if timezone was set
func (r *RecordExport) HasTimezone() bool {
	return r.hasTimezone
}

// RawTimezone returns raw value of timezone parameter
func (r *RecordExport) RawTimezone() string {
	return r.rawTimezone
}

// GetTimezone returns casted value of  timezone parameter
func (r *RecordExport) GetTimezone() string {
	return r.Timezone
}

// HasFilename returns true if filename was set
func (r *RecordExport) HasFilename() bool {
	return r.hasFilename
}

// RawFilename returns raw value of filename parameter
func (r *RecordExport) RawFilename() string {
	return r.rawFilename
}

// GetFilename returns casted value of  filename parameter
func (r *RecordExport) GetFilename() string {
	return r.Filename
}

// HasExt returns true if ext was set
func (r *RecordExport) HasExt() bool {
	return r.hasExt
}

// RawExt returns raw value of ext parameter
func (r *RecordExport) RawExt() string {
	return r.rawExt
}

// GetExt returns casted value of  ext parameter
func (r *RecordExport) GetExt() string {
	return r.Ext
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordExport) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordExport) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordExport) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordExport) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordExport) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordExport) GetModuleID() uint64 {
	return r.ModuleID
}

// HasProcedure returns true if procedure was set
func (r *RecordExec) HasProcedure() bool {
	return r.hasProcedure
}

// RawProcedure returns raw value of procedure parameter
func (r *RecordExec) RawProcedure() string {
	return r.rawProcedure
}

// GetProcedure returns casted value of  procedure parameter
func (r *RecordExec) GetProcedure() string {
	return r.Procedure
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordExec) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordExec) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordExec) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordExec) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordExec) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordExec) GetModuleID() uint64 {
	return r.ModuleID
}

// HasArgs returns true if args was set
func (r *RecordExec) HasArgs() bool {
	return r.hasArgs
}

// RawArgs returns raw value of args parameter
func (r *RecordExec) RawArgs() []string {
	return r.rawArgs
}

// GetArgs returns casted value of  args parameter
func (r *RecordExec) GetArgs() []ProcedureArg {
	return r.Args
}

// HasValues returns true if values was set
func (r *RecordCreate) HasValues() bool {
	return r.hasValues
}

// RawValues returns raw value of values parameter
func (r *RecordCreate) RawValues() string {
	return r.rawValues
}

// GetValues returns casted value of  values parameter
func (r *RecordCreate) GetValues() types.RecordValueSet {
	return r.Values
}

// HasRecords returns true if records was set
func (r *RecordCreate) HasRecords() bool {
	return r.hasRecords
}

// RawRecords returns raw value of records parameter
func (r *RecordCreate) RawRecords() string {
	return r.rawRecords
}

// GetRecords returns casted value of  records parameter
func (r *RecordCreate) GetRecords() types.RecordBulkSet {
	return r.Records
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordCreate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordCreate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordCreate) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordCreate) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordCreate) GetModuleID() uint64 {
	return r.ModuleID
}

// HasRecordID returns true if recordID was set
func (r *RecordRead) HasRecordID() bool {
	return r.hasRecordID
}

// RawRecordID returns raw value of recordID parameter
func (r *RecordRead) RawRecordID() string {
	return r.rawRecordID
}

// GetRecordID returns casted value of  recordID parameter
func (r *RecordRead) GetRecordID() uint64 {
	return r.RecordID
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordRead) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordRead) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordRead) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordRead) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordRead) GetModuleID() uint64 {
	return r.ModuleID
}

// HasRecordID returns true if recordID was set
func (r *RecordUpdate) HasRecordID() bool {
	return r.hasRecordID
}

// RawRecordID returns raw value of recordID parameter
func (r *RecordUpdate) RawRecordID() string {
	return r.rawRecordID
}

// GetRecordID returns casted value of  recordID parameter
func (r *RecordUpdate) GetRecordID() uint64 {
	return r.RecordID
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordUpdate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordUpdate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordUpdate) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordUpdate) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// HasValues returns true if values was set
func (r *RecordUpdate) HasValues() bool {
	return r.hasValues
}

// RawValues returns raw value of values parameter
func (r *RecordUpdate) RawValues() string {
	return r.rawValues
}

// GetValues returns casted value of  values parameter
func (r *RecordUpdate) GetValues() types.RecordValueSet {
	return r.Values
}

// HasRecords returns true if records was set
func (r *RecordUpdate) HasRecords() bool {
	return r.hasRecords
}

// RawRecords returns raw value of records parameter
func (r *RecordUpdate) RawRecords() string {
	return r.rawRecords
}

// GetRecords returns casted value of  records parameter
func (r *RecordUpdate) GetRecords() types.RecordBulkSet {
	return r.Records
}

// HasRecordIDs returns true if recordIDs was set
func (r *RecordBulkDelete) HasRecordIDs() bool {
	return r.hasRecordIDs
}

// RawRecordIDs returns raw value of recordIDs parameter
func (r *RecordBulkDelete) RawRecordIDs() []string {
	return r.rawRecordIDs
}

// GetRecordIDs returns casted value of  recordIDs parameter
func (r *RecordBulkDelete) GetRecordIDs() []string {
	return r.RecordIDs
}

// HasTruncate returns true if truncate was set
func (r *RecordBulkDelete) HasTruncate() bool {
	return r.hasTruncate
}

// RawTruncate returns raw value of truncate parameter
func (r *RecordBulkDelete) RawTruncate() string {
	return r.rawTruncate
}

// GetTruncate returns casted value of  truncate parameter
func (r *RecordBulkDelete) GetTruncate() bool {
	return r.Truncate
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordBulkDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordBulkDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordBulkDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordBulkDelete) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordBulkDelete) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordBulkDelete) GetModuleID() uint64 {
	return r.ModuleID
}

// HasRecordID returns true if recordID was set
func (r *RecordDelete) HasRecordID() bool {
	return r.hasRecordID
}

// RawRecordID returns raw value of recordID parameter
func (r *RecordDelete) RawRecordID() string {
	return r.rawRecordID
}

// GetRecordID returns casted value of  recordID parameter
func (r *RecordDelete) GetRecordID() uint64 {
	return r.RecordID
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordDelete) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordDelete) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordDelete) GetModuleID() uint64 {
	return r.ModuleID
}

// HasRecordID returns true if recordID was set
func (r *RecordUpload) HasRecordID() bool {
	return r.hasRecordID
}

// RawRecordID returns raw value of recordID parameter
func (r *RecordUpload) RawRecordID() string {
	return r.rawRecordID
}

// GetRecordID returns casted value of  recordID parameter
func (r *RecordUpload) GetRecordID() uint64 {
	return r.RecordID
}

// HasFieldName returns true if fieldName was set
func (r *RecordUpload) HasFieldName() bool {
	return r.hasFieldName
}

// RawFieldName returns raw value of fieldName parameter
func (r *RecordUpload) RawFieldName() string {
	return r.rawFieldName
}

// GetFieldName returns casted value of  fieldName parameter
func (r *RecordUpload) GetFieldName() string {
	return r.FieldName
}

// HasUpload returns true if upload was set
func (r *RecordUpload) HasUpload() bool {
	return r.hasUpload
}

// RawUpload returns raw value of upload parameter
func (r *RecordUpload) RawUpload() string {
	return r.rawUpload
}

// GetUpload returns casted value of  upload parameter
func (r *RecordUpload) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordUpload) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordUpload) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordUpload) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordUpload) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordUpload) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordUpload) GetModuleID() uint64 {
	return r.ModuleID
}

// HasRecordID returns true if recordID was set
func (r *RecordTriggerScript) HasRecordID() bool {
	return r.hasRecordID
}

// RawRecordID returns raw value of recordID parameter
func (r *RecordTriggerScript) RawRecordID() string {
	return r.rawRecordID
}

// GetRecordID returns casted value of  recordID parameter
func (r *RecordTriggerScript) GetRecordID() uint64 {
	return r.RecordID
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordTriggerScript) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordTriggerScript) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordTriggerScript) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordTriggerScript) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordTriggerScript) GetModuleID() uint64 {
	return r.ModuleID
}

// HasScript returns true if script was set
func (r *RecordTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *RecordTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *RecordTriggerScript) GetScript() string {
	return r.Script
}

// HasValues returns true if values was set
func (r *RecordTriggerScript) HasValues() bool {
	return r.hasValues
}

// RawValues returns raw value of values parameter
func (r *RecordTriggerScript) RawValues() string {
	return r.rawValues
}

// GetValues returns casted value of  values parameter
func (r *RecordTriggerScript) GetValues() types.RecordValueSet {
	return r.Values
}

// HasScript returns true if script was set
func (r *RecordTriggerScriptOnList) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *RecordTriggerScriptOnList) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *RecordTriggerScriptOnList) GetScript() string {
	return r.Script
}

// HasNamespaceID returns true if namespaceID was set
func (r *RecordTriggerScriptOnList) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *RecordTriggerScriptOnList) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *RecordTriggerScriptOnList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *RecordTriggerScriptOnList) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *RecordTriggerScriptOnList) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *RecordTriggerScriptOnList) GetModuleID() uint64 {
	return r.ModuleID
}
