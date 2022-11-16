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
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/types"
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
	TemplateList struct {
		// Query GET parameter
		//
		// Query
		Query string

		// Handle GET parameter
		//
		// Handle
		Handle string

		// Type GET parameter
		//
		// Type
		Type string

		// OwnerID GET parameter
		//
		// OwnerID
		OwnerID uint64 `json:",string"`

		// Partial GET parameter
		//
		// Show partial templates
		Partial bool

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted templates
		Deleted uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total counter
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	TemplateCreate struct {
		// Handle POST parameter
		//
		// Handle
		Handle string

		// Language POST parameter
		//
		// Language
		Language string

		// Type POST parameter
		//
		// Type
		Type string

		// Partial POST parameter
		//
		// Partial
		Partial bool

		// Meta POST parameter
		//
		// Meta
		Meta types.TemplateMeta

		// Template POST parameter
		//
		// Template
		Template string

		// OwnerID POST parameter
		//
		// OwnerID
		OwnerID uint64 `json:",string"`

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	TemplateRead struct {
		// TemplateID PATH parameter
		//
		// ID
		TemplateID uint64 `json:",string"`
	}

	TemplateUpdate struct {
		// TemplateID PATH parameter
		//
		// ID
		TemplateID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Language POST parameter
		//
		// Language
		Language string

		// Type POST parameter
		//
		// Type
		Type string

		// Partial POST parameter
		//
		// Partial
		Partial bool

		// Meta POST parameter
		//
		// Meta
		Meta types.TemplateMeta

		// Template POST parameter
		//
		// Template
		Template string

		// OwnerID POST parameter
		//
		// OwnerID
		OwnerID uint64 `json:",string"`

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	TemplateDelete struct {
		// TemplateID PATH parameter
		//
		// ID
		TemplateID uint64 `json:",string"`
	}

	TemplateUndelete struct {
		// TemplateID PATH parameter
		//
		// Template ID
		TemplateID uint64 `json:",string"`
	}

	TemplateRenderDrivers struct {
	}

	TemplateRender struct {
		// TemplateID PATH parameter
		//
		// Render template to use
		TemplateID uint64 `json:",string"`

		// Filename PATH parameter
		//
		// Filename to use
		Filename string

		// Ext PATH parameter
		//
		// Export format
		Ext string

		// Variables POST parameter
		//
		// Variables defined by import file
		Variables json.RawMessage

		// Options POST parameter
		//
		// Rendering options
		Options json.RawMessage
	}
)

// NewTemplateList request
func NewTemplateList() *TemplateList {
	return &TemplateList{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"handle":     r.Handle,
		"type":       r.Type,
		"ownerID":    r.OwnerID,
		"partial":    r.Partial,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetOwnerID() uint64 {
	return r.OwnerID
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetPartial() bool {
	return r.Partial
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r TemplateList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *TemplateList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
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
		if val, ok := tmp["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["partial"]; ok && len(val) > 0 {
			r.Partial, err = payload.ParseBool(val[0]), nil
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

	return err
}

// NewTemplateCreate request
func NewTemplateCreate() *TemplateCreate {
	return &TemplateCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":   r.Handle,
		"language": r.Language,
		"type":     r.Type,
		"partial":  r.Partial,
		"meta":     r.Meta,
		"template": r.Template,
		"ownerID":  r.OwnerID,
		"labels":   r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetLanguage() string {
	return r.Language
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetPartial() bool {
	return r.Partial
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetMeta() types.TemplateMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetTemplate() string {
	return r.Template
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetOwnerID() uint64 {
	return r.OwnerID
}

// Auditable returns all auditable/loggable parameters
func (r TemplateCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *TemplateCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["language"]; ok && len(val) > 0 {
				r.Language, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["type"]; ok && len(val) > 0 {
				r.Type, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["partial"]; ok && len(val) > 0 {
				r.Partial, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseTemplateMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseTemplateMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["template"]; ok && len(val) > 0 {
				r.Template, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownerID"]; ok && len(val) > 0 {
				r.OwnerID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["language"]; ok && len(val) > 0 {
			r.Language, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["partial"]; ok && len(val) > 0 {
			r.Partial, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTemplateMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTemplateMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["template"]; ok && len(val) > 0 {
			r.Template, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

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

	return err
}

// NewTemplateRead request
func NewTemplateRead() *TemplateRead {
	return &TemplateRead{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"templateID": r.TemplateID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRead) GetTemplateID() uint64 {
	return r.TemplateID
}

// Fill processes request and fills internal variables
func (r *TemplateRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "templateID")
		r.TemplateID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTemplateUpdate request
func NewTemplateUpdate() *TemplateUpdate {
	return &TemplateUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"templateID": r.TemplateID,
		"handle":     r.Handle,
		"language":   r.Language,
		"type":       r.Type,
		"partial":    r.Partial,
		"meta":       r.Meta,
		"template":   r.Template,
		"ownerID":    r.OwnerID,
		"labels":     r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetTemplateID() uint64 {
	return r.TemplateID
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetLanguage() string {
	return r.Language
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetPartial() bool {
	return r.Partial
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetMeta() types.TemplateMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetTemplate() string {
	return r.Template
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetOwnerID() uint64 {
	return r.OwnerID
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *TemplateUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["language"]; ok && len(val) > 0 {
				r.Language, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["type"]; ok && len(val) > 0 {
				r.Type, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["partial"]; ok && len(val) > 0 {
				r.Partial, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseTemplateMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseTemplateMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["template"]; ok && len(val) > 0 {
				r.Template, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownerID"]; ok && len(val) > 0 {
				r.OwnerID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["language"]; ok && len(val) > 0 {
			r.Language, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["partial"]; ok && len(val) > 0 {
			r.Partial, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTemplateMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTemplateMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["template"]; ok && len(val) > 0 {
			r.Template, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownerID"]; ok && len(val) > 0 {
			r.OwnerID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

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

		val = chi.URLParam(req, "templateID")
		r.TemplateID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTemplateDelete request
func NewTemplateDelete() *TemplateDelete {
	return &TemplateDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"templateID": r.TemplateID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateDelete) GetTemplateID() uint64 {
	return r.TemplateID
}

// Fill processes request and fills internal variables
func (r *TemplateDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "templateID")
		r.TemplateID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTemplateUndelete request
func NewTemplateUndelete() *TemplateUndelete {
	return &TemplateUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"templateID": r.TemplateID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateUndelete) GetTemplateID() uint64 {
	return r.TemplateID
}

// Fill processes request and fills internal variables
func (r *TemplateUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "templateID")
		r.TemplateID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTemplateRenderDrivers request
func NewTemplateRenderDrivers() *TemplateRenderDrivers {
	return &TemplateRenderDrivers{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRenderDrivers) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *TemplateRenderDrivers) Fill(req *http.Request) (err error) {

	return err
}

// NewTemplateRender request
func NewTemplateRender() *TemplateRender {
	return &TemplateRender{}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRender) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"templateID": r.TemplateID,
		"filename":   r.Filename,
		"ext":        r.Ext,
		"variables":  r.Variables,
		"options":    r.Options,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRender) GetTemplateID() uint64 {
	return r.TemplateID
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRender) GetFilename() string {
	return r.Filename
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRender) GetExt() string {
	return r.Ext
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRender) GetVariables() json.RawMessage {
	return r.Variables
}

// Auditable returns all auditable/loggable parameters
func (r TemplateRender) GetOptions() json.RawMessage {
	return r.Options
}

// Fill processes request and fills internal variables
func (r *TemplateRender) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["variables"]; ok && len(val) > 0 {
				r.Variables, err = json.RawMessage(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["options"]; ok && len(val) > 0 {
				r.Options, err = json.RawMessage(val[0]), nil
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

		if val, ok := req.Form["variables"]; ok && len(val) > 0 {
			r.Variables, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["options"]; ok && len(val) > 0 {
			r.Options, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "templateID")
		r.TemplateID, err = payload.ParseUint64(val), nil
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
