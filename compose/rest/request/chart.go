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
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Chart list request parameters
type ChartList struct {
}

func NewChartList() *ChartList {
	return &ChartList{}
}

func (cReq *ChartList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

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

	return err
}

var _ RequestFiller = NewChartList()

// Chart create request parameters
type ChartCreate struct {
	Config sqlxTypes.JSONText
	Name   string
}

func NewChartCreate() *ChartCreate {
	return &ChartCreate{}
}

func (cReq *ChartCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

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

	if val, ok := post["config"]; ok {

		if cReq.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["name"]; ok {

		cReq.Name = val
	}

	return err
}

var _ RequestFiller = NewChartCreate()

// Chart read request parameters
type ChartRead struct {
	ChartID uint64 `json:",string"`
}

func NewChartRead() *ChartRead {
	return &ChartRead{}
}

func (cReq *ChartRead) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

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

	cReq.ChartID = parseUInt64(chi.URLParam(r, "chartID"))

	return err
}

var _ RequestFiller = NewChartRead()

// Chart update request parameters
type ChartUpdate struct {
	ChartID uint64 `json:",string"`
	Config  sqlxTypes.JSONText
	Name    string
}

func NewChartUpdate() *ChartUpdate {
	return &ChartUpdate{}
}

func (cReq *ChartUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

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

	cReq.ChartID = parseUInt64(chi.URLParam(r, "chartID"))
	if val, ok := post["config"]; ok {

		if cReq.Config, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["name"]; ok {

		cReq.Name = val
	}

	return err
}

var _ RequestFiller = NewChartUpdate()

// Chart delete request parameters
type ChartDelete struct {
	ChartID uint64 `json:",string"`
}

func NewChartDelete() *ChartDelete {
	return &ChartDelete{}
}

func (cReq *ChartDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(cReq)

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

	cReq.ChartID = parseUInt64(chi.URLParam(r, "chartID"))

	return err
}

var _ RequestFiller = NewChartDelete()
