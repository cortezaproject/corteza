package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	RecordReport struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Metrics GET parameter
		//
		// Metrics (eg: 'SUM(money), MAX(calls)')
		Metrics string

		// Dimensions GET parameter
		//
		// Dimensions (eg: 'DATE(foo), status')
		Dimensions string

		// Filter GET parameter
		//
		// Filter (eg: 'DATE(foo) > 2010')
		Filter string
	}

	RecordList struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Query GET parameter
		//
		// Record filtering query
		Query string

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted records
		Deleted uint

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total records counter
		IncTotal bool

		// IncPageNavigation GET parameter
		//
		// Include page navigation
		IncPageNavigation bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	RecordImportInit struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Upload POST parameter
		//
		// File import
		Upload *multipart.FileHeader
	}

	RecordImportRun struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// SessionID PATH parameter
		//
		// Import session
		SessionID uint64 `json:",string"`

		// Fields POST parameter
		//
		// Fields defined by import file
		Fields json.RawMessage

		// OnError POST parameter
		//
		// What happens if record fails to import
		OnError string
	}

	RecordImportProgress struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// SessionID PATH parameter
		//
		// Import session
		SessionID uint64 `json:",string"`
	}

	RecordExport struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Filename PATH parameter
		//
		// Filename to use
		Filename string

		// Ext PATH parameter
		//
		// Export format
		Ext string

		// Filter GET parameter
		//
		// Filtering condition
		Filter string

		// Fields GET parameter
		//
		// Fields to export
		Fields []string

		// Timezone GET parameter
		//
		// Convert times to this timezone
		Timezone string
	}

	RecordExec struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Procedure PATH parameter
		//
		// Name of procedure to execute
		Procedure string

		// Args POST parameter
		//
		// Procedure arguments
		Args []ProcedureArg
	}

	RecordCreate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Values POST parameter
		//
		// Record values
		Values types.RecordValueSet

		// OwnedBy POST parameter
		//
		// Record Owner
		OwnedBy uint64 `json:",string"`

		// Records POST parameter
		//
		// Records
		Records types.RecordBulkSet

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	RecordRead struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// RecordID PATH parameter
		//
		// Record ID
		RecordID uint64 `json:",string"`
	}

	RecordUpdate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// RecordID PATH parameter
		//
		// Record ID
		RecordID uint64 `json:",string"`

		// Values POST parameter
		//
		// Record values
		Values types.RecordValueSet

		// OwnedBy POST parameter
		//
		// Record Owner
		OwnedBy uint64 `json:",string"`

		// Records POST parameter
		//
		// Records
		Records types.RecordBulkSet

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	RecordBulkDelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// RecordIDs POST parameter
		//
		// IDs of records to delete
		RecordIDs []string

		// Truncate POST parameter
		//
		// Remove ALL records of a specified module (pending implementation)
		Truncate bool
	}

	RecordDelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// RecordID PATH parameter
		//
		// Record ID
		RecordID uint64 `json:",string"`
	}

	RecordUpload struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// RecordID POST parameter
		//
		// Record ID
		RecordID uint64 `json:",string"`

		// FieldName POST parameter
		//
		// Field name
		FieldName string

		// Upload POST parameter
		//
		// File to upload
		Upload *multipart.FileHeader
	}

	RecordTriggerScript struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// RecordID PATH parameter
		//
		// ID
		RecordID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string

		// Values POST parameter
		//
		// Record values
		Values types.RecordValueSet
	}

	RecordTriggerScriptOnList struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string

		// Args POST parameter
		//
		// Arguments to pass to the script
		Args map[string]interface{}
	}
)

// NewRecordReport request
func NewRecordReport() *RecordReport {
	return &RecordReport{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"metrics":     r.Metrics,
		"dimensions":  r.Dimensions,
		"filter":      r.Filter,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) GetMetrics() string {
	return r.Metrics
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) GetDimensions() string {
	return r.Dimensions
}

// Auditable returns all auditable/loggable parameters
func (r RecordReport) GetFilter() string {
	return r.Filter
}

// Fill processes request and fills internal variables
func (r *RecordReport) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["metrics"]; ok && len(val) > 0 {
			r.Metrics, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["dimensions"]; ok && len(val) > 0 {
			r.Dimensions, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["filter"]; ok && len(val) > 0 {
			r.Filter, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordList request
func NewRecordList() *RecordList {
	return &RecordList{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":       r.NamespaceID,
		"moduleID":          r.ModuleID,
		"query":             r.Query,
		"labels":            r.Labels,
		"deleted":           r.Deleted,
		"limit":             r.Limit,
		"incTotal":          r.IncTotal,
		"incPageNavigation": r.IncPageNavigation,
		"pageCursor":        r.PageCursor,
		"sort":              r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetIncPageNavigation() bool {
	return r.IncPageNavigation
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r RecordList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *RecordList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := tmp["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incPageNavigation"]; ok && len(val) > 0 {
			r.IncPageNavigation, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordImportInit request
func NewRecordImportInit() *RecordImportInit {
	return &RecordImportInit{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportInit) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"upload":      r.Upload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportInit) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportInit) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportInit) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *RecordImportInit) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			// Ignoring upload as its handled in the POST params section
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordImportRun request
func NewRecordImportRun() *RecordImportRun {
	return &RecordImportRun{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"sessionID":   r.SessionID,
		"fields":      r.Fields,
		"onError":     r.OnError,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) GetSessionID() uint64 {
	return r.SessionID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) GetFields() json.RawMessage {
	return r.Fields
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportRun) GetOnError() string {
	return r.OnError
}

// Fill processes request and fills internal variables
func (r *RecordImportRun) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["fields"]; ok && len(val) > 0 {
				r.Fields, err = json.RawMessage(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["onError"]; ok && len(val) > 0 {
				r.OnError, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["fields"]; ok && len(val) > 0 {
			r.Fields, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["onError"]; ok && len(val) > 0 {
			r.OnError, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordImportProgress request
func NewRecordImportProgress() *RecordImportProgress {
	return &RecordImportProgress{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportProgress) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"sessionID":   r.SessionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportProgress) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportProgress) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordImportProgress) GetSessionID() uint64 {
	return r.SessionID
}

// Fill processes request and fills internal variables
func (r *RecordImportProgress) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "sessionID")
		r.SessionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordExport request
func NewRecordExport() *RecordExport {
	return &RecordExport{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"filename":    r.Filename,
		"ext":         r.Ext,
		"filter":      r.Filter,
		"fields":      r.Fields,
		"timezone":    r.Timezone,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetFilename() string {
	return r.Filename
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetExt() string {
	return r.Ext
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetFilter() string {
	return r.Filter
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetFields() []string {
	return r.Fields
}

// Auditable returns all auditable/loggable parameters
func (r RecordExport) GetTimezone() string {
	return r.Timezone
}

// Fill processes request and fills internal variables
func (r *RecordExport) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["filter"]; ok && len(val) > 0 {
			r.Filter, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["fields[]"]; ok {
			r.Fields, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["fields"]; ok {
			r.Fields, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["timezone"]; ok && len(val) > 0 {
			r.Timezone, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "filename")
		r.Filename, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "ext")
		r.Ext, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordExec request
func NewRecordExec() *RecordExec {
	return &RecordExec{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordExec) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"procedure":   r.Procedure,
		"args":        r.Args,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordExec) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordExec) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordExec) GetProcedure() string {
	return r.Procedure
}

// Auditable returns all auditable/loggable parameters
func (r RecordExec) GetArgs() []ProcedureArg {
	return r.Args
}

// Fill processes request and fills internal variables
func (r *RecordExec) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["args[]"]; ok && len(val) > 0  {
		//    r.Args, err = []ProcedureArg(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "procedure")
		r.Procedure, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordCreate request
func NewRecordCreate() *RecordCreate {
	return &RecordCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"values":      r.Values,
		"ownedBy":     r.OwnedBy,
		"records":     r.Records,
		"labels":      r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) GetValues() types.RecordValueSet {
	return r.Values
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) GetRecords() types.RecordBulkSet {
	return r.Records
}

// Auditable returns all auditable/loggable parameters
func (r RecordCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *RecordCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["ownedBy"]; ok && len(val) > 0 {
				r.OwnedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["values[]"]; ok && len(val) > 0  {
		//    r.Values, err = types.RecordValueSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["records[]"]; ok && len(val) > 0  {
		//    r.Records, err = types.RecordBulkSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordRead request
func NewRecordRead() *RecordRead {
	return &RecordRead{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"recordID":    r.RecordID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordRead) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordRead) GetRecordID() uint64 {
	return r.RecordID
}

// Fill processes request and fills internal variables
func (r *RecordRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "recordID")
		r.RecordID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordUpdate request
func NewRecordUpdate() *RecordUpdate {
	return &RecordUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"recordID":    r.RecordID,
		"values":      r.Values,
		"ownedBy":     r.OwnedBy,
		"records":     r.Records,
		"labels":      r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetRecordID() uint64 {
	return r.RecordID
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetValues() types.RecordValueSet {
	return r.Values
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetRecords() types.RecordBulkSet {
	return r.Records
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *RecordUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["ownedBy"]; ok && len(val) > 0 {
				r.OwnedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["values[]"]; ok && len(val) > 0  {
		//    r.Values, err = types.RecordValueSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["records[]"]; ok && len(val) > 0  {
		//    r.Records, err = types.RecordBulkSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "recordID")
		r.RecordID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordBulkDelete request
func NewRecordBulkDelete() *RecordBulkDelete {
	return &RecordBulkDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordBulkDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"recordIDs":   r.RecordIDs,
		"truncate":    r.Truncate,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordBulkDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordBulkDelete) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordBulkDelete) GetRecordIDs() []string {
	return r.RecordIDs
}

// Auditable returns all auditable/loggable parameters
func (r RecordBulkDelete) GetTruncate() bool {
	return r.Truncate
}

// Fill processes request and fills internal variables
func (r *RecordBulkDelete) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["truncate"]; ok && len(val) > 0 {
				r.Truncate, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		//if val, ok := req.Form["recordIDs[]"]; ok && len(val) > 0  {
		//    r.RecordIDs, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["truncate"]; ok && len(val) > 0 {
			r.Truncate, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordDelete request
func NewRecordDelete() *RecordDelete {
	return &RecordDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"recordID":    r.RecordID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordDelete) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordDelete) GetRecordID() uint64 {
	return r.RecordID
}

// Fill processes request and fills internal variables
func (r *RecordDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "recordID")
		r.RecordID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordUpload request
func NewRecordUpload() *RecordUpload {
	return &RecordUpload{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"recordID":    r.RecordID,
		"fieldName":   r.FieldName,
		"upload":      r.Upload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) GetRecordID() uint64 {
	return r.RecordID
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) GetFieldName() string {
	return r.FieldName
}

// Auditable returns all auditable/loggable parameters
func (r RecordUpload) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *RecordUpload) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["recordID"]; ok && len(val) > 0 {
				r.RecordID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["fieldName"]; ok && len(val) > 0 {
				r.FieldName, err = val[0], nil
				if err != nil {
					return err
				}
			}

			// Ignoring upload as its handled in the POST params section
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["recordID"]; ok && len(val) > 0 {
			r.RecordID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["fieldName"]; ok && len(val) > 0 {
			r.FieldName, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordTriggerScript request
func NewRecordTriggerScript() *RecordTriggerScript {
	return &RecordTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"recordID":    r.RecordID,
		"script":      r.Script,
		"values":      r.Values,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) GetRecordID() uint64 {
	return r.RecordID
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) GetScript() string {
	return r.Script
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScript) GetValues() types.RecordValueSet {
	return r.Values
}

// Fill processes request and fills internal variables
func (r *RecordTriggerScript) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["script"]; ok && len(val) > 0 {
				r.Script, err = val[0], nil
				if err != nil {
					return err
				}
			}

		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["script"]; ok && len(val) > 0 {
			r.Script, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["values[]"]; ok && len(val) > 0  {
		//    r.Values, err = types.RecordValueSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "recordID")
		r.RecordID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewRecordTriggerScriptOnList request
func NewRecordTriggerScriptOnList() *RecordTriggerScriptOnList {
	return &RecordTriggerScriptOnList{}
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScriptOnList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"moduleID":    r.ModuleID,
		"script":      r.Script,
		"args":        r.Args,
	}
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScriptOnList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScriptOnList) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScriptOnList) GetScript() string {
	return r.Script
}

// Auditable returns all auditable/loggable parameters
func (r RecordTriggerScriptOnList) GetArgs() map[string]interface{} {
	return r.Args
}

// Fill processes request and fills internal variables
func (r *RecordTriggerScriptOnList) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["script"]; ok && len(val) > 0 {
				r.Script, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["args[]"]; ok {
				r.Args, err = parseMapStringInterface(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["args"]; ok {
				r.Args, err = parseMapStringInterface(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["script"]; ok && len(val) > 0 {
			r.Script, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["args[]"]; ok {
			r.Args, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["args"]; ok {
			r.Args, err = parseMapStringInterface(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
