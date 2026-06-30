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
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/crusttech/human/server/system/types"
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
	DmlConnectionList struct {
		// Handle GET parameter
		//
		//
		Handle string

		// Type GET parameter
		//
		//
		Type string

		// ConnectionID GET parameter
		//
		//
		ConnectionID []string
	}

	DmlConnectionCreate struct {
		// Handle POST parameter
		//
		//
		Handle string

		// Label POST parameter
		//
		//
		Label string

		// Params POST parameter
		//
		//
		Params types.DmlConnectionParams
	}

	DmlConnectionRead struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`
	}

	DmlConnectionModels struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// Ident GET parameter
		//
		//
		Ident []string
	}

	DmlMappingCreate struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// NamespaceHandle POST parameter
		//
		//
		NamespaceHandle string

		// SourceIdent POST parameter
		//
		//
		SourceIdent string

		// ModuleHandle POST parameter
		//
		//
		ModuleHandle string

		// ModuleName POST parameter
		//
		//
		ModuleName string

		// Skip POST parameter
		//
		//
		Skip bool

		// Identifier POST parameter
		//
		//
		Identifier string

		// Columns POST parameter
		//
		//
		Columns types.DmlColumnMapSet
	}

	DmlMappingRead struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// MappingID PATH parameter
		//
		//
		MappingID uint64 `json:",string"`
	}

	DmlMappingUpdate struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// MappingID PATH parameter
		//
		//
		MappingID uint64 `json:",string"`

		// NamespaceHandle POST parameter
		//
		//
		NamespaceHandle string

		// SourceIdent POST parameter
		//
		//
		SourceIdent string

		// ModuleHandle POST parameter
		//
		//
		ModuleHandle string

		// ModuleName POST parameter
		//
		//
		ModuleName string

		// Skip POST parameter
		//
		//
		Skip bool

		// Identifier POST parameter
		//
		//
		Identifier string

		// Columns POST parameter
		//
		//
		Columns types.DmlColumnMapSet
	}

	DmlMappingDelete struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// MappingID PATH parameter
		//
		//
		MappingID uint64 `json:",string"`
	}

	DmlImportRun struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// MappingID PATH parameter
		//
		//
		MappingID uint64 `json:",string"`

		// Method POST parameter
		//
		//
		Method types.DmlImportMethod
	}

	DmlImportRunRead struct {
		// ConnectionID PATH parameter
		//
		//
		ConnectionID uint64 `json:",string"`

		// MappingID PATH parameter
		//
		//
		MappingID uint64 `json:",string"`

		// RunID PATH parameter
		//
		//
		RunID uint64 `json:",string"`
	}
)

// NewDmlConnectionList request
func NewDmlConnectionList() *DmlConnectionList {
	return &DmlConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":       r.Handle,
		"type":         r.Type,
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionList) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionList) GetConnectionID() []string {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DmlConnectionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["connectionID[]"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["connectionID"]; ok {
			r.ConnectionID, err = val, nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDmlConnectionCreate request
func NewDmlConnectionCreate() *DmlConnectionCreate {
	return &DmlConnectionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle": r.Handle,
		"label":  r.Label,
		"params": r.Params,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionCreate) GetLabel() string {
	return r.Label
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionCreate) GetParams() types.DmlConnectionParams {
	return r.Params
}

// Fill processes request and fills internal variables
func (r *DmlConnectionCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["label"]; ok && len(val) > 0 {
				r.Label, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["params[]"]; ok {
				r.Params, err = types.ParseDmlConnectionParams(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["params"]; ok {
				r.Params, err = types.ParseDmlConnectionParams(val)
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["label"]; ok && len(val) > 0 {
			r.Label, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["params[]"]; ok {
			r.Params, err = types.ParseDmlConnectionParams(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["params"]; ok {
			r.Params, err = types.ParseDmlConnectionParams(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDmlConnectionRead request
func NewDmlConnectionRead() *DmlConnectionRead {
	return &DmlConnectionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DmlConnectionRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlConnectionModels request
func NewDmlConnectionModels() *DmlConnectionModels {
	return &DmlConnectionModels{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionModels) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"ident":        r.Ident,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionModels) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlConnectionModels) GetIdent() []string {
	return r.Ident
}

// Fill processes request and fills internal variables
func (r *DmlConnectionModels) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["ident[]"]; ok {
			r.Ident, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["ident"]; ok {
			r.Ident, err = val, nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlMappingCreate request
func NewDmlMappingCreate() *DmlMappingCreate {
	return &DmlMappingCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID":    r.ConnectionID,
		"namespaceHandle": r.NamespaceHandle,
		"sourceIdent":     r.SourceIdent,
		"moduleHandle":    r.ModuleHandle,
		"moduleName":      r.ModuleName,
		"skip":            r.Skip,
		"identifier":      r.Identifier,
		"columns":         r.Columns,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetNamespaceHandle() string {
	return r.NamespaceHandle
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetSourceIdent() string {
	return r.SourceIdent
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetModuleHandle() string {
	return r.ModuleHandle
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetModuleName() string {
	return r.ModuleName
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetSkip() bool {
	return r.Skip
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetIdentifier() string {
	return r.Identifier
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingCreate) GetColumns() types.DmlColumnMapSet {
	return r.Columns
}

// Fill processes request and fills internal variables
func (r *DmlMappingCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["namespaceHandle"]; ok && len(val) > 0 {
				r.NamespaceHandle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["sourceIdent"]; ok && len(val) > 0 {
				r.SourceIdent, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["moduleHandle"]; ok && len(val) > 0 {
				r.ModuleHandle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["moduleName"]; ok && len(val) > 0 {
				r.ModuleName, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["skip"]; ok && len(val) > 0 {
				r.Skip, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["identifier"]; ok && len(val) > 0 {
				r.Identifier, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["columns[]"]; ok {
				r.Columns, err = types.ParseDmlColumnMapSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["columns"]; ok {
				r.Columns, err = types.ParseDmlColumnMapSet(val)
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

		if val, ok := req.Form["namespaceHandle"]; ok && len(val) > 0 {
			r.NamespaceHandle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["sourceIdent"]; ok && len(val) > 0 {
			r.SourceIdent, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["moduleHandle"]; ok && len(val) > 0 {
			r.ModuleHandle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["moduleName"]; ok && len(val) > 0 {
			r.ModuleName, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["skip"]; ok && len(val) > 0 {
			r.Skip, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["identifier"]; ok && len(val) > 0 {
			r.Identifier, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["columns[]"]; ok {
			r.Columns, err = types.ParseDmlColumnMapSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["columns"]; ok {
			r.Columns, err = types.ParseDmlColumnMapSet(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlMappingRead request
func NewDmlMappingRead() *DmlMappingRead {
	return &DmlMappingRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"mappingID":    r.MappingID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingRead) GetMappingID() uint64 {
	return r.MappingID
}

// Fill processes request and fills internal variables
func (r *DmlMappingRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "mappingID")
		r.MappingID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlMappingUpdate request
func NewDmlMappingUpdate() *DmlMappingUpdate {
	return &DmlMappingUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID":    r.ConnectionID,
		"mappingID":       r.MappingID,
		"namespaceHandle": r.NamespaceHandle,
		"sourceIdent":     r.SourceIdent,
		"moduleHandle":    r.ModuleHandle,
		"moduleName":      r.ModuleName,
		"skip":            r.Skip,
		"identifier":      r.Identifier,
		"columns":         r.Columns,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetMappingID() uint64 {
	return r.MappingID
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetNamespaceHandle() string {
	return r.NamespaceHandle
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetSourceIdent() string {
	return r.SourceIdent
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetModuleHandle() string {
	return r.ModuleHandle
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetModuleName() string {
	return r.ModuleName
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetSkip() bool {
	return r.Skip
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetIdentifier() string {
	return r.Identifier
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingUpdate) GetColumns() types.DmlColumnMapSet {
	return r.Columns
}

// Fill processes request and fills internal variables
func (r *DmlMappingUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["namespaceHandle"]; ok && len(val) > 0 {
				r.NamespaceHandle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["sourceIdent"]; ok && len(val) > 0 {
				r.SourceIdent, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["moduleHandle"]; ok && len(val) > 0 {
				r.ModuleHandle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["moduleName"]; ok && len(val) > 0 {
				r.ModuleName, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["skip"]; ok && len(val) > 0 {
				r.Skip, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["identifier"]; ok && len(val) > 0 {
				r.Identifier, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["columns[]"]; ok {
				r.Columns, err = types.ParseDmlColumnMapSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["columns"]; ok {
				r.Columns, err = types.ParseDmlColumnMapSet(val)
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

		if val, ok := req.Form["namespaceHandle"]; ok && len(val) > 0 {
			r.NamespaceHandle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["sourceIdent"]; ok && len(val) > 0 {
			r.SourceIdent, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["moduleHandle"]; ok && len(val) > 0 {
			r.ModuleHandle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["moduleName"]; ok && len(val) > 0 {
			r.ModuleName, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["skip"]; ok && len(val) > 0 {
			r.Skip, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["identifier"]; ok && len(val) > 0 {
			r.Identifier, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["columns[]"]; ok {
			r.Columns, err = types.ParseDmlColumnMapSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["columns"]; ok {
			r.Columns, err = types.ParseDmlColumnMapSet(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "mappingID")
		r.MappingID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlMappingDelete request
func NewDmlMappingDelete() *DmlMappingDelete {
	return &DmlMappingDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"mappingID":    r.MappingID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingDelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlMappingDelete) GetMappingID() uint64 {
	return r.MappingID
}

// Fill processes request and fills internal variables
func (r *DmlMappingDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "mappingID")
		r.MappingID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlImportRun request
func NewDmlImportRun() *DmlImportRun {
	return &DmlImportRun{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRun) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"mappingID":    r.MappingID,
		"method":       r.Method,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRun) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRun) GetMappingID() uint64 {
	return r.MappingID
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRun) GetMethod() types.DmlImportMethod {
	return r.Method
}

// Fill processes request and fills internal variables
func (r *DmlImportRun) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["method"]; ok && len(val) > 0 {
				r.Method, err = types.DmlImportMethod(val[0]), nil
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

		if val, ok := req.Form["method"]; ok && len(val) > 0 {
			r.Method, err = types.DmlImportMethod(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "mappingID")
		r.MappingID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDmlImportRunRead request
func NewDmlImportRunRead() *DmlImportRunRead {
	return &DmlImportRunRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRunRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"mappingID":    r.MappingID,
		"runID":        r.RunID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRunRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRunRead) GetMappingID() uint64 {
	return r.MappingID
}

// Auditable returns all auditable/loggable parameters
func (r DmlImportRunRead) GetRunID() uint64 {
	return r.RunID
}

// Fill processes request and fills internal variables
func (r *DmlImportRunRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "mappingID")
		r.MappingID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "runID")
		r.RunID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
