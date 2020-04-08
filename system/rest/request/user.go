package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `user.go`, `user.util.go` or `user_test.go` to
	implement your API calls, helper functions and tests. The file `user.go`
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

	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// UserList request parameters
type UserList struct {
	hasUserID bool
	rawUserID []string
	UserID    []string

	hasRoleID bool
	rawRoleID []string
	RoleID    []string

	hasQuery bool
	rawQuery string
	Query    string

	hasUsername bool
	rawUsername string
	Username    string

	hasEmail bool
	rawEmail string
	Email    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasKind bool
	rawKind string
	Kind    types.UserKind

	hasIncDeleted bool
	rawIncDeleted string
	IncDeleted    bool

	hasIncSuspended bool
	rawIncSuspended string
	IncSuspended    bool

	hasDeleted bool
	rawDeleted string
	Deleted    uint

	hasSuspended bool
	rawSuspended string
	Suspended    uint

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

// NewUserList request
func NewUserList() *UserList {
	return &UserList{}
}

// Auditable returns all auditable/loggable parameters
func (r UserList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID
	out["roleID"] = r.RoleID
	out["query"] = r.Query
	out["username"] = r.Username
	out["email"] = r.Email
	out["handle"] = r.Handle
	out["kind"] = r.Kind
	out["incDeleted"] = r.IncDeleted
	out["incSuspended"] = r.IncSuspended
	out["deleted"] = r.Deleted
	out["suspended"] = r.Suspended
	out["limit"] = r.Limit
	out["offset"] = r.Offset
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["sort"] = r.Sort

	return out
}

// Fill processes request and fills internal variables
func (r *UserList) Fill(req *http.Request) (err error) {
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

	if val, ok := urlQuery["userID[]"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseStrings(val)
	} else if val, ok = urlQuery["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseStrings(val)
	}

	if val, ok := urlQuery["roleID[]"]; ok {
		r.hasRoleID = true
		r.rawRoleID = val
		r.RoleID = parseStrings(val)
	} else if val, ok = urlQuery["roleID"]; ok {
		r.hasRoleID = true
		r.rawRoleID = val
		r.RoleID = parseStrings(val)
	}

	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}
	if val, ok := get["username"]; ok {
		r.hasUsername = true
		r.rawUsername = val
		r.Username = val
	}
	if val, ok := get["email"]; ok {
		r.hasEmail = true
		r.rawEmail = val
		r.Email = val
	}
	if val, ok := get["handle"]; ok {
		r.hasHandle = true
		r.rawHandle = val
		r.Handle = val
	}
	if val, ok := get["kind"]; ok {
		r.hasKind = true
		r.rawKind = val
		r.Kind = types.UserKind(val)
	}
	if val, ok := get["incDeleted"]; ok {
		r.hasIncDeleted = true
		r.rawIncDeleted = val
		r.IncDeleted = parseBool(val)
	}
	if val, ok := get["incSuspended"]; ok {
		r.hasIncSuspended = true
		r.rawIncSuspended = val
		r.IncSuspended = parseBool(val)
	}
	if val, ok := get["deleted"]; ok {
		r.hasDeleted = true
		r.rawDeleted = val
		r.Deleted = parseUint(val)
	}
	if val, ok := get["suspended"]; ok {
		r.hasSuspended = true
		r.rawSuspended = val
		r.Suspended = parseUint(val)
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

var _ RequestFiller = NewUserList()

// UserCreate request parameters
type UserCreate struct {
	hasEmail bool
	rawEmail string
	Email    string

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasKind bool
	rawKind string
	Kind    types.UserKind
}

// NewUserCreate request
func NewUserCreate() *UserCreate {
	return &UserCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["email"] = r.Email
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["kind"] = r.Kind

	return out
}

// Fill processes request and fills internal variables
func (r *UserCreate) Fill(req *http.Request) (err error) {
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

	if val, ok := post["email"]; ok {
		r.hasEmail = true
		r.rawEmail = val
		r.Email = val
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
	if val, ok := post["kind"]; ok {
		r.hasKind = true
		r.rawKind = val
		r.Kind = types.UserKind(val)
	}

	return err
}

var _ RequestFiller = NewUserCreate()

// UserUpdate request parameters
type UserUpdate struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasEmail bool
	rawEmail string
	Email    string

	hasName bool
	rawName string
	Name    string

	hasHandle bool
	rawHandle string
	Handle    string

	hasKind bool
	rawKind string
	Kind    types.UserKind
}

// NewUserUpdate request
func NewUserUpdate() *UserUpdate {
	return &UserUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r UserUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID
	out["email"] = r.Email
	out["name"] = r.Name
	out["handle"] = r.Handle
	out["kind"] = r.Kind

	return out
}

// Fill processes request and fills internal variables
func (r *UserUpdate) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))
	if val, ok := post["email"]; ok {
		r.hasEmail = true
		r.rawEmail = val
		r.Email = val
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
	if val, ok := post["kind"]; ok {
		r.hasKind = true
		r.rawKind = val
		r.Kind = types.UserKind(val)
	}

	return err
}

var _ RequestFiller = NewUserUpdate()

// UserRead request parameters
type UserRead struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserRead request
func NewUserRead() *UserRead {
	return &UserRead{}
}

// Auditable returns all auditable/loggable parameters
func (r UserRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserRead) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserRead()

// UserDelete request parameters
type UserDelete struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserDelete request
func NewUserDelete() *UserDelete {
	return &UserDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r UserDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserDelete) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserDelete()

// UserSuspend request parameters
type UserSuspend struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserSuspend request
func NewUserSuspend() *UserSuspend {
	return &UserSuspend{}
}

// Auditable returns all auditable/loggable parameters
func (r UserSuspend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserSuspend) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserSuspend()

// UserUnsuspend request parameters
type UserUnsuspend struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserUnsuspend request
func NewUserUnsuspend() *UserUnsuspend {
	return &UserUnsuspend{}
}

// Auditable returns all auditable/loggable parameters
func (r UserUnsuspend) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserUnsuspend) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserUnsuspend()

// UserUndelete request parameters
type UserUndelete struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserUndelete request
func NewUserUndelete() *UserUndelete {
	return &UserUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r UserUndelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserUndelete) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserUndelete()

// UserSetPassword request parameters
type UserSetPassword struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasPassword bool
	rawPassword string
	Password    string
}

// NewUserSetPassword request
func NewUserSetPassword() *UserSetPassword {
	return &UserSetPassword{}
}

// Auditable returns all auditable/loggable parameters
func (r UserSetPassword) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID
	out["password"] = "*masked*sensitive*data*"

	return out
}

// Fill processes request and fills internal variables
func (r *UserSetPassword) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))
	if val, ok := post["password"]; ok {
		r.hasPassword = true
		r.rawPassword = val
		r.Password = val
	}

	return err
}

var _ RequestFiller = NewUserSetPassword()

// UserMembershipList request parameters
type UserMembershipList struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserMembershipList request
func NewUserMembershipList() *UserMembershipList {
	return &UserMembershipList{}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserMembershipList) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))

	return err
}

var _ RequestFiller = NewUserMembershipList()

// UserMembershipAdd request parameters
type UserMembershipAdd struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserMembershipAdd request
func NewUserMembershipAdd() *UserMembershipAdd {
	return &UserMembershipAdd{}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipAdd) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserMembershipAdd) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewUserMembershipAdd()

// UserMembershipRemove request parameters
type UserMembershipRemove struct {
	hasRoleID bool
	rawRoleID string
	RoleID    uint64 `json:",string"`

	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`
}

// NewUserMembershipRemove request
func NewUserMembershipRemove() *UserMembershipRemove {
	return &UserMembershipRemove{}
}

// Auditable returns all auditable/loggable parameters
func (r UserMembershipRemove) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["roleID"] = r.RoleID
	out["userID"] = r.UserID

	return out
}

// Fill processes request and fills internal variables
func (r *UserMembershipRemove) Fill(req *http.Request) (err error) {
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

var _ RequestFiller = NewUserMembershipRemove()

// UserTriggerScript request parameters
type UserTriggerScript struct {
	hasUserID bool
	rawUserID string
	UserID    uint64 `json:",string"`

	hasScript bool
	rawScript string
	Script    string
}

// NewUserTriggerScript request
func NewUserTriggerScript() *UserTriggerScript {
	return &UserTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r UserTriggerScript) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["userID"] = r.UserID
	out["script"] = r.Script

	return out
}

// Fill processes request and fills internal variables
func (r *UserTriggerScript) Fill(req *http.Request) (err error) {
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

	r.hasUserID = true
	r.rawUserID = chi.URLParam(req, "userID")
	r.UserID = parseUInt64(chi.URLParam(req, "userID"))
	if val, ok := post["script"]; ok {
		r.hasScript = true
		r.rawScript = val
		r.Script = val
	}

	return err
}

var _ RequestFiller = NewUserTriggerScript()

// HasUserID returns true if userID was set
func (r *UserList) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserList) RawUserID() []string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserList) GetUserID() []string {
	return r.UserID
}

// HasRoleID returns true if roleID was set
func (r *UserList) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *UserList) RawRoleID() []string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *UserList) GetRoleID() []string {
	return r.RoleID
}

// HasQuery returns true if query was set
func (r *UserList) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *UserList) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *UserList) GetQuery() string {
	return r.Query
}

// HasUsername returns true if username was set
func (r *UserList) HasUsername() bool {
	return r.hasUsername
}

// RawUsername returns raw value of username parameter
func (r *UserList) RawUsername() string {
	return r.rawUsername
}

// GetUsername returns casted value of  username parameter
func (r *UserList) GetUsername() string {
	return r.Username
}

// HasEmail returns true if email was set
func (r *UserList) HasEmail() bool {
	return r.hasEmail
}

// RawEmail returns raw value of email parameter
func (r *UserList) RawEmail() string {
	return r.rawEmail
}

// GetEmail returns casted value of  email parameter
func (r *UserList) GetEmail() string {
	return r.Email
}

// HasHandle returns true if handle was set
func (r *UserList) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *UserList) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *UserList) GetHandle() string {
	return r.Handle
}

// HasKind returns true if kind was set
func (r *UserList) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *UserList) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *UserList) GetKind() types.UserKind {
	return r.Kind
}

// HasIncDeleted returns true if incDeleted was set
func (r *UserList) HasIncDeleted() bool {
	return r.hasIncDeleted
}

// RawIncDeleted returns raw value of incDeleted parameter
func (r *UserList) RawIncDeleted() string {
	return r.rawIncDeleted
}

// GetIncDeleted returns casted value of  incDeleted parameter
func (r *UserList) GetIncDeleted() bool {
	return r.IncDeleted
}

// HasIncSuspended returns true if incSuspended was set
func (r *UserList) HasIncSuspended() bool {
	return r.hasIncSuspended
}

// RawIncSuspended returns raw value of incSuspended parameter
func (r *UserList) RawIncSuspended() string {
	return r.rawIncSuspended
}

// GetIncSuspended returns casted value of  incSuspended parameter
func (r *UserList) GetIncSuspended() bool {
	return r.IncSuspended
}

// HasDeleted returns true if deleted was set
func (r *UserList) HasDeleted() bool {
	return r.hasDeleted
}

// RawDeleted returns raw value of deleted parameter
func (r *UserList) RawDeleted() string {
	return r.rawDeleted
}

// GetDeleted returns casted value of  deleted parameter
func (r *UserList) GetDeleted() uint {
	return r.Deleted
}

// HasSuspended returns true if suspended was set
func (r *UserList) HasSuspended() bool {
	return r.hasSuspended
}

// RawSuspended returns raw value of suspended parameter
func (r *UserList) RawSuspended() string {
	return r.rawSuspended
}

// GetSuspended returns casted value of  suspended parameter
func (r *UserList) GetSuspended() uint {
	return r.Suspended
}

// HasLimit returns true if limit was set
func (r *UserList) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *UserList) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *UserList) GetLimit() uint {
	return r.Limit
}

// HasOffset returns true if offset was set
func (r *UserList) HasOffset() bool {
	return r.hasOffset
}

// RawOffset returns raw value of offset parameter
func (r *UserList) RawOffset() string {
	return r.rawOffset
}

// GetOffset returns casted value of  offset parameter
func (r *UserList) GetOffset() uint {
	return r.Offset
}

// HasPage returns true if page was set
func (r *UserList) HasPage() bool {
	return r.hasPage
}

// RawPage returns raw value of page parameter
func (r *UserList) RawPage() string {
	return r.rawPage
}

// GetPage returns casted value of  page parameter
func (r *UserList) GetPage() uint {
	return r.Page
}

// HasPerPage returns true if perPage was set
func (r *UserList) HasPerPage() bool {
	return r.hasPerPage
}

// RawPerPage returns raw value of perPage parameter
func (r *UserList) RawPerPage() string {
	return r.rawPerPage
}

// GetPerPage returns casted value of  perPage parameter
func (r *UserList) GetPerPage() uint {
	return r.PerPage
}

// HasSort returns true if sort was set
func (r *UserList) HasSort() bool {
	return r.hasSort
}

// RawSort returns raw value of sort parameter
func (r *UserList) RawSort() string {
	return r.rawSort
}

// GetSort returns casted value of  sort parameter
func (r *UserList) GetSort() string {
	return r.Sort
}

// HasEmail returns true if email was set
func (r *UserCreate) HasEmail() bool {
	return r.hasEmail
}

// RawEmail returns raw value of email parameter
func (r *UserCreate) RawEmail() string {
	return r.rawEmail
}

// GetEmail returns casted value of  email parameter
func (r *UserCreate) GetEmail() string {
	return r.Email
}

// HasName returns true if name was set
func (r *UserCreate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *UserCreate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *UserCreate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *UserCreate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *UserCreate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *UserCreate) GetHandle() string {
	return r.Handle
}

// HasKind returns true if kind was set
func (r *UserCreate) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *UserCreate) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *UserCreate) GetKind() types.UserKind {
	return r.Kind
}

// HasUserID returns true if userID was set
func (r *UserUpdate) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserUpdate) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserUpdate) GetUserID() uint64 {
	return r.UserID
}

// HasEmail returns true if email was set
func (r *UserUpdate) HasEmail() bool {
	return r.hasEmail
}

// RawEmail returns raw value of email parameter
func (r *UserUpdate) RawEmail() string {
	return r.rawEmail
}

// GetEmail returns casted value of  email parameter
func (r *UserUpdate) GetEmail() string {
	return r.Email
}

// HasName returns true if name was set
func (r *UserUpdate) HasName() bool {
	return r.hasName
}

// RawName returns raw value of name parameter
func (r *UserUpdate) RawName() string {
	return r.rawName
}

// GetName returns casted value of  name parameter
func (r *UserUpdate) GetName() string {
	return r.Name
}

// HasHandle returns true if handle was set
func (r *UserUpdate) HasHandle() bool {
	return r.hasHandle
}

// RawHandle returns raw value of handle parameter
func (r *UserUpdate) RawHandle() string {
	return r.rawHandle
}

// GetHandle returns casted value of  handle parameter
func (r *UserUpdate) GetHandle() string {
	return r.Handle
}

// HasKind returns true if kind was set
func (r *UserUpdate) HasKind() bool {
	return r.hasKind
}

// RawKind returns raw value of kind parameter
func (r *UserUpdate) RawKind() string {
	return r.rawKind
}

// GetKind returns casted value of  kind parameter
func (r *UserUpdate) GetKind() types.UserKind {
	return r.Kind
}

// HasUserID returns true if userID was set
func (r *UserRead) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserRead) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserRead) GetUserID() uint64 {
	return r.UserID
}

// HasUserID returns true if userID was set
func (r *UserDelete) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserDelete) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserDelete) GetUserID() uint64 {
	return r.UserID
}

// HasUserID returns true if userID was set
func (r *UserSuspend) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserSuspend) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserSuspend) GetUserID() uint64 {
	return r.UserID
}

// HasUserID returns true if userID was set
func (r *UserUnsuspend) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserUnsuspend) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserUnsuspend) GetUserID() uint64 {
	return r.UserID
}

// HasUserID returns true if userID was set
func (r *UserUndelete) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserUndelete) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserUndelete) GetUserID() uint64 {
	return r.UserID
}

// HasUserID returns true if userID was set
func (r *UserSetPassword) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserSetPassword) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserSetPassword) GetUserID() uint64 {
	return r.UserID
}

// HasPassword returns true if password was set
func (r *UserSetPassword) HasPassword() bool {
	return r.hasPassword
}

// RawPassword returns raw value of password parameter
func (r *UserSetPassword) RawPassword() string {
	return r.rawPassword
}

// GetPassword returns casted value of  password parameter
func (r *UserSetPassword) GetPassword() string {
	return r.Password
}

// HasUserID returns true if userID was set
func (r *UserMembershipList) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserMembershipList) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserMembershipList) GetUserID() uint64 {
	return r.UserID
}

// HasRoleID returns true if roleID was set
func (r *UserMembershipAdd) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *UserMembershipAdd) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *UserMembershipAdd) GetRoleID() uint64 {
	return r.RoleID
}

// HasUserID returns true if userID was set
func (r *UserMembershipAdd) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserMembershipAdd) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserMembershipAdd) GetUserID() uint64 {
	return r.UserID
}

// HasRoleID returns true if roleID was set
func (r *UserMembershipRemove) HasRoleID() bool {
	return r.hasRoleID
}

// RawRoleID returns raw value of roleID parameter
func (r *UserMembershipRemove) RawRoleID() string {
	return r.rawRoleID
}

// GetRoleID returns casted value of  roleID parameter
func (r *UserMembershipRemove) GetRoleID() uint64 {
	return r.RoleID
}

// HasUserID returns true if userID was set
func (r *UserMembershipRemove) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserMembershipRemove) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserMembershipRemove) GetUserID() uint64 {
	return r.UserID
}

// HasUserID returns true if userID was set
func (r *UserTriggerScript) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *UserTriggerScript) RawUserID() string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *UserTriggerScript) GetUserID() uint64 {
	return r.UserID
}

// HasScript returns true if script was set
func (r *UserTriggerScript) HasScript() bool {
	return r.hasScript
}

// RawScript returns raw value of script parameter
func (r *UserTriggerScript) RawScript() string {
	return r.rawScript
}

// GetScript returns casted value of  script parameter
func (r *UserTriggerScript) GetScript() string {
	return r.Script
}
