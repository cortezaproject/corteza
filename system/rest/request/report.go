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
	"github.com/cortezaproject/corteza-server/pkg/report"
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
	ReportList struct {
		// Handle GET parameter
		//
		// Report handle
		Handle string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted reports
		Deleted uint

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

	ReportCreate struct {
		// Handle POST parameter
		//
		// Client handle
		Handle string

		// Meta POST parameter
		//
		// Additional info
		Meta *types.ReportMeta

		// Scenarios POST parameter
		//
		// Report scenarios
		Scenarios types.ReportScenarioSet

		// Sources POST parameter
		//
		// Report source definitions
		Sources types.ReportDataSourceSet

		// Blocks POST parameter
		//
		// Report blocks definition
		Blocks types.ReportBlockSet

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	ReportUpdate struct {
		// ReportID PATH parameter
		//
		// Report ID
		ReportID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Client handle
		Handle string

		// Meta POST parameter
		//
		// Additional info
		Meta *types.ReportMeta

		// Scenarios POST parameter
		//
		// Report scenarios
		Scenarios types.ReportScenarioSet

		// Sources POST parameter
		//
		// Report sources definition
		Sources types.ReportDataSourceSet

		// Blocks POST parameter
		//
		// Report blocks definition
		Blocks types.ReportBlockSet

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	ReportRead struct {
		// ReportID PATH parameter
		//
		// Report ID
		ReportID uint64 `json:",string"`
	}

	ReportDelete struct {
		// ReportID PATH parameter
		//
		// Report ID
		ReportID uint64 `json:",string"`
	}

	ReportUndelete struct {
		// ReportID PATH parameter
		//
		// Report ID
		ReportID uint64 `json:",string"`
	}

	ReportDescribe struct {
		// Sources POST parameter
		//
		// Report steps definition
		Sources types.ReportDataSourceSet

		// Steps POST parameter
		//
		// Report steps definition
		Steps report.StepDefinitionSet

		// Describe POST parameter
		//
		// The source descriptions to generate
		Describe []string
	}

	ReportRun struct {
		// ReportID PATH parameter
		//
		// Report ID
		ReportID uint64 `json:",string"`

		// Frames POST parameter
		//
		// Report data frame definitions
		Frames report.FrameDefinitionSet
	}
)

// NewReportList request
func NewReportList() *ReportList {
	return &ReportList{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":     r.Handle,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ReportList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ReportList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
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

	return err
}

// NewReportCreate request
func NewReportCreate() *ReportCreate {
	return &ReportCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":    r.Handle,
		"meta":      r.Meta,
		"scenarios": r.Scenarios,
		"sources":   r.Sources,
		"blocks":    r.Blocks,
		"labels":    r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) GetMeta() *types.ReportMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) GetScenarios() types.ReportScenarioSet {
	return r.Scenarios
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) GetSources() types.ReportDataSourceSet {
	return r.Sources
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) GetBlocks() types.ReportBlockSet {
	return r.Blocks
}

// Auditable returns all auditable/loggable parameters
func (r ReportCreate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ReportCreate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseReportMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseReportMeta(val)
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["scenarios[]"]; ok && len(val) > 0  {
		//    r.Scenarios, err = types.ReportScenarioSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["sources[]"]; ok && len(val) > 0  {
		//    r.Sources, err = types.ReportDataSourceSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["blocks[]"]; ok && len(val) > 0  {
		//    r.Blocks, err = types.ReportBlockSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

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

	return err
}

// NewReportUpdate request
func NewReportUpdate() *ReportUpdate {
	return &ReportUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reportID":  r.ReportID,
		"handle":    r.Handle,
		"meta":      r.Meta,
		"scenarios": r.Scenarios,
		"sources":   r.Sources,
		"blocks":    r.Blocks,
		"labels":    r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetReportID() uint64 {
	return r.ReportID
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetMeta() *types.ReportMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetScenarios() types.ReportScenarioSet {
	return r.Scenarios
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetSources() types.ReportDataSourceSet {
	return r.Sources
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetBlocks() types.ReportBlockSet {
	return r.Blocks
}

// Auditable returns all auditable/loggable parameters
func (r ReportUpdate) GetLabels() map[string]string {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ReportUpdate) Fill(req *http.Request) (err error) {

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

		if val, ok := req.Form["handle"]; ok && len(val) > 0 {
			r.Handle, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseReportMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseReportMeta(val)
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["scenarios[]"]; ok && len(val) > 0  {
		//    r.Scenarios, err = types.ReportScenarioSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["sources[]"]; ok && len(val) > 0  {
		//    r.Sources, err = types.ReportDataSourceSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["blocks[]"]; ok && len(val) > 0  {
		//    r.Blocks, err = types.ReportBlockSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

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

		val = chi.URLParam(req, "reportID")
		r.ReportID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReportRead request
func NewReportRead() *ReportRead {
	return &ReportRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reportID": r.ReportID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportRead) GetReportID() uint64 {
	return r.ReportID
}

// Fill processes request and fills internal variables
func (r *ReportRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "reportID")
		r.ReportID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReportDelete request
func NewReportDelete() *ReportDelete {
	return &ReportDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reportID": r.ReportID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportDelete) GetReportID() uint64 {
	return r.ReportID
}

// Fill processes request and fills internal variables
func (r *ReportDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "reportID")
		r.ReportID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReportUndelete request
func NewReportUndelete() *ReportUndelete {
	return &ReportUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reportID": r.ReportID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportUndelete) GetReportID() uint64 {
	return r.ReportID
}

// Fill processes request and fills internal variables
func (r *ReportUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "reportID")
		r.ReportID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReportDescribe request
func NewReportDescribe() *ReportDescribe {
	return &ReportDescribe{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportDescribe) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"sources":  r.Sources,
		"steps":    r.Steps,
		"describe": r.Describe,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportDescribe) GetSources() types.ReportDataSourceSet {
	return r.Sources
}

// Auditable returns all auditable/loggable parameters
func (r ReportDescribe) GetSteps() report.StepDefinitionSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r ReportDescribe) GetDescribe() []string {
	return r.Describe
}

// Fill processes request and fills internal variables
func (r *ReportDescribe) Fill(req *http.Request) (err error) {

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

		//if val, ok := req.Form["sources[]"]; ok && len(val) > 0  {
		//    r.Sources, err = types.ReportDataSourceSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["steps[]"]; ok && len(val) > 0  {
		//    r.Steps, err = report.StepDefinitionSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}

		//if val, ok := req.Form["describe[]"]; ok && len(val) > 0  {
		//    r.Describe, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	return err
}

// NewReportRun request
func NewReportRun() *ReportRun {
	return &ReportRun{}
}

// Auditable returns all auditable/loggable parameters
func (r ReportRun) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reportID": r.ReportID,
		"frames":   r.Frames,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReportRun) GetReportID() uint64 {
	return r.ReportID
}

// Auditable returns all auditable/loggable parameters
func (r ReportRun) GetFrames() report.FrameDefinitionSet {
	return r.Frames
}

// Fill processes request and fills internal variables
func (r *ReportRun) Fill(req *http.Request) (err error) {

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

		//if val, ok := req.Form["frames[]"]; ok && len(val) > 0  {
		//    r.Frames, err = report.FrameDefinitionSet(val), nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "reportID")
		r.ReportID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
