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
	"github.com/cortezaproject/corteza/server/automation/types"
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
	TriggerDefinitionList struct {
		// TriggerDefinitionID GET parameter
		//
		// Filter by trigger definition ID
		TriggerDefinitionID []string

		// Query GET parameter
		//
		// Filter trigger definitions
		Query string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted
		Deleted uint

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

	TriggerDefinitionCreate struct {
		// Handle POST parameter
		//
		// Trigger definition handle
		Handle string

		// SkipEventBus POST parameter
		//
		// Skip event bus
		SkipEventBus bool

		// InputSchema POST parameter
		//
		// Input schema
		InputSchema types.TriggerDefinitionSchema

		// OutputSchema POST parameter
		//
		// Output schema
		OutputSchema types.TriggerDefinitionSchema

		// Meta POST parameter
		//
		// Meta data
		Meta *types.TriggerDefinitionMeta
	}

	TriggerDefinitionUpdate struct {
		// TriggerDefinitionID PATH parameter
		//
		// Trigger Definition ID
		TriggerDefinitionID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Trigger definition handle
		Handle string

		// SkipEventBus POST parameter
		//
		// Skip event bus
		SkipEventBus bool

		// InputSchema POST parameter
		//
		// Input schema
		InputSchema types.TriggerDefinitionSchema

		// OutputSchema POST parameter
		//
		// Output schema
		OutputSchema types.TriggerDefinitionSchema

		// Meta POST parameter
		//
		// Meta data
		Meta *types.TriggerDefinitionMeta
	}

	TriggerDefinitionRead struct {
		// TriggerDefinitionID PATH parameter
		//
		// Trigger Definition ID
		TriggerDefinitionID uint64 `json:",string"`
	}

	TriggerDefinitionDelete struct {
		// TriggerDefinitionID PATH parameter
		//
		// Trigger Definition ID
		TriggerDefinitionID uint64 `json:",string"`
	}

	TriggerDefinitionUndelete struct {
		// TriggerDefinitionID PATH parameter
		//
		// Trigger Definition ID
		TriggerDefinitionID uint64 `json:",string"`
	}
)

// NewTriggerDefinitionList request
func NewTriggerDefinitionList() *TriggerDefinitionList {
	return &TriggerDefinitionList{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerDefinitionID": r.TriggerDefinitionID,
		"query":               r.Query,
		"deleted":             r.Deleted,
		"limit":               r.Limit,
		"pageCursor":          r.PageCursor,
		"sort":                r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) GetTriggerDefinitionID() []string {
	return r.TriggerDefinitionID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *TriggerDefinitionList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["triggerDefinitionID[]"]; ok {
			r.TriggerDefinitionID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["triggerDefinitionID"]; ok {
			r.TriggerDefinitionID, err = val, nil
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

// NewTriggerDefinitionCreate request
func NewTriggerDefinitionCreate() *TriggerDefinitionCreate {
	return &TriggerDefinitionCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":       r.Handle,
		"skipEventBus": r.SkipEventBus,
		"inputSchema":  r.InputSchema,
		"outputSchema": r.OutputSchema,
		"meta":         r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionCreate) GetSkipEventBus() bool {
	return r.SkipEventBus
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionCreate) GetInputSchema() types.TriggerDefinitionSchema {
	return r.InputSchema
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionCreate) GetOutputSchema() types.TriggerDefinitionSchema {
	return r.OutputSchema
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionCreate) GetMeta() *types.TriggerDefinitionMeta {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *TriggerDefinitionCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["skipEventBus"]; ok && len(val) > 0 {
				r.SkipEventBus, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["inputSchema[]"]; ok {
				r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["inputSchema"]; ok {
				r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["outputSchema[]"]; ok {
				r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["outputSchema"]; ok {
				r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseTriggerDefinitionMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseTriggerDefinitionMeta(val)
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

		if val, ok := req.Form["skipEventBus"]; ok && len(val) > 0 {
			r.SkipEventBus, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["inputSchema[]"]; ok {
			r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["inputSchema"]; ok {
			r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["outputSchema[]"]; ok {
			r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["outputSchema"]; ok {
			r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTriggerDefinitionMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTriggerDefinitionMeta(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewTriggerDefinitionUpdate request
func NewTriggerDefinitionUpdate() *TriggerDefinitionUpdate {
	return &TriggerDefinitionUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerDefinitionID": r.TriggerDefinitionID,
		"handle":              r.Handle,
		"skipEventBus":        r.SkipEventBus,
		"inputSchema":         r.InputSchema,
		"outputSchema":        r.OutputSchema,
		"meta":                r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) GetTriggerDefinitionID() uint64 {
	return r.TriggerDefinitionID
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) GetSkipEventBus() bool {
	return r.SkipEventBus
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) GetInputSchema() types.TriggerDefinitionSchema {
	return r.InputSchema
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) GetOutputSchema() types.TriggerDefinitionSchema {
	return r.OutputSchema
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUpdate) GetMeta() *types.TriggerDefinitionMeta {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *TriggerDefinitionUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["skipEventBus"]; ok && len(val) > 0 {
				r.SkipEventBus, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["inputSchema[]"]; ok {
				r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["inputSchema"]; ok {
				r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["outputSchema[]"]; ok {
				r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["outputSchema"]; ok {
				r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["meta[]"]; ok {
				r.Meta, err = types.ParseTriggerDefinitionMeta(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["meta"]; ok {
				r.Meta, err = types.ParseTriggerDefinitionMeta(val)
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

		if val, ok := req.Form["skipEventBus"]; ok && len(val) > 0 {
			r.SkipEventBus, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["inputSchema[]"]; ok {
			r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["inputSchema"]; ok {
			r.InputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["outputSchema[]"]; ok {
			r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["outputSchema"]; ok {
			r.OutputSchema, err = types.ParseTriggerDefinitionSchema(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = types.ParseTriggerDefinitionMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = types.ParseTriggerDefinitionMeta(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "triggerDefinitionID")
		r.TriggerDefinitionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTriggerDefinitionRead request
func NewTriggerDefinitionRead() *TriggerDefinitionRead {
	return &TriggerDefinitionRead{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerDefinitionID": r.TriggerDefinitionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionRead) GetTriggerDefinitionID() uint64 {
	return r.TriggerDefinitionID
}

// Fill processes request and fills internal variables
func (r *TriggerDefinitionRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "triggerDefinitionID")
		r.TriggerDefinitionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTriggerDefinitionDelete request
func NewTriggerDefinitionDelete() *TriggerDefinitionDelete {
	return &TriggerDefinitionDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerDefinitionID": r.TriggerDefinitionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionDelete) GetTriggerDefinitionID() uint64 {
	return r.TriggerDefinitionID
}

// Fill processes request and fills internal variables
func (r *TriggerDefinitionDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "triggerDefinitionID")
		r.TriggerDefinitionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewTriggerDefinitionUndelete request
func NewTriggerDefinitionUndelete() *TriggerDefinitionUndelete {
	return &TriggerDefinitionUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"triggerDefinitionID": r.TriggerDefinitionID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r TriggerDefinitionUndelete) GetTriggerDefinitionID() uint64 {
	return r.TriggerDefinitionID
}

// Fill processes request and fills internal variables
func (r *TriggerDefinitionUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "triggerDefinitionID")
		r.TriggerDefinitionID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
