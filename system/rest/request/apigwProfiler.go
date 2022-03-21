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
	ApigwProfilerAggregation struct {
		// Path GET parameter
		//
		// Filter by request path
		Path string

		// Before GET parameter
		//
		// Entries before specified route
		Before string

		// Sort GET parameter
		//
		// Sort items
		Sort string

		// Limit GET parameter
		//
		// Limit
		Limit uint
	}

	ApigwProfilerRoute struct {
		// RouteID PATH parameter
		//
		// Route ID
		RouteID string

		// Path GET parameter
		//
		// Filter by request path
		Path string

		// Before GET parameter
		//
		// Entries before specified hit ID
		Before string

		// Sort GET parameter
		//
		// Sort items
		Sort string

		// Limit GET parameter
		//
		// Limit
		Limit uint
	}

	ApigwProfilerHit struct {
		// HitID PATH parameter
		//
		// Hit ID
		HitID string
	}
)

// NewApigwProfilerAggregation request
func NewApigwProfilerAggregation() *ApigwProfilerAggregation {
	return &ApigwProfilerAggregation{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerAggregation) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"path":   r.Path,
		"before": r.Before,
		"sort":   r.Sort,
		"limit":  r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerAggregation) GetPath() string {
	return r.Path
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerAggregation) GetBefore() string {
	return r.Before
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerAggregation) GetSort() string {
	return r.Sort
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerAggregation) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *ApigwProfilerAggregation) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["path"]; ok && len(val) > 0 {
			r.Path, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["before"]; ok && len(val) > 0 {
			r.Before, err = val[0], nil
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
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewApigwProfilerRoute request
func NewApigwProfilerRoute() *ApigwProfilerRoute {
	return &ApigwProfilerRoute{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerRoute) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"routeID": r.RouteID,
		"path":    r.Path,
		"before":  r.Before,
		"sort":    r.Sort,
		"limit":   r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerRoute) GetRouteID() string {
	return r.RouteID
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerRoute) GetPath() string {
	return r.Path
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerRoute) GetBefore() string {
	return r.Before
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerRoute) GetSort() string {
	return r.Sort
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerRoute) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *ApigwProfilerRoute) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["path"]; ok && len(val) > 0 {
			r.Path, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["before"]; ok && len(val) > 0 {
			r.Before, err = val[0], nil
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
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "routeID")
		r.RouteID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApigwProfilerHit request
func NewApigwProfilerHit() *ApigwProfilerHit {
	return &ApigwProfilerHit{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerHit) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"hitID": r.HitID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwProfilerHit) GetHitID() string {
	return r.HitID
}

// Fill processes request and fills internal variables
func (r *ApigwProfilerHit) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "hitID")
		r.HitID, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
