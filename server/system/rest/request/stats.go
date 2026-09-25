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
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/go-chi/chi/v5"
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
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	StatsList struct {
		// From GET parameter
		//
		// Start of the reporting range (inclusive); defaults to 30 days before `to`
		From *time.Time

		// To GET parameter
		//
		// End of the reporting range (exclusive); defaults to now
		To *time.Time

		// Bucket GET parameter
		//
		// Series granularity (day, week, month); defaults by range length
		Bucket string
	}

	StatsDetail struct {
		// Resource PATH parameter
		//
		// Inventory resource name (users, roles, workflows, ...)
		Resource string

		// From GET parameter
		//
		// Start of the reporting range (inclusive); defaults to 30 days before `to`
		From *time.Time

		// To GET parameter
		//
		// End of the reporting range (exclusive); defaults to now
		To *time.Time

		// Bucket GET parameter
		//
		// Series granularity (day, week, month); defaults by range length
		Bucket string
	}
)

// NewStatsList request
func NewStatsList() *StatsList {
	return &StatsList{}
}

// Auditable returns all auditable/loggable parameters
func (r StatsList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"from":   r.From,
		"to":     r.To,
		"bucket": r.Bucket,
	}
}

// Auditable returns all auditable/loggable parameters
func (r StatsList) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r StatsList) GetTo() *time.Time {
	return r.To
}

// Auditable returns all auditable/loggable parameters
func (r StatsList) GetBucket() string {
	return r.Bucket
}

// Fill processes request and fills internal variables
func (r *StatsList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["from"]; ok && len(val) > 0 {
			r.From, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["to"]; ok && len(val) > 0 {
			r.To, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["bucket"]; ok && len(val) > 0 {
			r.Bucket, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewStatsDetail request
func NewStatsDetail() *StatsDetail {
	return &StatsDetail{}
}

// Auditable returns all auditable/loggable parameters
func (r StatsDetail) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource": r.Resource,
		"from":     r.From,
		"to":       r.To,
		"bucket":   r.Bucket,
	}
}

// Auditable returns all auditable/loggable parameters
func (r StatsDetail) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r StatsDetail) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r StatsDetail) GetTo() *time.Time {
	return r.To
}

// Auditable returns all auditable/loggable parameters
func (r StatsDetail) GetBucket() string {
	return r.Bucket
}

// Fill processes request and fills internal variables
func (r *StatsDetail) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["from"]; ok && len(val) > 0 {
			r.From, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["to"]; ok && len(val) > 0 {
			r.To, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["bucket"]; ok && len(val) > 0 {
			r.Bucket, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "resource")
		r.Resource, err = val, nil
		if err != nil {
			return err
		}

	}

	return err
}
