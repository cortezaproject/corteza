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
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	sqlxTypes "github.com/jmoiron/sqlx/types"
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
)

type (
	// Internal API interface
	ApplicationList struct {
		// Name GET parameter
		//
		// Application name
		Name string

		// Query GET parameter
		//
		// Filter applications
		Query string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted roles
		Deleted uint

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// Offset GET parameter
		//
		// Offset
		Offset uint

		// Page GET parameter
		//
		// Page number (1-based)
		Page uint

		// PerPage GET parameter
		//
		// Returned items per page (default 50)
		PerPage uint

		// Sort GET parameter
		//
		// Sort items
		Sort string
	}

	ApplicationCreate struct {
		// Name POST parameter
		//
		// Application name
		Name string

		// Enabled POST parameter
		//
		// Enabled
		Enabled bool

		// Unify POST parameter
		//
		// Unify properties
		Unify sqlxTypes.JSONText

		// Config POST parameter
		//
		// Arbitrary JSON holding application configuration
		Config sqlxTypes.JSONText
	}

	ApplicationUpdate struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`

		// Name POST parameter
		//
		// Email
		Name string

		// Enabled POST parameter
		//
		// Enabled
		Enabled bool

		// Unify POST parameter
		//
		// Unify properties
		Unify sqlxTypes.JSONText

		// Config POST parameter
		//
		// Arbitrary JSON holding application configuration
		Config sqlxTypes.JSONText
	}

	ApplicationRead struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`
	}

	ApplicationDelete struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`
	}

	ApplicationUndelete struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`
	}

	ApplicationTriggerScript struct {
		// ApplicationID PATH parameter
		//
		// ID
		ApplicationID uint64 `json:",string"`

		// Script POST parameter
		//
		// Script to execute
		Script string
	}
)

// NewApplicationList request
func NewApplicationList() *ApplicationList {
	return &ApplicationList{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":    r.Name,
		"query":   r.Query,
		"deleted": r.Deleted,
		"limit":   r.Limit,
		"offset":  r.Offset,
		"page":    r.Page,
		"perPage": r.PerPage,
		"sort":    r.Sort,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetDeleted() uint {
	return r.Deleted
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetOffset() uint {
	return r.Offset
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetPage() uint {
	return r.Page
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetPerPage() uint {
	return r.PerPage
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetSort() string {
	return r.Sort
}

// Fill processes request and fills internal variables
func (r *ApplicationList) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
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
		if val, ok := tmp["offset"]; ok && len(val) > 0 {
			r.Offset, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["page"]; ok && len(val) > 0 {
			r.Page, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["perPage"]; ok && len(val) > 0 {
			r.PerPage, err = payload.ParseUint(val[0]), nil
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

// NewApplicationCreate request
func NewApplicationCreate() *ApplicationCreate {
	return &ApplicationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":    r.Name,
		"enabled": r.Enabled,
		"unify":   r.Unify,
		"config":  r.Config,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetUnify() sqlxTypes.JSONText {
	return r.Unify
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Fill processes request and fills internal variables
func (r *ApplicationCreate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["unify"]; ok && len(val) > 0 {
			r.Unify, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewApplicationUpdate request
func NewApplicationUpdate() *ApplicationUpdate {
	return &ApplicationUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
		"name":          r.Name,
		"enabled":       r.Enabled,
		"unify":         r.Unify,
		"config":        r.Config,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetName() string {
	return r.Name
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetEnabled() bool {
	return r.Enabled
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetUnify() sqlxTypes.JSONText {
	return r.Unify
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Fill processes request and fills internal variables
func (r *ApplicationUpdate) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["name"]; ok && len(val) > 0 {
			r.Name, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["enabled"]; ok && len(val) > 0 {
			r.Enabled, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["unify"]; ok && len(val) > 0 {
			r.Unify, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["config"]; ok && len(val) > 0 {
			r.Config, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "applicationID")
		r.ApplicationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApplicationRead request
func NewApplicationRead() *ApplicationRead {
	return &ApplicationRead{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationRead) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Fill processes request and fills internal variables
func (r *ApplicationRead) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "applicationID")
		r.ApplicationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApplicationDelete request
func NewApplicationDelete() *ApplicationDelete {
	return &ApplicationDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationDelete) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Fill processes request and fills internal variables
func (r *ApplicationDelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "applicationID")
		r.ApplicationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApplicationUndelete request
func NewApplicationUndelete() *ApplicationUndelete {
	return &ApplicationUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUndelete) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Fill processes request and fills internal variables
func (r *ApplicationUndelete) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "applicationID")
		r.ApplicationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApplicationTriggerScript request
func NewApplicationTriggerScript() *ApplicationTriggerScript {
	return &ApplicationTriggerScript{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationTriggerScript) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
		"script":        r.Script,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationTriggerScript) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationTriggerScript) GetScript() string {
	return r.Script
}

// Fill processes request and fills internal variables
func (r *ApplicationTriggerScript) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["script"]; ok && len(val) > 0 {
			r.Script, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "applicationID")
		r.ApplicationID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
