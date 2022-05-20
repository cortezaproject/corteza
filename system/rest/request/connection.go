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
	ConnectionList struct {
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

	ConnectionCreate struct {
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

	ConnectionUpdate struct {
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
)

// NewConnectionList request
func NewConnectionList() *ConnectionList {
	return &ConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) Auditable() map[string]interface{} {
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
func (r ConnectionList) GetConnectionID() []string {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetLocation() string {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionList) GetLimit() uint {
	return r.Limit
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

// NewConnectionCreate request
func NewConnectionCreate() *ConnectionCreate {
	return &ConnectionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) Auditable() map[string]interface{} {
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
func (r ConnectionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetDsn() string {
	return r.Dsn
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetLocation() string {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetSensitive() bool {
	return r.Sensitive
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetCapabilities() types.ConnectionCapabilities {
	return r.Capabilities
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionCreate) GetLabels() map[string]string {
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

// NewConnectionUpdate request
func NewConnectionUpdate() *ConnectionUpdate {
	return &ConnectionUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) Auditable() map[string]interface{} {
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
func (r ConnectionUpdate) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetDsn() string {
	return r.Dsn
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetLocation() string {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetSensitive() bool {
	return r.Sensitive
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetCapabilities() types.ConnectionCapabilities {
	return r.Capabilities
}

// Auditable returns all auditable/loggable parameters
func (r ConnectionUpdate) GetLabels() map[string]string {
	return r.Labels
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
