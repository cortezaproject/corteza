package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `application.go`, `application.util.go` or `application_test.go` to
	implement your API calls, helper functions and tests. The file `application.go`
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
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// ApplicationList request parameters
type ApplicationList struct {
	hasName bool
	rawName string
	Name    string

	hasQuery bool
	rawQuery string
	Query    string

	hasDeleted bool
	rawDeleted string
	Deleted    uint

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

// NewApplicationList request
func NewApplicationList() *ApplicationList {
	return &ApplicationList{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["query"] = r.Query
	out["deleted"] = r.Deleted
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}
	if val, ok := get["deleted"]; ok {
		r.hasDeleted = true
		r.rawDeleted = val
		r.Deleted = parseUint(val)
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

var _ RequestFiller = NewApplicationList()

// ApplicationCreate request parameters
type ApplicationCreate struct {
	hasName bool
	rawName string
	Name    string

	hasEnabled bool
	rawEnabled string
	Enabled    bool

	hasUnify bool
	rawUnify string
	Unify    sqlxTypes.JSONText

	hasConfig bool
	rawConfig string
	Config    sqlxTypes.JSONText
}

// NewApplicationCreate request
func NewApplicationCreate() *ApplicationCreate {
	return &ApplicationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["enabled"] = r.Enabled
	out["unify"] = r.Unify
	out["config"] = r.Config

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := post["enabled"]; ok {
		r.hasEnabled = true
		r.rawEnabled = val
		r.Enabled = parseBool(val)
	}
	if val, ok := post["unify"]; ok {
		r.hasUnify = true
		r.rawUnify = val

		if r.Unify, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["config"]; ok {
		r.hasConfig = true
		r.rawConfig = val

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewApplicationCreate()

// ApplicationUpdate request parameters
type ApplicationUpdate struct {
	hasApplicationID bool
	rawApplicationID string
	ApplicationID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasEnabled bool
	rawEnabled string
	Enabled    bool

	hasUnify bool
	rawUnify string
	Unify    sqlxTypes.JSONText

	hasConfig bool
	rawConfig string
	Config    sqlxTypes.JSONText
}

// NewApplicationUpdate request
func NewApplicationUpdate() *ApplicationUpdate {
	return &ApplicationUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID
	out["name"] = r.Name
	out["enabled"] = r.Enabled
	out["unify"] = r.Unify
	out["config"] = r.Config

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationUpdate) Fill(req *http.Request) (err error) {
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

	r.hasApplicationID = true
	r.rawApplicationID = chi.URLParam(req, "applicationID")
	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))
	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := post["enabled"]; ok {
		r.hasEnabled = true
		r.rawEnabled = val
		r.Enabled = parseBool(val)
	}
	if val, ok := post["unify"]; ok {
		r.hasUnify = true
		r.rawUnify = val

		if r.Unify, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["config"]; ok {
		r.hasConfig = true
		r.rawConfig = val

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewApplicationUpdate()

// ApplicationRead request parameters
type ApplicationRead struct {
	hasApplicationID bool
	rawApplicationID string
	ApplicationID    uint64 `json:",string"`
}

// NewApplicationRead request
func NewApplicationRead() *ApplicationRead {
	return &ApplicationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationRead) Fill(req *http.Request) (err error) {
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

	r.hasApplicationID = true
	r.rawApplicationID = chi.URLParam(req, "applicationID")
	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))

	return err
}

var _ RequestFiller = NewApplicationRead()

// ApplicationDelete request parameters
type ApplicationDelete struct {
	hasApplicationID bool
	rawApplicationID string
	ApplicationID    uint64 `json:",string"`
}

// NewApplicationDelete request
func NewApplicationDelete() *ApplicationDelete {
	return &ApplicationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationDelete) Fill(req *http.Request) (err error) {
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

	r.hasApplicationID = true
	r.rawApplicationID = chi.URLParam(req, "applicationID")
	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))

	return err
}

var _ RequestFiller = NewApplicationDelete()

// ApplicationUndelete request parameters
type ApplicationUndelete struct {
	hasApplicationID bool
	rawApplicationID string
	ApplicationID    uint64 `json:",string"`
}

// NewApplicationUndelete request
func NewApplicationUndelete() *ApplicationUndelete {
	return &ApplicationUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUndelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationUndelete) Fill(req *http.Request) (err error) {
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

	r.hasApplicationID = true
	r.rawApplicationID = chi.URLParam(req, "applicationID")
	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))

	return err
}

var _ RequestFiller = NewApplicationUndelete()

// ApplicationTriggerScript request parameters
type ApplicationTriggerScript struct {
	hasApplicationID bool
	rawApplicationID string
	ApplicationID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string
}

// NewApplicationTriggerScript request
func NewApplicationTriggerScript() *ApplicationTriggerScript {
	return &ApplicationTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["applicationID"] = r.ApplicationID
	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
func (r *ApplicationTriggerScript) Fill(req *http.Request) (err error) {
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

	r.hasApplicationID = true
	r.rawApplicationID = chi.URLParam(req, "applicationID")
	r.ApplicationID = parseUInt64(chi.URLParam(req, "applicationID"))
	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewApplicationTriggerScript()

// HasName returns true if name was set
func (r *ApplicationList) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ApplicationList) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ApplicationList) GetName() string {
	return r.Name
}

// HasQuery returns true if query was set
func (r *ApplicationList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *ApplicationList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *ApplicationList) GetQuery() string {
	return r.Query
}

// HasDeleted returns true if deleted was set
func (r *ApplicationList) HasDeleted() bool {
	return r.hasDeleted
}

// RawDeleted returns raw value of deleted parameter
func (r *ApplicationList) RawDeleted() string {
	return r.rawDeleted
}

// GetDeleted returns casted value of  deleted parameter
func (r *ApplicationList) GetDeleted() uint {
	return r.Deleted
}

// HasLimit returns true if limit was set
func (r *ApplicationList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *ApplicationList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *ApplicationList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *ApplicationList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *ApplicationList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *ApplicationList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *ApplicationList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *ApplicationList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *ApplicationList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *ApplicationList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *ApplicationList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *ApplicationList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *ApplicationList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *ApplicationList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *ApplicationList) GetSort() string {
	return r.Sort
}

// HasName returns true if name was set
func (r *ApplicationCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ApplicationCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ApplicationCreate) GetName() string {
	return r.Name
}

// HasEnabled returns true if enabled was set
func (r *ApplicationCreate) HasEnabled() bool {
	return r.hasEnabled
}

// RawEnabled returns raw value of enabled parameter
func (r *ApplicationCreate) RawEnabled() string {
	return r.rawEnabled
}

// GetEnabled returns casted value of  enabled parameter
func (r *ApplicationCreate) GetEnabled() bool {
	return r.Enabled
}

// HasUnify returns true if unify was set
func (r *ApplicationCreate) HasUnify() bool {
	return r.hasUnify
}

// RawUnify returns raw value of unify parameter
func (r *ApplicationCreate) RawUnify() string {
	return r.rawUnify
}

// GetUnify returns casted value of  unify parameter
func (r *ApplicationCreate) GetUnify() sqlxTypes.JSONText {
	return r.Unify
}

// HasConfig returns true if config was set
func (r *ApplicationCreate) HasConfig() bool {
	return r.hasConfig
}

// RawConfig returns raw value of config parameter
func (r *ApplicationCreate) RawConfig() string {
	return r.rawConfig
}

// GetConfig returns casted value of  config parameter
func (r *ApplicationCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// HasApplicationID returns true if applicationID was set
func (r *ApplicationUpdate) HasApplicationID() bool {
	return r.hasApplicationID
}

// RawApplicationID returns raw value of applicationID parameter
func (r *ApplicationUpdate) RawApplicationID() string {
	return r.rawApplicationID
}

// GetApplicationID returns casted value of  applicationID parameter
func (r *ApplicationUpdate) GetApplicationID() uint64 {
	return r.ApplicationID
}

// HasName returns true if name was set
func (r *ApplicationUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ApplicationUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ApplicationUpdate) GetName() string {
	return r.Name
}

// HasEnabled returns true if enabled was set
func (r *ApplicationUpdate) HasEnabled() bool {
	return r.hasEnabled
}

// RawEnabled returns raw value of enabled parameter
func (r *ApplicationUpdate) RawEnabled() string {
	return r.rawEnabled
}

// GetEnabled returns casted value of  enabled parameter
func (r *ApplicationUpdate) GetEnabled() bool {
	return r.Enabled
}

// HasUnify returns true if unify was set
func (r *ApplicationUpdate) HasUnify() bool {
	return r.hasUnify
}

// RawUnify returns raw value of unify parameter
func (r *ApplicationUpdate) RawUnify() string {
	return r.rawUnify
}

// GetUnify returns casted value of  unify parameter
func (r *ApplicationUpdate) GetUnify() sqlxTypes.JSONText {
	return r.Unify
}

// HasConfig returns true if config was set
func (r *ApplicationUpdate) HasConfig() bool {
	return r.hasConfig
}

// RawConfig returns raw value of config parameter
func (r *ApplicationUpdate) RawConfig() string {
	return r.rawConfig
}

// GetConfig returns casted value of  config parameter
func (r *ApplicationUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// HasApplicationID returns true if applicationID was set
func (r *ApplicationRead) HasApplicationID() bool {
	return r.hasApplicationID
}

// RawApplicationID returns raw value of applicationID parameter
func (r *ApplicationRead) RawApplicationID() string {
	return r.rawApplicationID
}

// GetApplicationID returns casted value of  applicationID parameter
func (r *ApplicationRead) GetApplicationID() uint64 {
	return r.ApplicationID
}

// HasApplicationID returns true if applicationID was set
func (r *ApplicationDelete) HasApplicationID() bool {
	return r.hasApplicationID
}

// RawApplicationID returns raw value of applicationID parameter
func (r *ApplicationDelete) RawApplicationID() string {
	return r.rawApplicationID
}

// GetApplicationID returns casted value of  applicationID parameter
func (r *ApplicationDelete) GetApplicationID() uint64 {
	return r.ApplicationID
}

// HasApplicationID returns true if applicationID was set
func (r *ApplicationUndelete) HasApplicationID() bool {
	return r.hasApplicationID
}

// RawApplicationID returns raw value of applicationID parameter
func (r *ApplicationUndelete) RawApplicationID() string {
	return r.rawApplicationID
}

// GetApplicationID returns casted value of  applicationID parameter
func (r *ApplicationUndelete) GetApplicationID() uint64 {
	return r.ApplicationID
}

// HasApplicationID returns true if applicationID was set
func (r *ApplicationTriggerScript) HasApplicationID() bool {
	return r.hasApplicationID
}

// RawApplicationID returns raw value of applicationID parameter
func (r *ApplicationTriggerScript) RawApplicationID() string {
	return r.rawApplicationID
}

// GetApplicationID returns casted value of  applicationID parameter
func (r *ApplicationTriggerScript) GetApplicationID() uint64 {
	return r.ApplicationID
}

// HasScript returns true if script was set
func (r *ApplicationTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *ApplicationTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *ApplicationTriggerScript) GetScript() string {
	return r.Script
}
