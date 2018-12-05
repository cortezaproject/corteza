package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `workflow.go`, `workflow.util.go` or `workflow_test.go` to
	implement your API calls, helper functions and tests. The file `workflow.go`
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

	"github.com/crusttech/crust/crm/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Workflow list request parameters
type WorkflowList struct {
}

func NewWorkflowList() *WorkflowList {
	return &WorkflowList{}
}

func (w *WorkflowList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(w)

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

var _ RequestFiller = NewWorkflowList()

// Workflow create request parameters
type WorkflowCreate struct {
	Name    string
	Tasks   types.WorkflowTaskSet
	OnError types.WorkflowTaskSet
	Timeout int
}

func NewWorkflowCreate() *WorkflowCreate {
	return &WorkflowCreate{}
}

func (w *WorkflowCreate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(w)

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

		w.Name = val
	}
	if val, ok := post["timeout"]; ok {

		w.Timeout = parseInt(val)
	}

	return err
}

var _ RequestFiller = NewWorkflowCreate()

// Workflow get request parameters
type WorkflowGet struct {
	WorkflowID string
}

func NewWorkflowGet() *WorkflowGet {
	return &WorkflowGet{}
}

func (w *WorkflowGet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(w)

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

	w.WorkflowID = chi.URLParam(r, "workflowID")

	return err
}

var _ RequestFiller = NewWorkflowGet()

// Workflow update request parameters
type WorkflowUpdate struct {
	WorkflowID string
	Name       string
	Tasks      types.WorkflowTaskSet
	OnError    types.WorkflowTaskSet
	Timeout    int
}

func NewWorkflowUpdate() *WorkflowUpdate {
	return &WorkflowUpdate{}
}

func (w *WorkflowUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(w)

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

	w.WorkflowID = chi.URLParam(r, "workflowID")
	if val, ok := post["name"]; ok {

		w.Name = val
	}
	if val, ok := post["timeout"]; ok {

		w.Timeout = parseInt(val)
	}

	return err
}

var _ RequestFiller = NewWorkflowUpdate()

// Workflow delete request parameters
type WorkflowDelete struct {
	WorkflowID string
}

func NewWorkflowDelete() *WorkflowDelete {
	return &WorkflowDelete{}
}

func (w *WorkflowDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(w)

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

	w.WorkflowID = chi.URLParam(r, "workflowID")

	return err
}

var _ RequestFiller = NewWorkflowDelete()
