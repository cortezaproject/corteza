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

// ReminderList request parameters
type ReminderList struct {
	hasReminderID bool
	rawReminderID []string
	ReminderID    []string

	hasResource bool
	rawResource string
	Resource    string

	hasAssignedTo bool
	rawAssignedTo string
	AssignedTo    uint64 `json:",string"`

	hasScheduledFrom bool
	rawScheduledFrom string
	ScheduledFrom    *time.Time

	hasScheduledUntil bool
	rawScheduledUntil string
	ScheduledUntil    *time.Time

	hasScheduledOnly bool
	rawScheduledOnly string
	ScheduledOnly    bool

	hasExcludeDismissed bool
	rawExcludeDismissed string
	ExcludeDismissed    bool

	hasLimit bool
	rawLimit string
	Limit    uint

	hasOffset bool
	rawOffset string
	Offset    uint

	hasPage bool
	rawPage string
	Page    uint

	hasPerPage bool
	rawPerPage string
	PerPage    uint

	hasSort bool
	rawSort string
	Sort    string
}

// NewReminderList request
func NewReminderList() *ReminderList {
	return &ReminderList{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID
	out["resource"] = r.Resource
	out["assignedTo"] = r.AssignedTo
	out["scheduledFrom"] = r.ScheduledFrom
	out["scheduledUntil"] = r.ScheduledUntil
	out["scheduledOnly"] = r.ScheduledOnly
	out["excludeDismissed"] = r.ExcludeDismissed
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort

	return out
}

// Fill processes request and fills internal variables
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

	if val, ok := urlQuery["reminderID[]"]; ok {
		r.hasReminderID = true
		r.rawReminderID = val
		r.ReminderID = parseStrings(val)
	} else if val, ok = urlQuery["reminderID"]; ok {
		r.hasReminderID = true
		r.rawReminderID = val
		r.ReminderID = parseStrings(val)
	}

	if val, ok := get["resource"]; ok {
		r.hasResource = true
		r.rawResource = val
		r.Resource = val
	}
	if val, ok := get["assignedTo"]; ok {
		r.hasAssignedTo = true
		r.rawAssignedTo = val
		r.AssignedTo = parseUInt64(val)
	}
	if val, ok := get["scheduledFrom"]; ok {
		r.hasScheduledFrom = true
		r.rawScheduledFrom = val

		if r.ScheduledFrom, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := get["scheduledUntil"]; ok {
		r.hasScheduledUntil = true
		r.rawScheduledUntil = val

		if r.ScheduledUntil, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := get["scheduledOnly"]; ok {
		r.hasScheduledOnly = true
		r.rawScheduledOnly = val
		r.ScheduledOnly = parseBool(val)
	}
	if val, ok := get["excludeDismissed"]; ok {
		r.hasExcludeDismissed = true
		r.rawExcludeDismissed = val
		r.ExcludeDismissed = parseBool(val)
	}
	if val, ok := get["limit"]; ok {
		r.hasLimit = true
		r.rawLimit = val
		r.Limit = parseUint(val)
	}
	if val, ok := get["offset"]; ok {
		r.hasOffset = true
		r.rawOffset = val
		r.Offset = parseUint(val)
	}
	if val, ok := get["page"]; ok {
		r.hasPage = true
		r.rawPage = val
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.hasPerPage = true
		r.rawPerPage = val
		r.PerPage = parseUint(val)
	}
	if val, ok := get["sort"]; ok {
		r.hasSort = true
		r.rawSort = val
		r.Sort = val
	}

	return err
}

var _ RequestFiller = NewReminderList()

// ReminderCreate request parameters
type ReminderCreate struct {
	hasResource bool
	rawResource string
	Resource    string

	hasAssignedTo bool
	rawAssignedTo string
	AssignedTo    uint64 `json:",string"`

	hasPayload bool
	rawPayload string
	Payload    sqlxTypes.JSONText

	hasRemindAt bool
	rawRemindAt string
	RemindAt    *time.Time
}

// NewReminderCreate request
func NewReminderCreate() *ReminderCreate {
	return &ReminderCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["resource"] = r.Resource
	out["assignedTo"] = r.AssignedTo
	out["payload"] = r.Payload
	out["remindAt"] = r.RemindAt

	return out
}

// Fill processes request and fills internal variables
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
		r.hasResource = true
		r.rawResource = val
		r.Resource = val
	}
	if val, ok := post["assignedTo"]; ok {
		r.hasAssignedTo = true
		r.rawAssignedTo = val
		r.AssignedTo = parseUInt64(val)
	}
	if val, ok := post["payload"]; ok {
		r.hasPayload = true
		r.rawPayload = val

		if r.Payload, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["remindAt"]; ok {
		r.hasRemindAt = true
		r.rawRemindAt = val

		if r.RemindAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewReminderCreate()

// ReminderUpdate request parameters
type ReminderUpdate struct {
	hasReminderID bool
	rawReminderID string
	ReminderID    uint64 `json:",string"`

	hasResource bool
	rawResource string
	Resource    string

	hasAssignedTo bool
	rawAssignedTo string
	AssignedTo    uint64 `json:",string"`

	hasPayload bool
	rawPayload string
	Payload    sqlxTypes.JSONText

	hasRemindAt bool
	rawRemindAt string
	RemindAt    *time.Time
}

// NewReminderUpdate request
func NewReminderUpdate() *ReminderUpdate {
	return &ReminderUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID
	out["resource"] = r.Resource
	out["assignedTo"] = r.AssignedTo
	out["payload"] = r.Payload
	out["remindAt"] = r.RemindAt

	return out
}

// Fill processes request and fills internal variables
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

	r.hasReminderID = true
	r.rawReminderID = chi.URLParam(req, "reminderID")
	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))
	if val, ok := post["resource"]; ok {
		r.hasResource = true
		r.rawResource = val
		r.Resource = val
	}
	if val, ok := post["assignedTo"]; ok {
		r.hasAssignedTo = true
		r.rawAssignedTo = val
		r.AssignedTo = parseUInt64(val)
	}
	if val, ok := post["payload"]; ok {
		r.hasPayload = true
		r.rawPayload = val

		if r.Payload, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["remindAt"]; ok {
		r.hasRemindAt = true
		r.rawRemindAt = val

		if r.RemindAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewReminderUpdate()

// ReminderRead request parameters
type ReminderRead struct {
	hasReminderID bool
	rawReminderID string
	ReminderID    uint64 `json:",string"`
}

// NewReminderRead request
func NewReminderRead() *ReminderRead {
	return &ReminderRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasReminderID = true
	r.rawReminderID = chi.URLParam(req, "reminderID")
	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))

	return err
}

var _ RequestFiller = NewReminderRead()

// ReminderDelete request parameters
type ReminderDelete struct {
	hasReminderID bool
	rawReminderID string
	ReminderID    uint64 `json:",string"`
}

// NewReminderDelete request
func NewReminderDelete() *ReminderDelete {
	return &ReminderDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasReminderID = true
	r.rawReminderID = chi.URLParam(req, "reminderID")
	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))

	return err
}

var _ RequestFiller = NewReminderDelete()

// ReminderDismiss request parameters
type ReminderDismiss struct {
	hasReminderID bool
	rawReminderID string
	ReminderID    uint64 `json:",string"`
}

// NewReminderDismiss request
func NewReminderDismiss() *ReminderDismiss {
	return &ReminderDismiss{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderDismiss) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID

	return out
}

// Fill processes request and fills internal variables
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

	r.hasReminderID = true
	r.rawReminderID = chi.URLParam(req, "reminderID")
	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))

	return err
}

var _ RequestFiller = NewReminderDismiss()

// ReminderSnooze request parameters
type ReminderSnooze struct {
	hasReminderID bool
	rawReminderID string
	ReminderID    uint64 `json:",string"`

	hasRemindAt bool
	rawRemindAt string
	RemindAt    *time.Time
}

// NewReminderSnooze request
func NewReminderSnooze() *ReminderSnooze {
	return &ReminderSnooze{}
}

// Auditable returns all auditable/loggable parameters
func (r ReminderSnooze) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["reminderID"] = r.ReminderID
	out["remindAt"] = r.RemindAt

	return out
}

// Fill processes request and fills internal variables
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

	r.hasReminderID = true
	r.rawReminderID = chi.URLParam(req, "reminderID")
	r.ReminderID = parseUInt64(chi.URLParam(req, "reminderID"))
	if val, ok := post["remindAt"]; ok {
		r.hasRemindAt = true
		r.rawRemindAt = val

		if r.RemindAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewReminderSnooze()

// HasReminderID returns true if reminderID was set
func (r *ReminderList) HasReminderID() bool {
	return r.hasReminderID
}

// RawReminderID returns raw value of reminderID parameter
func (r *ReminderList) RawReminderID() []string {
	return r.rawReminderID
}

// GetReminderID returns casted value of  reminderID parameter
func (r *ReminderList) GetReminderID() []string {
	return r.ReminderID
}

// HasResource returns true if resource was set
func (r *ReminderList) HasResource() bool {
	return r.hasResource
}

// RawResource returns raw value of resource parameter
func (r *ReminderList) RawResource() string {
	return r.rawResource
}

// GetResource returns casted value of  resource parameter
func (r *ReminderList) GetResource() string {
	return r.Resource
}

// HasAssignedTo returns true if assignedTo was set
func (r *ReminderList) HasAssignedTo() bool {
	return r.hasAssignedTo
}

// RawAssignedTo returns raw value of assignedTo parameter
func (r *ReminderList) RawAssignedTo() string {
	return r.rawAssignedTo
}

// GetAssignedTo returns casted value of  assignedTo parameter
func (r *ReminderList) GetAssignedTo() uint64 {
	return r.AssignedTo
}

// HasScheduledFrom returns true if scheduledFrom was set
func (r *ReminderList) HasScheduledFrom() bool {
	return r.hasScheduledFrom
}

// RawScheduledFrom returns raw value of scheduledFrom parameter
func (r *ReminderList) RawScheduledFrom() string {
	return r.rawScheduledFrom
}

// GetScheduledFrom returns casted value of  scheduledFrom parameter
func (r *ReminderList) GetScheduledFrom() *time.Time {
	return r.ScheduledFrom
}

// HasScheduledUntil returns true if scheduledUntil was set
func (r *ReminderList) HasScheduledUntil() bool {
	return r.hasScheduledUntil
}

// RawScheduledUntil returns raw value of scheduledUntil parameter
func (r *ReminderList) RawScheduledUntil() string {
	return r.rawScheduledUntil
}

// GetScheduledUntil returns casted value of  scheduledUntil parameter
func (r *ReminderList) GetScheduledUntil() *time.Time {
	return r.ScheduledUntil
}

// HasScheduledOnly returns true if scheduledOnly was set
func (r *ReminderList) HasScheduledOnly() bool {
	return r.hasScheduledOnly
}

// RawScheduledOnly returns raw value of scheduledOnly parameter
func (r *ReminderList) RawScheduledOnly() string {
	return r.rawScheduledOnly
}

// GetScheduledOnly returns casted value of  scheduledOnly parameter
func (r *ReminderList) GetScheduledOnly() bool {
	return r.ScheduledOnly
}

// HasExcludeDismissed returns true if excludeDismissed was set
func (r *ReminderList) HasExcludeDismissed() bool {
	return r.hasExcludeDismissed
}

// RawExcludeDismissed returns raw value of excludeDismissed parameter
func (r *ReminderList) RawExcludeDismissed() string {
	return r.rawExcludeDismissed
}

// GetExcludeDismissed returns casted value of  excludeDismissed parameter
func (r *ReminderList) GetExcludeDismissed() bool {
	return r.ExcludeDismissed
}

// HasLimit returns true if limit was set
func (r *ReminderList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *ReminderList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *ReminderList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *ReminderList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *ReminderList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *ReminderList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *ReminderList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *ReminderList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *ReminderList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *ReminderList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *ReminderList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *ReminderList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *ReminderList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *ReminderList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *ReminderList) GetSort() string {
	return r.Sort
}

// HasResource returns true if resource was set
func (r *ReminderCreate) HasResource() bool {
	return r.hasResource
}

// RawResource returns raw value of resource parameter
func (r *ReminderCreate) RawResource() string {
	return r.rawResource
}

// GetResource returns casted value of  resource parameter
func (r *ReminderCreate) GetResource() string {
	return r.Resource
}

// HasAssignedTo returns true if assignedTo was set
func (r *ReminderCreate) HasAssignedTo() bool {
	return r.hasAssignedTo
}

// RawAssignedTo returns raw value of assignedTo parameter
func (r *ReminderCreate) RawAssignedTo() string {
	return r.rawAssignedTo
}

// GetAssignedTo returns casted value of  assignedTo parameter
func (r *ReminderCreate) GetAssignedTo() uint64 {
	return r.AssignedTo
}

// HasPayload returns true if payload was set
func (r *ReminderCreate) HasPayload() bool {
	return r.hasPayload
}

// RawPayload returns raw value of payload parameter
func (r *ReminderCreate) RawPayload() string {
	return r.rawPayload
}

// GetPayload returns casted value of  payload parameter
func (r *ReminderCreate) GetPayload() sqlxTypes.JSONText {
	return r.Payload
}

// HasRemindAt returns true if remindAt was set
func (r *ReminderCreate) HasRemindAt() bool {
	return r.hasRemindAt
}

// RawRemindAt returns raw value of remindAt parameter
func (r *ReminderCreate) RawRemindAt() string {
	return r.rawRemindAt
}

// GetRemindAt returns casted value of  remindAt parameter
func (r *ReminderCreate) GetRemindAt() *time.Time {
	return r.RemindAt
}

// HasReminderID returns true if reminderID was set
func (r *ReminderUpdate) HasReminderID() bool {
	return r.hasReminderID
}

// RawReminderID returns raw value of reminderID parameter
func (r *ReminderUpdate) RawReminderID() string {
	return r.rawReminderID
}

// GetReminderID returns casted value of  reminderID parameter
func (r *ReminderUpdate) GetReminderID() uint64 {
	return r.ReminderID
}

// HasResource returns true if resource was set
func (r *ReminderUpdate) HasResource() bool {
	return r.hasResource
}

// RawResource returns raw value of resource parameter
func (r *ReminderUpdate) RawResource() string {
	return r.rawResource
}

// GetResource returns casted value of  resource parameter
func (r *ReminderUpdate) GetResource() string {
	return r.Resource
}

// HasAssignedTo returns true if assignedTo was set
func (r *ReminderUpdate) HasAssignedTo() bool {
	return r.hasAssignedTo
}

// RawAssignedTo returns raw value of assignedTo parameter
func (r *ReminderUpdate) RawAssignedTo() string {
	return r.rawAssignedTo
}

// GetAssignedTo returns casted value of  assignedTo parameter
func (r *ReminderUpdate) GetAssignedTo() uint64 {
	return r.AssignedTo
}

// HasPayload returns true if payload was set
func (r *ReminderUpdate) HasPayload() bool {
	return r.hasPayload
}

// RawPayload returns raw value of payload parameter
func (r *ReminderUpdate) RawPayload() string {
	return r.rawPayload
}

// GetPayload returns casted value of  payload parameter
func (r *ReminderUpdate) GetPayload() sqlxTypes.JSONText {
	return r.Payload
}

// HasRemindAt returns true if remindAt was set
func (r *ReminderUpdate) HasRemindAt() bool {
	return r.hasRemindAt
}

// RawRemindAt returns raw value of remindAt parameter
func (r *ReminderUpdate) RawRemindAt() string {
	return r.rawRemindAt
}

// GetRemindAt returns casted value of  remindAt parameter
func (r *ReminderUpdate) GetRemindAt() *time.Time {
	return r.RemindAt
}

// HasReminderID returns true if reminderID was set
func (r *ReminderRead) HasReminderID() bool {
	return r.hasReminderID
}

// RawReminderID returns raw value of reminderID parameter
func (r *ReminderRead) RawReminderID() string {
	return r.rawReminderID
}

// GetReminderID returns casted value of  reminderID parameter
func (r *ReminderRead) GetReminderID() uint64 {
	return r.ReminderID
}

// HasReminderID returns true if reminderID was set
func (r *ReminderDelete) HasReminderID() bool {
	return r.hasReminderID
}

// RawReminderID returns raw value of reminderID parameter
func (r *ReminderDelete) RawReminderID() string {
	return r.rawReminderID
}

// GetReminderID returns casted value of  reminderID parameter
func (r *ReminderDelete) GetReminderID() uint64 {
	return r.ReminderID
}

// HasReminderID returns true if reminderID was set
func (r *ReminderDismiss) HasReminderID() bool {
	return r.hasReminderID
}

// RawReminderID returns raw value of reminderID parameter
func (r *ReminderDismiss) RawReminderID() string {
	return r.rawReminderID
}

// GetReminderID returns casted value of  reminderID parameter
func (r *ReminderDismiss) GetReminderID() uint64 {
	return r.ReminderID
}

// HasReminderID returns true if reminderID was set
func (r *ReminderSnooze) HasReminderID() bool {
	return r.hasReminderID
}

// RawReminderID returns raw value of reminderID parameter
func (r *ReminderSnooze) RawReminderID() string {
	return r.rawReminderID
}

// GetReminderID returns casted value of  reminderID parameter
func (r *ReminderSnooze) GetReminderID() uint64 {
	return r.ReminderID
}

// HasRemindAt returns true if remindAt was set
func (r *ReminderSnooze) HasRemindAt() bool {
	return r.hasRemindAt
}

// RawRemindAt returns raw value of remindAt parameter
func (r *ReminderSnooze) RawRemindAt() string {
	return r.rawRemindAt
}

// GetRemindAt returns casted value of  remindAt parameter
func (r *ReminderSnooze) GetRemindAt() *time.Time {
	return r.RemindAt
}
