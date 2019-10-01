package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `reminder.go`, `reminder.util.go` or `reminder_test.go` to
	implement your API calls, helper functions and tests. The file `reminder.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	sqlxTypes "github.com/jmoiron/sqlx/types"
	"time"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Reminder list request parameters
type ReminderList struct {
	Resource         string
	AssignedTo       uint64 `json:",string"`
	ScheduledFrom    *time.Time
	ScheduledUntil   *time.Time
	ScheduledOnly    bool
	ExcludeDismissed bool
	Page             uint
	PerPage          uint
}

func NewReminderList() *ReminderList {
	return &ReminderList{}
}

func (r ReminderList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource
	out["assignedTo"] = r.AssignedTo
	out["scheduledFrom"] = r.ScheduledFrom
	out["scheduledUntil"] = r.ScheduledUntil
	out["scheduledOnly"] = r.ScheduledOnly
	out["excludeDismissed"] = r.ExcludeDismissed
	out["page"] = r.Page
	out["perPage"] = r.PerPage

	return out
}

func (r *ReminderList) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["resource"]; ok {
		r.Resource = val
	}
	if val, ok := get["assignedTo"]; ok {
		r.AssignedTo = parseUInt64(val)
	}
	if val, ok := get["scheduledFrom"]; ok {

		if r.ScheduledFrom, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := get["scheduledUntil"]; ok {

		if r.ScheduledUntil, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := get["scheduledOnly"]; ok {
		r.ScheduledOnly = parseBool(val)
	}
	if val, ok := get["excludeDismissed"]; ok {
		r.ExcludeDismissed = parseBool(val)
	}
	if val, ok := get["page"]; ok {
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.PerPage = parseUint(val)
	}

	return err
}

var _ RequestFiller = NewReminderList()

// Reminder create request parameters
type ReminderCreate struct {
	Resource   string
	AssignedTo uint64 `json:",string"`
	Payload    sqlxTypes.JSONText
	RemindAt   *time.Time
}

func NewReminderCreate() *ReminderCreate {
	return &ReminderCreate{}
}

func (r ReminderCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource
	out["assignedTo"] = r.AssignedTo
	out["payload"] = r.Payload
	out["remindAt"] = r.RemindAt

	return out
}

func (r *ReminderCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := post["resource"]; ok {
		r.Resource = val
	}
	if val, ok := post["assignedTo"]; ok {
		r.AssignedTo = parseUInt64(val)
	}
	if val, ok := post["payload"]; ok {

		if r.Payload, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["remindAt"]; ok {

		if r.RemindAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewReminderCreate()

// Reminder update request parameters
type ReminderUpdate struct {
	ReminderID uint64 `json:",string"`
	Resource   string
	AssignedTo uint64 `json:",string"`
	Payload    sqlxTypes.JSONText
	RemindAt   *time.Time
}

func NewReminderUpdate() *ReminderUpdate {
	return &ReminderUpdate{}
}

func (r ReminderUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID
	out["resource"] = r.Resource
	out["assignedTo"] = r.AssignedTo
	out["payload"] = r.Payload
	out["remindAt"] = r.RemindAt

	return out
}

func (r *ReminderUpdate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))
	if val, ok := post["resource"]; ok {
		r.Resource = val
	}
	if val, ok := post["assignedTo"]; ok {
		r.AssignedTo = parseUInt64(val)
	}
	if val, ok := post["payload"]; ok {

		if r.Payload, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["remindAt"]; ok {

		if r.RemindAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewReminderUpdate()

// Reminder read request parameters
type ReminderRead struct {
	ReminderID uint64 `json:",string"`
}

func NewReminderRead() *ReminderRead {
	return &ReminderRead{}
}

func (r ReminderRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID

	return out
}

func (r *ReminderRead) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))

	return err
}

var _ RequestFiller = NewReminderRead()

// Reminder delete request parameters
type ReminderDelete struct {
	ReminderID uint64 `json:",string"`
}

func NewReminderDelete() *ReminderDelete {
	return &ReminderDelete{}
}

func (r ReminderDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID

	return out
}

func (r *ReminderDelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))

	return err
}

var _ RequestFiller = NewReminderDelete()

// Reminder dismiss request parameters
type ReminderDismiss struct {
	ReminderID uint64 `json:",string"`
}

func NewReminderDismiss() *ReminderDismiss {
	return &ReminderDismiss{}
}

func (r ReminderDismiss) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID

	return out
}

func (r *ReminderDismiss) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))

	return err
}

var _ RequestFiller = NewReminderDismiss()

// Reminder snooze request parameters
type ReminderSnooze struct {
	ReminderID uint64 `json:",string"`
	RemindAt   *time.Time
}

func NewReminderSnooze() *ReminderSnooze {
	return &ReminderSnooze{}
}

func (r ReminderSnooze) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID
	out["remindAt"] = r.RemindAt

	return out
}

func (r *ReminderSnooze) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))
	if val, ok := post["remindAt"]; ok {

		if r.RemindAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewReminderSnooze()
