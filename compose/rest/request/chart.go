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

// Chart list request parameters
type ChartList struct {
	Query       string
	Page        uint
	PerPage     uint
	NamespaceID uint64 `json:",string"`
}

func NewChartList() *ChartList {
	return &ChartList{}
}

func (r ChartList) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["query"] = r.Query
	out["page"] = r.Page
	out["perPage"] = r.PerPage
	out["namespaceID"] = r.NamespaceID

	return out
}

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
		r.Query = val
	}
	if val, ok := get["page"]; ok {
		r.Page = parseUint(val)
	}
	if val, ok := get["perPage"]; ok {
		r.PerPage = parseUint(val)
	}
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartList()

// Chart create request parameters
type ChartCreate struct {
	Config      sqlxTypes.JSONText
	Name        string
	NamespaceID uint64 `json:",string"`
}

func NewChartCreate() *ChartCreate {
	return &ChartCreate{}
}

func (r ChartCreate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["config"] = r.Config
	out["name"] = r.Name
	out["namespaceID"] = r.NamespaceID

	return out
}

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

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartCreate()

// Chart read request parameters
type ChartRead struct {
	ChartID     uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewChartRead() *ChartRead {
	return &ChartRead{}
}

func (r ChartRead) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["chartID"] = r.ChartID
	out["namespaceID"] = r.NamespaceID

	return out
}

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

	r.ChartID = parseUInt64(chi.URLParam(req, "chartID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartRead()

// Chart update request parameters
type ChartUpdate struct {
	ChartID     uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
	Config      sqlxTypes.JSONText
	Name        string
	UpdatedAt   *time.Time
}

func NewChartUpdate() *ChartUpdate {
	return &ChartUpdate{}
}

func (r ChartUpdate) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["chartID"] = r.ChartID
	out["namespaceID"] = r.NamespaceID
	out["config"] = r.Config
	out["name"] = r.Name
	out["updatedAt"] = r.UpdatedAt

	return out
}

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

	r.ChartID = parseUInt64(chi.URLParam(req, "chartID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))
	if val, ok := post["config"]; ok {

		if r.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["name"]; ok {
		r.Name = val
	}
	if val, ok := post["updatedAt"]; ok {

		if r.UpdatedAt, err = parseISODatePtrWithErr(val); err != nil {
			return err
		}
	}

	return err
}

var _ RequestFiller = NewChartUpdate()

// Chart delete request parameters
type ChartDelete struct {
	ChartID     uint64 `json:",string"`
	NamespaceID uint64 `json:",string"`
}

func NewChartDelete() *ChartDelete {
	return &ChartDelete{}
}

func (r ChartDelete) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["chartID"] = r.ChartID
	out["namespaceID"] = r.NamespaceID

	return out
}

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

	r.ChartID = parseUInt64(chi.URLParam(req, "chartID"))
	r.NamespaceID = parseUInt64(chi.URLParam(req, "namespaceID"))

	return err
}

var _ RequestFiller = NewChartDelete()
