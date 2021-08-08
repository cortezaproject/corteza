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
	"github.com/cortezaproject/corteza-server/system/types"
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
	ApigwFilterList struct {
		// RouteID GET parameter
		//
		// Filter by route ID
		RouteID uint64 `json:",string"`

		// Query GET parameter
		//
		// Filter filters
		Query string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted filters
		Deleted uint64 `json:",string"`

		// Disabled GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) disabled filters
		Disabled uint64 `json:",string"`

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

	ApigwFilterCreate struct {
		// RouteID POST parameter
		//
		// Route
		RouteID uint64 `json:",string"`

		// Weight POST parameter
		//
		// Filter priority
		Weight uint64 `json:",string"`

		// Kind POST parameter
		//
		// Filter kind
		Kind string

		// Ref POST parameter
		//
		// Filter ref
		Ref string

		// Params POST parameter
		//
		// Filter parameters
		Params types.ApigwFilterParams
	}

	ApigwFilterUpdate struct {
		// FilterID PATH parameter
		//
		// Filter ID
		FilterID uint64 `json:",string"`

		// RouteID POST parameter
		//
		// Route
		RouteID uint64 `json:",string"`

		// Weight POST parameter
		//
		// Filter priority
		Weight uint64 `json:",string"`

		// Kind POST parameter
		//
		// Filter kind
		Kind string

		// Ref POST parameter
		//
		// Filter ref
		Ref string

		// Params POST parameter
		//
		// Filter parameters
		Params types.ApigwFilterParams
	}

	ApigwFilterRead struct {
		// FilterID PATH parameter
		//
		// Filter ID
		FilterID uint64 `json:",string"`
	}

	ApigwFilterDelete struct {
		// FilterID PATH parameter
		//
		// Filter ID
		FilterID uint64 `json:",string"`
	}

	ApigwFilterUndelete struct {
		// FilterID PATH parameter
		//
		// Filter ID
		FilterID uint64 `json:",string"`
	}

	ApigwFilterDefFilter struct {
		// Kind GET parameter
		//
		// Filter filters by kind
		Kind string
	}

	ApigwFilterDefProxyAuth struct {
	}
)

// NewApigwFilterList request
func NewApigwFilterList() *ApigwFilterList {
	return &ApigwFilterList{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"routeID":    r.RouteID,
		"query":      r.Query,
		"deleted":    r.Deleted,
		"disabled":   r.Disabled,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetRouteID() uint64 {
	return r.RouteID
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetDeleted() uint64 {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetDisabled() uint64 {
	return r.Disabled
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ApigwFilterList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["routeID"]; ok && len(val) > 0 {
			r.RouteID, err = payload.ParseUint64(val[0]), nil
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
			r.Deleted, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["disabled"]; ok && len(val) > 0 {
			r.Disabled, err = payload.ParseUint64(val[0]), nil
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

// NewApigwFilterCreate request
func NewApigwFilterCreate() *ApigwFilterCreate {
	return &ApigwFilterCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"routeID": r.RouteID,
		"weight":  r.Weight,
		"kind":    r.Kind,
		"ref":     r.Ref,
		"params":  r.Params,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterCreate) GetRouteID() uint64 {
	return r.RouteID
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterCreate) GetWeight() uint64 {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterCreate) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterCreate) GetRef() string {
	return r.Ref
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterCreate) GetParams() types.ApigwFilterParams {
	return r.Params
}

// Fill processes request and fills internal variables
func (r *ApigwFilterCreate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["routeID"]; ok && len(val) > 0 {
			r.RouteID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["weight"]; ok && len(val) > 0 {
			r.Weight, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ref"]; ok && len(val) > 0 {
			r.Ref, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["params[]"]; ok {
			r.Params, err = types.ParseApigwfFilterParams(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["params"]; ok {
			r.Params, err = types.ParseApigwfFilterParams(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewApigwFilterUpdate request
func NewApigwFilterUpdate() *ApigwFilterUpdate {
	return &ApigwFilterUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"filterID": r.FilterID,
		"routeID":  r.RouteID,
		"weight":   r.Weight,
		"kind":     r.Kind,
		"ref":      r.Ref,
		"params":   r.Params,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) GetFilterID() uint64 {
	return r.FilterID
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) GetRouteID() uint64 {
	return r.RouteID
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) GetWeight() uint64 {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) GetRef() string {
	return r.Ref
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUpdate) GetParams() types.ApigwFilterParams {
	return r.Params
}

// Fill processes request and fills internal variables
func (r *ApigwFilterUpdate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["routeID"]; ok && len(val) > 0 {
			r.RouteID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["weight"]; ok && len(val) > 0 {
			r.Weight, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ref"]; ok && len(val) > 0 {
			r.Ref, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["params[]"]; ok {
			r.Params, err = types.ParseApigwfFilterParams(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["params"]; ok {
			r.Params, err = types.ParseApigwfFilterParams(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "filterID")
		r.FilterID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApigwFilterRead request
func NewApigwFilterRead() *ApigwFilterRead {
	return &ApigwFilterRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"filterID": r.FilterID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterRead) GetFilterID() uint64 {
	return r.FilterID
}

// Fill processes request and fills internal variables
func (r *ApigwFilterRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "filterID")
		r.FilterID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApigwFilterDelete request
func NewApigwFilterDelete() *ApigwFilterDelete {
	return &ApigwFilterDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"filterID": r.FilterID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterDelete) GetFilterID() uint64 {
	return r.FilterID
}

// Fill processes request and fills internal variables
func (r *ApigwFilterDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "filterID")
		r.FilterID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApigwFilterUndelete request
func NewApigwFilterUndelete() *ApigwFilterUndelete {
	return &ApigwFilterUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"filterID": r.FilterID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterUndelete) GetFilterID() uint64 {
	return r.FilterID
}

// Fill processes request and fills internal variables
func (r *ApigwFilterUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "filterID")
		r.FilterID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApigwFilterDefFilter request
func NewApigwFilterDefFilter() *ApigwFilterDefFilter {
	return &ApigwFilterDefFilter{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterDefFilter) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind": r.Kind,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterDefFilter) GetKind() string {
	return r.Kind
}

// Fill processes request and fills internal variables
func (r *ApigwFilterDefFilter) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewApigwFilterDefProxyAuth request
func NewApigwFilterDefProxyAuth() *ApigwFilterDefProxyAuth {
	return &ApigwFilterDefProxyAuth{}
}

// Auditable returns all auditable/loggable parameters
func (r ApigwFilterDefProxyAuth) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *ApigwFilterDefProxyAuth) Fill(req *http.Request) (err error) {

	return err
}
