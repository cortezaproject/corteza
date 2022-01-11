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
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	sqlxTypes "github.com/jmoiron/sqlx/types"
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
	PageList struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// SelfID GET parameter
		//
		// Parent page ID
		SelfID uint64 `json:",string"`

		// ModuleID GET parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query
		Query string

		// Handle GET parameter
		//
		// Search by handle
		Handle string

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	PageCreate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// SelfID POST parameter
		//
		// Parent Page ID
		SelfID uint64 `json:",string"`

		// ModuleID POST parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Description POST parameter
		//
		// Description
		Description string

		// Weight POST parameter
		//
		// Page tree weight
		Weight int

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Visible POST parameter
		//
		// Visible in navigation
		Visible bool

		// Blocks POST parameter
		//
		// Blocks JSON
		Blocks sqlxTypes.JSONText
	}

	PageRead struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`
	}

	PageTree struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`
	}

	PageUpdate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// SelfID POST parameter
		//
		// Parent Page ID
		SelfID uint64 `json:",string"`

		// ModuleID POST parameter
		//
		// Module ID (optional)
		ModuleID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Description POST parameter
		//
		// Description
		Description string

		// Weight POST parameter
		//
		// Page tree weight
		Weight int

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// Visible POST parameter
		//
		// Visible in navigation
		Visible bool

		// Blocks POST parameter
		//
		// Blocks JSON
		Blocks sqlxTypes.JSONText
	}

	PageReorder struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// SelfID PATH parameter
		//
		// Parent page ID
		SelfID uint64 `json:",string"`

		// PageIDs POST parameter
		//
		// Page ID order
		PageIDs []string
	}

	PageDelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`
	}

	PageUpload struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// Upload POST parameter
		//
		// File to upload
		Upload *multipart.FileHeader
	}

	PageTriggerScript struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string

		// Args POST parameter
		//
		// Arguments to pass to the script
		Args map[string]interface{}
	}

	PageListTranslations struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// ID
		PageID uint64 `json:",string"`
	}

	PageUpdateTranslations struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// ID
		PageID uint64 `json:",string"`

		// Translations POST parameter
		//
		// Resource translation to upsert
		Translations locale.ResourceTranslationSet
	}
)

// NewPageList request
func NewPageList() *PageList {
	return &PageList{}
}

// Auditable returns all auditable/loggable parameters
func (r PageList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"selfID":      r.SelfID,
		"moduleID":    r.ModuleID,
		"query":       r.Query,
		"handle":      r.Handle,
		"labels":      r.Labels,
		"limit":       r.Limit,
		"pageCursor":  r.PageCursor,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetSelfID() uint64 {
	return r.SelfID
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r PageList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *PageList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["selfID"]; ok && len(val) > 0 {
			r.SelfID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["moduleID"]; ok && len(val) > 0 {
			r.ModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
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

	}

	return err
}

// NewPageCreate request
func NewPageCreate() *PageCreate {
	return &PageCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"selfID":      r.SelfID,
		"moduleID":    r.ModuleID,
		"title":       r.Title,
		"handle":      r.Handle,
		"description": r.Description,
		"weight":      r.Weight,
		"labels":      r.Labels,
		"visible":     r.Visible,
		"blocks":      r.Blocks,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetSelfID() uint64 {
	return r.SelfID
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetWeight() int {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetVisible() bool {
	return r.Visible
}

// Auditable returns all auditable/loggable parameters
func (r PageCreate) GetBlocks() sqlxTypes.JSONText {
	return r.Blocks
}

// Fill processes request and fills internal variables
func (r *PageCreate) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["selfID"]; ok && len(val) > 0 {
				r.SelfID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["moduleID"]; ok && len(val) > 0 {
				r.ModuleID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["weight"]; ok && len(val) > 0 {
				r.Weight, err = payload.ParseInt(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["visible"]; ok && len(val) > 0 {
				r.Visible, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["blocks"]; ok && len(val) > 0 {
				r.Blocks, err = payload.ParseJSONTextWithErr(val[0])
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

		if val, ok := req.Form["selfID"]; ok && len(val) > 0 {
			r.SelfID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["moduleID"]; ok && len(val) > 0 {
			r.ModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["weight"]; ok && len(val) > 0 {
			r.Weight, err = payload.ParseInt(val[0]), nil
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

		if val, ok := req.Form["visible"]; ok && len(val) > 0 {
			r.Visible, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["blocks"]; ok && len(val) > 0 {
			r.Blocks, err = payload.ParseJSONTextWithErr(val[0])
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

	}

	return err
}

// NewPageRead request
func NewPageRead() *PageRead {
	return &PageRead{}
}

// Auditable returns all auditable/loggable parameters
func (r PageRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageRead) GetPageID() uint64 {
	return r.PageID
}

// Fill processes request and fills internal variables
func (r *PageRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageTree request
func NewPageTree() *PageTree {
	return &PageTree{}
}

// Auditable returns all auditable/loggable parameters
func (r PageTree) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageTree) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Fill processes request and fills internal variables
func (r *PageTree) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageUpdate request
func NewPageUpdate() *PageUpdate {
	return &PageUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"selfID":      r.SelfID,
		"moduleID":    r.ModuleID,
		"title":       r.Title,
		"handle":      r.Handle,
		"description": r.Description,
		"weight":      r.Weight,
		"labels":      r.Labels,
		"visible":     r.Visible,
		"blocks":      r.Blocks,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetSelfID() uint64 {
	return r.SelfID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetWeight() int {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetVisible() bool {
	return r.Visible
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdate) GetBlocks() sqlxTypes.JSONText {
	return r.Blocks
}

// Fill processes request and fills internal variables
func (r *PageUpdate) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

			if val, ok := req.MultipartForm.Value["selfID"]; ok && len(val) > 0 {
				r.SelfID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["moduleID"]; ok && len(val) > 0 {
				r.ModuleID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["weight"]; ok && len(val) > 0 {
				r.Weight, err = payload.ParseInt(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["visible"]; ok && len(val) > 0 {
				r.Visible, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["blocks"]; ok && len(val) > 0 {
				r.Blocks, err = payload.ParseJSONTextWithErr(val[0])
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

		if val, ok := req.Form["selfID"]; ok && len(val) > 0 {
			r.SelfID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["moduleID"]; ok && len(val) > 0 {
			r.ModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["weight"]; ok && len(val) > 0 {
			r.Weight, err = payload.ParseInt(val[0]), nil
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

		if val, ok := req.Form["visible"]; ok && len(val) > 0 {
			r.Visible, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["blocks"]; ok && len(val) > 0 {
			r.Blocks, err = payload.ParseJSONTextWithErr(val[0])
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

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageReorder request
func NewPageReorder() *PageReorder {
	return &PageReorder{}
}

// Auditable returns all auditable/loggable parameters
func (r PageReorder) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"selfID":      r.SelfID,
		"pageIDs":     r.PageIDs,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageReorder) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageReorder) GetSelfID() uint64 {
	return r.SelfID
}

// Auditable returns all auditable/loggable parameters
func (r PageReorder) GetPageIDs() []string {
	return r.PageIDs
}

// Fill processes request and fills internal variables
func (r *PageReorder) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

		//if val, ok := req.Form["pageIDs[]"]; ok && len(val) > 0  {
		//    r.PageIDs, err = val, nil
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

		val = chi.URLParam(req, "selfID")
		r.SelfID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageDelete request
func NewPageDelete() *PageDelete {
	return &PageDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r PageDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageDelete) GetPageID() uint64 {
	return r.PageID
}

// Fill processes request and fills internal variables
func (r *PageDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageUpload request
func NewPageUpload() *PageUpload {
	return &PageUpload{}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpload) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"upload":      r.Upload,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpload) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpload) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpload) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *PageUpload) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageTriggerScript request
func NewPageTriggerScript() *PageTriggerScript {
	return &PageTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r PageTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"script":      r.Script,
		"args":        r.Args,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageTriggerScript) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageTriggerScript) GetScript() string {
	return r.Script
}

// Auditable returns all auditable/loggable parameters
func (r PageTriggerScript) GetArgs() map[string]interface{} {
	return r.Args
}

// Fill processes request and fills internal variables
func (r *PageTriggerScript) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageListTranslations request
func NewPageListTranslations() *PageListTranslations {
	return &PageListTranslations{}
}

// Auditable returns all auditable/loggable parameters
func (r PageListTranslations) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageListTranslations) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageListTranslations) GetPageID() uint64 {
	return r.PageID
}

// Fill processes request and fills internal variables
func (r *PageListTranslations) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageUpdateTranslations request
func NewPageUpdateTranslations() *PageUpdateTranslations {
	return &PageUpdateTranslations{}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdateTranslations) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"translations": r.Translations,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdateTranslations) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdateTranslations) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageUpdateTranslations) GetTranslations() locale.ResourceTranslationSet {
	return r.Translations
}

// Fill processes request and fills internal variables
func (r *PageUpdateTranslations) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
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

		//if val, ok := req.Form["translations[]"]; ok && len(val) > 0  {
		//    r.Translations, err = locale.ResourceTranslationSet(val), nil
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

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
