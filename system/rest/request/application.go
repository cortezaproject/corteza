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
	"github.com/cortezaproject/corteza-server/pkg/label"
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

		// Labels GET parameter
		//
		// Labels
		Labels map[string]string

		// Flags GET parameter
		//
		// Flags
		Flags []string

		// IncFlags GET parameter
		//
		// Calculated (0, default), global (1) or return only (2) own flags
		IncFlags uint

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

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

		// Weight POST parameter
		//
		// Weight for sorting
		Weight int

		// Unify POST parameter
		//
		// Unify properties
		Unify sqlxTypes.JSONText

		// Config POST parameter
		//
		// Arbitrary JSON holding application configuration
		Config sqlxTypes.JSONText

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
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

		// Weight POST parameter
		//
		// Weight for sorting
		Weight int

		// Unify POST parameter
		//
		// Unify properties
		Unify sqlxTypes.JSONText

		// Config POST parameter
		//
		// Arbitrary JSON holding application configuration
		Config sqlxTypes.JSONText

		// Labels POST parameter
		//
		// Labels
		Labels map[string]string
	}

	ApplicationFlagCreate struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`

		// Flag PATH parameter
		//
		// Flag
		Flag string

		// OwnedBy PATH parameter
		//
		// Owner; 0 = everyone
		OwnedBy uint64 `json:",string"`
	}

	ApplicationFlagDelete struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`

		// Flag PATH parameter
		//
		// Flag
		Flag string

		// OwnedBy PATH parameter
		//
		// Owner; 0 = everyone
		OwnedBy uint64 `json:",string"`
	}

	ApplicationRead struct {
		// ApplicationID PATH parameter
		//
		// Application ID
		ApplicationID uint64 `json:",string"`

		// IncFlags GET parameter
		//
		// Calculated (0, default), global (1) or return only (2) own flags
		IncFlags uint
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

	ApplicationReorder struct {
		// ApplicationIDs POST parameter
		//
		// Application order
		ApplicationIDs []string
	}
)

// NewApplicationList request
func NewApplicationList() *ApplicationList {
	return &ApplicationList{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":       r.Name,
		"query":      r.Query,
		"deleted":    r.Deleted,
		"labels":     r.Labels,
		"flags":      r.Flags,
		"incFlags":   r.IncFlags,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
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
func (r ApplicationList) GetLabels() map[string]string {
	return r.Labels
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetFlags() []string {
	return r.Flags
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetIncFlags() uint {
	return r.IncFlags
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationList) GetPageCursor() string {
	return r.PageCursor
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
		if val, ok := tmp["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := tmp["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["flags[]"]; ok {
			r.Flags, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["flags"]; ok {
			r.Flags, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["incFlags"]; ok && len(val) > 0 {
			r.IncFlags, err = payload.ParseUint(val[0]), nil
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

// NewApplicationCreate request
func NewApplicationCreate() *ApplicationCreate {
	return &ApplicationCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"name":    r.Name,
		"enabled": r.Enabled,
		"weight":  r.Weight,
		"unify":   r.Unify,
		"config":  r.Config,
		"labels":  r.Labels,
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
func (r ApplicationCreate) GetWeight() int {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetUnify() sqlxTypes.JSONText {
	return r.Unify
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationCreate) GetLabels() map[string]string {
	return r.Labels
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

		if val, ok := req.Form["weight"]; ok && len(val) > 0 {
			r.Weight, err = payload.ParseInt(val[0]), nil
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

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
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
		"weight":        r.Weight,
		"unify":         r.Unify,
		"config":        r.Config,
		"labels":        r.Labels,
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
func (r ApplicationUpdate) GetWeight() int {
	return r.Weight
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetUnify() sqlxTypes.JSONText {
	return r.Unify
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetConfig() sqlxTypes.JSONText {
	return r.Config
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationUpdate) GetLabels() map[string]string {
	return r.Labels
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

		if val, ok := req.Form["weight"]; ok && len(val) > 0 {
			r.Weight, err = payload.ParseInt(val[0]), nil
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

		if val, ok := req.Form["labels[]"]; ok {
			r.Labels, err = label.ParseStrings(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["labels"]; ok {
			r.Labels, err = label.ParseStrings(val)
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

// NewApplicationFlagCreate request
func NewApplicationFlagCreate() *ApplicationFlagCreate {
	return &ApplicationFlagCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
		"flag":          r.Flag,
		"ownedBy":       r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagCreate) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagCreate) GetFlag() string {
	return r.Flag
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagCreate) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *ApplicationFlagCreate) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "flag")
		r.Flag, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "ownedBy")
		r.OwnedBy, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewApplicationFlagDelete request
func NewApplicationFlagDelete() *ApplicationFlagDelete {
	return &ApplicationFlagDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationID": r.ApplicationID,
		"flag":          r.Flag,
		"ownedBy":       r.OwnedBy,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagDelete) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagDelete) GetFlag() string {
	return r.Flag
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationFlagDelete) GetOwnedBy() uint64 {
	return r.OwnedBy
}

// Fill processes request and fills internal variables
func (r *ApplicationFlagDelete) Fill(req *http.Request) (err error) {
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

		val = chi.URLParam(req, "flag")
		r.Flag, err = val, nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "ownedBy")
		r.OwnedBy, err = payload.ParseUint64(val), nil
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
		"incFlags":      r.IncFlags,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationRead) GetApplicationID() uint64 {
	return r.ApplicationID
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationRead) GetIncFlags() uint {
	return r.IncFlags
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
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["incFlags"]; ok && len(val) > 0 {
			r.IncFlags, err = payload.ParseUint(val[0]), nil
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

// NewApplicationReorder request
func NewApplicationReorder() *ApplicationReorder {
	return &ApplicationReorder{}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationReorder) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"applicationIDs": r.ApplicationIDs,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ApplicationReorder) GetApplicationIDs() []string {
	return r.ApplicationIDs
}

// Fill processes request and fills internal variables
func (r *ApplicationReorder) Fill(req *http.Request) (err error) {
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

		//if val, ok := req.Form["applicationIDs[]"]; ok && len(val) > 0  {
		//    r.ApplicationIDs, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}
	}

	return err
}
