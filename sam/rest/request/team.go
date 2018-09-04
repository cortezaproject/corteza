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
	"net/http"
	"strings"
)

var _ = chi.URLParam
var _ = types.JSONText{}

// Team list request parameters
type TeamList struct {
	Query string
}

func NewTeamList() *TeamList {
	return &TeamList{}
}

func (t *TeamList) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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
	Members []uint64
}

func NewTeamCreate() *TeamCreate {
	return &TeamCreate{}
}

func (t *TeamCreate) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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

	return err
}

var _ RequestFiller = NewTeamCreate()

// Team edit request parameters
type TeamEdit struct {
	TeamID  uint64
	Name    string
	Members []uint64
}

func NewTeamEdit() *TeamEdit {
	return &TeamEdit{}
}

func (t *TeamEdit) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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

	return err
}

var _ RequestFiller = NewTeamEdit()

// Team read request parameters
type TeamRead struct {
	TeamID uint64
}

func NewTeamRead() *TeamRead {
	return &TeamRead{}
}

func (t *TeamRead) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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
	TeamID uint64
}

func NewTeamRemove() *TeamRemove {
	return &TeamRemove{}
}

func (t *TeamRemove) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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
	TeamID uint64
}

func NewTeamArchive() *TeamArchive {
	return &TeamArchive{}
}

func (t *TeamArchive) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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
	TeamID          uint64
	Organisation_id uint64
}

func NewTeamMove() *TeamMove {
	return &TeamMove{}
}

func (t *TeamMove) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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
	if val, ok := post["organisation_id"]; ok {

		t.Organisation_id = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewTeamMove()

// Team merge request parameters
type TeamMerge struct {
	TeamID      uint64
	Destination uint64
}

func NewTeamMerge() *TeamMerge {
	return &TeamMerge{}
}

func (t *TeamMerge) Fill(r *http.Request) error {
	var err error

	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(t)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	r.ParseForm()
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
