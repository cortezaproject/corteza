package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	ProjectTaskList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// ProjectID GET parameter
		//
		// Filter by project ID
		ProjectID uint64 `json:",string"`

		// Status GET parameter
		//
		// Filter by status
		Status string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// IncTotal GET parameter
		//
		// Include total counter
		IncTotal bool

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	ProjectTaskCreate struct {
		// ProjectID POST parameter
		//
		// Project this task belongs to
		ProjectID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// TaskName POST parameter
		//
		// Task name
		TaskName string

		// TaskType POST parameter
		//
		// Task type
		TaskType string

		// Status POST parameter
		//
		// Status
		Status string

		// Severity POST parameter
		//
		// Severity
		Severity string

		// Risk POST parameter
		//
		// Risk
		Risk string

		// Owner POST parameter
		//
		// Owner user ID
		Owner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// Backlog POST parameter
		//
		// Backlog IDs (comma-separated)
		Backlog string

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string

		// CompletedDate POST parameter
		//
		// Completed date (ISO)
		CompletedDate string
	}

	ProjectTaskRead struct {
		// TaskID PATH parameter
		//
		// Task ID
		TaskID uint64 `json:",string"`
	}

	ProjectTaskUpdate struct {
		// TaskID PATH parameter
		//
		// Task ID
		TaskID uint64 `json:",string"`

		// Title POST parameter
		//
		// Title
		Title string

		// Description POST parameter
		//
		// Description
		Description string

		// TaskName POST parameter
		//
		// Task name
		TaskName string

		// TaskType POST parameter
		//
		// Task type
		TaskType string

		// Status POST parameter
		//
		// Status
		Status string

		// Severity POST parameter
		//
		// Severity
		Severity string

		// Risk POST parameter
		//
		// Risk
		Risk string

		// Owner POST parameter
		//
		// Owner user ID
		Owner uint64 `json:",string"`

		// ChangeOwner POST parameter
		//
		// Change owner user ID
		ChangeOwner uint64 `json:",string"`

		// Backlog POST parameter
		//
		// Backlog IDs (comma-separated)
		Backlog string

		// DateDue POST parameter
		//
		// Due date (ISO)
		DateDue string

		// CompletedDate POST parameter
		//
		// Completed date (ISO)
		CompletedDate string
	}

	ProjectTaskDelete struct {
		// TaskID PATH parameter
		//
		// Task ID
		TaskID uint64 `json:",string"`
	}
)

// NewProjectTaskList request
func NewProjectTaskList() *ProjectTaskList {
	return &ProjectTaskList{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"projectID":  r.ProjectID,
		"status":     r.Status,
		"limit":      r.Limit,
		"incTotal":   r.IncTotal,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetIncTotal() bool {
	return r.IncTotal
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ProjectTaskList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incTotal"]; ok && len(val) > 0 {
			r.IncTotal, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewProjectTaskCreate request
func NewProjectTaskCreate() *ProjectTaskCreate {
	return &ProjectTaskCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"projectID":     r.ProjectID,
		"title":         r.Title,
		"description":   r.Description,
		"taskName":      r.TaskName,
		"taskType":      r.TaskType,
		"status":        r.Status,
		"severity":      r.Severity,
		"risk":          r.Risk,
		"owner":         r.Owner,
		"changeOwner":   r.ChangeOwner,
		"backlog":       r.Backlog,
		"dateDue":       r.DateDue,
		"completedDate": r.CompletedDate,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetProjectID() uint64 {
	return r.ProjectID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetTaskName() string {
	return r.TaskName
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetTaskType() string {
	return r.TaskType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetOwner() uint64 {
	return r.Owner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetBacklog() string {
	return r.Backlog
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetDateDue() string {
	return r.DateDue
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskCreate) GetCompletedDate() string {
	return r.CompletedDate
}

// Fill processes request and fills internal variables
func (r *ProjectTaskCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["projectID"]; ok && len(val) > 0 {
				r.ProjectID, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["taskName"]; ok && len(val) > 0 {
				r.TaskName, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["taskType"]; ok && len(val) > 0 {
				r.TaskType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["severity"]; ok && len(val) > 0 {
				r.Severity, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["risk"]; ok && len(val) > 0 {
				r.Risk, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["owner"]; ok && len(val) > 0 {
				r.Owner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeOwner"]; ok && len(val) > 0 {
				r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["backlog"]; ok && len(val) > 0 {
				r.Backlog, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["dateDue"]; ok && len(val) > 0 {
				r.DateDue, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["completedDate"]; ok && len(val) > 0 {
				r.CompletedDate, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["projectID"]; ok && len(val) > 0 {
			r.ProjectID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["taskName"]; ok && len(val) > 0 {
			r.TaskName, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["taskType"]; ok && len(val) > 0 {
			r.TaskType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["risk"]; ok && len(val) > 0 {
			r.Risk, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["owner"]; ok && len(val) > 0 {
			r.Owner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeOwner"]; ok && len(val) > 0 {
			r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["backlog"]; ok && len(val) > 0 {
			r.Backlog, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["dateDue"]; ok && len(val) > 0 {
			r.DateDue, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["completedDate"]; ok && len(val) > 0 {
			r.CompletedDate, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewProjectTaskRead request
func NewProjectTaskRead() *ProjectTaskRead {
	return &ProjectTaskRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"taskID": r.TaskID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskRead) GetTaskID() uint64 {
	return r.TaskID
}

// Fill processes request and fills internal variables
func (r *ProjectTaskRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "taskID")
		r.TaskID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectTaskUpdate request
func NewProjectTaskUpdate() *ProjectTaskUpdate {
	return &ProjectTaskUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"taskID":        r.TaskID,
		"title":         r.Title,
		"description":   r.Description,
		"taskName":      r.TaskName,
		"taskType":      r.TaskType,
		"status":        r.Status,
		"severity":      r.Severity,
		"risk":          r.Risk,
		"owner":         r.Owner,
		"changeOwner":   r.ChangeOwner,
		"backlog":       r.Backlog,
		"dateDue":       r.DateDue,
		"completedDate": r.CompletedDate,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetTaskID() uint64 {
	return r.TaskID
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetTitle() string {
	return r.Title
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetDescription() string {
	return r.Description
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetTaskName() string {
	return r.TaskName
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetTaskType() string {
	return r.TaskType
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetStatus() string {
	return r.Status
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetSeverity() string {
	return r.Severity
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetRisk() string {
	return r.Risk
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetOwner() uint64 {
	return r.Owner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetChangeOwner() uint64 {
	return r.ChangeOwner
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetBacklog() string {
	return r.Backlog
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetDateDue() string {
	return r.DateDue
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskUpdate) GetCompletedDate() string {
	return r.CompletedDate
}

// Fill processes request and fills internal variables
func (r *ProjectTaskUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["title"]; ok && len(val) > 0 {
				r.Title, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["description"]; ok && len(val) > 0 {
				r.Description, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["taskName"]; ok && len(val) > 0 {
				r.TaskName, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["taskType"]; ok && len(val) > 0 {
				r.TaskType, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["status"]; ok && len(val) > 0 {
				r.Status, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["severity"]; ok && len(val) > 0 {
				r.Severity, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["risk"]; ok && len(val) > 0 {
				r.Risk, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["owner"]; ok && len(val) > 0 {
				r.Owner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["changeOwner"]; ok && len(val) > 0 {
				r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["backlog"]; ok && len(val) > 0 {
				r.Backlog, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["dateDue"]; ok && len(val) > 0 {
				r.DateDue, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["completedDate"]; ok && len(val) > 0 {
				r.CompletedDate, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["title"]; ok && len(val) > 0 {
			r.Title, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["description"]; ok && len(val) > 0 {
			r.Description, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["taskName"]; ok && len(val) > 0 {
			r.TaskName, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["taskType"]; ok && len(val) > 0 {
			r.TaskType, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["status"]; ok && len(val) > 0 {
			r.Status, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["severity"]; ok && len(val) > 0 {
			r.Severity, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["risk"]; ok && len(val) > 0 {
			r.Risk, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["owner"]; ok && len(val) > 0 {
			r.Owner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["changeOwner"]; ok && len(val) > 0 {
			r.ChangeOwner, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["backlog"]; ok && len(val) > 0 {
			r.Backlog, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["dateDue"]; ok && len(val) > 0 {
			r.DateDue, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["completedDate"]; ok && len(val) > 0 {
			r.CompletedDate, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "taskID")
		r.TaskID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewProjectTaskDelete request
func NewProjectTaskDelete() *ProjectTaskDelete {
	return &ProjectTaskDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"taskID": r.TaskID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ProjectTaskDelete) GetTaskID() uint64 {
	return r.TaskID
}

// Fill processes request and fills internal variables
func (r *ProjectTaskDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "taskID")
		r.TaskID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
