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
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-chi/chi/v5"
	sqlxTypes "github.com/jmoiron/sqlx/types"
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
	NotificationList struct {
		// NotificationID GET parameter
		//
		// Filter by notification ID
		NotificationID []string

		// Kind GET parameter
		//
		// Only notifications of a specific kind
		Kind string

		// Read GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) read notifications
		Read uint

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted notifications
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

	NotificationCreate struct {
		// Kind POST parameter
		//
		// Kind
		Kind string

		// Config POST parameter
		//
		// Configuration data for the notification
		Config sqlxTypes.JSONText

		// Recipient POST parameter
		//
		// Recipient
		Recipient uint64 `json:",string"`
	}

	NotificationUpdate struct {
		// NotificationID PATH parameter
		//
		// Notification ID
		NotificationID uint64 `json:",string"`

		// Kind POST parameter
		//
		// Kind
		Kind string

		// Config POST parameter
		//
		// Configuration data for the notification
		Config sqlxTypes.JSONText

		// Recipient POST parameter
		//
		// Recipient
		Recipient uint64 `json:",string"`
	}

	NotificationRead struct {
		// NotificationID PATH parameter
		//
		// Notification ID
		NotificationID uint64 `json:",string"`
	}

	NotificationDelete struct {
		// NotificationID PATH parameter
		//
		// Notification ID
		NotificationID uint64 `json:",string"`
	}

	NotificationMarkAsRead struct {
		// NotificationID PATH parameter
		//
		// Notification ID
		NotificationID uint64 `json:",string"`
	}

	NotificationMarkAllAsRead struct {
	}
)

// NewNotificationList request
func NewNotificationList() *NotificationList {
	return &NotificationList{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"notificationID": r.NotificationID,
		"kind":           r.Kind,
		"read":           r.Read,
		"deleted":        r.Deleted,
		"limit":          r.Limit,
		"pageCursor":     r.PageCursor,
		"sort":           r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetNotificationID() []string {
	return r.NotificationID
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetRead() uint {
	return r.Read
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r NotificationList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *NotificationList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["notificationID[]"]; ok {
			r.NotificationID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["notificationID"]; ok {
			r.NotificationID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["read"]; ok && len(val) > 0 {
			r.Read, err = payload.ParseUint(val[0]), nil
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

// NewNotificationCreate request
func NewNotificationCreate() *NotificationCreate {
	return &NotificationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"kind":      r.Kind,
		"config":    r.Config,
		"recipient": r.Recipient,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationCreate) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r NotificationCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r NotificationCreate) GetRecipient() uint64 {
	return r.Recipient
}

// Fill processes request and fills internal variables
func (r *NotificationCreate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["kind"]; ok && len(val) > 0 {
				r.Kind, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config"]; ok && len(val) > 0 {
				r.Config, err = payload.ParseJSONTextWithErr(val[0])
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["recipient"]; ok && len(val) > 0 {
				r.Recipient, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["recipient"]; ok && len(val) > 0 {
			r.Recipient, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewNotificationUpdate request
func NewNotificationUpdate() *NotificationUpdate {
	return &NotificationUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"notificationID": r.NotificationID,
		"kind":           r.Kind,
		"config":         r.Config,
		"recipient":      r.Recipient,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationUpdate) GetNotificationID() uint64 {
	return r.NotificationID
}

// Auditable returns all auditable/loggable parameters
func (r NotificationUpdate) GetKind() string {
	return r.Kind
}

// Auditable returns all auditable/loggable parameters
func (r NotificationUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r NotificationUpdate) GetRecipient() uint64 {
	return r.Recipient
}

// Fill processes request and fills internal variables
func (r *NotificationUpdate) Fill(req *http.Request) (err error) {

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

			if val, ok := req.MultipartForm.Value["kind"]; ok && len(val) > 0 {
				r.Kind, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["config"]; ok && len(val) > 0 {
				r.Config, err = payload.ParseJSONTextWithErr(val[0])
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["recipient"]; ok && len(val) > 0 {
				r.Recipient, err = payload.ParseUint64(val[0]), nil
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

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["recipient"]; ok && len(val) > 0 {
			r.Recipient, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "notificationID")
		r.NotificationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNotificationRead request
func NewNotificationRead() *NotificationRead {
	return &NotificationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"notificationID": r.NotificationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationRead) GetNotificationID() uint64 {
	return r.NotificationID
}

// Fill processes request and fills internal variables
func (r *NotificationRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "notificationID")
		r.NotificationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNotificationDelete request
func NewNotificationDelete() *NotificationDelete {
	return &NotificationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"notificationID": r.NotificationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationDelete) GetNotificationID() uint64 {
	return r.NotificationID
}

// Fill processes request and fills internal variables
func (r *NotificationDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "notificationID")
		r.NotificationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNotificationMarkAsRead request
func NewNotificationMarkAsRead() *NotificationMarkAsRead {
	return &NotificationMarkAsRead{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationMarkAsRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"notificationID": r.NotificationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationMarkAsRead) GetNotificationID() uint64 {
	return r.NotificationID
}

// Fill processes request and fills internal variables
func (r *NotificationMarkAsRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "notificationID")
		r.NotificationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewNotificationMarkAllAsRead request
func NewNotificationMarkAllAsRead() *NotificationMarkAllAsRead {
	return &NotificationMarkAllAsRead{}
}

// Auditable returns all auditable/loggable parameters
func (r NotificationMarkAllAsRead) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *NotificationMarkAllAsRead) Fill(req *http.Request) (err error) {

	return err
}
