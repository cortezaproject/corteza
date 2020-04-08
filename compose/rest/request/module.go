package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
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

	"github.com/cortezaproject/corteza-server/compose/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"time"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// ModuleList request parameters
type ModuleList struct {
	hasQuery bool
	rawQuery string
	Query    string

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

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

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewModuleList request
func NewModuleList() *ModuleList {
	return &ModuleList{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ModuleList) Fill(req *http.Request) (err error) {
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

	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}
	if val, ok := get["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := get["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
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
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleList()

// ModuleCreate request parameters
type ModuleCreate struct {
	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasFields bool
	rawFields string
	Fields    types.ModuleFieldSet

	hasMeta bool
	rawMeta string
	Meta    sqlxTypes.JSONText

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewModuleCreate request
func NewModuleCreate() *ModuleCreate {
	return &ModuleCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["handle"] = r.Handle
	out["fields"] = r.Fields
	out["meta"] = r.Meta
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ModuleCreate) Fill(req *http.Request) (err error) {
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
	if val, ok := post["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
	}
	if val, ok := post["meta"]; ok {
		r.hasMeta = true
		r.rawMeta = val

		if r.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleCreate()

// ModuleRead request parameters
type ModuleRead struct {
	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewModuleRead request
func NewModuleRead() *ModuleRead {
	return &ModuleRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["moduleID"] = r.ModuleID
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ModuleRead) Fill(req *http.Request) (err error) {
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

	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleRead()

// ModuleUpdate request parameters
type ModuleUpdate struct {
	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasFields bool
	rawFields string
	Fields    types.ModuleFieldSet

	hasMeta bool
	rawMeta string
	Meta    sqlxTypes.JSONText

	hasUpdatedAt bool
	rawUpdatedAt string
	UpdatedAt    *time.Time
}

// NewModuleUpdate request
func NewModuleUpdate() *ModuleUpdate {
	return &ModuleUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["moduleID"] = r.ModuleID
	out["namespaceID"] = r.NamespaceID
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["fields"] = r.Fields
	out["meta"] = r.Meta
	out["updatedAt"] = r.UpdatedAt

	return out
}

// Fill processes request and fills internal variables
func (r *ModuleUpdate) Fill(req *http.Request) (err error) {
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

	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := post["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
	}
	if val, ok := post["meta"]; ok {
		r.hasMeta = true
		r.rawMeta = val

		if r.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["updatedAt"]; ok {
		r.hasUpdatedAt = true
		r.rawUpdatedAt = val

		if r.UpdatedAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleUpdate()

// ModuleDelete request parameters
type ModuleDelete struct {
	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewModuleDelete request
func NewModuleDelete() *ModuleDelete {
	return &ModuleDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["moduleID"] = r.ModuleID
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ModuleDelete) Fill(req *http.Request) (err error) {
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

	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleDelete()

// ModuleTriggerScript request parameters
type ModuleTriggerScript struct {
	hasModuleID bool
	rawModuleID string
	ModuleID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string
}

// NewModuleTriggerScript request
func NewModuleTriggerScript() *ModuleTriggerScript {
	return &ModuleTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r ModuleTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["moduleID"] = r.ModuleID
	out["namespaceID"] = r.NamespaceID
	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
func (r *ModuleTriggerScript) Fill(req *http.Request) (err error) {
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

	r.hasModuleID = true
	r.rawModuleID = chi.URLParam(req, "moduleID")
	r.ModuleID = parseUInt64(chi.URLParam(req, "moduleID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewModuleTriggerScript()

// HasQuery returns true if query was set
func (r *ModuleList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *ModuleList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *ModuleList) GetQuery() string {
	return r.Query
}

// HasName returns true if name was set
func (r *ModuleList) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ModuleList) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ModuleList) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *ModuleList) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *ModuleList) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *ModuleList) GetHandle() string {
	return r.Handle
}

// HasLimit returns true if limit was set
func (r *ModuleList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *ModuleList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *ModuleList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *ModuleList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *ModuleList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *ModuleList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *ModuleList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *ModuleList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *ModuleList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *ModuleList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *ModuleList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *ModuleList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *ModuleList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *ModuleList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *ModuleList) GetSort() string {
	return r.Sort
}

// HasNamespaceID returns true if namespaceID was set
func (r *ModuleList) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ModuleList) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ModuleList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasName returns true if name was set
func (r *ModuleCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ModuleCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ModuleCreate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *ModuleCreate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *ModuleCreate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *ModuleCreate) GetHandle() string {
	return r.Handle
}

// HasFields returns true if fields was set
func (r *ModuleCreate) HasFields() bool {
	return r.hasFields
}

// RawFields returns raw value of fields parameter
func (r *ModuleCreate) RawFields() string {
	return r.rawFields
}

// GetFields returns casted value of  fields parameter
func (r *ModuleCreate) GetFields() types.ModuleFieldSet {
	return r.Fields
}

// HasMeta returns true if meta was set
func (r *ModuleCreate) HasMeta() bool {
	return r.hasMeta
}

// RawMeta returns raw value of meta parameter
func (r *ModuleCreate) RawMeta() string {
	return r.rawMeta
}

// GetMeta returns casted value of  meta parameter
func (r *ModuleCreate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// HasNamespaceID returns true if namespaceID was set
func (r *ModuleCreate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ModuleCreate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ModuleCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *ModuleRead) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *ModuleRead) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *ModuleRead) GetModuleID() uint64 {
	return r.ModuleID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ModuleRead) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ModuleRead) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ModuleRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *ModuleUpdate) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *ModuleUpdate) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *ModuleUpdate) GetModuleID() uint64 {
	return r.ModuleID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ModuleUpdate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ModuleUpdate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ModuleUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasName returns true if name was set
func (r *ModuleUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ModuleUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ModuleUpdate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *ModuleUpdate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *ModuleUpdate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *ModuleUpdate) GetHandle() string {
	return r.Handle
}

// HasFields returns true if fields was set
func (r *ModuleUpdate) HasFields() bool {
	return r.hasFields
}

// RawFields returns raw value of fields parameter
func (r *ModuleUpdate) RawFields() string {
	return r.rawFields
}

// GetFields returns casted value of  fields parameter
func (r *ModuleUpdate) GetFields() types.ModuleFieldSet {
	return r.Fields
}

// HasMeta returns true if meta was set
func (r *ModuleUpdate) HasMeta() bool {
	return r.hasMeta
}

// RawMeta returns raw value of meta parameter
func (r *ModuleUpdate) RawMeta() string {
	return r.rawMeta
}

// GetMeta returns casted value of  meta parameter
func (r *ModuleUpdate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// HasUpdatedAt returns true if updatedAt was set
func (r *ModuleUpdate) HasUpdatedAt() bool {
	return r.hasUpdatedAt
}

// RawUpdatedAt returns raw value of updatedAt parameter
func (r *ModuleUpdate) RawUpdatedAt() string {
	return r.rawUpdatedAt
}

// GetUpdatedAt returns casted value of  updatedAt parameter
func (r *ModuleUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// HasModuleID returns true if moduleID was set
func (r *ModuleDelete) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *ModuleDelete) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *ModuleDelete) GetModuleID() uint64 {
	return r.ModuleID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ModuleDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ModuleDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ModuleDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasModuleID returns true if moduleID was set
func (r *ModuleTriggerScript) HasModuleID() bool {
	return r.hasModuleID
}

// RawModuleID returns raw value of moduleID parameter
func (r *ModuleTriggerScript) RawModuleID() string {
	return r.rawModuleID
}

// GetModuleID returns casted value of  moduleID parameter
func (r *ModuleTriggerScript) GetModuleID() uint64 {
	return r.ModuleID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ModuleTriggerScript) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ModuleTriggerScript) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ModuleTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasScript returns true if script was set
func (r *ModuleTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *ModuleTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *ModuleTriggerScript) GetScript() string {
	return r.Script
}
