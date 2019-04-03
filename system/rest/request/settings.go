package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `settings.go`, `settings.util.go` or `settings_test.go` to
	implement your API calls, helper functions and tests. The file `settings.go`
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
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Settings list request parameters
type SettingsList struct {
	Prefix string
}

func NewSettingsList() *SettingsList {
	return &SettingsList{}
}

func (seReq *SettingsList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(seReq)

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

	if val, ok := get["prefix"]; ok {

		seReq.Prefix = val
	}

	return err
}

var _ RequestFiller = NewSettingsList()

// Settings update request parameters
type SettingsUpdate struct {
	Values sqlxTypes.JSONText
}

func NewSettingsUpdate() *SettingsUpdate {
	return &SettingsUpdate{}
}

func (seReq *SettingsUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(seReq)

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

	if val, ok := post["values"]; ok {

		if seReq.Values, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewSettingsUpdate()

// Settings get request parameters
type SettingsGet struct {
	OwnerID uint64 `json:",string"`
	Key     string
}

func NewSettingsGet() *SettingsGet {
	return &SettingsGet{}
}

func (seReq *SettingsGet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(seReq)

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

	if val, ok := get["ownerID"]; ok {

		seReq.OwnerID = parseUInt64(val)
	}
	seReq.Key = chi.URLParam(r, "key")

	return err
}

var _ RequestFiller = NewSettingsGet()

// Settings set request parameters
type SettingsSet struct {
	Key     string
	OwnerID uint64 `json:",string"`
	Value   sqlxTypes.JSONText
}

func NewSettingsSet() *SettingsSet {
	return &SettingsSet{}
}

func (seReq *SettingsSet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(seReq)

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

	seReq.Key = chi.URLParam(r, "key")
	if val, ok := post["ownerID"]; ok {

		seReq.OwnerID = parseUInt64(val)
	}
	if val, ok := post["value"]; ok {

		if seReq.Value, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewSettingsSet()
