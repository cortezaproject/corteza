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
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/crusttech/human/server/system/types"
	"github.com/go-chi/chi/v5"
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
	ConnectionList struct {
		// Handle GET parameter
		//
		// Filter by handle
		Handle string

		// Status GET parameter
		//
		// Filter by status
		Status []string

		// Query GET parameter
		//
		// Search query
		Query string

		// Tags GET parameter
		//
		// Filter by tags
		Tags []string

		// Source GET parameter
		//
		// Filter by source ("catalog", "local", or "" for all)
		Source string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted connections
		Deleted uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

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

	ConnectionCreate struct {
		// Handle POST parameter
		//
		// Connection handle
		Handle string

		// Meta POST parameter
		//
		// Meta
		Meta types.ConnectionMeta

		// Service POST parameter
		//
		// Service config
		Service types.ConnectionService

		// Resources POST parameter
		//
		// Resources
		Resources json.RawMessage

		// StandardOperations POST parameter
		//
		// Standard operations
		StandardOperations json.RawMessage

		// Operations POST parameter
		//
		// Operations
		Operations json.RawMessage

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	ConnectionUpdate struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Connection handle
		Handle string

		// Meta POST parameter
		//
		// Meta
		Meta types.ConnectionMeta

		// Service POST parameter
		//
		// Service config
		Service types.ConnectionService

		// Resources POST parameter
		//
		// Resources
		Resources json.RawMessage

		// StandardOperations POST parameter
		//
		// Standard operations
		StandardOperations json.RawMessage

		// Operations POST parameter
		//
		// Operations
		Operations json.RawMessage

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	ConnectionRead struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	ConnectionDelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	ConnectionUndelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	ConnectionEnable struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	ConnectionGenerate struct {
		// Prompt POST parameter
		//
		// Description or API docs
		Prompt string

		// Context POST parameter
		//
		// Optional additional context
		Context string
	}

	ConnectionImport struct {
		// CatalogID POST parameter
		//
		// Catalog ID
		CatalogID string
	}

	ConnectionConfigure struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`

		// CatalogID POST parameter
		//
		// Catalog ID to auto-import on first configure
		CatalogID string

		// Name POST parameter
		//
		// Connection name
		Name string

		// Config POST parameter
		//
		// Connection config
		Config types.ConfiguredConnectionConfig

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	ConnectionUpdateConfiguration struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`

		// ConfiguredConnectionID PATH parameter
		//
		// Configured Connection ID
		ConfiguredConnectionID uint64 `json:",string"`

		// Name POST parameter
		//
		// Connection name
		Name string

		// Config POST parameter
		//
		// Connection config
		Config types.ConfiguredConnectionConfig

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}
)

// NewConnectionList request
func NewConnectionList() *ConnectionList {
	return &ConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":     r.Handle,
		"status":     r.Status,
		"query":      r.Query,
		"tags":       r.Tags,
		"source":     r.Source,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetStatus() []string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetTags() []string {
	return r.Tags
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetSource() string {
	return r.Source
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ConnectionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status[]"]; ok {
			r.Status, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["status"]; ok {
			r.Status, err = val, nil
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
		if val, ok := tmp["tags[]"]; ok {
			r.Tags, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["tags"]; ok {
			r.Tags, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["source"]; ok && len(val) > 0 {
			r.Source, err = val[0], nil
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

// NewConnectionCreate request
func NewConnectionCreate() *ConnectionCreate {
	return &ConnectionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":             r.Handle,
		"meta":               r.Meta,
		"service":            r.Service,
		"resources":          r.Resources,
		"standardOperations": r.StandardOperations,
		"operations":         r.Operations,
		"labels":             r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetMeta() types.ConnectionMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetService() types.ConnectionService {
	return r.Service
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetResources() json.RawMessage {
	return r.Resources
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetStandardOperations() json.RawMessage {
	return r.StandardOperations
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetOperations() json.RawMessage {
	return r.Operations
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ConnectionCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["service[]"]; ok {
				r.Service, err = types.ParseConnectionService(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["service"]; ok {
				r.Service, err = types.ParseConnectionService(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["resources"]; ok && len(val) > 0 {
				r.Resources, err = json.RawMessage(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["standardOperations"]; ok && len(val) > 0 {
				r.StandardOperations, err = json.RawMessage(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["operations"]; ok && len(val) > 0 {
				r.Operations, err = json.RawMessage(val[0]), nil
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["service[]"]; ok {
			r.Service, err = types.ParseConnectionService(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["service"]; ok {
			r.Service, err = types.ParseConnectionService(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resources"]; ok && len(val) > 0 {
			r.Resources, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["standardOperations"]; ok && len(val) > 0 {
			r.StandardOperations, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["operations"]; ok && len(val) > 0 {
			r.Operations, err = json.RawMessage(val[0]), nil
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

// NewConnectionUpdate request
func NewConnectionUpdate() *ConnectionUpdate {
	return &ConnectionUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID":       r.ConnectionID,
		"handle":             r.Handle,
		"meta":               r.Meta,
		"service":            r.Service,
		"resources":          r.Resources,
		"standardOperations": r.StandardOperations,
		"operations":         r.Operations,
		"labels":             r.Labels,
		"updatedAt":          r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetMeta() types.ConnectionMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetService() types.ConnectionService {
	return r.Service
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetResources() json.RawMessage {
	return r.Resources
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetStandardOperations() json.RawMessage {
	return r.StandardOperations
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetOperations() json.RawMessage {
	return r.Operations
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *ConnectionUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseConnectionMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["service[]"]; ok {
				r.Service, err = types.ParseConnectionService(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["service"]; ok {
				r.Service, err = types.ParseConnectionService(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["resources"]; ok && len(val) > 0 {
				r.Resources, err = json.RawMessage(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["standardOperations"]; ok && len(val) > 0 {
				r.StandardOperations, err = json.RawMessage(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["operations"]; ok && len(val) > 0 {
				r.Operations, err = json.RawMessage(val[0]), nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseConnectionMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["service[]"]; ok {
			r.Service, err = types.ParseConnectionService(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["service"]; ok {
			r.Service, err = types.ParseConnectionService(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["resources"]; ok && len(val) > 0 {
			r.Resources, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["standardOperations"]; ok && len(val) > 0 {
			r.StandardOperations, err = json.RawMessage(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["operations"]; ok && len(val) > 0 {
			r.Operations, err = json.RawMessage(val[0]), nil
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

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewConnectionRead request
func NewConnectionRead() *ConnectionRead {
	return &ConnectionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *ConnectionRead) Fill(req *http.Request) (err error) {

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

// NewConnectionDelete request
func NewConnectionDelete() *ConnectionDelete {
	return &ConnectionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionDelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *ConnectionDelete) Fill(req *http.Request) (err error) {

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

// NewConnectionUndelete request
func NewConnectionUndelete() *ConnectionUndelete {
	return &ConnectionUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUndelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *ConnectionUndelete) Fill(req *http.Request) (err error) {

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

// NewConnectionEnable request
func NewConnectionEnable() *ConnectionEnable {
	return &ConnectionEnable{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionEnable) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionEnable) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *ConnectionEnable) Fill(req *http.Request) (err error) {

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

// NewConnectionGenerate request
func NewConnectionGenerate() *ConnectionGenerate {
	return &ConnectionGenerate{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionGenerate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"prompt":  r.Prompt,
		"context": r.Context,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionGenerate) GetPrompt() string {
	return r.Prompt
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionGenerate) GetContext() string {
	return r.Context
}

// Fill processes request and fills internal variables
func (r *ConnectionGenerate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["prompt"]; ok && len(val) > 0 {
				r.Prompt, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["context"]; ok && len(val) > 0 {
				r.Context, err = val[0], nil
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

		if val, ok := req.Form["prompt"]; ok && len(val) > 0 {
			r.Prompt, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["context"]; ok && len(val) > 0 {
			r.Context, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewConnectionImport request
func NewConnectionImport() *ConnectionImport {
	return &ConnectionImport{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionImport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"catalogID": r.CatalogID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionImport) GetCatalogID() string {
	return r.CatalogID
}

// Fill processes request and fills internal variables
func (r *ConnectionImport) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["catalogID"]; ok && len(val) > 0 {
				r.CatalogID, err = val[0], nil
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

		if val, ok := req.Form["catalogID"]; ok && len(val) > 0 {
			r.CatalogID, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewConnectionConfigure request
func NewConnectionConfigure() *ConnectionConfigure {
	return &ConnectionConfigure{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionConfigure) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"catalogID":    r.CatalogID,
		"name":         r.Name,
		"config":       r.Config,
		"labels":       r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionConfigure) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionConfigure) GetCatalogID() string {
	return r.CatalogID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionConfigure) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionConfigure) GetConfig() types.ConfiguredConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionConfigure) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ConnectionConfigure) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["catalogID"]; ok && len(val) > 0 {
				r.CatalogID, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["catalogID"]; ok && len(val) > 0 {
				r.CatalogID, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseConfiguredConnectionConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseConfiguredConnectionConfig(val)
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

		if val, ok := req.Form["catalogID"]; ok && len(val) > 0 {
			r.CatalogID, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["catalogID"]; ok && len(val) > 0 {
			r.CatalogID, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseConfiguredConnectionConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseConfiguredConnectionConfig(val)
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

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewConnectionUpdateConfiguration request
func NewConnectionUpdateConfiguration() *ConnectionUpdateConfiguration {
	return &ConnectionUpdateConfiguration{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdateConfiguration) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID":           r.ConnectionID,
		"configuredConnectionID": r.ConfiguredConnectionID,
		"name":                   r.Name,
		"config":                 r.Config,
		"labels":                 r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdateConfiguration) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdateConfiguration) GetConfiguredConnectionID() uint64 {
	return r.ConfiguredConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdateConfiguration) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdateConfiguration) GetConfig() types.ConfiguredConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdateConfiguration) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ConnectionUpdateConfiguration) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseConfiguredConnectionConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseConfiguredConnectionConfig(val)
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseConfiguredConnectionConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseConfiguredConnectionConfig(val)
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

		val = chi.URLParam(req, "connectionID")
		r.ConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "configuredConnectionID")
		r.ConfiguredConnectionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
