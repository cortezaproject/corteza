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
	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
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
	NgAutomationList struct {
		// AutomationID GET parameter
		//
		// Filter by automation ID
		AutomationID []string

		// ProjectID GET parameter
		//
		// Filter by project ID
		ProjectID uint64 `json:",string"`

		// Query GET parameter
		//
		// Filter automation
		Query string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted automation
		Deleted uint

		// Disabled GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) disabled automation
		Disabled uint

		// Labels GET parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total rows counter
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	NgAutomationCreate struct {
		// Handle POST parameter
		//
		// NgAutomation name
		Handle string

		// ProjectID POST parameter
		//
		// Project this automation belongs to
		ProjectID uint64 `json:",string"`

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

		// Meta POST parameter
		//
		// NgAutomation meta data
		Meta *types.NgAutomationMeta

		// Enabled POST parameter
		//
		// Is automation enabled
		Enabled bool

		// Scope POST parameter
		//
		// NgAutomation meta data
		Scope *expr.Vars

		// Triggers POST parameter
		//
		// NgAutomation steps definition
		Triggers types.NgAutomationTriggerSet

		// Steps POST parameter
		//
		// NgAutomation steps definition
		Steps types.NgAutomationStepSet

		// Paths POST parameter
		//
		// NgAutomation step paths definition
		Paths types.NgAutomationPathSet

		// RunAs POST parameter
		//
		// Is automation enabled
		RunAs uint64 `json:",string"`

		// OwnedBy POST parameter
		//
		// Owner of the automation
		OwnedBy uint64 `json:",string"`
	}

	NgAutomationUpdate struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`

		// Handle POST parameter
		//
		// NgAutomation name
		Handle string

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

		// Meta POST parameter
		//
		// NgAutomation meta data
		Meta *types.NgAutomationMeta

		// Enabled POST parameter
		//
		// Is automation enabled
		Enabled bool

		// Scope POST parameter
		//
		// NgAutomation meta data
		Scope *expr.Vars

		// Triggers POST parameter
		//
		// NgAutomation steps definition
		Triggers types.NgAutomationTriggerSet

		// Steps POST parameter
		//
		// NgAutomation steps definition
		Steps types.NgAutomationStepSet

		// Paths POST parameter
		//
		// NgAutomation step paths definition
		Paths types.NgAutomationPathSet

		// RunAs POST parameter
		//
		// Is automation enabled
		RunAs uint64 `json:",string"`

		// OwnedBy POST parameter
		//
		// Owner of the automation
		OwnedBy uint64 `json:",string"`

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	NgAutomationRead struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`
	}

	NgAutomationDelete struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`
	}

	NgAutomationUndelete struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`
	}

	NgAutomationTest struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`

		// Scope POST parameter
		//
		// NgAutomation meta data
		Scope *expr.Vars

		// RunAs POST parameter
		//
		// Is automation enabled
		RunAs bool
	}

	NgAutomationExec struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`

		// Input POST parameter
		//
		// Input
		Input *expr.Vars

		// Trace POST parameter
		//
		// Trace ngAutomation execution
		Trace bool

		// Wait POST parameter
		//
		// Wait for ngAutomation to complete
		Wait bool

		// Async POST parameter
		//
		// Execute step and return immediately
		Async bool
	}

	NgAutomationExecutions struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`
	}

	NgAutomationExecutionTrace struct {
		// AutomationID PATH parameter
		//
		// NgAutomation ID
		AutomationID uint64 `json:",string"`

		// ExecutionID PATH parameter
		//
		// Execution ID
		ExecutionID uint64 `json:",string"`
	}

	NgAutomationAllExecutions struct {
		// AutomationID GET parameter
		//
		// Filter by automation ID
		AutomationID []string

		// EventType GET parameter
		//
		// Filter by event type
		EventType string

		// ResourceType GET parameter
		//
		// Filter by resource type
		ResourceType string

		// Status GET parameter
		//
		// Filter by status
		Status []string
	}
)

// NewNgAutomationList request
func NewNgAutomationList() *NgAutomationList {
	return &NgAutomationList{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
		"projectID":    r.ProjectID,
		"query":        r.Query,
		"deleted":      r.Deleted,
		"disabled":     r.Disabled,
		"labels":       r.Labels,
		"limit":        r.Limit,
		"incTotal":     r.IncTotal,
		"pageCursor":   r.PageCursor,
		"sort":         r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetAutomationID() []string {
	return r.AutomationID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetDisabled() uint {
	return r.Disabled
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *NgAutomationList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["automationID[]"]; ok {
			r.AutomationID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["automationID"]; ok {
			r.AutomationID, err = val, nil
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
		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
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
		if val, ok := tmp["disabled"]; ok && len(val) > 0 {
			r.Disabled, err = payload.ParseUint(val[0]), nil
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
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
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

// NewNgAutomationCreate request
func NewNgAutomationCreate() *NgAutomationCreate {
	return &NgAutomationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":    r.Handle,
		"projectID": r.ProjectID,
		"labels":    r.Labels,
		"meta":      r.Meta,
		"enabled":   r.Enabled,
		"scope":     r.Scope,
		"triggers":  r.Triggers,
		"steps":     r.Steps,
		"paths":     r.Paths,
		"runAs":     r.RunAs,
		"ownedBy":   r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetMeta() *types.NgAutomationMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetScope() *expr.Vars {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetTriggers() types.NgAutomationTriggerSet {
	return r.Triggers
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetSteps() types.NgAutomationStepSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetPaths() types.NgAutomationPathSet {
	return r.Paths
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetRunAs() uint64 {
	return r.RunAs
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *NgAutomationCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["projectID"]; ok && len(val) > 0 {
				r.ProjectID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseNgAutomationMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseNgAutomationMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["enabled"]; ok && len(val) > 0 {
				r.Enabled, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scope[]"]; ok {
				r.Scope, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["scope"]; ok {
				r.Scope, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["triggers[]"]; ok {
				r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["triggers"]; ok {
				r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["steps[]"]; ok {
				r.Steps, err = types.ParseNgAutomationStepSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["steps"]; ok {
				r.Steps, err = types.ParseNgAutomationStepSet(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["paths[]"]; ok {
				r.Paths, err = types.ParseNgAutomationPathSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["paths"]; ok {
				r.Paths, err = types.ParseNgAutomationPathSet(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["runAs"]; ok && len(val) > 0 {
				r.RunAs, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownedBy"]; ok && len(val) > 0 {
				r.OwnedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}
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

		if val, ok := req.Form["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseNgAutomationMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseNgAutomationMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope[]"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scope"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["triggers[]"]; ok {
			r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["triggers"]; ok {
			r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["steps[]"]; ok {
			r.Steps, err = types.ParseNgAutomationStepSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["steps"]; ok {
			r.Steps, err = types.ParseNgAutomationStepSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["paths[]"]; ok {
			r.Paths, err = types.ParseNgAutomationPathSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["paths"]; ok {
			r.Paths, err = types.ParseNgAutomationPathSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["runAs"]; ok && len(val) > 0 {
			r.RunAs, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewNgAutomationUpdate request
func NewNgAutomationUpdate() *NgAutomationUpdate {
	return &NgAutomationUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
		"handle":       r.Handle,
		"labels":       r.Labels,
		"meta":         r.Meta,
		"enabled":      r.Enabled,
		"scope":        r.Scope,
		"triggers":     r.Triggers,
		"steps":        r.Steps,
		"paths":        r.Paths,
		"runAs":        r.RunAs,
		"ownedBy":      r.OwnedBy,
		"updatedAt":    r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetAutomationID() uint64 {
	return r.AutomationID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetMeta() *types.NgAutomationMeta {
	return r.Meta
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetScope() *expr.Vars {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetTriggers() types.NgAutomationTriggerSet {
	return r.Triggers
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetSteps() types.NgAutomationStepSet {
	return r.Steps
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetPaths() types.NgAutomationPathSet {
	return r.Paths
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetRunAs() uint64 {
	return r.RunAs
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *NgAutomationUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["handle"]; ok && len(val) > 0 {
				r.Handle, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["labels[]"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["labels"]; ok {
				r.Labels, err = label.ParseStrings(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseNgAutomationMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseNgAutomationMeta(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["enabled"]; ok && len(val) > 0 {
				r.Enabled, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scope[]"]; ok {
				r.Scope, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["scope"]; ok {
				r.Scope, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["triggers[]"]; ok {
				r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["triggers"]; ok {
				r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["steps[]"]; ok {
				r.Steps, err = types.ParseNgAutomationStepSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["steps"]; ok {
				r.Steps, err = types.ParseNgAutomationStepSet(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["paths[]"]; ok {
				r.Paths, err = types.ParseNgAutomationPathSet(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["paths"]; ok {
				r.Paths, err = types.ParseNgAutomationPathSet(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["runAs"]; ok && len(val) > 0 {
				r.RunAs, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["ownedBy"]; ok && len(val) > 0 {
				r.OwnedBy, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["updatedAt"]; ok && len(val) > 0 {
				r.UpdatedAt, err = payload.ParseISODatePtrWithErr(val[0])
				if err != nil {
					return err
				}
			}
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

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseNgAutomationMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseNgAutomationMeta(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scope[]"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scope"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["triggers[]"]; ok {
			r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["triggers"]; ok {
			r.Triggers, err = types.ParseNgAutomationTriggerSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["steps[]"]; ok {
			r.Steps, err = types.ParseNgAutomationStepSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["steps"]; ok {
			r.Steps, err = types.ParseNgAutomationStepSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["paths[]"]; ok {
			r.Paths, err = types.ParseNgAutomationPathSet(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["paths"]; ok {
			r.Paths, err = types.ParseNgAutomationPathSet(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["runAs"]; ok && len(val) > 0 {
			r.RunAs, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["ownedBy"]; ok && len(val) > 0 {
			r.OwnedBy, err = payload.ParseUint64(val[0]), nil
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

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationRead request
func NewNgAutomationRead() *NgAutomationRead {
	return &NgAutomationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationRead) GetAutomationID() uint64 {
	return r.AutomationID
}

// Fill processes request and fills internal variables
func (r *NgAutomationRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationDelete request
func NewNgAutomationDelete() *NgAutomationDelete {
	return &NgAutomationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationDelete) GetAutomationID() uint64 {
	return r.AutomationID
}

// Fill processes request and fills internal variables
func (r *NgAutomationDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationUndelete request
func NewNgAutomationUndelete() *NgAutomationUndelete {
	return &NgAutomationUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationUndelete) GetAutomationID() uint64 {
	return r.AutomationID
}

// Fill processes request and fills internal variables
func (r *NgAutomationUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationTest request
func NewNgAutomationTest() *NgAutomationTest {
	return &NgAutomationTest{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationTest) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
		"scope":        r.Scope,
		"runAs":        r.RunAs,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationTest) GetAutomationID() uint64 {
	return r.AutomationID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationTest) GetScope() *expr.Vars {
	return r.Scope
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationTest) GetRunAs() bool {
	return r.RunAs
}

// Fill processes request and fills internal variables
func (r *NgAutomationTest) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["scope[]"]; ok {
				r.Scope, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["scope"]; ok {
				r.Scope, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["runAs"]; ok && len(val) > 0 {
				r.RunAs, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["scope[]"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scope"]; ok {
			r.Scope, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["runAs"]; ok && len(val) > 0 {
			r.RunAs, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationExec request
func NewNgAutomationExec() *NgAutomationExec {
	return &NgAutomationExec{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExec) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
		"input":        r.Input,
		"trace":        r.Trace,
		"wait":         r.Wait,
		"async":        r.Async,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExec) GetAutomationID() uint64 {
	return r.AutomationID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExec) GetInput() *expr.Vars {
	return r.Input
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExec) GetTrace() bool {
	return r.Trace
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExec) GetWait() bool {
	return r.Wait
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExec) GetAsync() bool {
	return r.Async
}

// Fill processes request and fills internal variables
func (r *NgAutomationExec) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["input[]"]; ok {
				r.Input, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["input"]; ok {
				r.Input, err = types.ParseWorkflowVariables(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["trace"]; ok && len(val) > 0 {
				r.Trace, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["wait"]; ok && len(val) > 0 {
				r.Wait, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["async"]; ok && len(val) > 0 {
				r.Async, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["input[]"]; ok {
			r.Input, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["input"]; ok {
			r.Input, err = types.ParseWorkflowVariables(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["trace"]; ok && len(val) > 0 {
			r.Trace, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["wait"]; ok && len(val) > 0 {
			r.Wait, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["async"]; ok && len(val) > 0 {
			r.Async, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationExecutions request
func NewNgAutomationExecutions() *NgAutomationExecutions {
	return &NgAutomationExecutions{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExecutions) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExecutions) GetAutomationID() uint64 {
	return r.AutomationID
}

// Fill processes request and fills internal variables
func (r *NgAutomationExecutions) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationExecutionTrace request
func NewNgAutomationExecutionTrace() *NgAutomationExecutionTrace {
	return &NgAutomationExecutionTrace{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExecutionTrace) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
		"executionID":  r.ExecutionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExecutionTrace) GetAutomationID() uint64 {
	return r.AutomationID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationExecutionTrace) GetExecutionID() uint64 {
	return r.ExecutionID
}

// Fill processes request and fills internal variables
func (r *NgAutomationExecutionTrace) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "automationID")
		r.AutomationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "executionID")
		r.ExecutionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNgAutomationAllExecutions request
func NewNgAutomationAllExecutions() *NgAutomationAllExecutions {
	return &NgAutomationAllExecutions{}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationAllExecutions) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"automationID": r.AutomationID,
		"eventType":    r.EventType,
		"resourceType": r.ResourceType,
		"status":       r.Status,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationAllExecutions) GetAutomationID() []string {
	return r.AutomationID
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationAllExecutions) GetEventType() string {
	return r.EventType
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationAllExecutions) GetResourceType() string {
	return r.ResourceType
}

// Auditable returns all auditable/loggable parameters
func (r NgAutomationAllExecutions) GetStatus() []string {
	return r.Status
}

// Fill processes request and fills internal variables
func (r *NgAutomationAllExecutions) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["automationID[]"]; ok {
			r.AutomationID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["automationID"]; ok {
			r.AutomationID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["eventType"]; ok && len(val) > 0 {
			r.EventType, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["resourceType"]; ok && len(val) > 0 {
			r.ResourceType, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status[]"]; ok {
			r.Status, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["status"]; ok {
			r.Status, err = val, nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
