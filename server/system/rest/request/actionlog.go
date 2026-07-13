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
	ActionlogList struct {
		// From GET parameter
		//
		// From
		From *time.Time

		// To GET parameter
		//
		// To
		To *time.Time

		// BeforeActionID GET parameter
		//
		// Entries before specified action ID
		BeforeActionID uint64 `json:",string"`

		// Resource GET parameter
		//
		// Resource
		Resource string

		// ProjectID GET parameter
		//
		// Scope to a project
		ProjectID uint64 `json:",string"`

		// Action GET parameter
		//
		// Action
		Action string

		// ActorID GET parameter
		//
		// Filter by one or more actors
		ActorID []string

		// Origin GET parameter
		//
		// Origin
		Origin string

		// Limit GET parameter
		//
		// Limit
		Limit uint
	}

	ActionlogReport struct {
		// Dimensions GET parameter
		//
		// Group rows by one or more dimensions
		Dimensions []string

		// Metrics GET parameter
		//
		// Metrics to compute per group
		Metrics []string

		// From GET parameter
		//
		// From
		From *time.Time

		// To GET parameter
		//
		// To
		To *time.Time

		// Resource GET parameter
		//
		// Resource
		Resource string

		// Action GET parameter
		//
		// Action
		Action string

		// ActorID GET parameter
		//
		// Filter by one or more actors
		ActorID []string

		// Origin GET parameter
		//
		// Origin
		Origin string

		// Limit GET parameter
		//
		// Limit
		Limit uint
	}
)

// NewActionlogList request
func NewActionlogList() *ActionlogList {
	return &ActionlogList{}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"from":           r.From,
		"to":             r.To,
		"beforeActionID": r.BeforeActionID,
		"resource":       r.Resource,
		"projectID":      r.ProjectID,
		"action":         r.Action,
		"actorID":        r.ActorID,
		"origin":         r.Origin,
		"limit":          r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetTo() *time.Time {
	return r.To
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetBeforeActionID() uint64 {
	return r.BeforeActionID
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetAction() string {
	return r.Action
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetActorID() []string {
	return r.ActorID
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetOrigin() string {
	return r.Origin
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogList) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *ActionlogList) Fill(req *http.Request) (err error) {

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
		if val, ok := tmp["beforeActionID"]; ok && len(val) > 0 {
			r.BeforeActionID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["action"]; ok && len(val) > 0 {
			r.Action, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["actorID[]"]; ok {
			r.ActorID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["actorID"]; ok {
			r.ActorID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["origin"]; ok && len(val) > 0 {
			r.Origin, err = val[0], nil
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

// NewActionlogReport request
func NewActionlogReport() *ActionlogReport {
	return &ActionlogReport{}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"dimensions": r.Dimensions,
		"metrics":    r.Metrics,
		"from":       r.From,
		"to":         r.To,
		"resource":   r.Resource,
		"action":     r.Action,
		"actorID":    r.ActorID,
		"origin":     r.Origin,
		"limit":      r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetDimensions() []string {
	return r.Dimensions
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetMetrics() []string {
	return r.Metrics
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetFrom() *time.Time {
	return r.From
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetTo() *time.Time {
	return r.To
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetAction() string {
	return r.Action
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetActorID() []string {
	return r.ActorID
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetOrigin() string {
	return r.Origin
}

// Auditable returns all auditable/loggable parameters
func (r ActionlogReport) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *ActionlogReport) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["dimensions[]"]; ok {
			r.Dimensions, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["dimensions"]; ok {
			r.Dimensions, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["metrics[]"]; ok {
			r.Metrics, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["metrics"]; ok {
			r.Metrics, err = val, nil
			if err != nil {
				return err
			}
		}
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
		if val, ok := tmp["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["action"]; ok && len(val) > 0 {
			r.Action, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["actorID[]"]; ok {
			r.ActorID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["actorID"]; ok {
			r.ActorID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["origin"]; ok && len(val) > 0 {
			r.Origin, err = val[0], nil
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
