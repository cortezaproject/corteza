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
	ReminderList struct {
		// ReminderID GET parameter
		//
		// Filter by reminder ID
		ReminderID []string

		// Resource GET parameter
		//
		// Only reminders of a specific resource
		Resource string

		// AssignedTo GET parameter
		//
		// Only reminders for a given user
		AssignedTo uint64 `json:",string"`

		// ScheduledFrom GET parameter
		//
		// Only reminders from this time (included)
		ScheduledFrom *time.Time

		// ScheduledUntil GET parameter
		//
		// Only reminders up to this time (included)
		ScheduledUntil *time.Time

		// ScheduledOnly GET parameter
		//
		// Only scheduled reminders
		ScheduledOnly bool

		// ExcludeDismissed GET parameter
		//
		// Filter out dismissed reminders
		ExcludeDismissed bool

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// Offset GET parameter
		//
		// Offset
		Offset uint

		// Page GET parameter
		//
		// Page number (1-based)
		Page uint

		// PerPage GET parameter
		//
		// Returned items per page (default 50)
		PerPage uint

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	ReminderCreate struct {
		// Resource POST parameter
		//
		// Resource
		Resource string

		// AssignedTo POST parameter
		//
		// Assigned To
		AssignedTo uint64 `json:",string"`

		// Payload POST parameter
		//
		// Payload
		Payload sqlxTypes.JSONText

		// RemindAt POST parameter
		//
		// Remind At
		RemindAt *time.Time
	}

	ReminderUpdate struct {
		// ReminderID PATH parameter
		//
		// Reminder ID
		ReminderID uint64 `json:",string"`

		// Resource POST parameter
		//
		// Resource
		Resource string

		// AssignedTo POST parameter
		//
		// Assigned To
		AssignedTo uint64 `json:",string"`

		// Payload POST parameter
		//
		// Payload
		Payload sqlxTypes.JSONText

		// RemindAt POST parameter
		//
		// Remind At
		RemindAt *time.Time
	}

	ReminderRead struct {
		// ReminderID PATH parameter
		//
		// Reminder ID
		ReminderID uint64 `json:",string"`
	}

	ReminderDelete struct {
		// ReminderID PATH parameter
		//
		// Reminder ID
		ReminderID uint64 `json:",string"`
	}

	ReminderDismiss struct {
		// ReminderID PATH parameter
		//
		// reminder ID
		ReminderID uint64 `json:",string"`
	}

	ReminderSnooze struct {
		// ReminderID PATH parameter
		//
		// reminder ID
		ReminderID uint64 `json:",string"`

		// RemindAt POST parameter
		//
		// New Remind At Time
		RemindAt *time.Time
	}
)

// NewReminderList request
func NewReminderList() *ReminderList {
	return &ReminderList{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reminderID":       r.ReminderID,
		"resource":         r.Resource,
		"assignedTo":       r.AssignedTo,
		"scheduledFrom":    r.ScheduledFrom,
		"scheduledUntil":   r.ScheduledUntil,
		"scheduledOnly":    r.ScheduledOnly,
		"excludeDismissed": r.ExcludeDismissed,
		"limit":            r.Limit,
		"offset":           r.Offset,
		"page":             r.Page,
		"perPage":          r.PerPage,
		"sort":             r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetReminderID() []string {
	return r.ReminderID
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetAssignedTo() uint64 {
	return r.AssignedTo
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetScheduledFrom() *time.Time {
	return r.ScheduledFrom
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetScheduledUntil() *time.Time {
	return r.ScheduledUntil
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetScheduledOnly() bool {
	return r.ScheduledOnly
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetExcludeDismissed() bool {
	return r.ExcludeDismissed
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetOffset() uint {
	return r.Offset
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetPage() uint {
	return r.Page
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetPerPage() uint {
	return r.PerPage
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ReminderList) Fill(req *http.Request) (err error) {
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

		if val, ok := tmp["reminderID[]"]; ok {
			r.ReminderID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["reminderID"]; ok {
			r.ReminderID, err = val, nil
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
		if val, ok := tmp["assignedTo"]; ok && len(val) > 0 {
			r.AssignedTo, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["scheduledFrom"]; ok && len(val) > 0 {
			r.ScheduledFrom, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["scheduledUntil"]; ok && len(val) > 0 {
			r.ScheduledUntil, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["scheduledOnly"]; ok && len(val) > 0 {
			r.ScheduledOnly, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["excludeDismissed"]; ok && len(val) > 0 {
			r.ExcludeDismissed, err = payload.ParseBool(val[0]), nil
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
		if val, ok := tmp["offset"]; ok && len(val) > 0 {
			r.Offset, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["page"]; ok && len(val) > 0 {
			r.Page, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["perPage"]; ok && len(val) > 0 {
			r.PerPage, err = payload.ParseUint(val[0]), nil
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

// NewReminderCreate request
func NewReminderCreate() *ReminderCreate {
	return &ReminderCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"resource":   r.Resource,
		"assignedTo": r.AssignedTo,
		"payload":    r.Payload,
		"remindAt":   r.RemindAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderCreate) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ReminderCreate) GetAssignedTo() uint64 {
	return r.AssignedTo
}

// Auditable returns all auditable/loggable parameters
func (r ReminderCreate) GetPayload() sqlxTypes.JSONText {
	return r.Payload
}

// Auditable returns all auditable/loggable parameters
func (r ReminderCreate) GetRemindAt() *time.Time {
	return r.RemindAt
}

// Fill processes request and fills internal variables
func (r *ReminderCreate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["assignedTo"]; ok && len(val) > 0 {
			r.AssignedTo, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["payload"]; ok && len(val) > 0 {
			r.Payload, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["remindAt"]; ok && len(val) > 0 {
			r.RemindAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewReminderUpdate request
func NewReminderUpdate() *ReminderUpdate {
	return &ReminderUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reminderID": r.ReminderID,
		"resource":   r.Resource,
		"assignedTo": r.AssignedTo,
		"payload":    r.Payload,
		"remindAt":   r.RemindAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) GetReminderID() uint64 {
	return r.ReminderID
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) GetResource() string {
	return r.Resource
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) GetAssignedTo() uint64 {
	return r.AssignedTo
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) GetPayload() sqlxTypes.JSONText {
	return r.Payload
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) GetRemindAt() *time.Time {
	return r.RemindAt
}

// Fill processes request and fills internal variables
func (r *ReminderUpdate) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["resource"]; ok && len(val) > 0 {
			r.Resource, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["assignedTo"]; ok && len(val) > 0 {
			r.AssignedTo, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["payload"]; ok && len(val) > 0 {
			r.Payload, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["remindAt"]; ok && len(val) > 0 {
			r.RemindAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "reminderID")
		r.ReminderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReminderRead request
func NewReminderRead() *ReminderRead {
	return &ReminderRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reminderID": r.ReminderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderRead) GetReminderID() uint64 {
	return r.ReminderID
}

// Fill processes request and fills internal variables
func (r *ReminderRead) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "reminderID")
		r.ReminderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReminderDelete request
func NewReminderDelete() *ReminderDelete {
	return &ReminderDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reminderID": r.ReminderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderDelete) GetReminderID() uint64 {
	return r.ReminderID
}

// Fill processes request and fills internal variables
func (r *ReminderDelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "reminderID")
		r.ReminderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReminderDismiss request
func NewReminderDismiss() *ReminderDismiss {
	return &ReminderDismiss{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderDismiss) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reminderID": r.ReminderID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderDismiss) GetReminderID() uint64 {
	return r.ReminderID
}

// Fill processes request and fills internal variables
func (r *ReminderDismiss) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "reminderID")
		r.ReminderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewReminderSnooze request
func NewReminderSnooze() *ReminderSnooze {
	return &ReminderSnooze{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderSnooze) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"reminderID": r.ReminderID,
		"remindAt":   r.RemindAt,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderSnooze) GetReminderID() uint64 {
	return r.ReminderID
}

// Auditable returns all auditable/loggable parameters
func (r ReminderSnooze) GetRemindAt() *time.Time {
	return r.RemindAt
}

// Fill processes request and fills internal variables
func (r *ReminderSnooze) Fill(req *http.Request) (err error) {
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

		if val, ok := req.Form["remindAt"]; ok && len(val) > 0 {
			r.RemindAt, err = payload.ParseISODatePtrWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "reminderID")
		r.ReminderID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
