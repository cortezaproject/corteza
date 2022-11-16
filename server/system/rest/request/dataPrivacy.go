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
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/payload"
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
	DataPrivacyConnectionList struct {
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
		Deleted filter.State
	}
)

// NewDataPrivacyConnectionList request
func NewDataPrivacyConnectionList() *DataPrivacyConnectionList {
	return &DataPrivacyConnectionList{}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"connectionID": r.ConnectionID,
		"handle":       r.Handle,
		"type":         r.Type,
		"deleted":      r.Deleted,
	}
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetConnectionID() []string {
	return r.ConnectionID
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetType() string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r DataPrivacyConnectionList) GetDeleted() filter.State {
	return r.Deleted
}

// Fill processes request and fills internal variables
func (r *DataPrivacyConnectionList) Fill(req *http.Request) (err error) {

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
			r.Deleted, err = payload.ParseFilterState(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
