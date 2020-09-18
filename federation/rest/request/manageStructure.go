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
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
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

		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

		// ComposeModuleID POST parameter
		//
		// Compose module id
		ComposeModuleID uint64 `json:",string"`

		// Fields POST parameter
		//
		// Exposed module fields
		Fields types.ModuleFieldMappingList
	}

	ManageStructureRemove struct {
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
		"nodeID":          r.NodeID,
		"moduleID":        r.ModuleID,
		"composeModuleID": r.ComposeModuleID,
		"fields":          r.Fields,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetComposeModuleID() uint64 {
	return r.ComposeModuleID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureCreateExposed) GetFields() types.ModuleFieldMappingList {
	return r.Fields
}

// Fill processes request and fills internal variables
func (r *ManageStructureCreateExposed) Fill(req *http.Request) (err error) {
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

		// if val, ok := req.Form["fields"]; ok && len(val) > 0 {
		// 	r.Fields, err = types.ModuleFieldMappingList(val[0]), nil
		// 	if err != nil {
		// 		return err
		// 	}
		// }
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

// NewManageStructureRemove request
func NewManageStructureRemove() *ManageStructureRemove {
	return &ManageStructureRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureRemove) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":   r.NodeID,
		"moduleID": r.ModuleID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureRemove) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r ManageStructureRemove) GetModuleID() uint64 {
	return r.ModuleID
}

// Fill processes request and fills internal variables
func (r *ManageStructureRemove) Fill(req *http.Request) (err error) {
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

// Fill processes request and fills internal variables
func (r *ManageStructureListAll) Fill(req *http.Request) (err error) {
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
