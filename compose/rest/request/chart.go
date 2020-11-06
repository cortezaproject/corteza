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
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
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
	ChartList struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// Query GET parameter
		//
		// Search query to match against charts
		Query string

		// Handle GET parameter
		//
		// Search charts by handle
		Handle string

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

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

	ChartCreate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// Config POST parameter
		//
		// Chart JSON
		Config sqlxTypes.JSONText

		// Name POST parameter
		//
		// Chart name
		Name string

		// Handle POST parameter
		//
		// Chart handle
		Handle string

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	ChartRead struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ChartID PATH parameter
		//
		// Chart ID
		ChartID uint64 `json:",string"`
	}

	ChartUpdate struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ChartID PATH parameter
		//
		// Chart ID
		ChartID uint64 `json:",string"`

		// Config POST parameter
		//
		// Chart JSON
		Config sqlxTypes.JSONText

		// Name POST parameter
		//
		// Chart name
		Name string

		// Handle POST parameter
		//
		// Chart handle
		Handle string

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	ChartDelete struct {
		// NamespaceID PATH parameter
		//
		// Namespace ID
		NamespaceID uint64 `json:",string"`

		// ChartID PATH parameter
		//
		// Chart ID
		ChartID uint64 `json:",string"`
	}
)

// NewChartList request
func NewChartList() *ChartList {
	return &ChartList{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"query":       r.Query,
		"handle":      r.Handle,
		"labels":      r.Labels,
		"limit":       r.Limit,
		"pageCursor":  r.PageCursor,
		"sort":        r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ChartList) Fill(req *http.Request) (err error) {
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
		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
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

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChartCreate request
func NewChartCreate() *ChartCreate {
	return &ChartCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"config":      r.Config,
		"name":        r.Name,
		"handle":      r.Handle,
		"labels":      r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ChartCreate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
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

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChartRead request
func NewChartRead() *ChartRead {
	return &ChartRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"chartID":     r.ChartID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChartRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ChartRead) GetChartID() uint64 {
	return r.ChartID
}

// Fill processes request and fills internal variables
func (r *ChartRead) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "chartID")
		r.ChartID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChartUpdate request
func NewChartUpdate() *ChartUpdate {
	return &ChartUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"chartID":     r.ChartID,
		"config":      r.Config,
		"name":        r.Name,
		"handle":      r.Handle,
		"labels":      r.Labels,
		"updatedAt":   r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetChartID() uint64 {
	return r.ChartID
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *ChartUpdate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
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

		if val, ok := req.Form["updatedAt"]; ok && len(val) > 0 {
			r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "chartID")
		r.ChartID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChartDelete request
func NewChartDelete() *ChartDelete {
	return &ChartDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"namespaceID": r.NamespaceID,
		"chartID":     r.ChartID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChartDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// Auditable returns all auditable/loggable parameters
func (r ChartDelete) GetChartID() uint64 {
	return r.ChartID
}

// Fill processes request and fills internal variables
func (r *ChartDelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "namespaceID")
		r.NamespaceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "chartID")
		r.ChartID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
