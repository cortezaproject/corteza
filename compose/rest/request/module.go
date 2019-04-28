package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
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

	"github.com/crusttech/crust/compose/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"time"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Module list request parameters
type ModuleList struct {
	Query       string
	Page        uint
	PerPage     uint
	NamespaceID uint64 `json:",string"`
}

func NewModuleList() *ModuleList {
	return &ModuleList{}
}

func (mReq *ModuleList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

		mReq.Query = val
	}
	if val, ok := get["page"]; ok {

		mReq.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {

		mReq.PerPage = parseUint(val)
	}
	mReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleList()

// Module create request parameters
type ModuleCreate struct {
	Name        string
	Fields      types.ModuleFieldSet
	Meta        sqlxTypes.JSONText
	NamespaceID uint64 `json:",string"`
}

func NewModuleCreate() *ModuleCreate {
	return &ModuleCreate{}
}

func (mReq *ModuleCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

		mReq.Name = val
	}
	if val, ok := post["meta"]; ok {

		if mReq.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	mReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleCreate()

// Module read request parameters
type ModuleRead struct {
	ModuleID    uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewModuleRead() *ModuleRead {
	return &ModuleRead{}
}

func (mReq *ModuleRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	mReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleRead()

// Module update request parameters
type ModuleUpdate struct {
	ModuleID    uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	Name        string
	Fields      types.ModuleFieldSet
	Meta        sqlxTypes.JSONText
	UpdatedAt   *time.Time
}

func NewModuleUpdate() *ModuleUpdate {
	return &ModuleUpdate{}
}

func (mReq *ModuleUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	mReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))
	if val, ok := post["name"]; ok {

		mReq.Name = val
	}
	if val, ok := post["meta"]; ok {

		if mReq.Meta, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["updatedAt"]; ok {

		if mReq.UpdatedAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewModuleUpdate()

// Module delete request parameters
type ModuleDelete struct {
	ModuleID    uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewModuleDelete() *ModuleDelete {
	return &ModuleDelete{}
}

func (mReq *ModuleDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(mReq)

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

	mReq.ModuleID = parseUInt64(chi.URLParam(r, "moduleID"))
	mReq.NamespaceID = parseUInt64(chi.URLParam(r, "namespaceID"))

	return err
}

var _ RequestFiller = NewModuleDelete()
