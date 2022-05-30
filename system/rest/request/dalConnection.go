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
	"github.com/cortezaproject/corteza-server/pkg/geolocation"
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

		// Type GET parameter
		//
		// Search type to match against connections
		Type string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted connections
		Deleted uint
	}

	DalConnectionCreate struct {
		// Handle POST parameter
		//
		// handle
		Handle string

		// Name POST parameter
		//
		// name
		Name string

		// Type POST parameter
		//
		// type
		Type string

		// Location POST parameter
		//
		// location
		Location geolocation.Full

		// Ownership POST parameter
		//
		// ownership
		Ownership string

		// SensitivityLevel POST parameter
		//
		// sensitivityLevel
		SensitivityLevel uint64 `json:",string"`

		// Config POST parameter
		//
		// config
		Config types.ConnectionConfig

		// Capabilities POST parameter
		//
		// capabilities
		Capabilities types.ConnectionCapabilities
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

		// Name POST parameter
		//
		// name
		Name string

		// Type POST parameter
		//
		// type
		Type string

		// Location POST parameter
		//
		// location
		Location geolocation.Full

		// Ownership POST parameter
		//
		// ownership
		Ownership string

		// SensitivityLevel POST parameter
		//
		// sensitivityLevel
		SensitivityLevel uint64 `json:",string"`

		// Config POST parameter
		//
		// config
		Config types.ConnectionConfig

		// Capabilities POST parameter
		//
		// capabilities
		Capabilities types.ConnectionCapabilities
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
		"type":         r.Type,
		"deleted":      r.Deleted,
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
func (r DalConnectionList) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionList) GetDeleted() uint {
	return r.Deleted
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
		if val, ok := tmp["type"]; ok && len(val) > 0 {
			r.Type, err = val[0], nil
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
		"handle":           r.Handle,
		"name":             r.Name,
		"type":             r.Type,
		"location":         r.Location,
		"ownership":        r.Ownership,
		"sensitivityLevel": r.SensitivityLevel,
		"config":           r.Config,
		"capabilities":     r.Capabilities,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetLocation() geolocation.Full {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetSensitivityLevel() uint64 {
	return r.SensitivityLevel
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionCreate) GetCapabilities() types.ConnectionCapabilities {
	return r.Capabilities
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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["location[]"]; ok {
				r.Location, err = geolocation.Parse(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["location"]; ok {
				r.Location, err = geolocation.Parse(val)
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

			if val, ok := req.MultipartForm.Value["sensitivityLevel"]; ok && len(val) > 0 {
				r.SensitivityLevel, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
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

		if val, ok := req.Form["location[]"]; ok {
			r.Location, err = geolocation.Parse(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["location"]; ok {
			r.Location, err = geolocation.Parse(val)
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

		if val, ok := req.Form["sensitivityLevel"]; ok && len(val) > 0 {
			r.SensitivityLevel, err = payload.ParseUint64(val[0]), nil
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
		"connectionID":     r.ConnectionID,
		"handle":           r.Handle,
		"name":             r.Name,
		"type":             r.Type,
		"location":         r.Location,
		"ownership":        r.Ownership,
		"sensitivityLevel": r.SensitivityLevel,
		"config":           r.Config,
		"capabilities":     r.Capabilities,
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
func (r DalConnectionUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetLocation() geolocation.Full {
	return r.Location
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetOwnership() string {
	return r.Ownership
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetSensitivityLevel() uint64 {
	return r.SensitivityLevel
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetConfig() types.ConnectionConfig {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r DalConnectionUpdate) GetCapabilities() types.ConnectionCapabilities {
	return r.Capabilities
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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["location[]"]; ok {
				r.Location, err = geolocation.Parse(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["location"]; ok {
				r.Location, err = geolocation.Parse(val)
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

			if val, ok := req.MultipartForm.Value["sensitivityLevel"]; ok && len(val) > 0 {
				r.SensitivityLevel, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
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

		if val, ok := req.Form["location[]"]; ok {
			r.Location, err = geolocation.Parse(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["location"]; ok {
			r.Location, err = geolocation.Parse(val)
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

		if val, ok := req.Form["sensitivityLevel"]; ok && len(val) > 0 {
			r.SensitivityLevel, err = payload.ParseUint64(val[0]), nil
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
