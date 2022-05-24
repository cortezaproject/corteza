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
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/types"
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
	DalConnectionList struct {
		// ConnectionID GET parameter
		//
		// Filter by connection ID
		ConnectionID []string

		// Handle GET parameter
		//
		// Search handle to match against connections
		Handle string

		// Location GET parameter
		//
		// Search location to match against connections
		Location string

		// Ownership GET parameter
		//
		// Search ownership to match against connections
		Ownership string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted connections
		Deleted uint

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

	DalConnectionCreate struct {
		// Handle POST parameter
		//
		// handle
		Handle string

		// Dsn POST parameter
		//
		// dsn
		Dsn string

		// Location POST parameter
		//
		// location
		Location string

		// Ownership POST parameter
		//
		// ownership
		Ownership string

		// Sensitive POST parameter
		//
		// sensitive
		Sensitive bool

		// Config POST parameter
		//
		// config
		Config types.ConnectionConfig

		// Capabilities POST parameter
		//
		// capabilities
		Capabilities types.ConnectionCapabilities

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	DalConnectionUpdate struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`

		// Handle POST parameter
		//
		// handle
		Handle string

		// Dsn POST parameter
		//
		// dsn
		Dsn string

		// Location POST parameter
		//
		// location
		Location string

		// Ownership POST parameter
		//
		// ownership
		Ownership string

		// Sensitive POST parameter
		//
		// sensitive
		Sensitive bool

		// Config POST parameter
		//
		// config
		Config types.ConnectionConfig

		// Capabilities POST parameter
		//
		// capabilities
		Capabilities types.ConnectionCapabilities

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	DalConnectionReadPrimary struct {
	}

	DalConnectionRead struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	DalConnectionDelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	DalConnectionUndelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}
)

// NewDalConnectionList request
func NewDalConnectionList() *DalConnectionList {
	return &DalConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"handle":       r.Handle,
		"location":     r.Location,
		"ownership":    r.Ownership,
		"deleted":      r.Deleted,
		"labels":       r.Labels,
		"limit":        r.Limit,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetConnectionID() []string {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetLocation() string {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *DalConnectionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

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
		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["location"]; ok && len(val) > 0 {
			r.Location, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["ownership"]; ok && len(val) > 0 {
			r.Ownership, err = val[0], nil
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

// NewDalConnectionCreate request
func NewDalConnectionCreate() *DalConnectionCreate {
	return &DalConnectionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":       r.Handle,
		"dsn":          r.Dsn,
		"location":     r.Location,
		"ownership":    r.Ownership,
		"sensitive":    r.Sensitive,
		"config":       r.Config,
		"capabilities": r.Capabilities,
		"labels":       r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetDsn() string {
	return r.Dsn
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetLocation() string {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetSensitive() bool {
	return r.Sensitive
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetCapabilities() types.ConnectionCapabilities {
	return r.Capabilities
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *DalConnectionCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["dsn"]; ok && len(val) > 0 {
				r.Dsn, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["location"]; ok && len(val) > 0 {
				r.Location, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownership"]; ok && len(val) > 0 {
				r.Ownership, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["sensitive"]; ok && len(val) > 0 {
				r.Sensitive, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["capabilities[]"]; ok {
				r.Capabilities, err = types.ParseConnectionCapabilities(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["capabilities"]; ok {
				r.Capabilities, err = types.ParseConnectionCapabilities(val)
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

		if val, ok := req.Form["dsn"]; ok && len(val) > 0 {
			r.Dsn, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["location"]; ok && len(val) > 0 {
			r.Location, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownership"]; ok && len(val) > 0 {
			r.Ownership, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["sensitive"]; ok && len(val) > 0 {
			r.Sensitive, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["capabilities[]"]; ok {
			r.Capabilities, err = types.ParseConnectionCapabilities(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["capabilities"]; ok {
			r.Capabilities, err = types.ParseConnectionCapabilities(val)
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

// NewDalConnectionUpdate request
func NewDalConnectionUpdate() *DalConnectionUpdate {
	return &DalConnectionUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"handle":       r.Handle,
		"dsn":          r.Dsn,
		"location":     r.Location,
		"ownership":    r.Ownership,
		"sensitive":    r.Sensitive,
		"config":       r.Config,
		"capabilities": r.Capabilities,
		"labels":       r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetDsn() string {
	return r.Dsn
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetLocation() string {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetSensitive() bool {
	return r.Sensitive
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetCapabilities() types.ConnectionCapabilities {
	return r.Capabilities
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *DalConnectionUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["dsn"]; ok && len(val) > 0 {
				r.Dsn, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["location"]; ok && len(val) > 0 {
				r.Location, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownership"]; ok && len(val) > 0 {
				r.Ownership, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["sensitive"]; ok && len(val) > 0 {
				r.Sensitive, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config[]"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["config"]; ok {
				r.Config, err = types.ParseConnectionConfig(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["capabilities[]"]; ok {
				r.Capabilities, err = types.ParseConnectionCapabilities(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["capabilities"]; ok {
				r.Capabilities, err = types.ParseConnectionCapabilities(val)
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

		if val, ok := req.Form["dsn"]; ok && len(val) > 0 {
			r.Dsn, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["location"]; ok && len(val) > 0 {
			r.Location, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownership"]; ok && len(val) > 0 {
			r.Ownership, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["sensitive"]; ok && len(val) > 0 {
			r.Sensitive, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config[]"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["config"]; ok {
			r.Config, err = types.ParseConnectionConfig(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["capabilities[]"]; ok {
			r.Capabilities, err = types.ParseConnectionCapabilities(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["capabilities"]; ok {
			r.Capabilities, err = types.ParseConnectionCapabilities(val)
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

// NewDalConnectionReadPrimary request
func NewDalConnectionReadPrimary() *DalConnectionReadPrimary {
	return &DalConnectionReadPrimary{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionReadPrimary) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *DalConnectionReadPrimary) Fill(req *http.Request) (err error) {

	return err
}

// NewDalConnectionRead request
func NewDalConnectionRead() *DalConnectionRead {
	return &DalConnectionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DalConnectionRead) Fill(req *http.Request) (err error) {

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

// NewDalConnectionDelete request
func NewDalConnectionDelete() *DalConnectionDelete {
	return &DalConnectionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionDelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DalConnectionDelete) Fill(req *http.Request) (err error) {

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

// NewDalConnectionUndelete request
func NewDalConnectionUndelete() *DalConnectionUndelete {
	return &DalConnectionUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUndelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *DalConnectionUndelete) Fill(req *http.Request) (err error) {

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
