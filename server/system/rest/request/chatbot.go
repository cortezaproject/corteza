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
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/crusttech/human/server/system/types"
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
	ChatbotList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// Handle GET parameter
		//
		// Search by handle
		Handle string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted chatbots
		Deleted uint

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
		// Include total counter
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

	ChatbotCreate struct {
		// Handle POST parameter
		//
		// Chatbot handle
		Handle string

		// Name POST parameter
		//
		// Chatbot display name
		Name string

		// Enabled POST parameter
		//
		// Whether the chatbot is enabled
		Enabled bool

		// SessionTTL POST parameter
		//
		// Session TTL (duration string)
		SessionTTL string

		// AllowedOrigins POST parameter
		//
		// Allowed origins
		AllowedOrigins types.ChatbotAllowedOrigins

		// Handoff POST parameter
		//
		// Chatbot handoff settings
		Handoff types.ChatbotHandoff

		// Styling POST parameter
		//
		// Chatbot styling
		Styling types.ChatbotStyling

		// Scenarios POST parameter
		//
		// Chatbot scenarios
		Scenarios types.ChatbotScenarios

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue
	}

	ChatbotRead struct {
		// ChatbotID PATH parameter
		//
		// Chatbot ID
		ChatbotID uint64 `json:",string"`
	}

	ChatbotUpdate struct {
		// ChatbotID PATH parameter
		//
		// Chatbot ID
		ChatbotID uint64 `json:",string"`

		// Handle POST parameter
		//
		// Chatbot handle
		Handle string

		// Name POST parameter
		//
		// Chatbot display name
		Name string

		// Enabled POST parameter
		//
		// Whether the chatbot is enabled
		Enabled bool

		// SessionTTL POST parameter
		//
		// Session TTL (duration string)
		SessionTTL string

		// AllowedOrigins POST parameter
		//
		// Allowed origins
		AllowedOrigins types.ChatbotAllowedOrigins

		// Handoff POST parameter
		//
		// Chatbot handoff settings
		Handoff types.ChatbotHandoff

		// Styling POST parameter
		//
		// Chatbot styling
		Styling types.ChatbotStyling

		// Scenarios POST parameter
		//
		// Chatbot scenarios
		Scenarios types.ChatbotScenarios

		// Labels POST parameter
		//
		// Labels
		Labels map[string]labelTypes.LabelValue

		// UpdatedAt POST parameter
		//
		// Last update (or creation) date
		UpdatedAt *time.Time
	}

	ChatbotDelete struct {
		// ChatbotID PATH parameter
		//
		// Chatbot ID
		ChatbotID uint64 `json:",string"`
	}

	ChatbotUndelete struct {
		// ChatbotID PATH parameter
		//
		// Chatbot ID
		ChatbotID uint64 `json:",string"`
	}

	ChatbotRegenerateWidgetKey struct {
		// ChatbotID PATH parameter
		//
		// Chatbot ID
		ChatbotID uint64 `json:",string"`
	}

	ChatbotUploadAsset struct {
		// ChatbotID PATH parameter
		//
		// Chatbot ID
		ChatbotID uint64 `json:",string"`

		// Upload POST parameter
		//
		// Asset file to upload
		Upload *multipart.FileHeader
	}
)

// NewChatbotList request
func NewChatbotList() *ChatbotList {
	return &ChatbotList{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"handle":     r.Handle,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ChatbotList) Fill(req *http.Request) (err error) {

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

// NewChatbotCreate request
func NewChatbotCreate() *ChatbotCreate {
	return &ChatbotCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"handle":         r.Handle,
		"name":           r.Name,
		"enabled":        r.Enabled,
		"sessionTTL":     r.SessionTTL,
		"allowedOrigins": r.AllowedOrigins,
		"handoff":        r.Handoff,
		"styling":        r.Styling,
		"scenarios":      r.Scenarios,
		"labels":         r.Labels,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetSessionTTL() string {
	return r.SessionTTL
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetAllowedOrigins() types.ChatbotAllowedOrigins {
	return r.AllowedOrigins
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetHandoff() types.ChatbotHandoff {
	return r.Handoff
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetStyling() types.ChatbotStyling {
	return r.Styling
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetScenarios() types.ChatbotScenarios {
	return r.Scenarios
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotCreate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Fill processes request and fills internal variables
func (r *ChatbotCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["sessionTTL"]; ok && len(val) > 0 {
				r.SessionTTL, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["allowedOrigins[]"]; ok {
				r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["allowedOrigins"]; ok {
				r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["handoff[]"]; ok {
				r.Handoff, err = types.ParseChatbotHandoff(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["handoff"]; ok {
				r.Handoff, err = types.ParseChatbotHandoff(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["styling[]"]; ok {
				r.Styling, err = types.ParseChatbotStyling(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["styling"]; ok {
				r.Styling, err = types.ParseChatbotStyling(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scenarios[]"]; ok {
				r.Scenarios, err = types.ParseChatbotScenarios(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["scenarios"]; ok {
				r.Scenarios, err = types.ParseChatbotScenarios(val)
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
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

		if val, ok := req.Form["sessionTTL"]; ok && len(val) > 0 {
			r.SessionTTL, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["allowedOrigins[]"]; ok {
			r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["allowedOrigins"]; ok {
			r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handoff[]"]; ok {
			r.Handoff, err = types.ParseChatbotHandoff(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["handoff"]; ok {
			r.Handoff, err = types.ParseChatbotHandoff(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["styling[]"]; ok {
			r.Styling, err = types.ParseChatbotStyling(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["styling"]; ok {
			r.Styling, err = types.ParseChatbotStyling(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scenarios[]"]; ok {
			r.Scenarios, err = types.ParseChatbotScenarios(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scenarios"]; ok {
			r.Scenarios, err = types.ParseChatbotScenarios(val)
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

	return err
}

// NewChatbotRead request
func NewChatbotRead() *ChatbotRead {
	return &ChatbotRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID": r.ChatbotID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotRead) GetChatbotID() uint64 {
	return r.ChatbotID
}

// Fill processes request and fills internal variables
func (r *ChatbotRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "chatbotID")
		r.ChatbotID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotUpdate request
func NewChatbotUpdate() *ChatbotUpdate {
	return &ChatbotUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID":      r.ChatbotID,
		"handle":         r.Handle,
		"name":           r.Name,
		"enabled":        r.Enabled,
		"sessionTTL":     r.SessionTTL,
		"allowedOrigins": r.AllowedOrigins,
		"handoff":        r.Handoff,
		"styling":        r.Styling,
		"scenarios":      r.Scenarios,
		"labels":         r.Labels,
		"updatedAt":      r.UpdatedAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetChatbotID() uint64 {
	return r.ChatbotID
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetHandle() string {
	return r.Handle
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetSessionTTL() string {
	return r.SessionTTL
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetAllowedOrigins() types.ChatbotAllowedOrigins {
	return r.AllowedOrigins
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetHandoff() types.ChatbotHandoff {
	return r.Handoff
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetStyling() types.ChatbotStyling {
	return r.Styling
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetScenarios() types.ChatbotScenarios {
	return r.Scenarios
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetLabels() map[string]labelTypes.LabelValue {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// Fill processes request and fills internal variables
func (r *ChatbotUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["name"]; ok && len(val) > 0 {
				r.Name, err = val[0], nil
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

			if val, ok := req.MultipartForm.Value["sessionTTL"]; ok && len(val) > 0 {
				r.SessionTTL, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["allowedOrigins[]"]; ok {
				r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["allowedOrigins"]; ok {
				r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["handoff[]"]; ok {
				r.Handoff, err = types.ParseChatbotHandoff(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["handoff"]; ok {
				r.Handoff, err = types.ParseChatbotHandoff(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["styling[]"]; ok {
				r.Styling, err = types.ParseChatbotStyling(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["styling"]; ok {
				r.Styling, err = types.ParseChatbotStyling(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["scenarios[]"]; ok {
				r.Scenarios, err = types.ParseChatbotScenarios(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["scenarios"]; ok {
				r.Scenarios, err = types.ParseChatbotScenarios(val)
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

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
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

		if val, ok := req.Form["sessionTTL"]; ok && len(val) > 0 {
			r.SessionTTL, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["allowedOrigins[]"]; ok {
			r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["allowedOrigins"]; ok {
			r.AllowedOrigins, err = types.ParseChatbotAllowedOrigins(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handoff[]"]; ok {
			r.Handoff, err = types.ParseChatbotHandoff(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["handoff"]; ok {
			r.Handoff, err = types.ParseChatbotHandoff(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["styling[]"]; ok {
			r.Styling, err = types.ParseChatbotStyling(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["styling"]; ok {
			r.Styling, err = types.ParseChatbotStyling(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["scenarios[]"]; ok {
			r.Scenarios, err = types.ParseChatbotScenarios(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["scenarios"]; ok {
			r.Scenarios, err = types.ParseChatbotScenarios(val)
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

		val = chi.URLParam(req, "chatbotID")
		r.ChatbotID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotDelete request
func NewChatbotDelete() *ChatbotDelete {
	return &ChatbotDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID": r.ChatbotID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotDelete) GetChatbotID() uint64 {
	return r.ChatbotID
}

// Fill processes request and fills internal variables
func (r *ChatbotDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "chatbotID")
		r.ChatbotID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotUndelete request
func NewChatbotUndelete() *ChatbotUndelete {
	return &ChatbotUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID": r.ChatbotID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUndelete) GetChatbotID() uint64 {
	return r.ChatbotID
}

// Fill processes request and fills internal variables
func (r *ChatbotUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "chatbotID")
		r.ChatbotID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotRegenerateWidgetKey request
func NewChatbotRegenerateWidgetKey() *ChatbotRegenerateWidgetKey {
	return &ChatbotRegenerateWidgetKey{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotRegenerateWidgetKey) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID": r.ChatbotID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotRegenerateWidgetKey) GetChatbotID() uint64 {
	return r.ChatbotID
}

// Fill processes request and fills internal variables
func (r *ChatbotRegenerateWidgetKey) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "chatbotID")
		r.ChatbotID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewChatbotUploadAsset request
func NewChatbotUploadAsset() *ChatbotUploadAsset {
	return &ChatbotUploadAsset{}
}

// Auditable returns all auditable/loggable parameters
func (r ChatbotUploadAsset) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"chatbotID": r.ChatbotID,
		"upload":    r.Upload,
	}
}

func (r ChatbotUploadAsset) GetChatbotID() uint64 {
	return r.ChatbotID
}

func (r ChatbotUploadAsset) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Fill processes request and fills internal variables
func (r *ChatbotUploadAsset) Fill(req *http.Request) (err error) {
	{
		var val string
		// path params
		val = chi.URLParam(req, "chatbotID")
		r.ChatbotID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}
		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}
	}

	return err
}
