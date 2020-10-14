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
)

type (
	// Internal API interface
	SyncDataReadExposedAll struct {
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

	SyncDataReadExposed struct {
		// ModuleID PATH parameter
		//
		// Module ID
		ModuleID uint64 `json:",string"`

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

// NewSyncDataReadExposedAll request
func NewSyncDataReadExposedAll() *SyncDataReadExposedAll {
	return &SyncDataReadExposedAll{}
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposedAll) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposedAll) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposedAll) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposedAll) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposedAll) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *SyncDataReadExposedAll) Fill(req *http.Request) (err error) {
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

	return err
}

// NewSyncDataReadExposed request
func NewSyncDataReadExposed() *SyncDataReadExposed {
	return &SyncDataReadExposed{}
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposed) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"moduleID":   r.ModuleID,
		"query":      r.Query,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposed) GetModuleID() uint64 {
	return r.ModuleID
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposed) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposed) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposed) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r SyncDataReadExposed) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *SyncDataReadExposed) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "moduleID")
		r.ModuleID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
