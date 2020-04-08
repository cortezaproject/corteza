package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `role.go`, `role.util.go` or `role_test.go` to
	implement your API calls, helper functions and tests. The file `role.go`
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
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// RoleList request parameters
type RoleList struct {
	hasQuery bool
	rawQuery string
	Query    string

	hasDeleted bool
	rawDeleted string
	Deleted    uint

	hasArchived bool
	rawArchived string
	Archived    uint

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

// NewRoleList request
func NewRoleList() *RoleList {
	return &RoleList{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["deleted"] = r.Deleted
	out["archived"] = r.Archived
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort

	return out
}

// Fill processes request and fills internal variables
func (r *RoleList) Fill(req *http.Request) (err error) {
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
	if val, ok := get["deleted"]; ok {
		r.hasDeleted = true
		r.rawDeleted = val
		r.Deleted = parseUint(val)
	}
	if val, ok := get["archived"]; ok {
		r.hasArchived = true
		r.rawArchived = val
		r.Archived = parseUint(val)
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

var _ RequestFiller = NewRoleList()

// RoleCreate request parameters
type RoleCreate struct {
	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasMembers bool
	rawMembers []string
	Members    []string
}

// NewRoleCreate request
func NewRoleCreate() *RoleCreate {
	return &RoleCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["handle"] = r.Handle
	out["members"] = r.Members

	return out
}

// Fill processes request and fills internal variables
func (r *RoleCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := req.Form["members"]; ok {
		r.hasMembers = true
		r.rawMembers = val
		r.Members = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewRoleCreate()

// RoleUpdate request parameters
type RoleUpdate struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasMembers bool
	rawMembers []string
	Members    []string
}

// NewRoleUpdate request
func NewRoleUpdate() *RoleUpdate {
	return &RoleUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["members"] = r.Members

	return out
}

// Fill processes request and fills internal variables
func (r *RoleUpdate) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
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

	if val, ok := req.Form["members"]; ok {
		r.hasMembers = true
		r.rawMembers = val
		r.Members = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewRoleUpdate()

// RoleRead request parameters
type RoleRead struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewRoleRead request
func NewRoleRead() *RoleRead {
	return &RoleRead{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleRead) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleRead()

// RoleDelete request parameters
type RoleDelete struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewRoleDelete request
func NewRoleDelete() *RoleDelete {
	return &RoleDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleDelete) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleDelete()

// RoleArchive request parameters
type RoleArchive struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewRoleArchive request
func NewRoleArchive() *RoleArchive {
	return &RoleArchive{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleArchive) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleArchive) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleArchive()

// RoleUnarchive request parameters
type RoleUnarchive struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewRoleUnarchive request
func NewRoleUnarchive() *RoleUnarchive {
	return &RoleUnarchive{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUnarchive) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleUnarchive) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleUnarchive()

// RoleUndelete request parameters
type RoleUndelete struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewRoleUndelete request
func NewRoleUndelete() *RoleUndelete {
	return &RoleUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleUndelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleUndelete) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleUndelete()

// RoleMove request parameters
type RoleMove struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasOrganisationID bool
	rawOrganisationID string
	OrganisationID    uint64 `json:",string"`
}

// NewRoleMove request
func NewRoleMove() *RoleMove {
	return &RoleMove{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["organisationID"] = r.OrganisationID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleMove) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["organisationID"]; ok {
		r.hasOrganisationID = true
		r.rawOrganisationID = val
		r.OrganisationID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewRoleMove()

// RoleMerge request parameters
type RoleMerge struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasDestination bool
	rawDestination string
	Destination    uint64 `json:",string"`
}

// NewRoleMerge request
func NewRoleMerge() *RoleMerge {
	return &RoleMerge{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMerge) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["destination"] = r.Destination

	return out
}

// Fill processes request and fills internal variables
func (r *RoleMerge) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["destination"]; ok {
		r.hasDestination = true
		r.rawDestination = val
		r.Destination = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewRoleMerge()

// RoleMemberList request parameters
type RoleMemberList struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`
}

// NewRoleMemberList request
func NewRoleMemberList() *RoleMemberList {
	return &RoleMemberList{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleMemberList) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleMemberList()

// RoleMemberAdd request parameters
type RoleMemberAdd struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewRoleMemberAdd request
func NewRoleMemberAdd() *RoleMemberAdd {
	return &RoleMemberAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberAdd) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleMemberAdd) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewRoleMemberAdd()

// RoleMemberRemove request parameters
type RoleMemberRemove struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewRoleMemberRemove request
func NewRoleMemberRemove() *RoleMemberRemove {
	return &RoleMemberRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleMemberRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *RoleMemberRemove) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewRoleMemberRemove()

// RoleTriggerScript request parameters
type RoleTriggerScript struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string
}

// NewRoleTriggerScript request
func NewRoleTriggerScript() *RoleTriggerScript {
	return &RoleTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r RoleTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
func (r *RoleTriggerScript) Fill(req *http.Request) (err error) {
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

	r.hasRoleID = true
	r.rawRoleID = chi.URLParam(req, "roleID")
	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewRoleTriggerScript()

// HasQuery returns true if query was set
func (r *RoleList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *RoleList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *RoleList) GetQuery() string {
	return r.Query
}

// HasDeleted returns true if deleted was set
func (r *RoleList) HasDeleted() bool {
	return r.hasDeleted
}

// RawDeleted returns raw value of deleted parameter
func (r *RoleList) RawDeleted() string {
	return r.rawDeleted
}

// GetDeleted returns casted value of  deleted parameter
func (r *RoleList) GetDeleted() uint {
	return r.Deleted
}

// HasArchived returns true if archived was set
func (r *RoleList) HasArchived() bool {
	return r.hasArchived
}

// RawArchived returns raw value of archived parameter
func (r *RoleList) RawArchived() string {
	return r.rawArchived
}

// GetArchived returns casted value of  archived parameter
func (r *RoleList) GetArchived() uint {
	return r.Archived
}

// HasLimit returns true if limit was set
func (r *RoleList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *RoleList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *RoleList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *RoleList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *RoleList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *RoleList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *RoleList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *RoleList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *RoleList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *RoleList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *RoleList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *RoleList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *RoleList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *RoleList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *RoleList) GetSort() string {
	return r.Sort
}

// HasName returns true if name was set
func (r *RoleCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *RoleCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *RoleCreate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *RoleCreate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *RoleCreate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *RoleCreate) GetHandle() string {
	return r.Handle
}

// HasMembers returns true if members was set
func (r *RoleCreate) HasMembers() bool {
	return r.hasMembers
}

// RawMembers returns raw value of members parameter
func (r *RoleCreate) RawMembers() []string {
	return r.rawMembers
}

// GetMembers returns casted value of  members parameter
func (r *RoleCreate) GetMembers() []string {
	return r.Members
}

// HasRoleID returns true if roleID was set
func (r *RoleUpdate) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleUpdate) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleUpdate) GetRoleID() uint64 {
	return r.RoleID
}

// HasName returns true if name was set
func (r *RoleUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *RoleUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *RoleUpdate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *RoleUpdate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *RoleUpdate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *RoleUpdate) GetHandle() string {
	return r.Handle
}

// HasMembers returns true if members was set
func (r *RoleUpdate) HasMembers() bool {
	return r.hasMembers
}

// RawMembers returns raw value of members parameter
func (r *RoleUpdate) RawMembers() []string {
	return r.rawMembers
}

// GetMembers returns casted value of  members parameter
func (r *RoleUpdate) GetMembers() []string {
	return r.Members
}

// HasRoleID returns true if roleID was set
func (r *RoleRead) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleRead) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleRead) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *RoleDelete) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleDelete) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleDelete) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *RoleArchive) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleArchive) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleArchive) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *RoleUnarchive) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleUnarchive) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleUnarchive) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *RoleUndelete) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleUndelete) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleUndelete) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *RoleMove) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleMove) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleMove) GetRoleID() uint64 {
	return r.RoleID
}

// HasOrganisationID returns true if organisationID was set
func (r *RoleMove) HasOrganisationID() bool {
	return r.hasOrganisationID
}

// RawOrganisationID returns raw value of organisationID parameter
func (r *RoleMove) RawOrganisationID() string {
	return r.rawOrganisationID
}

// GetOrganisationID returns casted value of  organisationID parameter
func (r *RoleMove) GetOrganisationID() uint64 {
	return r.OrganisationID
}

// HasRoleID returns true if roleID was set
func (r *RoleMerge) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleMerge) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleMerge) GetRoleID() uint64 {
	return r.RoleID
}

// HasDestination returns true if destination was set
func (r *RoleMerge) HasDestination() bool {
	return r.hasDestination
}

// RawDestination returns raw value of destination parameter
func (r *RoleMerge) RawDestination() string {
	return r.rawDestination
}

// GetDestination returns casted value of  destination parameter
func (r *RoleMerge) GetDestination() uint64 {
	return r.Destination
}

// HasRoleID returns true if roleID was set
func (r *RoleMemberList) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleMemberList) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleMemberList) GetRoleID() uint64 {
	return r.RoleID
}

// HasRoleID returns true if roleID was set
func (r *RoleMemberAdd) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleMemberAdd) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleMemberAdd) GetRoleID() uint64 {
	return r.RoleID
}

// HasUserID returns true if userID was set
func (r *RoleMemberAdd) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *RoleMemberAdd) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *RoleMemberAdd) GetUserID() uint64 {
	return r.UserID
}

// HasRoleID returns true if roleID was set
func (r *RoleMemberRemove) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleMemberRemove) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleMemberRemove) GetRoleID() uint64 {
	return r.RoleID
}

// HasUserID returns true if userID was set
func (r *RoleMemberRemove) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *RoleMemberRemove) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *RoleMemberRemove) GetUserID() uint64 {
	return r.UserID
}

// HasRoleID returns true if roleID was set
func (r *RoleTriggerScript) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *RoleTriggerScript) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *RoleTriggerScript) GetRoleID() uint64 {
	return r.RoleID
}

// HasScript returns true if script was set
func (r *RoleTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *RoleTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *RoleTriggerScript) GetScript() string {
	return r.Script
}
