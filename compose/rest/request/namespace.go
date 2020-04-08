package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `namespace.go`, `namespace.util.go` or `namespace_test.go` to
	implement your API calls, helper functions and tests. The file `namespace.go`
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

// NamespaceList request parameters
type NamespaceList struct {
	hasQuery bool
	rawQuery string
	Query    string

	hasSlug bool
	rawSlug string
	Slug    string

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

// NewNamespaceList request
func NewNamespaceList() *NamespaceList {
	return &NamespaceList{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["slug"] = r.Slug
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort

	return out
}

// Fill processes request and fills internal variables
func (r *NamespaceList) Fill(req *http.Request) (err error) {
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
	if val, ok := get["slug"]; ok {
		r.hasSlug = true
		r.rawSlug = val
		r.Slug = val
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

var _ RequestFiller = NewNamespaceList()

// NamespaceCreate request parameters
type NamespaceCreate struct {
	hasName bool
	rawName string
	Name    string

	hasSlug bool
	rawSlug string
	Slug    string

	hasEnabled bool
	rawEnabled string
	Enabled    bool

	hasMeta bool
	rawMeta string
	Meta    sqlxTypes.JSONText
}

// NewNamespaceCreate request
func NewNamespaceCreate() *NamespaceCreate {
	return &NamespaceCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["slug"] = r.Slug
	out["enabled"] = r.Enabled
	out["meta"] = r.Meta

	return out
}

// Fill processes request and fills internal variables
func (r *NamespaceCreate) Fill(req *http.Request) (err error) {
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
	if val, ok := post["slug"]; ok {
		r.hasSlug = true
		r.rawSlug = val
		r.Slug = val
	}
	if val, ok := post["enabled"]; ok {
		r.hasEnabled = true
		r.rawEnabled = val
		r.Enabled = parseBool(val)
	}
	if val, ok := post["meta"]; ok {
		r.hasMeta = true
		r.rawMeta = val

		if r.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewNamespaceCreate()

// NamespaceRead request parameters
type NamespaceRead struct {
	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewNamespaceRead request
func NewNamespaceRead() *NamespaceRead {
	return &NamespaceRead{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *NamespaceRead) Fill(req *http.Request) (err error) {
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

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewNamespaceRead()

// NamespaceUpdate request parameters
type NamespaceUpdate struct {
	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasSlug bool
	rawSlug string
	Slug    string

	hasEnabled bool
	rawEnabled string
	Enabled    bool

	hasMeta bool
	rawMeta string
	Meta    sqlxTypes.JSONText

	hasUpdatedAt bool
	rawUpdatedAt string
	UpdatedAt    *time.Time
}

// NewNamespaceUpdate request
func NewNamespaceUpdate() *NamespaceUpdate {
	return &NamespaceUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID
	out["name"] = r.Name
	out["slug"] = r.Slug
	out["enabled"] = r.Enabled
	out["meta"] = r.Meta
	out["updatedAt"] = r.UpdatedAt

	return out
}

// Fill processes request and fills internal variables
func (r *NamespaceUpdate) Fill(req *http.Request) (err error) {
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

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["name"]; ok {
		r.hasName = true
		r.rawName = val
		r.Name = val
	}
	if val, ok := post["slug"]; ok {
		r.hasSlug = true
		r.rawSlug = val
		r.Slug = val
	}
	if val, ok := post["enabled"]; ok {
		r.hasEnabled = true
		r.rawEnabled = val
		r.Enabled = parseBool(val)
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

var _ RequestFiller = NewNamespaceUpdate()

// NamespaceDelete request parameters
type NamespaceDelete struct {
	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewNamespaceDelete request
func NewNamespaceDelete() *NamespaceDelete {
	return &NamespaceDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *NamespaceDelete) Fill(req *http.Request) (err error) {
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

	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewNamespaceDelete()

// NamespaceTriggerScript request parameters
type NamespaceTriggerScript struct {
	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string
}

// NewNamespaceTriggerScript request
func NewNamespaceTriggerScript() *NamespaceTriggerScript {
	return &NamespaceTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r NamespaceTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["namespaceID"] = r.NamespaceID
	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
func (r *NamespaceTriggerScript) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewNamespaceTriggerScript()

// HasQuery returns true if query was set
func (r *NamespaceList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *NamespaceList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *NamespaceList) GetQuery() string {
	return r.Query
}

// HasSlug returns true if slug was set
func (r *NamespaceList) HasSlug() bool {
	return r.hasSlug
}

// RawSlug returns raw value of slug parameter
func (r *NamespaceList) RawSlug() string {
	return r.rawSlug
}

// GetSlug returns casted value of  slug parameter
func (r *NamespaceList) GetSlug() string {
	return r.Slug
}

// HasLimit returns true if limit was set
func (r *NamespaceList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *NamespaceList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *NamespaceList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *NamespaceList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *NamespaceList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *NamespaceList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *NamespaceList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *NamespaceList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *NamespaceList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *NamespaceList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *NamespaceList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *NamespaceList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *NamespaceList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *NamespaceList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *NamespaceList) GetSort() string {
	return r.Sort
}

// HasName returns true if name was set
func (r *NamespaceCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *NamespaceCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *NamespaceCreate) GetName() string {
	return r.Name
}

// HasSlug returns true if slug was set
func (r *NamespaceCreate) HasSlug() bool {
	return r.hasSlug
}

// RawSlug returns raw value of slug parameter
func (r *NamespaceCreate) RawSlug() string {
	return r.rawSlug
}

// GetSlug returns casted value of  slug parameter
func (r *NamespaceCreate) GetSlug() string {
	return r.Slug
}

// HasEnabled returns true if enabled was set
func (r *NamespaceCreate) HasEnabled() bool {
	return r.hasEnabled
}

// RawEnabled returns raw value of enabled parameter
func (r *NamespaceCreate) RawEnabled() string {
	return r.rawEnabled
}

// GetEnabled returns casted value of  enabled parameter
func (r *NamespaceCreate) GetEnabled() bool {
	return r.Enabled
}

// HasMeta returns true if meta was set
func (r *NamespaceCreate) HasMeta() bool {
	return r.hasMeta
}

// RawMeta returns raw value of meta parameter
func (r *NamespaceCreate) RawMeta() string {
	return r.rawMeta
}

// GetMeta returns casted value of  meta parameter
func (r *NamespaceCreate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// HasNamespaceID returns true if namespaceID was set
func (r *NamespaceRead) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *NamespaceRead) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *NamespaceRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasNamespaceID returns true if namespaceID was set
func (r *NamespaceUpdate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *NamespaceUpdate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *NamespaceUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasName returns true if name was set
func (r *NamespaceUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *NamespaceUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *NamespaceUpdate) GetName() string {
	return r.Name
}

// HasSlug returns true if slug was set
func (r *NamespaceUpdate) HasSlug() bool {
	return r.hasSlug
}

// RawSlug returns raw value of slug parameter
func (r *NamespaceUpdate) RawSlug() string {
	return r.rawSlug
}

// GetSlug returns casted value of  slug parameter
func (r *NamespaceUpdate) GetSlug() string {
	return r.Slug
}

// HasEnabled returns true if enabled was set
func (r *NamespaceUpdate) HasEnabled() bool {
	return r.hasEnabled
}

// RawEnabled returns raw value of enabled parameter
func (r *NamespaceUpdate) RawEnabled() string {
	return r.rawEnabled
}

// GetEnabled returns casted value of  enabled parameter
func (r *NamespaceUpdate) GetEnabled() bool {
	return r.Enabled
}

// HasMeta returns true if meta was set
func (r *NamespaceUpdate) HasMeta() bool {
	return r.hasMeta
}

// RawMeta returns raw value of meta parameter
func (r *NamespaceUpdate) RawMeta() string {
	return r.rawMeta
}

// GetMeta returns casted value of  meta parameter
func (r *NamespaceUpdate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// HasUpdatedAt returns true if updatedAt was set
func (r *NamespaceUpdate) HasUpdatedAt() bool {
	return r.hasUpdatedAt
}

// RawUpdatedAt returns raw value of updatedAt parameter
func (r *NamespaceUpdate) RawUpdatedAt() string {
	return r.rawUpdatedAt
}

// GetUpdatedAt returns casted value of  updatedAt parameter
func (r *NamespaceUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// HasNamespaceID returns true if namespaceID was set
func (r *NamespaceDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *NamespaceDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *NamespaceDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasNamespaceID returns true if namespaceID was set
func (r *NamespaceTriggerScript) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *NamespaceTriggerScript) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *NamespaceTriggerScript) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasScript returns true if script was set
func (r *NamespaceTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *NamespaceTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *NamespaceTriggerScript) GetScript() string {
	return r.Script
}
