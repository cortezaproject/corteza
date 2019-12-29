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

// Role list request parameters
type RoleList struct {
	Query    string
	Deleted  uint
	Archived uint
	Page     uint
	PerPage  uint
	Sort     string
}

func NewRoleList() *RoleList {
	return &RoleList{}
}

func (r RoleList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["deleted"] = r.Deleted
	out["archived"] = r.Archived
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort

	return out
}

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
		r.Query = val
	}
	if val, ok := get["deleted"]; ok {
		r.Deleted = parseUint(val)
	}
	if val, ok := get["archived"]; ok {
		r.Archived = parseUint(val)
	}
	if val, ok := get["page"]; ok {
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.PerPage = parseUint(val)
	}
	if val, ok := get["sort"]; ok {
		r.Sort = val
	}

	return err
}

var _ RequestFiller = NewRoleList()

// Role create request parameters
type RoleCreate struct {
	Name    string
	Handle  string
	Members []string
}

func NewRoleCreate() *RoleCreate {
	return &RoleCreate{}
}

func (r RoleCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["name"] = r.Name
	out["handle"] = r.Handle
	out["members"] = r.Members

	return out
}

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
		r.Name = val
	}
	if val, ok := post["handle"]; ok {
		r.Handle = val
	}

	if val, ok := req.Form["members"]; ok {
		r.Members = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewRoleCreate()

// Role update request parameters
type RoleUpdate struct {
	RoleID  uint64 `json:",string"`
	Name    string
	Handle  string
	Members []string
}

func NewRoleUpdate() *RoleUpdate {
	return &RoleUpdate{}
}

func (r RoleUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["members"] = r.Members

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["handle"]; ok {
		r.Handle = val
	}

	if val, ok := req.Form["members"]; ok {
		r.Members = parseStrings(val)
	}

	return err
}

var _ RequestFiller = NewRoleUpdate()

// Role read request parameters
type RoleRead struct {
	RoleID uint64 `json:",string"`
}

func NewRoleRead() *RoleRead {
	return &RoleRead{}
}

func (r RoleRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleRead()

// Role delete request parameters
type RoleDelete struct {
	RoleID uint64 `json:",string"`
}

func NewRoleDelete() *RoleDelete {
	return &RoleDelete{}
}

func (r RoleDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleDelete()

// Role archive request parameters
type RoleArchive struct {
	RoleID uint64 `json:",string"`
}

func NewRoleArchive() *RoleArchive {
	return &RoleArchive{}
}

func (r RoleArchive) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleArchive()

// Role unarchive request parameters
type RoleUnarchive struct {
	RoleID uint64 `json:",string"`
}

func NewRoleUnarchive() *RoleUnarchive {
	return &RoleUnarchive{}
}

func (r RoleUnarchive) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleUnarchive()

// Role undelete request parameters
type RoleUndelete struct {
	RoleID uint64 `json:",string"`
}

func NewRoleUndelete() *RoleUndelete {
	return &RoleUndelete{}
}

func (r RoleUndelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleUndelete()

// Role move request parameters
type RoleMove struct {
	RoleID         uint64 `json:",string"`
	OrganisationID uint64 `json:",string"`
}

func NewRoleMove() *RoleMove {
	return &RoleMove{}
}

func (r RoleMove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["organisationID"] = r.OrganisationID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["organisationID"]; ok {
		r.OrganisationID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewRoleMove()

// Role merge request parameters
type RoleMerge struct {
	RoleID      uint64 `json:",string"`
	Destination uint64 `json:",string"`
}

func NewRoleMerge() *RoleMerge {
	return &RoleMerge{}
}

func (r RoleMerge) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["destination"] = r.Destination

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["destination"]; ok {
		r.Destination = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewRoleMerge()

// Role memberList request parameters
type RoleMemberList struct {
	RoleID uint64 `json:",string"`
}

func NewRoleMemberList() *RoleMemberList {
	return &RoleMemberList{}
}

func (r RoleMemberList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))

	return err
}

var _ RequestFiller = NewRoleMemberList()

// Role memberAdd request parameters
type RoleMemberAdd struct {
	RoleID uint64 `json:",string"`
	UserID uint64 `json:",string"`
}

func NewRoleMemberAdd() *RoleMemberAdd {
	return &RoleMemberAdd{}
}

func (r RoleMemberAdd) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["userID"] = r.UserID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewRoleMemberAdd()

// Role memberRemove request parameters
type RoleMemberRemove struct {
	RoleID uint64 `json:",string"`
	UserID uint64 `json:",string"`
}

func NewRoleMemberRemove() *RoleMemberRemove {
	return &RoleMemberRemove{}
}

func (r RoleMemberRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["userID"] = r.UserID

	return out
}

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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewRoleMemberRemove()

// Role fireTrigger request parameters
type RoleFireTrigger struct {
	RoleID uint64 `json:",string"`
	Script string
}

func NewRoleFireTrigger() *RoleFireTrigger {
	return &RoleFireTrigger{}
}

func (r RoleFireTrigger) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["script"] = r.Script

	return out
}

func (r *RoleFireTrigger) Fill(req *http.Request) (err error) {
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

	r.RoleID = parseUInt64(chi.URLParam(req, "roleID"))
	if val, ok := post["script"]; ok {
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewRoleFireTrigger()
