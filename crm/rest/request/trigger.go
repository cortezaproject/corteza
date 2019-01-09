package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `trigger.go`, `trigger.util.go` or `trigger_test.go` to
	implement your API calls, helper functions and tests. The file `trigger.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Trigger list request parameters
type TriggerList struct {
	ModuleID uint64 `json:",string"`
}

func NewTriggerList() *TriggerList {
	return &TriggerList{}
}

func (t *TriggerList) Fill(r *http.Request) (err error) {
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

	if val, ok := get["moduleID"]; ok {

		t.ModuleID = parseUInt64(val)
	}

	return err
}

var _ RequestFiller = NewTriggerList()

// Trigger create request parameters
type TriggerCreate struct {
	ModuleID uint64 `json:",string"`
	Name     string
	Actions  []string
	Enabled  bool
	Source   string
}

func NewTriggerCreate() *TriggerCreate {
	return &TriggerCreate{}
}

func (t *TriggerCreate) Fill(r *http.Request) (err error) {
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

	if val, ok := post["moduleID"]; ok {

		t.ModuleID = parseUInt64(val)
	}
	if val, ok := post["name"]; ok {

		t.Name = val
	}
	if val, ok := post["enabled"]; ok {

		t.Enabled = parseBool(val)
	}
	if val, ok := post["source"]; ok {

		t.Source = val
	}

	return err
}

var _ RequestFiller = NewTriggerCreate()

// Trigger read request parameters
type TriggerRead struct {
	TriggerID uint64 `json:",string"`
}

func NewTriggerRead() *TriggerRead {
	return &TriggerRead{}
}

func (t *TriggerRead) Fill(r *http.Request) (err error) {
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

	t.TriggerID = parseUInt64(chi.URLParam(r, "triggerID"))

	return err
}

var _ RequestFiller = NewTriggerRead()

// Trigger edit request parameters
type TriggerEdit struct {
	TriggerID uint64 `json:",string"`
	ModuleID  uint64 `json:",string"`
	Name      string
	Actions   []string
	Enabled   bool
	Source    string
}

func NewTriggerEdit() *TriggerEdit {
	return &TriggerEdit{}
}

func (t *TriggerEdit) Fill(r *http.Request) (err error) {
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

	t.TriggerID = parseUInt64(chi.URLParam(r, "triggerID"))
	if val, ok := post["moduleID"]; ok {

		t.ModuleID = parseUInt64(val)
	}
	if val, ok := post["name"]; ok {

		t.Name = val
	}
	if val, ok := post["enabled"]; ok {

		t.Enabled = parseBool(val)
	}
	if val, ok := post["source"]; ok {

		t.Source = val
	}

	return err
}

var _ RequestFiller = NewTriggerEdit()

// Trigger delete request parameters
type TriggerDelete struct {
	TriggerID uint64 `json:",string"`
}

func NewTriggerDelete() *TriggerDelete {
	return &TriggerDelete{}
}

func (t *TriggerDelete) Fill(r *http.Request) (err error) {
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

	t.TriggerID = parseUInt64(chi.URLParam(r, "triggerID"))

	return err
}

var _ RequestFiller = NewTriggerDelete()
