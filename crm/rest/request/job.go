package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `job.go`, `job.util.go` or `job_test.go` to
	implement your API calls, helper functions and tests. The file `job.go`
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
	sqlxTypes "github.com/jmoiron/sqlx/types"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Job list request parameters
type JobList struct {
	Status  string
	Page    int
	PerPage int
}

func NewJobList() *JobList {
	return &JobList{}
}

func (j *JobList) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(j)

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

	if val, ok := get["status"]; ok {

		j.Status = val
	}
	if val, ok := get["page"]; ok {

		j.Page = parseInt(val)
	}
	if val, ok := get["perPage"]; ok {

		j.PerPage = parseInt(val)
	}

	return err
}

var _ RequestFiller = NewJobList()

// Job run request parameters
type JobRun struct {
	WorkflowID string
	StartAt    string
	Parameters types.JobParameterSet
}

func NewJobRun() *JobRun {
	return &JobRun{}
}

func (j *JobRun) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(j)

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

	if val, ok := post["workflowID"]; ok {

		j.WorkflowID = val
	}
	if val, ok := post["startAt"]; ok {

		j.StartAt = val
	}

	return err
}

var _ RequestFiller = NewJobRun()

// Job get request parameters
type JobGet struct {
	JobID string
}

func NewJobGet() *JobGet {
	return &JobGet{}
}

func (j *JobGet) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(j)

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

	j.JobID = chi.URLParam(r, "jobID")

	return err
}

var _ RequestFiller = NewJobGet()

// Job logs request parameters
type JobLogs struct {
	JobID   string
	Page    int
	PerPage int
}

func NewJobLogs() *JobLogs {
	return &JobLogs{}
}

func (j *JobLogs) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(j)

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

	j.JobID = chi.URLParam(r, "jobID")
	j.Page = parseInt(chi.URLParam(r, "page"))
	j.PerPage = parseInt(chi.URLParam(r, "perPage"))

	return err
}

var _ RequestFiller = NewJobLogs()

// Job update request parameters
type JobUpdate struct {
	JobID      string
	Status     string
	Log        sqlxTypes.JSONText
	WorkflowID string
	StartAt    string
	Parameters types.JobParameterSet
}

func NewJobUpdate() *JobUpdate {
	return &JobUpdate{}
}

func (j *JobUpdate) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(j)

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

	j.JobID = chi.URLParam(r, "jobID")
	if val, ok := post["status"]; ok {

		j.Status = val
	}
	if val, ok := post["log"]; ok {

		if j.Log, err = parseJSONTextWithErr(val); err != nil {
			return err
		}
	}
	if val, ok := post["workflowID"]; ok {

		j.WorkflowID = val
	}
	if val, ok := post["startAt"]; ok {

		j.StartAt = val
	}

	return err
}

var _ RequestFiller = NewJobUpdate()

// Job delete request parameters
type JobDelete struct {
	JobID string
}

func NewJobDelete() *JobDelete {
	return &JobDelete{}
}

func (j *JobDelete) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(j)

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

	j.JobID = chi.URLParam(r, "jobID")

	return err
}

var _ RequestFiller = NewJobDelete()
