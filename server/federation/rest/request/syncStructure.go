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
	SyncStructureReadExposedInternal struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// LastSync GET parameter
		//
		// Last sync timestamp
		LastSync uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query
		Query string

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

	SyncStructureReadExposedSocial struct {
		// NodeID PATH parameter
		//
		// Node ID
		NodeID uint64 `json:",string"`

		// LastSync GET parameter
		//
		// Last sync timestamp
		LastSync uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query
		Query string

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
)

// NewSyncStructureReadExposedInternal request
func NewSyncStructureReadExposedInternal() *SyncStructureReadExposedInternal {
	return &SyncStructureReadExposedInternal{}
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":     r.NodeID,
		"lastSync":   r.LastSync,
		"query":      r.Query,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) GetLastSync() uint64 {
	return r.LastSync
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedInternal) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *SyncStructureReadExposedInternal) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["lastSync"]; ok && len(val) > 0 {
			r.LastSync, err = payload.ParseUint64(val[0]), nil
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

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewSyncStructureReadExposedSocial request
func NewSyncStructureReadExposedSocial() *SyncStructureReadExposedSocial {
	return &SyncStructureReadExposedSocial{}
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"nodeID":     r.NodeID,
		"lastSync":   r.LastSync,
		"query":      r.Query,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) GetNodeID() uint64 {
	return r.NodeID
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) GetLastSync() uint64 {
	return r.LastSync
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r SyncStructureReadExposedSocial) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *SyncStructureReadExposedSocial) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["lastSync"]; ok && len(val) > 0 {
			r.LastSync, err = payload.ParseUint64(val[0]), nil
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

		val = chi.URLParam(req, "nodeID")
		r.NodeID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
