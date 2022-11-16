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
	DalSensitivityLevelList struct {
		// SensitivityLevelID GET parameter
		//
		// Filter by sensitivity level ID
		SensitivityLevelID []string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted sensitivity levels
		Deleted uint

		// IncTotal GET parameter
		//
		// Include total counter
		IncTotal bool
	}

	DalSensitivityLevelCreate struct {
		// Handle POST parameter
		//
		//
		Handle string

		// Level POST parameter
		//
		//
		Level int

		// Meta POST parameter
		//
		//
		Meta types.DalSensitivityLevelMeta
	}

	DalSensitivityLevelUpdate struct {
		// SensitivityLevelID PATH parameter
		//
		// Connection ID
		SensitivityLevelID uint64 `json:",string"`

		// Handle POST parameter
		//
		//
		Handle string

		// Level POST parameter
		//
		//
		Level int

		// Meta POST parameter
		//
		//
		Meta types.DalSensitivityLevelMeta
	}

	DalSensitivityLevelRead struct {
		// SensitivityLevelID PATH parameter
		//
		// Connection ID
		SensitivityLevelID uint64 `json:",string"`
	}

	DalSensitivityLevelDelete struct {
		// SensitivityLevelID PATH parameter
		//
		// Connection ID
		SensitivityLevelID uint64 `json:",string"`
	}

	DalSensitivityLevelUndelete struct {
		// SensitivityLevelID PATH parameter
		//
		// Connection ID
		SensitivityLevelID uint64 `json:",string"`
	}
)

// NewDalSensitivityLevelList request
func NewDalSensitivityLevelList() *DalSensitivityLevelList {
	return &DalSensitivityLevelList{}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
		"deleted":            r.Deleted,
		"incTotal":           r.IncTotal,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelList) GetSensitivityLevelID() []string {
	return r.SensitivityLevelID
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelList) GetIncTotal() bool {
	return r.IncTotal
}

// Fill processes request and fills internal variables
func (r *DalSensitivityLevelList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["sensitivityLevelID[]"]; ok {
			r.SensitivityLevelID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["sensitivityLevelID"]; ok {
			r.SensitivityLevelID, err = val, nil
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
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDalSensitivityLevelCreate request
func NewDalSensitivityLevelCreate() *DalSensitivityLevelCreate {
	return &DalSensitivityLevelCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle": r.Handle,
		"level":  r.Level,
		"meta":   r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelCreate) GetLevel() int {
	return r.Level
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelCreate) GetMeta() types.DalSensitivityLevelMeta {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *DalSensitivityLevelCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["level"]; ok && len(val) > 0 {
				r.Level, err = payload.ParseInt(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
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

		if val, ok := req.Form["level"]; ok && len(val) > 0 {
			r.Level, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewDalSensitivityLevelUpdate request
func NewDalSensitivityLevelUpdate() *DalSensitivityLevelUpdate {
	return &DalSensitivityLevelUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
		"handle":             r.Handle,
		"level":              r.Level,
		"meta":               r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUpdate) GetSensitivityLevelID() uint64 {
	return r.SensitivityLevelID
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUpdate) GetLevel() int {
	return r.Level
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUpdate) GetMeta() types.DalSensitivityLevelMeta {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *DalSensitivityLevelUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["level"]; ok && len(val) > 0 {
				r.Level, err = payload.ParseInt(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
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

		if val, ok := req.Form["level"]; ok && len(val) > 0 {
			r.Level, err = payload.ParseInt(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseDalSensitivityLevelMeta(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "sensitivityLevelID")
		r.SensitivityLevelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDalSensitivityLevelRead request
func NewDalSensitivityLevelRead() *DalSensitivityLevelRead {
	return &DalSensitivityLevelRead{}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelRead) GetSensitivityLevelID() uint64 {
	return r.SensitivityLevelID
}

// Fill processes request and fills internal variables
func (r *DalSensitivityLevelRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sensitivityLevelID")
		r.SensitivityLevelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDalSensitivityLevelDelete request
func NewDalSensitivityLevelDelete() *DalSensitivityLevelDelete {
	return &DalSensitivityLevelDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelDelete) GetSensitivityLevelID() uint64 {
	return r.SensitivityLevelID
}

// Fill processes request and fills internal variables
func (r *DalSensitivityLevelDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sensitivityLevelID")
		r.SensitivityLevelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewDalSensitivityLevelUndelete request
func NewDalSensitivityLevelUndelete() *DalSensitivityLevelUndelete {
	return &DalSensitivityLevelUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sensitivityLevelID": r.SensitivityLevelID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DalSensitivityLevelUndelete) GetSensitivityLevelID() uint64 {
	return r.SensitivityLevelID
}

// Fill processes request and fills internal variables
func (r *DalSensitivityLevelUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "sensitivityLevelID")
		r.SensitivityLevelID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
