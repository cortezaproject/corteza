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
	Query string
}

func NewRoleList() *RoleList {
	return &RoleList{}
}

func (roReq *RoleList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["query"]; ok {

		roReq.Query = val
	}

	return err
}

var _ RequestFiller = NewRoleList()

// Role create request parameters
type RoleCreate struct {
	Name    string
	Members []string
}

func NewRoleCreate() *RoleCreate {
	return &RoleCreate{}
}

func (roReq *RoleCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := post["name"]; ok {

		roReq.Name = val
	}

	return err
}

var _ RequestFiller = NewRoleCreate()

// Role update request parameters
type RoleUpdate struct {
	RoleID  uint64 `json:",string"`
	Name    string
	Members []string
}

func NewRoleUpdate() *RoleUpdate {
	return &RoleUpdate{}
}

func (roReq *RoleUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))
	if val, ok := post["name"]; ok {

		roReq.Name = val
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

func (roReq *RoleRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

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

func (roReq *RoleDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

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

func (roReq *RoleArchive) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

	return err
}

var _ RequestFiller = NewRoleArchive()

// Role move request parameters
type RoleMove struct {
	RoleID         uint64 `json:",string"`
	OrganisationID uint64 `json:",string"`
}

func NewRoleMove() *RoleMove {
	return &RoleMove{}
}

func (roReq *RoleMove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))
	if val, ok := post["organisationID"]; ok {

		roReq.OrganisationID = parseUInt64(val)
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

func (roReq *RoleMerge) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))
	if val, ok := post["destination"]; ok {

		roReq.Destination = parseUInt64(val)
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

func (roReq *RoleMemberList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))

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

func (roReq *RoleMemberAdd) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))
	roReq.UserID = parseUInt64(chi.URLParam(r, "userID"))

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

func (roReq *RoleMemberRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(roReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	roReq.RoleID = parseUInt64(chi.URLParam(r, "roleID"))
	roReq.UserID = parseUInt64(chi.URLParam(r, "userID"))

	return err
}

var _ RequestFiller = NewRoleMemberRemove()
