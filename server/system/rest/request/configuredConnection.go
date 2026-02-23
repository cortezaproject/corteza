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
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
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
	ConfiguredConnectionList struct {
		// ConnectionID GET parameter
		//
		// Filter by connection ID
		ConnectionID uint64 `json:",string"`

		// Status GET parameter
		//
		// Filter by status
		Status []string

		// Query GET parameter
		//
		// Search query
		Query string

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

	ConfiguredConnectionInstall struct {
		// ConnectionID POST parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`

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

	ConfiguredConnectionRead struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}

	ConfiguredConnectionDelete struct {
		// ConnectionID PATH parameter
		//
		// Connection ID
		ConnectionID uint64 `json:",string"`
	}
)

// NewConfiguredConnectionList request
func NewConfiguredConnectionList() *ConfiguredConnectionList {
	return &ConfiguredConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"status":       r.Status,
		"query":        r.Query,
		"deleted":      r.Deleted,
		"labels":       r.Labels,
		"limit":        r.Limit,
		"incTotal":     r.IncTotal,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetStatus() []string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ConfiguredConnectionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["connectionID"]; ok && len(val) > 0 {
			r.ConnectionID, err = payload.ParseUint64(val[0]), nil
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

// NewConfiguredConnectionInstall request
func NewConfiguredConnectionInstall() *ConfiguredConnectionInstall {
	return &ConfiguredConnectionInstall{}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionInstall) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"name":         r.Name,
		"config":       r.Config,
		"labels":       r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionInstall) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionInstall) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionInstall) GetConfig() types.ConfiguredConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionInstall) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ConfiguredConnectionInstall) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["connectionID"]; ok && len(val) > 0 {
				r.ConnectionID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["connectionID"]; ok && len(val) > 0 {
			r.ConnectionID, err = payload.ParseUint64(val[0]), nil
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

	return err
}

// NewConfiguredConnectionRead request
func NewConfiguredConnectionRead() *ConfiguredConnectionRead {
	return &ConfiguredConnectionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionRead) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *ConfiguredConnectionRead) Fill(req *http.Request) (err error) {

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

// NewConfiguredConnectionDelete request
func NewConfiguredConnectionDelete() *ConfiguredConnectionDelete {
	return &ConfiguredConnectionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ConfiguredConnectionDelete) GetConnectionID() uint64 {
	return r.ConnectionID
}

// Fill processes request and fills internal variables
func (r *ConfiguredConnectionDelete) Fill(req *http.Request) (err error) {

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
