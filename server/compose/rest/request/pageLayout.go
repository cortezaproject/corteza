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
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-chi/chi/v5"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
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
	PageLayoutListNamespace struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID GET parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// ModuleID GET parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// ParentID GET parameter
		//
		// Parent ID
		ParentID uint64 `json:",string"`

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

	PageLayoutList struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// ModuleID GET parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// ParentID GET parameter
		//
		// Parent ID
		ParentID uint64 `json:",string"`

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

	PageLayoutCreate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// ParentID POST parameter
		//
		// ParentID
		ParentID uint64 `json:",string"`

		// Weight POST parameter
		//
		// Weight
		Weight int

		// ModuleID POST parameter
		//
		// ModuleID
		ModuleID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Meta POST parameter
		//
		// Meta
		Meta types.PageLayoutMeta

		// Config POST parameter
		//
		// Config
		Config sqlxTypes.JSONText

		// Blocks POST parameter
		//
		// Blocks
		Blocks sqlxTypes.JSONText

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// OwnedBy POST parameter
		//
		// OwnedBy
		OwnedBy uint64 `json:",string"`
	}

	PageLayoutRead struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageLayoutID PATH parameter
		//
		// Page layout ID
		PageLayoutID uint64 `json:",string"`
	}

	PageLayoutUpdate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageLayoutID PATH parameter
		//
		// Page layout ID
		PageLayoutID uint64 `json:",string"`

		// ParentID POST parameter
		//
		// ParentID
		ParentID uint64 `json:",string"`

		// Weight POST parameter
		//
		// Weight
		Weight int

		// ModuleID POST parameter
		//
		// ModuleID
		ModuleID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Handle
		Handle string

		// Meta POST parameter
		//
		// Meta
		Meta types.PageLayoutMeta

		// Config POST parameter
		//
		// Config
		Config sqlxTypes.JSONText

		// Blocks POST parameter
		//
		// Blocks
		Blocks sqlxTypes.JSONText

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// OwnedBy POST parameter
		//
		// OwnedBy
		OwnedBy uint64 `json:",string"`

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	PageLayoutReorder struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageIDs POST parameter
		//
		// Page ID order
		PageIDs []string
	}

	PageLayoutDelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageLayoutID PATH parameter
		//
		// Page layout ID
		PageLayoutID uint64 `json:",string"`

		// Strategy GET parameter
		//
		// Page delete strategy (abort, force, rebase, cascade)
		Strategy string
	}

	PageLayoutUndelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageLayoutID PATH parameter
		//
		// Page layout ID
		PageLayoutID uint64 `json:",string"`
	}

	PageLayoutListTranslations struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageLayoutID PATH parameter
		//
		// ID
		PageLayoutID uint64 `json:",string"`
	}

	PageLayoutUpdateTranslations struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// PageID PATH parameter
		//
		// Page ID
		PageID uint64 `json:",string"`

		// PageLayoutID PATH parameter
		//
		// ID
		PageLayoutID uint64 `json:",string"`

		// Translations POST parameter
		//
		// Resource translation to upsert
		Translations locale.ResourceTranslationSet
	}
)

// NewPageLayoutListNamespace request
func NewPageLayoutListNamespace() *PageLayoutListNamespace {
	return &PageLayoutListNamespace{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"moduleID":    r.ModuleID,
		"parentID":    r.ParentID,
		"query":       r.Query,
		"handle":      r.Handle,
		"labels":      r.Labels,
		"limit":       r.Limit,
		"pageCursor":  r.PageCursor,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetParentID() uint64 {
	return r.ParentID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListNamespace) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *PageLayoutListNamespace) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["pageID"]; ok && len(val) > 0 {
			r.PageID, err = payload.ParseUint64(val[0]), nil
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
		if val, ok := tmp["parentID"]; ok && len(val) > 0 {
			r.ParentID, err = payload.ParseUint64(val[0]), nil
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

// NewPageLayoutList request
func NewPageLayoutList() *PageLayoutList {
	return &PageLayoutList{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"moduleID":    r.ModuleID,
		"parentID":    r.ParentID,
		"query":       r.Query,
		"handle":      r.Handle,
		"labels":      r.Labels,
		"limit":       r.Limit,
		"pageCursor":  r.PageCursor,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetParentID() uint64 {
	return r.ParentID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *PageLayoutList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["moduleID"]; ok && len(val) > 0 {
			r.ModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["parentID"]; ok && len(val) > 0 {
			r.ParentID, err = payload.ParseUint64(val[0]), nil
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

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutCreate request
func NewPageLayoutCreate() *PageLayoutCreate {
	return &PageLayoutCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"parentID":    r.ParentID,
		"weight":      r.Weight,
		"moduleID":    r.ModuleID,
		"handle":      r.Handle,
		"meta":        r.Meta,
		"config":      r.Config,
		"blocks":      r.Blocks,
		"labels":      r.Labels,
		"ownedBy":     r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetParentID() uint64 {
	return r.ParentID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetWeight() int {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetMeta() types.PageLayoutMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetBlocks() sqlxTypes.JSONText {
	return r.Blocks
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *PageLayoutCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["parentID"]; ok && len(val) > 0 {
				r.ParentID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["moduleID"]; ok && len(val) > 0 {
				r.ModuleID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParsePageLayoutMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParsePageLayoutMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config"]; ok && len(val) > 0 {
				r.Config, err = payload.ParseJSONTextWithErr(val[0])
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

			if val, ok := req.MultipartForm.Value["ownedBy"]; ok && len(val) > 0 {
				r.OwnedBy, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["parentID"]; ok && len(val) > 0 {
			r.ParentID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["moduleID"]; ok && len(val) > 0 {
			r.ModuleID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParsePageLayoutMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParsePageLayoutMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
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

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
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

// NewPageLayoutRead request
func NewPageLayoutRead() *PageLayoutRead {
	return &PageLayoutRead{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"pageLayoutID": r.PageLayoutID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutRead) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutRead) GetPageLayoutID() uint64 {
	return r.PageLayoutID
}

// Fill processes request and fills internal variables
func (r *PageLayoutRead) Fill(req *http.Request) (err error) {

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

		val = chi.URLParam(req, "pageLayoutID")
		r.PageLayoutID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutUpdate request
func NewPageLayoutUpdate() *PageLayoutUpdate {
	return &PageLayoutUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"pageLayoutID": r.PageLayoutID,
		"parentID":     r.ParentID,
		"weight":       r.Weight,
		"moduleID":     r.ModuleID,
		"handle":       r.Handle,
		"meta":         r.Meta,
		"config":       r.Config,
		"blocks":       r.Blocks,
		"labels":       r.Labels,
		"ownedBy":      r.OwnedBy,
		"updatedAt":    r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetPageLayoutID() uint64 {
	return r.PageLayoutID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetParentID() uint64 {
	return r.ParentID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetWeight() int {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetMeta() types.PageLayoutMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetBlocks() sqlxTypes.JSONText {
	return r.Blocks
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *PageLayoutUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["parentID"]; ok && len(val) > 0 {
				r.ParentID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["moduleID"]; ok && len(val) > 0 {
				r.ModuleID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParsePageLayoutMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParsePageLayoutMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config"]; ok && len(val) > 0 {
				r.Config, err = payload.ParseJSONTextWithErr(val[0])
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

			if val, ok := req.MultipartForm.Value["ownedBy"]; ok && len(val) > 0 {
				r.OwnedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["updatedAt"]; ok && len(val) > 0 {
				r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
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

		if val, ok := req.Form["parentID"]; ok && len(val) > 0 {
			r.ParentID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["moduleID"]; ok && len(val) > 0 {
			r.ModuleID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParsePageLayoutMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParsePageLayoutMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
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

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["updatedAt"]; ok && len(val) > 0 {
			r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
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

		val = chi.URLParam(req, "pageLayoutID")
		r.PageLayoutID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutReorder request
func NewPageLayoutReorder() *PageLayoutReorder {
	return &PageLayoutReorder{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutReorder) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"pageID":      r.PageID,
		"pageIDs":     r.PageIDs,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutReorder) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutReorder) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutReorder) GetPageIDs() []string {
	return r.PageIDs
}

// Fill processes request and fills internal variables
func (r *PageLayoutReorder) Fill(req *http.Request) (err error) {

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

		val = chi.URLParam(req, "pageID")
		r.PageID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutDelete request
func NewPageLayoutDelete() *PageLayoutDelete {
	return &PageLayoutDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"pageLayoutID": r.PageLayoutID,
		"strategy":     r.Strategy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutDelete) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutDelete) GetPageLayoutID() uint64 {
	return r.PageLayoutID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutDelete) GetStrategy() string {
	return r.Strategy
}

// Fill processes request and fills internal variables
func (r *PageLayoutDelete) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["strategy"]; ok && len(val) > 0 {
			r.Strategy, err = val[0], nil
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

		val = chi.URLParam(req, "pageLayoutID")
		r.PageLayoutID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutUndelete request
func NewPageLayoutUndelete() *PageLayoutUndelete {
	return &PageLayoutUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"pageLayoutID": r.PageLayoutID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUndelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUndelete) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUndelete) GetPageLayoutID() uint64 {
	return r.PageLayoutID
}

// Fill processes request and fills internal variables
func (r *PageLayoutUndelete) Fill(req *http.Request) (err error) {

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

		val = chi.URLParam(req, "pageLayoutID")
		r.PageLayoutID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutListTranslations request
func NewPageLayoutListTranslations() *PageLayoutListTranslations {
	return &PageLayoutListTranslations{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListTranslations) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"pageLayoutID": r.PageLayoutID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListTranslations) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListTranslations) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutListTranslations) GetPageLayoutID() uint64 {
	return r.PageLayoutID
}

// Fill processes request and fills internal variables
func (r *PageLayoutListTranslations) Fill(req *http.Request) (err error) {

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

		val = chi.URLParam(req, "pageLayoutID")
		r.PageLayoutID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewPageLayoutUpdateTranslations request
func NewPageLayoutUpdateTranslations() *PageLayoutUpdateTranslations {
	return &PageLayoutUpdateTranslations{}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdateTranslations) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID":  r.NamespaceID,
		"pageID":       r.PageID,
		"pageLayoutID": r.PageLayoutID,
		"translations": r.Translations,
	}
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdateTranslations) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdateTranslations) GetPageID() uint64 {
	return r.PageID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdateTranslations) GetPageLayoutID() uint64 {
	return r.PageLayoutID
}

// Auditable returns all auditable/loggable parameters
func (r PageLayoutUpdateTranslations) GetTranslations() locale.ResourceTranslationSet {
	return r.Translations
}

// Fill processes request and fills internal variables
func (r *PageLayoutUpdateTranslations) Fill(req *http.Request) (err error) {

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

		val = chi.URLParam(req, "pageLayoutID")
		r.PageLayoutID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
