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
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
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
	ManageStructureReadExposed struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`
	}

	ManageStructureCreateExposed struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ComposeModuleID POST parameter
		//
		// Compose module id
		ComposeModuleID uint64 `json:",string"`

		// ComposeNamespaceID POST parameter
		//
		// Compose namespace id
		ComposeNamespaceID uint64 `json:",string"`

		// Name POST parameter
		//
		// Module name
		Name string

		// Handle POST parameter
		//
		// Module handle
		Handle string

		// Fields POST parameter
		//
		// Exposed module fields
		Fields types.ModuleFieldSet
	}

	ManageStructureUpdateExposed struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// ComposeModuleID POST parameter
		//
		// Compose module id
		ComposeModuleID uint64 `json:",string"`

		// ComposeNamespaceID POST parameter
		//
		// Compose namespace id
		ComposeNamespaceID uint64 `json:",string"`

		// Name POST parameter
		//
		// Module name
		Name string

		// Handle POST parameter
		//
		// Module handle
		Handle string

		// Fields POST parameter
		//
		// Exposed module fields
		Fields types.ModuleFieldSet
	}

	ManageStructureRemoveExposed struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`
	}

	ManageStructureReadShared struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`
	}

	ManageStructureCreateMappings struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// ComposeModuleID POST parameter
		//
		// Compose module id
		ComposeModuleID uint64 `json:",string"`

		// ComposeNamespaceID POST parameter
		//
		// Compose namespace id
		ComposeNamespaceID uint64 `json:",string"`

		// Fields POST parameter
		//
		// Exposed module fields
		Fields types.ModuleFieldMappingSet
	}

	ManageStructureReadMappings struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// ComposeModuleID GET parameter
		//
		// Compose module id
		ComposeModuleID uint64 `json:",string"`
	}

	ManageStructureListAll struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// Shared GET parameter
		//
		// List shared modules
		Shared bool

		// Exposed GET parameter
		//
		// List exposed modules
		Exposed bool

		// Mapped GET parameter
		//
		// List mapped modules
		Mapped bool
	}
)

// NewManageStructureReadExposed request
func NewManageStructureReadExposed() *ManageStructureReadExposed {
	return &ManageStructureReadExposed{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadExposed) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":   r.NodeID,
		"moduleID": r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadExposed) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadExposed) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ManageStructureReadExposed) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureCreateExposed request
func NewManageStructureCreateExposed() *ManageStructureCreateExposed {
	return &ManageStructureCreateExposed{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":             r.NodeID,
		"composeModuleID":    r.ComposeModuleID,
		"composeNamespaceID": r.ComposeNamespaceID,
		"name":               r.Name,
		"handle":             r.Handle,
		"fields":             r.Fields,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetComposeModuleID() uint64 {
	return r.ComposeModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetComposeNamespaceID() uint64 {
	return r.ComposeNamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetFields() types.ModuleFieldSet {
	return r.Fields
}

// Fill processes request and fills internal variables
func (r *ManageStructureCreateExposed) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["composeModuleID"]; ok && len(val) > 0 {
				r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["composeNamespaceID"]; ok && len(val) > 0 {
				r.ComposeNamespaceID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
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

		if val, ok := req.Form["composeModuleID"]; ok && len(val) > 0 {
			r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["composeNamespaceID"]; ok && len(val) > 0 {
			r.ComposeNamespaceID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["fields[]"]; ok && len(val) > 0  {
		//    r.Fields, err = types.ModuleFieldSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureUpdateExposed request
func NewManageStructureUpdateExposed() *ManageStructureUpdateExposed {
	return &ManageStructureUpdateExposed{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":             r.NodeID,
		"moduleID":           r.ModuleID,
		"composeModuleID":    r.ComposeModuleID,
		"composeNamespaceID": r.ComposeNamespaceID,
		"name":               r.Name,
		"handle":             r.Handle,
		"fields":             r.Fields,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetComposeModuleID() uint64 {
	return r.ComposeModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetComposeNamespaceID() uint64 {
	return r.ComposeNamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureUpdateExposed) GetFields() types.ModuleFieldSet {
	return r.Fields
}

// Fill processes request and fills internal variables
func (r *ManageStructureUpdateExposed) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["composeModuleID"]; ok && len(val) > 0 {
				r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["composeNamespaceID"]; ok && len(val) > 0 {
				r.ComposeNamespaceID, err = payload.ParseUint64(val[0]), nil
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

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
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

		if val, ok := req.Form["composeModuleID"]; ok && len(val) > 0 {
			r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["composeNamespaceID"]; ok && len(val) > 0 {
			r.ComposeNamespaceID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["fields[]"]; ok && len(val) > 0  {
		//    r.Fields, err = types.ModuleFieldSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureRemoveExposed request
func NewManageStructureRemoveExposed() *ManageStructureRemoveExposed {
	return &ManageStructureRemoveExposed{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureRemoveExposed) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":   r.NodeID,
		"moduleID": r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureRemoveExposed) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureRemoveExposed) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ManageStructureRemoveExposed) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureReadShared request
func NewManageStructureReadShared() *ManageStructureReadShared {
	return &ManageStructureReadShared{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadShared) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":   r.NodeID,
		"moduleID": r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadShared) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadShared) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ManageStructureReadShared) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureCreateMappings request
func NewManageStructureCreateMappings() *ManageStructureCreateMappings {
	return &ManageStructureCreateMappings{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateMappings) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":             r.NodeID,
		"moduleID":           r.ModuleID,
		"composeModuleID":    r.ComposeModuleID,
		"composeNamespaceID": r.ComposeNamespaceID,
		"fields":             r.Fields,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateMappings) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateMappings) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateMappings) GetComposeModuleID() uint64 {
	return r.ComposeModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateMappings) GetComposeNamespaceID() uint64 {
	return r.ComposeNamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateMappings) GetFields() types.ModuleFieldMappingSet {
	return r.Fields
}

// Fill processes request and fills internal variables
func (r *ManageStructureCreateMappings) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["composeModuleID"]; ok && len(val) > 0 {
				r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["composeNamespaceID"]; ok && len(val) > 0 {
				r.ComposeNamespaceID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["composeModuleID"]; ok && len(val) > 0 {
			r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["composeNamespaceID"]; ok && len(val) > 0 {
			r.ComposeNamespaceID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["fields[]"]; ok && len(val) > 0  {
		//    r.Fields, err = types.ModuleFieldMappingSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureReadMappings request
func NewManageStructureReadMappings() *ManageStructureReadMappings {
	return &ManageStructureReadMappings{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadMappings) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":          r.NodeID,
		"moduleID":        r.ModuleID,
		"composeModuleID": r.ComposeModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadMappings) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadMappings) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureReadMappings) GetComposeModuleID() uint64 {
	return r.ComposeModuleID
}

// Fill processes request and fills internal variables
func (r *ManageStructureReadMappings) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["composeModuleID"]; ok && len(val) > 0 {
			r.ComposeModuleID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewManageStructureListAll request
func NewManageStructureListAll() *ManageStructureListAll {
	return &ManageStructureListAll{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureListAll) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":  r.NodeID,
		"shared":  r.Shared,
		"exposed": r.Exposed,
		"mapped":  r.Mapped,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureListAll) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureListAll) GetShared() bool {
	return r.Shared
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureListAll) GetExposed() bool {
	return r.Exposed
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureListAll) GetMapped() bool {
	return r.Mapped
}

// Fill processes request and fills internal variables
func (r *ManageStructureListAll) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["shared"]; ok && len(val) > 0 {
			r.Shared, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["exposed"]; ok && len(val) > 0 {
			r.Exposed, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["mapped"]; ok && len(val) > 0 {
			r.Mapped, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
