package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `team.go`, `team.util.go` or `team_test.go` to
	implement your API calls, helper functions and tests. The file `team.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}
var _ = multipart.FileHeader{}

// Team list request parameters
type TeamList struct {
	Query string
}

func NewTeamList() *TeamList {
	return &TeamList{}
}

func (t *TeamList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

		t.Query = val
	}

	return err
}

var _ RequestFiller = NewTeamList()

// Team create request parameters
type TeamCreate struct {
	Name    string
	Members []uint64 `json:",string"`
}

func NewTeamCreate() *TeamCreate {
	return &TeamCreate{}
}

func (t *TeamCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

		t.Name = val
	}
	t.Members = parseUInt64A(r.Form["members"])

	return err
}

var _ RequestFiller = NewTeamCreate()

// Team edit request parameters
type TeamEdit struct {
	TeamID  uint64 `json:",string"`
	Name    string
	Members []uint64 `json:",string"`
}

func NewTeamEdit() *TeamEdit {
	return &TeamEdit{}
}

func (t *TeamEdit) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))
	if val, ok := post["name"]; ok {

		t.Name = val
	}
	t.Members = parseUInt64A(r.Form["members"])

	return err
}

var _ RequestFiller = NewTeamEdit()

// Team read request parameters
type TeamRead struct {
	TeamID uint64 `json:",string"`
}

func NewTeamRead() *TeamRead {
	return &TeamRead{}
}

func (t *TeamRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))

	return err
}

var _ RequestFiller = NewTeamRead()

// Team remove request parameters
type TeamRemove struct {
	TeamID uint64 `json:",string"`
}

func NewTeamRemove() *TeamRemove {
	return &TeamRemove{}
}

func (t *TeamRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))

	return err
}

var _ RequestFiller = NewTeamRemove()

// Team archive request parameters
type TeamArchive struct {
	TeamID uint64 `json:",string"`
}

func NewTeamArchive() *TeamArchive {
	return &TeamArchive{}
}

func (t *TeamArchive) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))

	return err
}

var _ RequestFiller = NewTeamArchive()

// Team move request parameters
type TeamMove struct {
	TeamID         uint64 `json:",string"`
	OrganisationID uint64 `json:",string"`
}

func NewTeamMove() *TeamMove {
	return &TeamMove{}
}

func (t *TeamMove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))
	if val, ok := post["organisationID"]; ok {

		t.OrganisationID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewTeamMove()

// Team merge request parameters
type TeamMerge struct {
	TeamID      uint64 `json:",string"`
	Destination uint64 `json:",string"`
}

func NewTeamMerge() *TeamMerge {
	return &TeamMerge{}
}

func (t *TeamMerge) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))
	if val, ok := post["destination"]; ok {

		t.Destination = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewTeamMerge()

// Team memberAdd request parameters
type TeamMemberAdd struct {
	TeamID uint64 `json:",string"`
	UserID uint64 `json:",string"`
}

func NewTeamMemberAdd() *TeamMemberAdd {
	return &TeamMemberAdd{}
}

func (t *TeamMemberAdd) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))
	if val, ok := post["userID"]; ok {

		t.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewTeamMemberAdd()

// Team memberRemove request parameters
type TeamMemberRemove struct {
	TeamID uint64 `json:",string"`
	UserID uint64 `json:",string"`
}

func NewTeamMemberRemove() *TeamMemberRemove {
	return &TeamMemberRemove{}
}

func (t *TeamMemberRemove) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

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

	t.TeamID = parseUInt64(chi.URLParam(r, "teamID"))
	if val, ok := post["userID"]; ok {

		t.UserID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewTeamMemberRemove()
