package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `chart.go`, `chart.util.go` or `chart_test.go` to
	implement your API calls, helper functions and tests. The file `chart.go`
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

// ChartList request parameters
type ChartList struct {
	hasQuery bool
	rawQuery string
	Query    string

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

// NewChartList request
func NewChartList() *ChartList {
	return &ChartList{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
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
func (r *ChartList) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewChartList()

// ChartCreate request parameters
type ChartCreate struct {
	hasConfig bool
	rawConfig string
	Config    sqlxTypes.JSONText

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewChartCreate request
func NewChartCreate() *ChartCreate {
	return &ChartCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["config"] = r.Config
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ChartCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["config"]; ok {
		r.hasConfig = true
		r.rawConfig = val

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
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
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartCreate()

// ChartRead request parameters
type ChartRead struct {
	hasChartID bool
	rawChartID string
	ChartID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewChartRead request
func NewChartRead() *ChartRead {
	return &ChartRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["chartID"] = r.ChartID
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ChartRead) Fill(req *http.Request) (err error) {
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

	r.hasChartID = true
	r.rawChartID = chi.URLParam(req, "chartID")
	r.ChartID = parseUInt64(chi.URLParam(req, "chartID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartRead()

// ChartUpdate request parameters
type ChartUpdate struct {
	hasChartID bool
	rawChartID string
	ChartID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`

	hasConfig bool
	rawConfig string
	Config    sqlxTypes.JSONText

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasUpdatedAt bool
	rawUpdatedAt string
	UpdatedAt    *time.Time
}

// NewChartUpdate request
func NewChartUpdate() *ChartUpdate {
	return &ChartUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["chartID"] = r.ChartID
	out["namespaceID"] = r.NamespaceID
	out["config"] = r.Config
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["updatedAt"] = r.UpdatedAt

	return out
}

// Fill processes request and fills internal variables
func (r *ChartUpdate) Fill(req *http.Request) (err error) {
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

	r.hasChartID = true
	r.rawChartID = chi.URLParam(req, "chartID")
	r.ChartID = parseUInt64(chi.URLParam(req, "chartID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["config"]; ok {
		r.hasConfig = true
		r.rawConfig = val

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
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
	if val, ok := post["updatedAt"]; ok {
		r.hasUpdatedAt = true
		r.rawUpdatedAt = val

		if r.UpdatedAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewChartUpdate()

// ChartDelete request parameters
type ChartDelete struct {
	hasChartID bool
	rawChartID string
	ChartID    uint64 `json:",string"`

	hasNamespaceID bool
	rawNamespaceID string
	NamespaceID    uint64 `json:",string"`
}

// NewChartDelete request
func NewChartDelete() *ChartDelete {
	return &ChartDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ChartDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["chartID"] = r.ChartID
	out["namespaceID"] = r.NamespaceID

	return out
}

// Fill processes request and fills internal variables
func (r *ChartDelete) Fill(req *http.Request) (err error) {
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

	r.hasChartID = true
	r.rawChartID = chi.URLParam(req, "chartID")
	r.ChartID = parseUInt64(chi.URLParam(req, "chartID"))
	r.hasNamespaceID = true
	r.rawNamespaceID = chi.URLParam(req, "namespaceID")
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartDelete()

// HasQuery returns true if query was set
func (r *ChartList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *ChartList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *ChartList) GetQuery() string {
	return r.Query
}

// HasHandle returns true if handle was set
func (r *ChartList) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *ChartList) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *ChartList) GetHandle() string {
	return r.Handle
}

// HasLimit returns true if limit was set
func (r *ChartList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *ChartList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *ChartList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *ChartList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *ChartList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *ChartList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *ChartList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *ChartList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *ChartList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *ChartList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *ChartList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *ChartList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *ChartList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *ChartList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *ChartList) GetSort() string {
	return r.Sort
}

// HasNamespaceID returns true if namespaceID was set
func (r *ChartList) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ChartList) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ChartList) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasConfig returns true if config was set
func (r *ChartCreate) HasConfig() bool {
	return r.hasConfig
}

// RawConfig returns raw value of config parameter
func (r *ChartCreate) RawConfig() string {
	return r.rawConfig
}

// GetConfig returns casted value of  config parameter
func (r *ChartCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// HasName returns true if name was set
func (r *ChartCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ChartCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ChartCreate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *ChartCreate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *ChartCreate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *ChartCreate) GetHandle() string {
	return r.Handle
}

// HasNamespaceID returns true if namespaceID was set
func (r *ChartCreate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ChartCreate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ChartCreate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasChartID returns true if chartID was set
func (r *ChartRead) HasChartID() bool {
	return r.hasChartID
}

// RawChartID returns raw value of chartID parameter
func (r *ChartRead) RawChartID() string {
	return r.rawChartID
}

// GetChartID returns casted value of  chartID parameter
func (r *ChartRead) GetChartID() uint64 {
	return r.ChartID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ChartRead) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ChartRead) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ChartRead) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasChartID returns true if chartID was set
func (r *ChartUpdate) HasChartID() bool {
	return r.hasChartID
}

// RawChartID returns raw value of chartID parameter
func (r *ChartUpdate) RawChartID() string {
	return r.rawChartID
}

// GetChartID returns casted value of  chartID parameter
func (r *ChartUpdate) GetChartID() uint64 {
	return r.ChartID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ChartUpdate) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ChartUpdate) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ChartUpdate) GetNamespaceID() uint64 {
	return r.NamespaceID
}

// HasConfig returns true if config was set
func (r *ChartUpdate) HasConfig() bool {
	return r.hasConfig
}

// RawConfig returns raw value of config parameter
func (r *ChartUpdate) RawConfig() string {
	return r.rawConfig
}

// GetConfig returns casted value of  config parameter
func (r *ChartUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// HasName returns true if name was set
func (r *ChartUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *ChartUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *ChartUpdate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *ChartUpdate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *ChartUpdate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *ChartUpdate) GetHandle() string {
	return r.Handle
}

// HasUpdatedAt returns true if updatedAt was set
func (r *ChartUpdate) HasUpdatedAt() bool {
	return r.hasUpdatedAt
}

// RawUpdatedAt returns raw value of updatedAt parameter
func (r *ChartUpdate) RawUpdatedAt() string {
	return r.rawUpdatedAt
}

// GetUpdatedAt returns casted value of  updatedAt parameter
func (r *ChartUpdate) GetUpdatedAt() *time.Time {
	return r.UpdatedAt
}

// HasChartID returns true if chartID was set
func (r *ChartDelete) HasChartID() bool {
	return r.hasChartID
}

// RawChartID returns raw value of chartID parameter
func (r *ChartDelete) RawChartID() string {
	return r.rawChartID
}

// GetChartID returns casted value of  chartID parameter
func (r *ChartDelete) GetChartID() uint64 {
	return r.ChartID
}

// HasNamespaceID returns true if namespaceID was set
func (r *ChartDelete) HasNamespaceID() bool {
	return r.hasNamespaceID
}

// RawNamespaceID returns raw value of namespaceID parameter
func (r *ChartDelete) RawNamespaceID() string {
	return r.rawNamespaceID
}

// GetNamespaceID returns casted value of  namespaceID parameter
func (r *ChartDelete) GetNamespaceID() uint64 {
	return r.NamespaceID
}
